package search

type BleveSearch struct {}

func (b *BleveSearch) Search(index string, query any) (any, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BleveSearch) Index(index string, doc map[string]any) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (b *BleveSearch) Update(index string, id string, doc map[string]any) error {
	panic("not implemented") // TODO: Implement
}

func (b *BleveSearch) Delete(index string, id string) error {
	panic("not implemented") // TODO: Implement
}

