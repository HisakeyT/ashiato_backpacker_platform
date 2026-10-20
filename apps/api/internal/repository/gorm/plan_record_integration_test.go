package gormrepo

import (
	"os"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/database"
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestPlanRecordTablesAndCascade(t *testing.T) {
	if os.Getenv("PLAN_RECORD_DB_TEST") != "1" {
		t.Skip("set PLAN_RECORD_DB_TEST=1 with a migrated local MySQL database")
	}
	db, err := database.New()
	if err != nil {
		t.Fatal(err)
	}
	tx := db.Begin()
	defer tx.Rollback()

	user := model.User{Nickname: "test", Email: "plan-record-test@example.invalid", PasswordHash: "test"}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	day := time.Date(2026, 10, 1, 0, 0, 0, 0, time.Local)
	travel := model.Travel{UserID: user.ID, Title: "test", StartDate: day, EndDate: day, ScheduleMode: model.ScheduleModePlanning}
	if err := tx.Omit("Plans").Create(&travel).Error; err != nil {
		t.Fatal(err)
	}
	plan := model.TravelPlan{TravelID: travel.ID, Place: "test", Content: "test", SortOrder: 1}
	if err := tx.Create(&plan).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewPlanRecordRepository(tx)
	for _, kind := range []model.RecordKind{model.RecordSightseeing, model.RecordStay, model.RecordMeal, model.RecordOther} {
		first := &model.PlanRecord{TravelPlanID: plan.ID, Name: "first"}
		if err := repo.Create(kind, first); err != nil {
			t.Fatalf("%s create: %v", kind, err)
		}
		second := &model.PlanRecord{TravelPlanID: plan.ID, Name: "second"}
		if err := repo.Create(kind, second); err != nil {
			t.Fatalf("%s second create: %v", kind, err)
		}
		list, err := repo.List(kind, plan.ID)
		if err != nil || len(list) != 2 || list[0].SortOrder != 1 || list[1].SortOrder != 2 {
			t.Fatalf("%s list = %+v, err = %v", kind, list, err)
		}
		list[0].LocationText = ptr("location")
		if err := repo.Update(kind, &list[0]); err != nil {
			t.Fatalf("%s update: %v", kind, err)
		}
		found, err := repo.Find(kind, first.ID)
		if err != nil || found.LocationText == nil || *found.LocationText != "location" {
			t.Fatalf("%s find = %+v, err = %v", kind, found, err)
		}
	}
	if err := tx.Delete(&plan).Error; err != nil {
		t.Fatal(err)
	}
	for _, kind := range []model.RecordKind{model.RecordSightseeing, model.RecordStay, model.RecordMeal, model.RecordOther} {
		list, err := repo.List(kind, plan.ID)
		if err != nil || len(list) != 0 {
			t.Fatalf("%s cascade list = %+v, err = %v", kind, list, err)
		}
	}
}

func ptr(value string) *string { return &value }
