//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

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
