package domain

import "github.com/google/uuid"

// Product представляет продукт с его деталями.
type Product struct {
	ID          uuid.UUID
	Weight      float64
	Name        string
	Description string
	Barcode     string // Штрихкод. Здесь хранится только название файла. Сам файл хранится на диске сервера. sdasad
	Params      map[string]any
}
