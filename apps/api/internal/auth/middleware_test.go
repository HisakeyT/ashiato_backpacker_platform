package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/gin-gonic/gin"
)

type versionOneLookup struct{}

func (versionOneLookup) FindByID(id uint) (*model.User, error) {
	return &model.User{ID: id, AuthVersion: 1}, nil
}

type versionLookup struct{ versions map[uint]uint }

func (v versionLookup) FindByID(id uint) (*model.User, error) {
	return &model.User{ID: id, AuthVersion: v.versions[id]}, nil
}

func TestAuthMiddlewareRejectsOldGenerationOnlyForChangedUser(t *testing.T) {
	manager := NewJWTManager("test-secret")
	first, err := manager.GenerateToken(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := manager.GenerateToken(2, 1)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.New()
	router.Use(AuthMiddleware(manager, versionLookup{versions: map[uint]uint{1: 2, 2: 1}}))
	router.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	for _, tc := range []struct {
		token string
		want  int
	}{{first, http.StatusUnauthorized}, {second, http.StatusOK}} {
		request := httptest.NewRequest(http.MethodGet, "/test", nil)
		request.Header.Set("Authorization", "Bearer "+tc.token)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != tc.want {
			t.Fatalf("status = %d, want %d", response.Code, tc.want)
		}
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	userID := uint(123)

	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager(t)
	token := generateTestToken(t, jwtManager, userID)

	r := gin.New()

	r.Use(AuthMiddleware(jwtManager, versionOneLookup{}))

	r.GET("/test", func(c *gin.Context) {
		getUserID, exists := c.Get("userID")
		if !exists {
			t.Fatal("user_id was not found in context")
		}

		if getUserID != userID {
			t.Errorf("expected userID 123, got %v", userID)
		}

		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager(t)
	router := newTestRouter(t, jwtManager)

	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidAuthorizationScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager(t)
	router := newTestRouter(t, jwtManager)

	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)

	req.Header.Set("Authorization", "Basic test-token")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := newTestJWTManager(t)

	router := newTestRouter(t, jwtManager)

	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)

	req.Header.Set("Authorization", "Bearer")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

// 共通関数化
func newTestJWTManager(t *testing.T) *JWTManager {
	t.Helper()

	return NewJWTManager("test-secret")
}

func generateTestToken(t *testing.T, jwtManager *JWTManager, userID uint) string {
	t.Helper()

	token, err := jwtManager.GenerateToken(userID, 1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	return token
}

func newTestRouter(t *testing.T, jwtManager *JWTManager) *gin.Engine {
	r := gin.New()

	r.Use(AuthMiddleware(jwtManager, versionOneLookup{}))

	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return r
}
