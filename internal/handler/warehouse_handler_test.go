package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PIRSON21/mediasoft-intership2025/internal/dto"
	"github.com/PIRSON21/mediasoft-intership2025/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWarehouses(t *testing.T) {
	cases := []struct {
		Name             string
		ReturnWarehouses []*dto.WarehouseAtListResponse
		ReturnError      error
		StatusCode       int
		JSON             bool
		ResponseBody     string
	}{
		{
			Name:             "Success",
			ReturnWarehouses: []*dto.WarehouseAtListResponse{{ID: "17b79680-4657-4ef4-9c3d-554a83c31828", Address: "Warehouse 1"}},
			ReturnError:      nil,
			StatusCode:       http.StatusOK,
			JSON:             true,
			ResponseBody:     `[{"id":"17b79680-4657-4ef4-9c3d-554a83c31828","address":"Warehouse 1"}]`,
		},
		{
			Name:             "Error",
			ReturnWarehouses: nil,
			ReturnError:      errors.New("internal server error"),
			StatusCode:       http.StatusInternalServerError,
			JSON:             true,
			ResponseBody:     `{"error":"error while getting warehouses: \"internal server error\""}`,
		},
		{
			Name:             "Empty List",
			ReturnWarehouses: []*dto.WarehouseAtListResponse{},
			ReturnError:      nil,
			StatusCode:       http.StatusOK,
			JSON:             true,
			ResponseBody:     `[]`,
		},
		{
			Name: "Some warehouses",
			ReturnWarehouses: []*dto.WarehouseAtListResponse{
				{ID: "1", Address: "Address 1"},
				{ID: "2", Address: "Address 2"},
			},
			ReturnError:  nil,
			StatusCode:   http.StatusOK,
			JSON:         true,
			ResponseBody: `[{"id":"1","address":"Address 1"},{"id":"2","address":"Address 2"}]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			mockService := NewMockWarehouseService(t)
			mockService.On("GetWarehouses", context.Background()).
				Return(tc.ReturnWarehouses, tc.ReturnError).
				Once()

			logger.CreateNOPLogger()

			handler := NewWarehouseHandler(mockService)
			req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)

			rr := httptest.NewRecorder()

			handler.GetWarehouses(rr, req)
			require.Equal(t, tc.StatusCode, rr.Code)

			if tc.JSON {
				assert.JSONEq(t, tc.ResponseBody, rr.Body.String())
			} else {
				assert.Equal(t, tc.ResponseBody, rr.Body.String())
			}
			require.True(t, mockService.AssertExpectations(t))
		})
	}
}
