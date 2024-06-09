package viewmodel

// FruitViewModel represents the ViewModel in the MVVM architecture.
type FruitViewModel struct {
	Names       []string // List of fruit names for the datalist.
	Description string   // Description of the selected fruit.
}

func NewFruitViewModel() *FruitViewModel {
	return &FruitViewModel{
		Names: []string{"Apple", "Banana", "Cherry"},
	}
}

func NewFruitViewModelWithDescription(fruit string) *FruitViewModel {
	vm := NewFruitViewModel()
	vm.Description = getFruitDescription(fruit)
	return vm
}

// getFruitDescription returns a description based on the fruit name.
func getFruitDescription(fruit string) string {
	descriptions := map[string]string{
		"Apple":  "A crisp, sweet fruit grown in cooler climates.",
		"Banana": "A long, yellow tropical fruit that is high in potassium.",
		"Cherry": "A small, round, deep red fruit with a pit.",
	}
	return descriptions[fruit]
}
