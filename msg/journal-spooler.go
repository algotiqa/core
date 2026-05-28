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
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

//=============================================================================
//===
//=== Spooler entry
//===
//=============================================================================

type SpoolerEntry struct {
	Exchange string
	Id       string
	Payload  []byte
}

//=============================================================================
//===
//=== Journal Spooler
//===
//=============================================================================

type JournalSpooler struct {
	queue           chan *SpoolerEntry
	queueSize       int
	compactMessages int
	writtenMessages int
	mu              sync.Mutex
}

//=============================================================================

func NewJournalSpooler(queueSize int, compactMessages int) (*JournalSpooler, error) {
	s := &JournalSpooler{
		queue          : make(chan *SpoolerEntry, queueSize),
		queueSize      : queueSize,
		compactMessages: compactMessages,
	}
	go s.worker()

	entries,err := s.compact()
	if err != nil {
		return nil, err
	}

	return s,s.recover(entries)
}

//=============================================================================
// Given that UUId and JSON conversion always work, the only possible error cases are:
//  - Queue full
//  - Error writing to the journal

func (s *JournalSpooler) Submit(id string, payload []byte, exchange string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.queue) == s.queueSize {
		return errors.New("queue is full")
	}

	if err := journal.Write(id, payload, exchange); err != nil {
		return fmt.Errorf("journal write failed: %w", err)
	}

	se := &SpoolerEntry{
		Exchange: exchange,
		Id      : id,
		Payload : payload,
	}

	s.queue <- se
	return nil
}

//=============================================================================
//===
//=== Worker
//===
//=============================================================================

func (s *JournalSpooler) worker() {
	for {
		select {
			case entry, ok := <- s.queue:
				if ok {
					//--- Run task
					s.run(entry)
				} else {
					//--- Channel closed. Exit from goroutine
					return
				}

			default:
				time.Sleep(time.Millisecond * 50)
		}
	}
}

//=============================================================================

func (s *JournalSpooler) run(se *SpoolerEntry) {
	errorFlag := false

	for {
		err := publish(se.Id, se.Payload, se.Exchange)
		if err == nil {
			err = journal.Ack(se.Id)
			if err != nil {
				slog.Error("Failed to write ACK to journal", "exchange", se.Exchange, "id", se.Id, "error", err.Error())
			} else {
				if errorFlag {
					slog.Info("Message retried with success", "exchange", se.Exchange, "id", se.Id)
				}
				break
			}
		}
		errorFlag = true
		time.Sleep(time.Millisecond * 1000)
	}

	s.writtenMessages++

	if s.writtenMessages >= s.compactMessages {
		_,_=s.compact()
	}
}

//=============================================================================

func (s *JournalSpooler) compact() ([]*JournalEntry, error) {
	start       := time.Now()
	entries,err := journal.Compact()
	duration    := time.Now().Sub(start)

	if err != nil {
		slog.Error("Could not compact journal: " + err.Error())
	} else {
		s.writtenMessages = len(entries)
		slog.Info("Journal compacted successfully", "count", s.writtenMessages, "time", duration.Seconds())
	}

	return entries, err
}

//=============================================================================

func (s *JournalSpooler) recover(entries []*JournalEntry) error {
	for _, entry := range entries {
		err := s.Submit(entry.Id, entry.Payload, entry.Exchange)
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================
