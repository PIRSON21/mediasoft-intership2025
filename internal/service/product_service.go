package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/PIRSON21/mediasoft-go/internal/domain"
	"github.com/PIRSON21/mediasoft-go/internal/dto"
	"github.com/PIRSON21/mediasoft-go/internal/repository"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"go.uber.org/zap"
)

type ProductService struct {
	host string
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository, host string) *ProductService {
	hostURL := createHostURL(host)
	return &ProductService{repo: repo, host: hostURL}
}

func createHostURL(host string) string {
	if strings.Index(host, ":") == 0 {
		host = "http://localhost" + host
	}

	if !(strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://")) {
		host = "http://" + host
	}

	return host
}

func (s *ProductService) GetProducts(ctx context.Context) ([]*dto.ProductAtListResponse, error) {
	log := logger.GetLogger().With(zap.String("op", "service.ProductService.GetProduct"))

	products, err := s.repo.GetProducts(ctx)
	if err != nil {
		log.Error("error while getting products from repo", zap.String("err", err.Error()))
		return nil, err
	}

	productsResponse := s.createProductsResponse(products)

	return productsResponse, nil
}

func (s *ProductService) createProductsResponse(products []*domain.Product) []*dto.ProductAtListResponse {
	var response []*dto.ProductAtListResponse
	for _, v := range products {
		params := copyMap(v.Params)
		response = append(response, &dto.ProductAtListResponse{
			ID:          v.ID,
			Weight:      v.Weight,
			Name:        v.Name,
			Description: v.Description,
			Barcode:     s.host + "/static/" + v.Barcode,
			Params:      params,
		})
	}

	return response
}

func copyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	newMap := make(map[string]any, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

func (s *ProductService) AddProduct(ctx context.Context, request *dto.ProductRequest) error {
	log := logger.GetLogger().With(zap.String("op", "service.ProductService.AddProduct"))

	filename, err := createFile(request.Barcode)
	if err != nil {
		log.Error("error while creating file", zap.String("err", err.Error()))
		return err
	}

	product := parseProductFromRequest(request, filename)

	err = s.repo.AddProduct(ctx, product)
	if err != nil {
		log.Error("error while adding product to repository", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func createFile(photo *dto.Photo) (string, error) {
	defer photo.File.Close()

	// TODO: можно сделать генерацию названия файла.
	savePath := filepath.Join("static", photo.Handler.Filename)
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, photo.File)
	if err != nil {
		return "", err
	}

	return photo.Handler.Filename, nil
}

func parseProductFromRequest(req *dto.ProductRequest, filename string) *domain.Product {
	return &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Weight:      *req.Weight,
		Params:      req.Params,
		Barcode:     filename,
	}
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID int, productReq *dto.ProductRequest) error {
	log := logger.GetLogger().With(zap.String("op", "service.ProductService.UpdateProduct"))

	var fileName string
	var err error
	if productReq.Barcode != nil {
		fileName, err = createFile(productReq.Barcode)
		if err != nil {
			log.Error("error while creating file", zap.String("err", err.Error()))
			return err
		}
	}

	product := parseProductFromUpdateRequest(productReq, fileName)
	product.ID = productID

	err = s.repo.UpdateProduct(ctx, product)
	if err != nil {
		log.Error("error while updating product at repository", zap.String("err", err.Error()))
		return err
	}

	return nil
}

func parseProductFromUpdateRequest(req *dto.ProductRequest, filename string) *domain.Product {
	var product domain.Product

	if req.Name != "" {
		product.Name = req.Name
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.Weight != nil {
		product.Weight = *req.Weight
	}

	if len(req.Params) != 0 {
		product.Params = copyMap(req.Params)
	}

	if filename != "" {
		product.Barcode = filename
	}

	return &product
}
