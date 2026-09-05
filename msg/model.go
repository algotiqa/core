//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package msg

//=============================================================================

const (
	TypeCreate     = 1
	TypeUpdate     = 2
	TypeDelete     = 3
	TypeActivate   = 4
	TypeDeactivate = 5
	TypeNewJob     = 6
	TypeChange     = 7
	TypeRestart    = 8

	//--- Queue: Inventory

	SourceTradingSystem  = "trading-system"
	SourceDataProduct    = "data-product"
	SourceBrokerProduct  = "broker-product"
	SourcePortfolio      = "portfolio"
	SourceTradingSession = "trading-session"
	SourceAccount        = "account"

	//--- Queue: Data Collector

	SourceUploadJob      = "upload-job"
	SourceRollRecalcJob  = "rollRecalc-job"

	//--- Queue: Runtime system

	SourceTrade          = "trade"

	//--- Queue: System adapter

	SourceConnection     = "connection"
	SourceSystem         = "system"

	//--- Queue: Event store

	SourceEvent          = "event"
)

//=============================================================================

const (
	ExInventory            = "algotiqa.inventory"
	QuInventoryToPortfolio = "algotiqa.inventory:portfolio"
	QuInventoryToCollector = "algotiqa.inventory:collector"
	QuInventoryToStorage   = "algotiqa.inventory:storage"

	ExCollector            = "algotiqa.collector"
	QuCollectorToInternal  = "algotiqa.collector:internal"

	ExRuntime              = "algotiqa.runtime"
	QuRuntimeToPortfolio   = "algotiqa.runtime:portfolio"

	ExSystem               = "algotiqa.system"
	QuSystemToInventory    = "algotiqa.system:inventory"
	QuSystemToCollector    = "algotiqa.system:collector"
	QuSystemToPortfolio    = "algotiqa.system:portfolio"

	ExEvent                = "algotiqa.event"
	QuAllToEvent           = "algotiqa.all:event"
)

//=============================================================================

type Message struct {
	Source string
	Type   int
	Entity []byte
}

//=============================================================================
