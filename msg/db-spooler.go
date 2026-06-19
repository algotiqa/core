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
