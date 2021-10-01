package api

type Database interface {
	//Read(path string) ([]contracts.Stock, error)
	Read(path string) ([][]string, error)
	Write(path string, data [][]string) error
}
