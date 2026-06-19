//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package dbms

import "time"

//=============================================================================

type OutboxMessage struct {
	Id        uint      `gorm:"primaryKey"`
	Timestamp time.Time
	Exchange  string
	Uuid      string
	Payload   []byte
	Size      int
}

//=============================================================================

func (OutboxMessage) TableName() string { return "outbox" }

//=============================================================================
