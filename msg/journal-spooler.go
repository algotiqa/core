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
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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
	compactMessages int
	writtenMessages int
}

//=============================================================================

func NewJournalSpooler(queueSize int, compactMessages int) (*JournalSpooler, error) {
	s := &JournalSpooler{
		queue          : make(chan *SpoolerEntry, queueSize),
		compactMessages: compactMessages,
	}
	go s.worker()

	return s,s.recover()
}

//=============================================================================

func (s *JournalSpooler) Submit(se *SpoolerEntry) {
	s.queue <- se
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
		err := s.publish(se)
		if err == nil {
			if errorFlag {
				slog.Info("Message retried with success", "exchange", se.Exchange, "id", se.Id)
			}
			break
		}
		errorFlag = true
		time.Sleep(time.Millisecond * 1000)
	}

	s.writtenMessages++

	if s.writtenMessages >= s.compactMessages {
		err := journal.Compact()
		if err != nil {
			slog.Error("Could not compact journal: " + err.Error())
		} else {
			s.writtenMessages = 0
		}
	}
}

//=============================================================================

func (s *JournalSpooler) publish(se *SpoolerEntry) error {
	if channel.IsClosed() {
		slog.Warn("publish: Channel is closed. Reconnecting...")
		err := connect()
		if err != nil {
			slog.Error("publish: Reconnect failure", "error", err.Error())
			return err
		}
		slog.Info("publish: Reconnected successfully")
	}

	dc,err := channel.PublishWithDeferredConfirm(se.Exchange, "", false, false,
		amqp.Publishing{
			MessageId   : se.Id,
			ContentType :"application/json",
			Body        : se.Payload,
			DeliveryMode: amqp.Persistent, // Ensure RabbitMQ persists it to disk too
		})

	if err != nil {
		slog.Error("Cannot publish a message to exchange", "exchange", se.Exchange, "id", se.Id, "error", err.Error())
		return err
	}

	if !dc.Wait() {
		slog.Error("Messaging system didn't ACK the message", "exchange", se.Exchange, "id", se.Id,)
		return err
	}

	if err = journal.Ack(se.Id); err != nil {
		slog.Error("Failed to write ACK to journal", "exchange", se.Exchange, "id", se.Id, "error", err.Error())
		return err
	}

	return nil
}

//=============================================================================

func (s *JournalSpooler) recover() error {
	entries,err := journal.Recover()

	if err == nil {
		for _, entry := range entries {
			se := &SpoolerEntry{
				Exchange: entry.Exchange,
				Id:       entry.Id,
				Payload:  entry.Payload,
			}
			s.Submit(se)
		}
	}

	return err
}

//=============================================================================
