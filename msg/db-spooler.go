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

	"github.com/algotiqa/core/dbms"
	"gorm.io/gorm"
)

//=============================================================================

func InitDbSpooler(pollInterval int) {
	slog.Info("Starting DB spooler...")

	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	go func() {
		for range ticker.C {
			runDbSpooler()
		}
	}()
}

//=============================================================================

func addOutboxMessage(tx *gorm.DB, id string, payload []byte, exchange string) (err error) {
	om := &dbms.OutboxMessage{
		Timestamp: time.Now(),
		Exchange : exchange,
		Uuid     : id,
		Payload  : payload,
		Size     : len(payload),
	}

	return dbms.AddOutboxMessage(tx, om)
}

//=============================================================================

func runDbSpooler() {
	list,err := getOutboxMessages()
	if err != nil {
		slog.Error("runDbSpooler: Error getting outbox messages", "error", err)
	}

	for _, om := range *list {
		err = publish(om.Uuid, om.Payload, om.Exchange)
		if err != nil {
			slog.Error("runDbSpooler: Error sending message", "error", err, "id", om.Id)
			return
		}

		err = deleteMessage(om.Id)
		if err != nil {
			slog.Error("runDbSpooler: Error deleting message", "error", err, "id", om.Id)
			return
		}
	}

	if len(*list) != 0 {
		slog.Info("runDbSpooler: Sent messages", "count", len(*list))
	}
}

//=============================================================================

func getOutboxMessages() (*[]dbms.OutboxMessage, error) {
	var list *[]dbms.OutboxMessage

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		var err error
		list, err = dbms.GetOutboxMessages(tx)
		return err
	})

	return list, err
}

//=============================================================================

func deleteMessage(id uint) error {
	return dbms.RunInTransaction(func(tx *gorm.DB) error {
		return dbms.DeleteOutboxMessage(tx, id)
	})
}

//=============================================================================
