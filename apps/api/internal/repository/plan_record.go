package repository

import "github.com/HisakeyT/backpacker-platform/internal/model"

type PlanRecordRepository interface {
	List(kind model.RecordKind, planID uint) ([]model.PlanRecord, error)
	Find(kind model.RecordKind, id uint) (*model.PlanRecord, error)
	Create(kind model.RecordKind, record *model.PlanRecord) error
	Update(kind model.RecordKind, record *model.PlanRecord) error
	Delete(kind model.RecordKind, id uint) error
}
