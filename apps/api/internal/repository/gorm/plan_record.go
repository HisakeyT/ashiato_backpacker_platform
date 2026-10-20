package gormrepo

import (
	"errors"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PlanRecordRepository struct{ db *gorm.DB }

func NewPlanRecordRepository(db *gorm.DB) *PlanRecordRepository {
	return &PlanRecordRepository{db: db}
}

func recordTable(kind model.RecordKind) (string, error) {
	if !kind.Valid() {
		return "", model.ErrInvalidPlanRecord
	}
	return string(kind), nil
}

func (r *PlanRecordRepository) List(kind model.RecordKind, planID uint) ([]model.PlanRecord, error) {
	table, err := recordTable(kind)
	if err != nil {
		return nil, err
	}
	var records []model.PlanRecord
	q := r.db.Table(table).Where("travel_plan_id = ?", planID).Order("sort_order ASC, id ASC")
	if err := q.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *PlanRecordRepository) Find(kind model.RecordKind, id uint) (*model.PlanRecord, error) {
	table, err := recordTable(kind)
	if err != nil {
		return nil, err
	}
	var record model.PlanRecord
	if err := r.db.Table(table).Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrPlanRecordNotFound
		}
		return nil, err
	}
	return &record, nil
}

func (r *PlanRecordRepository) Create(kind model.RecordKind, record *model.PlanRecord) error {
	table, err := recordTable(kind)
	if err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		var plan model.TravelPlan
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&plan, record.TravelPlanID).Error; err != nil {
			return err
		}
		var maxOrder int
		if err := tx.Table(table).Where("travel_plan_id = ?", record.TravelPlanID).
			Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder).Error; err != nil {
			return err
		}
		record.SortOrder = maxOrder + 1
		q := tx.Table(table)
		if kind != model.RecordMeal {
			q = q.Omit("MealType")
		}
		return q.Create(record).Error
	})
}

func (r *PlanRecordRepository) Update(kind model.RecordKind, record *model.PlanRecord) error {
	table, err := recordTable(kind)
	if err != nil {
		return err
	}
	values := map[string]any{
		"name": record.Name, "location_text": record.LocationText, "map_url": record.MapURL,
		"meal_type": record.MealType,
	}
	if kind != model.RecordMeal {
		delete(values, "meal_type")
	}
	return r.db.Table(table).Where("id = ?", record.ID).Updates(values).Error
}

func (r *PlanRecordRepository) Delete(kind model.RecordKind, id uint) error {
	table, err := recordTable(kind)
	if err != nil {
		return err
	}
	return r.db.Table(table).Where("id = ?", id).Delete(&model.PlanRecord{}).Error
}
