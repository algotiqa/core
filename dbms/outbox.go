//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package dbms

import (
	"github.com/algotiqa/core/req"
	"gorm.io/gorm"
)

//=============================================================================

func GetOutboxMessages(tx *gorm.DB) (*[]OutboxMessage, error) {
	var list []OutboxMessage
	res := tx.Find(&list).Order("id").Limit(1000)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddOutboxMessage(tx *gorm.DB, om *OutboxMessage) error {
	return tx.Create(om).Error
}

//=============================================================================

func DeleteOutboxMessage(tx *gorm.DB, id uint) error {
	return tx.Delete(&OutboxMessage{}, id).Error
}

//=============================================================================
