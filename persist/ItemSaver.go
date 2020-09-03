package persist

import (
	"Crawler/engine"
	"context"
	"errors"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func ItemSaver() chan engine.Item {
	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			itemCount++
			err := Save(item)
			if err != nil {
				log.Printf("Item Saver: error saving item %v: %v", item, err)
			}
		}
	}()
	return out
}

func Save(item engine.Item) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}

	if item.Type == "" {
		return errors.New("must supply Type")
	}

	indexService := client.Index().Index("dating_profile").
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}

	_, err = indexService.Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}
