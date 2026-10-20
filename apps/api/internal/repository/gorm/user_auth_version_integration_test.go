package gormrepo

import (
	"os"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/database"
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestBumpAuthVersionInTransaction(t *testing.T) {
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
		{Nickname: "one", Email: "auth-version-one@example.invalid", PasswordHash: "test", AuthVersion: 1},
		{Nickname: "two", Email: "auth-version-two@example.invalid", PasswordHash: "test", AuthVersion: 1},
	}
	if err := tx.Create(&users).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewUserRepository(tx)
	if err := repo.BumpAuthVersion(users[0].ID); err != nil {
		t.Fatal(err)
	}
	users[0].Nickname = "renamed"
	if err := repo.Update(&users[0]); err != nil {
		t.Fatal(err)
	}
	for index, want := range []uint{2, 1} {
		user, err := repo.FindByID(users[index].ID)
		if err != nil || user.AuthVersion != want {
			t.Fatalf("user %d = %+v, err = %v", index, user, err)
		}
	}
	if err := repo.BumpAuthVersion(0); err != model.ErrUserNotFound {
		t.Fatalf("missing user: %v", err)
	}
}
