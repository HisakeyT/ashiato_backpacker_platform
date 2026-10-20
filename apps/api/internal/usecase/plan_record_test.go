package usecase

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
)

type recordTravelRepo struct{ travel *model.Travel }

func (r recordTravelRepo) FindByID(uint) (*model.Travel, error) { return r.travel, nil }

type recordPlanRepo struct {
	repository.TravelPlanRepository
	plan *model.TravelPlan
}

func (r recordPlanRepo) FindByID(uint) (*model.TravelPlan, error) { return r.plan, nil }

type recordRepo struct {
	repository.PlanRecordRepository
	created *model.PlanRecord
	found   *model.PlanRecord
	listed  bool
}

func (r *recordRepo) Create(_ model.RecordKind, record *model.PlanRecord) error {
	r.created = record
	return nil
}
func (r *recordRepo) Find(_ model.RecordKind, _ uint) (*model.PlanRecord, error) {
	return r.found, nil
}
func (r *recordRepo) List(_ model.RecordKind, _ uint) ([]model.PlanRecord, error) {
	r.listed = true
	return nil, nil
}

func TestPlanRecordAuthorizationAndValidation(t *testing.T) {
	travel := &model.Travel{ID: 1, UserID: 2, IsPublic: false}
	plan := &model.TravelPlan{ID: 3, TravelID: 1}
	records := &recordRepo{}
	uc := NewPlanRecordUsecase(recordTravelRepo{travel}, recordPlanRepo{plan: plan}, records)

	if _, err := uc.List(0, 1, 3, model.RecordSightseeing, true); !errors.Is(err, model.ErrTravelNotFound) {
		t.Fatalf("private public list error = %v", err)
	}
	if records.listed {
		t.Fatal("private records were queried")
	}
	if _, err := uc.Create(4, 1, 3, model.RecordSightseeing, RecordInput{Name: "寺"}); !errors.Is(err, model.ErrUserNotAuthorized) {
		t.Fatalf("foreign create error = %v", err)
	}
	plan.TravelID = 9
	if _, err := uc.List(2, 1, 3, model.RecordSightseeing, false); !errors.Is(err, model.ErrTravelPlanNotFound) {
		t.Fatalf("foreign plan error = %v", err)
	}
	plan.TravelID = 1
	badURL := "javascript:alert(1)"
	if _, err := uc.Create(2, 1, 3, model.RecordSightseeing, RecordInput{Name: "寺", MapURL: &badURL}); !errors.Is(err, model.ErrInvalidPlanRecord) {
		t.Fatalf("unsafe URL error = %v", err)
	}
	if records.created != nil {
		t.Fatal("invalid record was saved")
	}
	if _, err := uc.Create(2, 1, 3, model.RecordSightseeing, RecordInput{Name: "寺"}); err != nil {
		t.Fatalf("valid create: %v", err)
	}
	if records.created == nil || records.created.Cost != nil {
		t.Fatalf("saved record = %+v", records.created)
	}
}

func TestPlanRecordRejectsForeignRecordAndInvalidMealType(t *testing.T) {
	records := &recordRepo{found: &model.PlanRecord{ID: 8, TravelPlanID: 99, Name: "昼食"}}
	uc := NewPlanRecordUsecase(
		recordTravelRepo{&model.Travel{ID: 1, UserID: 2}},
		recordPlanRepo{plan: &model.TravelPlan{ID: 3, TravelID: 1}}, records,
	)
	if err := uc.Delete(2, 1, 3, 8, model.RecordMeal); !errors.Is(err, model.ErrPlanRecordNotFound) {
		t.Fatalf("foreign record error = %v", err)
	}
	badType := "brunch"
	if _, err := uc.Create(2, 1, 3, model.RecordMeal, RecordInput{Name: "昼食", MealType: &badType}); !errors.Is(err, model.ErrInvalidPlanRecord) {
		t.Fatalf("invalid meal type error = %v", err)
	}
}
