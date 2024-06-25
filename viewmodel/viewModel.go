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
	View          []service.GeneralEntity
	Cursor        string
	PageSize      int
	HasNextPage   bool
	HasPrevPage   bool
	CurrentPage   int
	Pages         int
	SortKey       string
	SortDirection string
}

func NewTableViewModel(c *datastore.Client) *TableViewModel {
	return &TableViewModel{
		client:   c,
		PageSize: 50,
		Cursor:   "",
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
	vm.Reset()
	return nil

}

func (vm *TableViewModel) ToggleSortDirection() {
	newSortDirection := ""
	if vm.SortDirection == "" {
		newSortDirection = "desc"
	} else if vm.SortDirection == "desc" {
		newSortDirection = "asc"
	}
	vm.Reset()
	vm.SortDirection = newSortDirection

}

func (vm *TableViewModel) RowCount() int {
	return len(vm.Entities)
}

func (vm *TableViewModel) Reset() {
	vm.Cursor = ""
	vm.SortKey = ""
	vm.SortDirection = ""
	vm.Entities = nil
	vm.CurrentPage = 0
	vm.Pages = 0
	vm.HasPrevPage = false
	vm.HasNextPage = true

}

func (vm *TableViewModel) DebugInfo() {

	fmt.Println("Selected", vm.Selected)
	fmt.Println("Cursor", vm.Cursor)
	fmt.Println("CurrentPage", vm.CurrentPage)
	fmt.Println("PageSize", vm.PageSize)
	fmt.Println("Pages", vm.Pages)
	fmt.Println("SortKey", vm.SortKey)
	fmt.Println("SortDirection", vm.SortDirection)

}

func (vm *TableViewModel) GetNewPage(ctx context.Context) error {
	fmt.Println("Getting new page")
	if vm.Selected == "" {
		return fmt.Errorf("No kind selected")
	}

	entities, nextCursor, err := service.GetAllEntities(ctx, vm.client, vm.Selected, vm.SortKey, vm.SortDirection, vm.PageSize, vm.Cursor)

	if err != nil {
		return err
	}

	vm.Entities = append(vm.Entities, entities...)
	vm.Cursor = nextCursor

	vm.CurrentPage += 1
	vm.Pages += 1
	vm.HasPrevPage = vm.CurrentPage > 1
	vm.HasNextPage = nextCursor != ""
	return nil
}
