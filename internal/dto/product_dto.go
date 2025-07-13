package dto

import "mime/multipart"

// ProductAtListResponse представляет продукт в списке с его деталями.
type ProductAtListResponse struct {
	ID          string         `json:"id" example:"12345"`
	Name        string         `json:"name" example:"Product Name"`
	Weight      float64        `json:"weight" example:"1.5"`
	Description string         `json:"desc" example:"This is a product description."`
	Params      map[string]any `json:"params,omitempty" example:"{\"color\": \"red\", \"size\": \"M\"}"`
	Barcode     string         `json:"barcode_url" example:"http://localhost:8080/static/photo.png"` // Ссылка на доступ к штрихкоду.
}

// ProductRequest представляет запрос на создание или обновление продукта.
type ProductRequest struct {
	Name        string         `json:"name"`
	Weight      *float64       `json:"weight"`
	Description string         `json:"desc"`
	Params      map[string]any `json:"params"`
	Barcode     *Photo         `json:"barcode"` // Штрихкод в байтах
}

// Photo представляет файл изображения штрихкода.
type Photo struct {
	File    multipart.File
	Handler *multipart.FileHeader
}
