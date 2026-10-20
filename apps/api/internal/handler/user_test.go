package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepository struct {
	user *model.User
	err  error
}

type versionOneUserLookup struct{}

func (versionOneUserLookup) FindByID(id uint) (*model.User, error) {
	return &model.User{ID: id, AuthVersion: 1}, nil
}

func (f *fakeRepository) Create(user *model.User) error {
	return nil
}

func (f *fakeRepository) FindByID(id uint) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeRepository) FindByEmail(email string) (*model.User, error) {
	return f.user, f.err
}

func (f *fakeRepository) Update(user *model.User) error {
	return nil
}

func TestHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    string
		repository repository.UserRepository
		wantStatus int
	}{
		{
			name:       "正常に登録できる",
			request:    `{"nickname":"test","email":"test@example.com","password":"password","password_confirmation":"password"}`,
			repository: &fakeRepository{},
			wantStatus: http.StatusCreated,
		},
		{
			name:    "メールアドレスが重複している",
			request: `{"nickname":"test","email":"test@example.com","password":"password","password_confirmation":"password"}`,
			repository: &fakeRepository{
				user: &model.User{
					ID:    1,
					Email: "test@example.com",
				},
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "パスワードが一致しない",
			request:    `{"nickname":"test","email":"test@example.com","password":"password","password_confirmation":"wrong-password"}`,
			repository: &fakeRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "不正なJSON",
			request:    `{"nickname":}`,
			repository: &fakeRepository{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			jwtManager := auth.NewJWTManager("secret")
			useCase := usecase.NewUserUsecase(tt.repository, jwtManager)
			handler := NewHandler(useCase)

			router.POST("/register", handler.Register)

			req := httptest.NewRequest(
				http.MethodPost,
				"/register",
				strings.NewReader(tt.request),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	email := "test@example.com"
	nickName := "test"

	password := "password"
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		request      string
		repository   repository.UserRepository
		wantStatus   int
		wantResponse *LoginResponse
	}{
		{
			name:    "正常にログインできる",
			request: `{"email":"test@example.com","password":"password"}`,
			repository: &fakeRepository{
				user: &model.User{
					ID:           1,
					AuthVersion:  1,
					Nickname:     nickName,
					Email:        email,
					PasswordHash: string(passwordHash),
				},
			},
			wantStatus: http.StatusOK,
			wantResponse: &LoginResponse{
				Token: "dummy-token",
				User: UserResponse{
					ID:       1,
					Nickname: nickName,
					Email:    email,
				},
			},
		},
		{
			name:    "認証情報が不正",
			request: `{"email":"test@example.com","password":"wrong-password"}`,
			repository: &fakeRepository{
				user: &model.User{
					ID:           1,
					Nickname:     nickName,
					Email:        email,
					PasswordHash: string(passwordHash),
				},
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "不正なJSON",
			request:    `{"email":}`,
			repository: &fakeRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "ユーザーが存在しない",
			request: `{"email":"unknown@example.com","password":"password"}`,
			repository: &fakeRepository{
				err: model.ErrUserNotFound,
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			jwtManager := auth.NewJWTManager("secret")
			useCase := usecase.NewUserUsecase(tt.repository, jwtManager)
			handler := NewHandler(useCase)

			router.POST("/login", handler.Login)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(tt.request),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}

			if tt.wantResponse != nil {
				response := decodeLoginResponse(t, rec)

				if response.Token == "" {
					t.Error("expected token to be generated")
				}

				if response.User != tt.wantResponse.User {
					t.Errorf(
						"expected user %+v, got %+v",
						tt.wantResponse.User,
						response.User,
					)
				}
			}
		})
	}
}

func decodeLoginResponse(t *testing.T, rec *httptest.ResponseRecorder) LoginResponse {
	t.Helper()

	var response LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return response
}

func TestHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	repository := &fakeRepository{
		user: &model.User{
			ID:       123,
			Nickname: "test",
			Email:    "test@example.com",
		},
	}

	jwtManager := auth.NewJWTManager("secret")
	useCase := usecase.NewUserUsecase(repository, jwtManager)
	handler := NewHandler(useCase)

	router.Use(auth.AuthMiddleware(jwtManager, versionOneUserLookup{}))
	router.GET("/users/me", handler.Me)

	token, err := jwtManager.GenerateToken(123, 1)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/me",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response UserResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != 123 {
		t.Errorf("expected user ID %d, got %d", 123, response.ID)
	}

	if response.Nickname != "test" {
		t.Errorf("expected nickname %q, got %q", "test", response.Nickname)
	}

	if response.Email != "test@example.com" {
		t.Errorf("expected email %q, got %q", "test@example.com", response.Email)
	}
}

func TestHandler_Me_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repository := &fakeRepository{
		err: model.ErrUserNotFound,
	}

	jwtManager := auth.NewJWTManager("secret")
	useCase := usecase.NewUserUsecase(repository, jwtManager)
	handler := NewHandler(useCase)

	router := gin.New()
	router.Use(auth.AuthMiddleware(jwtManager, versionOneUserLookup{}))
	router.GET("/users/me", handler.Me)

	token, err := jwtManager.GenerateToken(123, 1)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/me",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestHandler_UpdateMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepository := &mock.MockUserRepository{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return &model.User{
				ID:       1,
				Nickname: "old-nickname",
				Email:    "test@example.com",
			}, nil
		},
		UpdateFunc: func(user *model.User) error {
			return nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := usecase.NewUserUsecase(mockRepository, jwtManager)
	handler := NewHandler(useCase)

	router := gin.New()
	router.PATCH("/users/me", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.UpdateMe(c)
	})

	body := `{"nickname":"new-nickname"}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/users/me",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response UserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("expected ID 1, got %d", response.ID)
	}

	if response.Nickname != "new-nickname" {
		t.Errorf("expected nickname new-nickname, got %s", response.Nickname)
	}

	if response.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", response.Email)
	}
}
