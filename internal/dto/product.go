package dto

import "mime/multipart"

type ProductAtListResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Weight      float64        `json:"weight"`
	Description string         `json:"desc"`
	Params      map[string]any `json:"params,omitempty"`
	Barcode     string         `json:"barcode_url"` // Ссылка на доступ к штрихкоду.
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
