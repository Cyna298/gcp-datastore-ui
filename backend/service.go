package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// Initializes a Datastore client
func NewDatastoreClient(ctx context.Context) (*datastore.Client, error) {
	if os.Getenv("DATASTORE_DATASET") == "" {
		log.Fatal("DATASTORE_DATASET environment variable must be set to the project ID")
	}
	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		log.Fatal("DATASTORE_EMULATOR_HOST environment variable must be set to the emulator host")
	}
	if os.Getenv("DATASTORE_EMULATOR_HOST_PATH") == "" {
		log.Fatal("DATASTORE_EMULATOR_HOST_PATH environment variable must be set to the emulator host path")
	}
	if os.Getenv("DATASTORE_HOST") == "" {
		log.Fatal("DATASTORE_HOST environment variable must be set to the datastore host")
	}
	if os.Getenv("DATASTORE_PROJECT_ID") == "" {
		log.Fatal("DATASTORE_PROJECT_ID environment variable must be set to the project ID")
	}

	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GetAllKinds retrieves all kinds from Datastore
func GetAllKinds(ctx context.Context, client *datastore.Client) ([]string, error) {
	query := datastore.NewQuery("__kind__").KeysOnly()
	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var kinds []string
	for _, key := range keys {
		kinds = append(kinds, key.Name)
	}
	return kinds, nil
}

type OutputProperty struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	TypeOf  string `json:"type"`
	Indexed bool   `json:"indexed"`
}

type GeneralEntity map[string]OutputProperty

func (x *GeneralEntity) Load(ps []datastore.Property) error {
	*x = make(map[string]OutputProperty)

	for _, p := range ps {
		(*x)[p.Name] = OutputProperty{
			Name:    p.Name,
			Value:   fmt.Sprintf("%v", p.Value),
			TypeOf:  fmt.Sprintf("%T", p.Value),
			Indexed: !p.NoIndex,
		}
	}

	return nil
}

func (x *GeneralEntity) Save() ([]datastore.Property, error) {

	return []datastore.Property{
		{
			Name:  "I",
			Value: x,
		},
	}, nil
}

// GetAllEntities retrieves entities of a specific kind from Datastore
func GetAllEntities(ctx context.Context, client *datastore.Client, kind string, sortKey string, sortDirection string, limit int, cursorStr string) ([]GeneralEntity, string, error) {
	query := datastore.NewQuery(kind).Limit(limit)
	if sortKey != "" {
		if sortDirection == "desc" {
			query = query.Order("-" + sortKey)

		} else {
			query = query.Order(sortKey)
		}
	}
	if cursorStr != "" {
		cursor, err := datastore.DecodeCursor(cursorStr)
		if err != nil {
			return nil, "", err
		}
		query = query.Start(cursor)
	}

	var entities []GeneralEntity
	it := client.Run(ctx, query)

	for {
		var entity GeneralEntity = make(map[string]OutputProperty)
		key, err := it.Next(&entity)

		entity["key"] = OutputProperty{
			Name:    "key",
			Value:   key.String(),
			TypeOf:  "string",
			Indexed: true,
		}

		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("err", err)
			return nil, "", err
		}

		if err != nil {
			fmt.Println("err", err)
			return nil, "", err
		}
		entities = append(entities, entity)
	}

	nextCursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}
	if limit > len(entities) {
		return entities, "", nil
	}

	return entities, nextCursor.String(), nil
}
