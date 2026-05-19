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

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/algotiqa/core"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

//=============================================================================
// Note: 1000 messages could eat (more or less) 700MB of memory. This is due to the
//       data from agents (each trading system with full history could be 700KB)

var url     string
var channel *amqp.Channel
var journal *Journal
var spooler *JournalSpooler
var config  *core.Journal

//=============================================================================

func InitMessaging(cfg *core.Messaging) {
	slog.Info("Starting messaging...")

	config = cfg.Journal

	var err error

	//--- Create journal

	journal,err = NewJournal(config.Directory)
	if err != nil {
		core.ExitWithMessage("Failed to create message journal: " + err.Error())
	}

	//--- Connect to messaging system

	url = "amqp://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Address + "/"

	err = connect()
	if err != nil {
		core.ExitWithMessage("Failed to connect to the messaging system or to get a channel: " + err.Error())
	}

	createExchange(ExInventory)
	createQueue(QuInventoryToPortfolio)
	bindQueue(ExInventory, QuInventoryToPortfolio)
	createQueue(QuInventoryToCollector)
	bindQueue(ExInventory, QuInventoryToCollector)
	createQueue(QuInventoryToStorage)
	bindQueue(ExInventory, QuInventoryToStorage)

	createExchange(ExCollector)
	createQueue(QuCollectorToInternal)
	bindQueue(ExCollector, QuCollectorToInternal)

	createExchange(ExRuntime)
	createQueue(QuRuntimeToPortfolio)
	bindQueue(ExRuntime, QuRuntimeToPortfolio)

	createExchange(ExSystem)
	createQueue(QuSystemToCollector)
	bindQueue(ExSystem, QuSystemToCollector)
	createQueue(QuSystemToInventory)
	bindQueue(ExSystem, QuSystemToInventory)
	createQueue(QuSystemToPortfolio)
	bindQueue(ExSystem, QuSystemToPortfolio)

	createExchange(ExEvent)
	createQueue(QuAllToEvent)
	bindQueue(ExEvent, QuAllToEvent)

	//--- Create spooler (here, because the channel must be ready)

	spooler,err = NewJournalSpooler(config.QueueSize, config.CompactMessages)
	if err != nil {
		core.ExitWithMessage("Failed to create journal spooler: " + err.Error())
	}

	if config.DbSpoolInterval > 0 {
		InitDbSpooler(config.DbSpoolInterval)
	}
}

//=============================================================================

func SendMessage(exchange string, source string, msgType int, entity any, tx *gorm.DB) error {
	uuidVal,err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("error generating UUID: %w", err)
	}

	id := uuidVal.String()

	body, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("error marshalling entity: %w", err)
	}

	message := &Message{
		Source: source,
		Type  : msgType,
		Entity: body,
	}

	body, err = json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %w", err)
	}

	if tx == nil {
		return spooler.Submit(id, body, exchange)
	}

	return addOutboxMessage(tx, id, body, exchange)
}

//=============================================================================

func ReceiveMessages(queue string, handler func(m *Message) bool) {
	for {
		messages, err := channel.Consume(queue, "", false, false, false, false, nil)

		if err != nil {
			core.ExitWithMessage("ReceiveMessages: Cannot create the consumer channel for '" + queue + "' : " + err.Error())
		}

		for d := range messages {
			msg := Message{}
			err = json.Unmarshal(d.Body, &msg)

			if err != nil {
				slog.Error("ReceiveMessages: Error unmarshalling message. Rejecting.", "error", err.Error())
				err = d.Reject(false)
				if err != nil {
					slog.Error("ReceiveMessages: Cannot reject message!", "error", err.Error())
				}
				continue
			}

			if handler(&msg) {
				err = d.Ack(false)
			} else {
				err = d.Nack(false, true)
			}

			if err != nil {
				slog.Error("ReceiveMessages: Cannot [N]acknowledge message!", "error", err.Error())
			}
		}

		slog.Warn("ReceiveMessages: Exited from for loop. Reconnecting...")

		if channel.IsClosed() {
			err = connect()
			if err != nil {
				core.ExitWithMessage("ReceiveMessages: Cannot reconnect to the channel: " + err.Error())
			} else {
				slog.Info("ReceiveMessages: Successfully reconnected")
			}
		}
	}
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func createExchange(name string) {
	err := channel.ExchangeDeclare(name, "fanout", true, false, false, false, nil)

	if err != nil {
		core.ExitWithMessage("Cannot create the '" + name + "' exchange in the messaging system: " + err.Error())
	}
}

//=============================================================================

func createQueue(name string) {
	_, err := channel.QueueDeclare(name, true, false, false, false, nil)

	if err != nil {
		core.ExitWithMessage("Cannot create the '" + name + "' queue in the messaging system: " + err.Error())
	}
}

//=============================================================================

func bindQueue(exchange, queue string) {
	err := channel.QueueBind(queue, "", exchange, false, nil)

	if err != nil {
		core.ExitWithMessage("Cannot bind queue '" + queue + "' to the exchange: " + err.Error())
	}
}

//=============================================================================

func connect() error {
	conn, err := amqp.Dial(url)
	if err == nil {
		channel, err = conn.Channel()
		if err == nil {
			//--- Put the channel into confirm mode
			err = channel.Confirm(false)
		}
	}

	return err
}

//=============================================================================

func publish(id string, payload []byte, exchange string) error {
	if channel.IsClosed() {
		slog.Warn("publish: Channel is closed. Reconnecting...")
		err := connect()
		if err != nil {
			slog.Error("publish: Reconnect failure", "error", err.Error())
			return err
		}
		slog.Info("publish: Reconnected successfully")
	}

	dc,err := channel.PublishWithDeferredConfirm(exchange, "", false, false,
		amqp.Publishing{
			MessageId   : id,
			ContentType :"application/json",
			Body        : payload,
			DeliveryMode: amqp.Persistent, // Ensure RabbitMQ persists it to disk too
		})

	if err != nil {
		slog.Error("Cannot publish a message to exchange", "exchange", exchange, "id", id, "error", err.Error())
		return err
	}

	if !dc.Wait() {
		slog.Error("Messaging system didn't ACK the message", "exchange", exchange, "id", id,)
		return err
	}

	return nil
}

//=============================================================================
