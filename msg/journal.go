//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package msg

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

//=============================================================================

type LogStatus int

const (
	StatusPending LogStatus = 0
	StatusAcked   LogStatus = 1
)

const JournalFile = "journal.bin"
const JournalTemp = "journal.temp"

//=============================================================================

type JournalEntry struct {
	Id        string
	Timestamp time.Time
	Status    LogStatus
	Exchange  string
	Payload   []byte
}

//=============================================================================

type Journal struct {
	dir     string
	file    *os.File
	writer  *bufio.Writer
	encoder *gob.Encoder
	mu      sync.Mutex
}

//=============================================================================

func NewJournal(journalDir string) (*Journal, error) {
	j := &Journal{
		dir : journalDir,
	}

	return j.init()
}

//=============================================================================

func (j *Journal) init() (*Journal,error) {
	var err error
	j.file, err = j.createFile(JournalFile)
	if err == nil {
		j.writer  = bufio.NewWriter(j.file)
		j.encoder = gob.NewEncoder(j.writer)
	}
	return j, err
}

//=============================================================================
// Write appends a pending message to the log and forces a disk sync

func (j *Journal) Write(id string, payload []byte, exchange string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	entry := &JournalEntry{
		Id       : id,
		Timestamp: time.Now(),
		Status   : StatusPending,
		Exchange : exchange,
		Payload  : payload,
	}

	return write(j.file, j.writer, j.encoder, entry, true)
}

//=============================================================================
// Ack appends a confirmation entry to the log

func (j *Journal) Ack(id string) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	entry := &JournalEntry{
		Id:        id,
		Status:    StatusAcked,
		Timestamp: time.Now(),
	}

	return write(j.file, j.writer, j.encoder, entry, true)
}

//=============================================================================

func (j *Journal) Close() error {
	err := j.writer.Flush()
	if err == nil {
		err = j.file.Close()
	}
	return err
}

//=============================================================================
// Recover parses the journal file and returns all unacknowledged messages

func (j *Journal) Recover() ([]*JournalEntry, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	if _, err := j.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	pendingMap := make(map[string]*JournalEntry)
	reader     := bufio.NewReader(j.file)
	decoder    := gob.NewDecoder(reader)

	for {
		var entry JournalEntry
		if err := decoder.Decode(&entry); err != nil {
			if errors.Is(err, io.EOF) {
				// We reached the end of the file successfully
				break
			}

			return nil, errors.New("Failed to decode journal entry: "+err.Error())
		}

		if entry.Status == StatusPending {
			pendingMap[entry.Id] = &entry
		} else if entry.Status == StatusAcked {
			delete(pendingMap, entry.Id)
		}
	}

	var unresolved []*JournalEntry
	for _, entry := range pendingMap {
		unresolved = append(unresolved, entry)
	}

	return unresolved, nil
}

//=============================================================================
// Compacts the journal file, removing all acknowledged messages

func (j *Journal) Compact() error {
	j.mu.Lock()
	defer j.mu.Unlock()

	//--- Remove any old file (just in case)

	err := j.deleteFile(JournalTemp)
	if err != nil {
		return err
	}

	//--- Retrieve current not-ack entries

	entries, err := j.Recover()
	if err != nil {
		return err
	}

	//--- Create temp destination

	file, err := j.createFile(JournalTemp)
	if err != nil {
		return err
	}

	writer  := bufio.NewWriter(file)
	encoder := gob.NewEncoder(writer)

	//--- Write entries into new temp file

	for _, entry := range entries {
		err = write(file, writer, encoder, entry, false)
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = j.deleteFile(JournalFile)
	if err != nil {
		return err
	}

	err = j.renameFile(JournalTemp, JournalFile)
	if err != nil {
		panic("Could not rename journal: " + err.Error())
	}

	j.file, err = j.createFile(JournalFile)
	if err != nil {
		panic("Could not recreate journal: " + err.Error())
	}

	j.writer = bufio.NewWriter(j.file)
	j.encoder= gob.NewEncoder(j.writer)

	return nil
}

//=============================================================================

func (j *Journal) createFile(name string) (*os.File, error) {
	return  os.OpenFile(j.dir +"/"+ name, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
}

//=============================================================================

func (j *Journal) deleteFile(name string) error {
	return os.Remove(j.dir +"/"+ name)
}

//=============================================================================

func (j *Journal) renameFile(oldName,newName string) error {
	return os.Rename(j.dir +"/"+ oldName, j.dir +"/"+ newName)
}

//=============================================================================

func write(file *os.File, writer *bufio.Writer, encoder *gob.Encoder, entry *JournalEntry, sync bool) error {
	err := encoder.Encode(entry)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	//--- Force the operating system to flush weights to stable storage
	if sync {
		return file.Sync()
	}

	return nil
}

//=============================================================================
