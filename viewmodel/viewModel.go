package viewmodel

import (
	"backend/service"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
)

type TableViewModel struct {
	Kinds         []string // List of fruit names for the datalist.
	Selected      string
	client        *datastore.Client
	Headers       []service.TableHeader
	Entities      []service.GeneralEntity
	Cursor        string
	PageSize      int
	SortKey       string
	SortDirection string
}

func NewTableViewModel(c *datastore.Client) *TableViewModel {
	return &TableViewModel{
		client:   c,
		PageSize: 100,
	}

}

func (vm *TableViewModel) UpdateKinds(ctx context.Context) error {

	kinds, err := service.GetAllKinds(ctx, vm.client)
	if err != nil {
		log.Fatal(err)
		return err
	}
	vm.Kinds = kinds
	return nil
}

func (vm *TableViewModel) SelectKind(kind string) error {
	vm.Selected = kind
	return nil

}

func (vm *TableViewModel) ToggleSortDirection() {
	if vm.SortDirection == "" {
		vm.SortDirection = "desc"
	} else if vm.SortDirection == "desc" {
		vm.SortDirection = "asc"
	} else {
		vm.SortDirection = ""
	}

}

func (vm *TableViewModel) GetData(ctx context.Context, cursor string, sortKey string) error {
	if vm.Selected == "" {
		return fmt.Errorf("No kind selected")
	}

	if sortKey != "" {
		vm.ToggleSortDirection()

	}

	entities, nextCursor, err := service.GetAllEntities(ctx, vm.client, vm.Selected, sortKey, vm.SortDirection, vm.PageSize, cursor)

	if err != nil {
		return err
	}

	headers := make(map[string]service.TableHeader)

	for _, e := range entities {
		for _, x := range e {
			if _, ok := headers[x.Name]; !ok {
				headers[x.Name] = service.TableHeader{
					Name: x.Name,
					Type: x.TypeOf,
				}

			}

		}
	}
	headerValues := make([]service.TableHeader, len(headers))
	i := 0
	for _, e := range headers {
		if e.Name == "key" {
			headerValues[0] = e
		} else {
			headerValues[i+1] = e
			i += 1
		}
	}

	vm.Entities = entities
	vm.Headers = headerValues
	vm.Cursor = nextCursor

	vm.SortKey = sortKey
	return nil
}
