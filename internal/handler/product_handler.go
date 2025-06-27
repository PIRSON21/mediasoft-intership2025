package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PIRSON21/mediasoft-go/internal/dto"
	custErr "github.com/PIRSON21/mediasoft-go/internal/errors"
	"github.com/PIRSON21/mediasoft-go/internal/middleware"
	"github.com/PIRSON21/mediasoft-go/internal/service"
	"github.com/PIRSON21/mediasoft-go/pkg/logger"
	"github.com/PIRSON21/mediasoft-go/pkg/render"
	"go.uber.org/zap"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

func (h *ProductHandler) ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProducts(w, r)
	case http.MethodPost:
		h.AddProduct(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.ProductHandler.GetProducts"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	productResponse, err := h.service.GetProducts(r.Context())
	if err != nil {
		log.Error("error while getting products", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while getting products")
		return
	}

	err = render.JSON(w, http.StatusOK, productResponse)
	if err != nil {
		log.Error("error while rendering products", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while marshalling products")
		return
	}
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(
		zap.String("op", "handler.ProductHandler.AddProduct"),
		zap.String("request-id", middleware.GetRequestID(r.Context())),
	)

	productRequest, err := parseProduct(r)
	if err != nil {
		log.Error("err while parsing product", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, "err while parsing product")
		return
	}

	validErr := validateCreateProduct(productRequest)
	if validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	err = h.service.AddProduct(r.Context(), productRequest)
	if err != nil {
		if errors.Is(err, custErr.ErrProductAlreadyExists) {
			custErr.UnnamedError(w, http.StatusConflict, "product with this name already exists")
			return
		}
		log.Error("err while adding product", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, "err while adding product")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func parseProduct(r *http.Request) (*dto.ProductRequest, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}

	var product dto.ProductRequest
	product.Name = r.FormValue("name")
	product.Description = r.FormValue("description")
	weightStr := r.FormValue("weight")
	if weightStr != "" {
		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing weight: %w", err)
		}
		product.Weight = &weight
	}
	params := r.FormValue("params")
	if params != "" {
		err = json.NewDecoder(strings.NewReader(r.FormValue("params"))).Decode(&product.Params)
		if err != nil {
			return nil, fmt.Errorf("error while parsing params: %w", err)
		}
	}

	file, handler, err := r.FormFile("barcode")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			return nil, err
		}
	}

	if file != nil && handler != nil {
		photo := dto.Photo{
			File:    file,
			Handler: handler,
		}

		product.Barcode = &photo
	}

	return &product, nil
}

func validateCreateProduct(product *dto.ProductRequest) map[string]string {
	validErr := make(map[string]string, 0)

	if len(product.Name) == 0 {
		validErr["name"] = "name cannot be empty"
	}

	if product.Weight == nil {
		validErr["weight"] = "weight cannot be empty"
	} else if *product.Weight <= 0 {
		validErr["weight"] = "weight is incorrect"
	}

	if product.Barcode == nil {
		validErr["barcode"] = "there must be barcode"
	}

	if len(validErr) == 0 {
		return nil
	}

	return validErr
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger().With(zap.String("op", "handler.ProductHandler.UpdateProduct"))

	productID, err := parseProductID(r)
	if err != nil {
		custErr.UnnamedError(w, http.StatusBadRequest, "wrong product ID")
		return
	}

	product, err := parseProduct(r)
	if err != nil {
		log.Error("error while parsing product", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while parsing product")
		return
	}

	validErr := validateUpdateProduct(product)
	if validErr != nil {
		render.JSON(w, http.StatusBadRequest, validErr)
		return
	}

	err = h.service.UpdateProduct(r.Context(), productID, product)
	if err != nil {
		log.Error("error while updating product", zap.String("err", err.Error()))
		custErr.UnnamedError(w, http.StatusInternalServerError, "error while updating product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseProductID(r *http.Request) (int, error) {
	splits := strings.Split(r.URL.Path, "/")
	productID, err := strconv.Atoi(splits[len(splits)-1])
	if err != nil {
		return 0, err
	}

	if productID <= 0 {
		return 0, fmt.Errorf("id cannot be less than 1")
	}

	return productID, nil
}

func validateUpdateProduct(product *dto.ProductRequest) map[string]string {
	validErr := make(map[string]string, 0)

	if product.Weight != nil && *product.Weight <= 0 {
		validErr["weight"] = "weight must be greater than 0"
	}

	if len(validErr) == 0 {
		return nil
	}

	return validErr
}
