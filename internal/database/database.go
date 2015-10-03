package database

type DB interface {
	Encode(url string) (string, error)
	Decode(url string) (string, error)
}


