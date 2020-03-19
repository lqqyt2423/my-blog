package article

// Article interface
type Article interface {
	Init(string)
	Get(name string) ([]byte, error)
	GetAll() ([]byte, error)
	Search(q string) ([]byte, error)
}
