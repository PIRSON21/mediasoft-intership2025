package dto

import "mime/multipart"

type ProductAtListResponse struct {
	ID          string         `json:"id" example:"12345"`
	Name        string         `json:"name" example:"Product Name"`
	Weight      float64        `json:"weight" example:"1.5"`
	Description string         `json:"desc" example:"This is a product description."`
	Params      map[string]any `json:"params,omitempty" example:"{\"color\": \"red\", \"size\": \"M\"}"`
	Barcode     string         `json:"barcode_url" example:"http://localhost:8080/static/photo.png"` // Ссылка на доступ к штрихкоду.
}

type ProductRequest struct {
	Name        string         `json:"name"`
	Weight      *float64       `json:"weight"`
	Description string         `json:"desc"`
	Params      map[string]any `json:"params"`
	Barcode     *Photo         `json:"barcode"` // Штрихкод в байтах
}

type Photo struct {
	File    multipart.File
	Handler *multipart.FileHeader
}
