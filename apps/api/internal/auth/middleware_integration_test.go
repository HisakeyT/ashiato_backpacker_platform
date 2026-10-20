package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/database"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	gormrepo "github.com/HisakeyT/backpacker-platform/internal/repository/gorm"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareAfterVersionBumpInMySQL(t *testing.T) {
	if os.Getenv("AUTH_VERSION_DB_TEST") != "1" {
		t.Skip("set AUTH_VERSION_DB_TEST=1 with a migrated local MySQL database")
	}
	db, err := database.New()
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer tx.Rollback()
	users := []model.User{
		{Nickname: "one", Email: "auth-token-one@example.invalid", PasswordHash: "test", AuthVersion: 1},
		{Nickname: "two", Email: "auth-token-two@example.invalid", PasswordHash: "test", AuthVersion: 1},
	}
	if err := tx.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	repo := gormrepo.NewUserRepository(tx)
	manager := NewJWTManager("test-secret")
	router := gin.New()
	router.Use(AuthMiddleware(manager, repo))
	router.GET("/private", func(c *gin.Context) { c.Status(http.StatusOK) })
	first, err := manager.GenerateToken(users[0].ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := manager.GenerateToken(users[1].ID, 1)
	if err != nil {
		t.Fatal(err)
	}
	check := func(token string, want int) {
		request := httptest.NewRequest(http.MethodGet, "/private", nil)
		request.Header.Set("Authorization", "Bearer "+token)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != want {
			t.Fatalf("status = %d, want %d", response.Code, want)
		}
	}
	check(first, http.StatusOK)
	if err := repo.BumpAuthVersion(users[0].ID); err != nil {
		t.Fatal(err)
	}
	check(first, http.StatusUnauthorized)
	check(second, http.StatusOK)
}
