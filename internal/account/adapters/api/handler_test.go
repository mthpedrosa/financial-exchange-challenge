package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/adapters/api"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
	mocks "github.com/mthpedrosa/financial-exchange-challenge/mocks/app"
	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestHandler(t *testing.T) (api.Account, *mocks.Account, *echo.Echo) {
	mockService := mocks.NewAccount(t)
	handler := api.NewAccountHandler(mockService)
	e := echo.New()
	return handler, mockService, e
}

func TestAccountHandler_Create(t *testing.T) {
	handler, mockService, e := setupTestHandler(t)
	reqBody := dto.CreateAccountRequest{Email: "new@example.com", Name: "New User"}

	t.Run("should create account successfully", func(t *testing.T) {
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		newID := uuid.NewString()
		mockService.On("Create", mock.Anything, reqBody).
			Return(dto.CreateAcountResponse{ID: newID}, nil).
			Once()

		err := handler.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response dto.CreateAcountResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, newID, response.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request on invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("should return bad request on validation failure", func(t *testing.T) {
		invalidReq := dto.CreateAccountRequest{Email: "", Name: ""}
		jsonBody, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Create(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("should return conflict on email already exists", func(t *testing.T) {
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("Create", mock.Anything, reqBody).
			Return(dto.CreateAcountResponse{}, ierr.ErrConflict).
			Once()

		err := handler.Create(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.True(t, errors.Is(err, ierr.ErrConflict))
	})
}

func TestAccountHandler_FindByID(t *testing.T) {
	handler, mockService, e := setupTestHandler(t)
	accountID := uuid.NewString()
	account := entity.Account{ID: accountID, Email: "test@example.com", Name: "Test User"}

	t.Run("should find account successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+accountID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		mockService.On("FindByID", mock.Anything, accountID).
			Return(account, nil).
			Once()

		err := handler.FindByID(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response entity.Account
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, account, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request on empty ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		err := handler.FindByID(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("should return not found on account not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+accountID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		mockService.On("FindByID", mock.Anything, accountID).
			Return(entity.Account{}, ierr.ErrNotFound).
			Once()

		err := handler.FindByID(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.True(t, errors.Is(err, ierr.ErrNotFound))
	})
}

func TestAccountHandler_GetAccounts(t *testing.T) {
	handler, mockService, e := setupTestHandler(t)
	accounts := []dto.AccountListDTO{{ID: uuid.NewString(), Name: "Test User"}}

	t.Run("should get accounts successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?email=test@example.com", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("GetAccounts", mock.Anything, dto.AccountFilter{Email: "test@example.com"}).
			Return(accounts, nil).
			Once()

		err := handler.GetAccounts(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response []dto.AccountListDTO
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, accounts, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return error on service failure", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService.On("GetAccounts", mock.Anything, dto.AccountFilter{}).
			Return(nil, errors.New("db error")).
			Once()

		err := handler.GetAccounts(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.EqualError(t, err, "db error")
	})
}

func TestAccountHandler_DeleteByID(t *testing.T) {
	handler, mockService, e := setupTestHandler(t)
	accountID := uuid.NewString()

	t.Run("should delete account successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/"+accountID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		mockService.On("DeleteByID", mock.Anything, accountID).
			Return(nil).
			Once()

		err := handler.DeleteByID(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "", rec.Body.String()) // NoContent
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request on empty ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		err := handler.DeleteByID(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, err.Error(), "cannot be empty")
	})
}

func TestAccountHandler_Update(t *testing.T) {
	handler, mockService, e := setupTestHandler(t)
	accountID := uuid.NewString()
	request := dto.UpdateAccountRequest{Email: "updated@example.com", Name: "Updated User"}
	updatedAccount := entity.Account{ID: accountID, Email: "updated@example.com", Name: "Updated User"}
	requestInvalid := dto.UpdateAccountRequest{Email: "invalid-email", Name: "Updated User"}

	t.Run("should update account successfully", func(t *testing.T) {
		jsonBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPut, "/"+accountID, strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		mockService.On("Update", mock.Anything, accountID, request).
			Return(updatedAccount, nil).
			Once()

		err := handler.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response entity.Account
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, updatedAccount, response)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request on invalid email validation", func(t *testing.T) {
		jsonBody, _ := json.Marshal(requestInvalid)
		req := httptest.NewRequest(http.MethodPut, "/"+accountID, strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		err := handler.Update(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("should return bad request on empty ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		err := handler.Update(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("should return bad request on invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/"+accountID, strings.NewReader("{invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		err := handler.Update(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("should return not found on account not found", func(t *testing.T) {
		jsonBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPut, "/"+accountID, strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		mockService.On("Update", mock.Anything, accountID, request).
			Return(entity.Account{}, ierr.ErrNotFound).
			Once()

		err := handler.Update(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.True(t, errors.Is(err, ierr.ErrNotFound))
	})

	t.Run("should return internal server error on generic error", func(t *testing.T) {
		jsonBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPut, "/"+accountID, strings.NewReader(string(jsonBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(accountID)

		mockService.On("Update", mock.Anything, accountID, request).
			Return(entity.Account{}, errors.New("database connection failed")).
			Once()

		err := handler.Update(c)

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, err.Error(), "an unexpected error occurred")
	})
}
