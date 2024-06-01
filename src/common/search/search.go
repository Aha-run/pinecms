package search


type ISearch interface {
	Search(index string, query any) (any, error)
	Index(index string, doc map[string]any) (string, error)
	Update(index, id string, doc map[string]any) error
	Delete(index string, id string) error
}
