package dbi

//Interface minimal database interface
//go:generate counterfeiter . Interface
type Interface interface {
	CreateTable(name string) error
	Index(tableName string) (map[string]string, error)
	Insert(tableName string, document map[string]string) error
}
