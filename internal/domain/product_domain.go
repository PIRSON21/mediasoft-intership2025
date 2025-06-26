package domain

type Product struct {
	ID          int
	Weight      int
	Name        string
	Description string
	Barcode     string // Штрихкод. Здесь хранится только название файла. Сам файл хранится на диске сервера. sdasad
	Params      map[string]any
}
