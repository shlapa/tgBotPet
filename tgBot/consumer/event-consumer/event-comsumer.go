package event_consumer

import (
	"log"
	"tgBot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	bathSize  int
}

func New(fetcher events.Fetcher, processor events.Processor, bathSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		bathSize:  bathSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvent, err := c.fetcher.Fetch(c.bathSize)
		if err != nil {
			log.Printf("Error consumer: %s", err.Error())

			continue
		}

		if len(gotEvent) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}
		if err := c.handleEvents(gotEvent); err != nil {
			log.Print("Error handling events")

			continue
		}
	}
}

func (c Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("cant handle event: %s", err.Error())

			continue
		}
	}
	return nil
}
