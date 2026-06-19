//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
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

	return write(j.file, j.writer, entry, true)
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

	return write(j.file, j.writer, entry, true)
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

func (j *Journal) Recover() ([]*JournalEntry, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	return j.recover()
}

//=============================================================================
// Compacts the journal file, removing all acknowledged messages

func (j *Journal) Compact() ([]*JournalEntry, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	//--- Remove any old file (just in case)

	exists,err := j.existsFile(JournalTemp)
	if err != nil {
		return nil, err
	}

	if exists {
		err = j.deleteFile(JournalTemp)
		if err != nil {
			return nil, err
		}
	}

	//--- Retrieve current not-ack entries

	entries, err := j.recover()
	if err != nil {
		return nil, err
	}

	//--- Create temp destination

	file, err := j.createFile(JournalTemp)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)

	//--- Write entries into new temp file

	for _, entry := range entries {
		err = write(file, writer, entry, false)
		if err != nil {
			return nil, err
		}
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	err = j.deleteFile(JournalFile)
	if err != nil {
		return nil, err
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

	return entries,nil
}

//=============================================================================
// Recover parses the journal file and returns all unacknowledged messages

func (j *Journal) recover() ([]*JournalEntry, error) {
	if _, err := j.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	pendingMap := make(map[string]*JournalEntry)
	reader     := bufio.NewReader(j.file)

	for {
		var entry JournalEntry
		decoder := gob.NewDecoder(reader)
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

func (j *Journal) existsFile(name string) (bool,error) {
	_, err := os.Stat(name)
	if err == nil {
		return true,nil // File exists
	}

	if errors.Is(err, os.ErrNotExist) {
		return false,nil // File explicitly does not exist
	}

	// The file might exist, but we got a different error (e.g., permission denied)
	return false,err
}

//=============================================================================

func write(file *os.File, writer *bufio.Writer, entry *JournalEntry, sync bool) error {
	encoder := gob.NewEncoder(writer)

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
