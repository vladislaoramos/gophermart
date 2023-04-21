package app

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/gophermart/internal/entity"
	mocks "github.com/vladislaoramos/gophermart/internal/mocks/usecase"
	"github.com/vladislaoramos/gophermart/internal/usecase"
	"github.com/vladislaoramos/gophermart/pkg/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPingHandler(t *testing.T) {
	mockLoyalSystem := &mocks.LoyalSystem{}
	mockLoyalSystem.On("PingRepo", mock.Anything).Return(nil)

	handler := pingHandler(mockLoyalSystem)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	mockLoyalSystem.AssertExpectations(t)
}

func TestUploadOrderHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//t.Parallel()
		user := entity.User{
			ID:    1,
			Login: "test",
		}

		mockLoyalSystem := mocks.NewLoyalSystem(t)
		mockLogger := logger.New("debug", os.Stdout)

		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		claims := map[string]interface{}{"user_id": user.ID}
		tt, tokenString, _ := tokenAuth.Encode(claims)

		handler := uploadOrder(mockLoyalSystem, mockLogger)

		orderNum := 378282246310005
		reqBody, _ := json.Marshal(orderNum)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
		req.Header.Add("Authorization", "Bearer "+tokenString)

		ctx := jwtauth.NewContext(req.Context(), tt, nil)

		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockLoyalSystem.On("UploadOrder", mock.Anything, user.ID, orderNum).Return(true, nil)

		handler(w, req)

		mockLoyalSystem.AssertExpectations(t)
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRegistrationUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		userAuth := entity.UserAuth{
			Login:    "test",
			Password: "test",
		}

		mockLoyalSystem := mocks.NewLoyalSystem(t)
		mockLogger := logger.New("debug", os.Stdout)

		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		handler := registrationUser(mockLoyalSystem, mockLogger, tokenAuth)

		reqBody, _ := json.Marshal(userAuth)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		user := entity.User{
			ID:    1,
			Login: "Test User",
		}

		mockLoyalSystem.On("CreateUser", mock.Anything, userAuth).Return(user, nil)

		handler(w, req)

		mockLoyalSystem.AssertExpectations(t)
		require.Equal(t, http.StatusOK, w.Code)
		require.NotEmpty(t, w.Header().Get("Authorization"))

	})

	t.Run("conflict", func(t *testing.T) {
		t.Parallel()

		userAuth := entity.UserAuth{
			Login:    "test",
			Password: "test",
		}

		mockLoyalSystem := mocks.NewLoyalSystem(t)
		mockLogger := logger.New("debug", os.Stdout)

		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		handler := registrationUser(mockLoyalSystem, mockLogger, tokenAuth)

		reqBody, _ := json.Marshal(userAuth)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		mockLoyalSystem.On("CreateUser", mock.Anything, userAuth).Return(entity.User{}, usecase.ErrConflict)

		handler(w, req)

		mockLoyalSystem.AssertExpectations(t)
		require.Equal(t, http.StatusConflict, w.Code)
		require.Empty(t, w.Header().Get("Authorization"))

	})
}

func TestLoginUser(t *testing.T) {
	t.Run("bad_request", func(t *testing.T) {
		t.Parallel()

		mockLoyalSystem := mocks.NewLoyalSystem(t)
		mockLogger := logger.New("debug", os.Stdout)

		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		handler := loginUser(mockLoyalSystem, mockLogger, tokenAuth)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("bad json"))
		w := httptest.NewRecorder()

		handler(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("auth_error", func(t *testing.T) {
		t.Parallel()

		mockLoyalSystem := mocks.NewLoyalSystem(t)
		mockLogger := logger.New("debug", os.Stdout)

		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		handler := loginUser(mockLoyalSystem, mockLogger, tokenAuth)

		userAuth := entity.UserAuth{
			Login:    "test",
			Password: "test",
		}

		reqBody, _ := json.Marshal(userAuth)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		mockLoyalSystem.On("CheckUser", mock.Anything, userAuth).Return(entity.User{}, usecase.ErrUnauthorized)

		handler(w, req)

		mockLoyalSystem.AssertExpectations(t)
		require.Equal(t, http.StatusUnauthorized, w.Code)
		require.Empty(t, w.Header().Get("Authorization"))
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockLoyalSystem := mocks.NewLoyalSystem(t)
		mockLogger := logger.New("debug", os.Stdout)

		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		handler := loginUser(mockLoyalSystem, mockLogger, tokenAuth)

		userAuth := entity.UserAuth{
			Login:    "test",
			Password: "test",
		}

		reqBody, _ := json.Marshal(userAuth)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		user := entity.User{
			ID:    1,
			Login: "Test User",
		}

		mockLoyalSystem.On("CheckUser", mock.Anything, userAuth).Return(user, nil)

		handler(w, req)

		mockLoyalSystem.AssertExpectations(t)
		require.Equal(t, http.StatusOK, w.Code)
		require.NotEmpty(t, w.Header().Get("Authorization"))
	})
}
