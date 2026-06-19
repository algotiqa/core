//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package msg

import (
	"time"

	"gorm.io/gorm"
)

//=============================================================================
//===
//=== Events
//===
//=============================================================================

type EventLevel int8

const (
	EventLevelInfo    = 0
	EventLevelWarning = 1
	EventLevelError   = 2
)

//=============================================================================

type Event struct {
	Username   string
	Level      EventLevel
	EventDate  time.Time
	Code       string
	Title      string
	Message    string
	Parameters map[string]any
}

//=============================================================================

func SendEventByCode(username string, code string, params map[string]any, tx *gorm.DB) error {
	e := Event{
		Username  : username,
		EventDate : time.Now(),
		Code      : code,
		Parameters: params,
	}

	return SendMessage(ExEvent, SourceEvent, TypeCreate, e, tx)
}

//=============================================================================

func SendEvent(username string, level EventLevel, title, message string, params map[string]any, tx *gorm.DB) error {
	e := Event{
		Username  : username,
		Level     : level,
		EventDate : time.Now(),
		Title     : title,
		Message   : message,
		Parameters: params,
	}

	return SendMessage(ExEvent, SourceEvent, TypeCreate, e, tx)
}

//=============================================================================
