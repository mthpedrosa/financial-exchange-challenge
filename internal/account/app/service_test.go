package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/app"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/dto"
	"github.com/mthpedrosa/financial-exchange-challenge/internal/account/domain/entity"
	mocks "github.com/mthpedrosa/financial-exchange-challenge/mocks/domain/port"

	"github.com/mthpedrosa/financial-exchange-challenge/pkg/ierr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountApp_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("should create account successfully when email is new", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		newID := uuid.NewString()
		testRequest := dto.CreateAccountRequest{Email: "new@example.com", Name: "New User"}

		mockRepo.On("GetAccounts", mock.Anything, mock.AnythingOfType("entity.AccountFilter")).
			Return([]entity.Account{}, nil).
			Once()
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("entity.Account")).
			Return(newID, nil).
			Once()

		response, err := accountService.Create(ctx, testRequest)

		assert.NoError(t, err)
		assert.Equal(t, newID, response.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return conflict error when email already exists", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)
		testRequest := dto.CreateAccountRequest{Email: "exists@example.com", Name: "Existing User"}

		mockRepo.On("GetAccounts", mock.Anything, mock.Anything).
			Return([]entity.Account{{ID: uuid.NewString()}}, nil).
			Once()

		_, err := accountService.Create(ctx, testRequest)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ierr.ErrConflict))
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)
		testRequest := dto.CreateAccountRequest{Email: "fail@example.com", Name: "Fail User"}

		mockRepo.On("GetAccounts", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error")).
			Once()

		_, err := accountService.Create(ctx, testRequest)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestAccountApp_FindByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return account successfully when account exists", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		accountID := uuid.NewString()
		expectedAccount := entity.Account{ID: accountID, Name: "Found User"}

		mockRepo.On("FindByID", mock.Anything, accountID).
			Return(expectedAccount, nil).
			Once()

		result, err := accountService.FindByID(ctx, accountID)

		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return not found error when account does not exist", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)
		nonExistentID := uuid.NewString()

		mockRepo.On("FindByID", mock.Anything, nonExistentID).
			Return(entity.Account{}, ierr.ErrNotFound).
			Once()

		_, err := accountService.FindByID(ctx, nonExistentID)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ierr.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestAccountApp_GetAccounts(t *testing.T) {
	ctx := context.Background()

	t.Run("should return accounts successfully", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		filters := dto.AccountFilter{Email: "test@example.com"}
		accounts := []entity.Account{{ID: uuid.NewString(), Email: "test@example.com", Name: "Test User"}}
		expectedDTOs := []dto.AccountListDTO{{ID: accounts[0].ID, Name: accounts[0].Name}}

		mockRepo.On("GetAccounts", mock.Anything, mock.AnythingOfType("entity.AccountFilter")).
			Return(accounts, nil).
			Once()

		result, err := accountService.GetAccounts(ctx, filters)

		assert.NoError(t, err)
		assert.Equal(t, expectedDTOs, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		filters := dto.AccountFilter{}

		mockRepo.On("GetAccounts", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error")).
			Once()

		_, err := accountService.GetAccounts(ctx, filters)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestAccountApp_DeleteByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete account successfully", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		accountID := uuid.NewString()

		mockRepo.On("DeleteByID", mock.Anything, accountID).
			Return(nil).
			Once()

		err := accountService.DeleteByID(ctx, accountID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		accountID := uuid.NewString()

		mockRepo.On("DeleteByID", mock.Anything, accountID).
			Return(errors.New("db error")).
			Once()

		err := accountService.DeleteByID(ctx, accountID)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestAccountApp_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("should update account successfully", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		accountID := uuid.NewString()
		request := dto.UpdateAccountRequest{Email: "updated@example.com", Name: "Updated User"}
		existingAccount := entity.Account{ID: accountID, Email: "old@example.com"}

		mockRepo.On("FindByID", ctx, accountID).
			Return(existingAccount, nil).
			Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("entity.Account")).
			Return(nil).
			Once()

		result, err := accountService.Update(ctx, accountID, request)

		assert.NoError(t, err)
		assert.Equal(t, accountID, result.ID)
		assert.Equal(t, request.Email, result.Email)
		assert.Equal(t, request.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return not found error when account does not exist", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		accountID := uuid.NewString()
		request := dto.UpdateAccountRequest{}

		mockRepo.On("FindByID", ctx, accountID).
			Return(entity.Account{}, ierr.ErrNotFound).
			Once()

		_, err := accountService.Update(ctx, accountID, request)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ierr.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error on repository failure", func(t *testing.T) {
		mockRepo := mocks.NewAccountRepository(t)
		accountService := app.NewAccountApp(mockRepo)

		accountID := uuid.NewString()
		request := dto.UpdateAccountRequest{Email: "updated@example.com", Name: "Updated User"}
		existingAccount := entity.Account{ID: accountID, Email: "old@example.com"}

		mockRepo.On("FindByID", ctx, accountID).
			Return(existingAccount, nil).
			Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("entity.Account")).
			Return(errors.New("db error")).
			Once()

		_, err := accountService.Update(ctx, accountID, request)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}
