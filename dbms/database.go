//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package dbms

import (
	"log/slog"
	"time"

	"github.com/algotiqa/core"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

//=============================================================================

var dbms *gorm.DB

//=============================================================================

func InitDatabase(cfg *core.Database) {

	slog.Info("Starting database...")
	url := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Address + ")/" + cfg.Name + "?charset=utf8mb4&parseTime=True"

	dialector := mysql.New(mysql.Config{
		DSN:                       url,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		core.ExitWithMessage("Failed to connect to the database: " + err.Error())
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			core.ExitWithMessage("Failed to connect to create database pool: " + err.Error())
		} else {
			sqlDB.SetConnMaxLifetime(time.Minute * 3)
			sqlDB.SetMaxOpenConns(50)
			sqlDB.SetMaxIdleConns(10)
		}

		dbms = db
	}
}

//=============================================================================

func RunInTransaction(f func(tx *gorm.DB) error) error {
	return dbms.Transaction(f)
}

//=============================================================================
