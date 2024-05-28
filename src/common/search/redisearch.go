package search

type RediSearch struct{}

func (r *RediSearch) Search(index string, query any) (any, error) {
	panic("not implemented") // TODO: Implement
}

func (r *RediSearch) Index(index string, doc map[string]any) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (r *RediSearch) Update(index string, id string, doc map[string]any) error {
	panic("not implemented") // TODO: Implement
}

func (r *RediSearch) Delete(index string, id string) error {
	panic("not implemented") // TODO: Implement
}
