package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// Initializes a Datastore client
func NewDatastoreClient(ctx context.Context) (*datastore.Client, error) {
	projectId := "my-project-id" //change this to your project id
	os.Setenv("DATASTORE_DATASET", projectId)
	os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")
	os.Setenv("DATASTORE_EMULATOR_HOST_PATH", "localhost:8081/datastore")
	os.Setenv("DATASTORE_HOST", "http://localhost:8081")
	os.Setenv("DATASTORE_PROJECT_ID", projectId)
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

		fmt.Println("p", p)
		fmt.Println("type", reflect.TypeOf(p.Value))
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
func GetAllEntities(ctx context.Context, client *datastore.Client, kind string, limit int, cursorStr string) ([]GeneralEntity, string, error) {
	query := datastore.NewQuery(kind).Limit(limit)
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
		var entity GeneralEntity
		_, err := it.Next(&entity)

		if err == iterator.Done {
			break
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

	return entities, nextCursor.String(), nil
}
