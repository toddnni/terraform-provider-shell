package shell

type Config struct {
	WorkingDirectory string
	CreateCommand    string
	CreateParameters []interface{}
	ReadCommand      string
	ReadParameters   []interface{}
	DeleteCommand    string
	DeleteParameters []interface{}
	UniqueParameters map[string]struct{}
}
