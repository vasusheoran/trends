package database

type DataStore interface {
	Read(path string) ([][]string, error)
	Write(path string, data [][]string) error
}
