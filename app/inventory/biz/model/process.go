package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	DraftStatusDraft     int32 = 1
	DraftStatusSubmitted int32 = 2
	DraftStatusDone      int32 = 3
)

type Process struct {
	ID        uint `gorm:"primarykey;index:idx_process_updated_id,priority:3;index:idx_process_owner_updated_id,priority:4;index:idx_process_owner_status_updated_id,priority:5;index:idx_process_status_updated_id,priority:4;index:idx_process_item_updated_id,priority:4;index:idx_process_item_status_updated_id,priority:5;index:idx_process_name_id,priority:3"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_process_updated_id,priority:2;index:idx_process_owner_updated_id,priority:3;index:idx_process_owner_status_updated_id,priority:4;index:idx_process_status_updated_id,priority:3;index:idx_process_item_updated_id,priority:3;index:idx_process_item_status_updated_id,priority:4"`
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_process_updated_id,priority:1;index:idx_process_owner_updated_id,priority:1;index:idx_process_owner_status_updated_id,priority:1;index:idx_process_status_updated_id,priority:1;index:idx_process_item_updated_id,priority:1;index:idx_process_item_status_updated_id,priority:1;index:idx_process_name_id,priority:1"`

	ItemID      uint   `gorm:"not null;index:idx_process_item_updated_id,priority:2;index:idx_process_item_status_updated_id,priority:2"`
	Item        Item   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OwnerUserID int64  `gorm:"not null;index:idx_process_owner_updated_id,priority:2;index:idx_process_owner_status_updated_id,priority:2"`
	Name        string `gorm:"type:varchar(100);not null;index:idx_process_name_id,priority:2,length:64"`
	Description string `gorm:"type:varchar(255);not null;default:''"`
	Status      int32  `gorm:"not null;default:1;index:idx_process_owner_status_updated_id,priority:3;index:idx_process_status_updated_id,priority:2;index:idx_process_item_status_updated_id,priority:3"`

	Items []ProcessItem `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ProcessItem struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ProcessID     uint `gorm:"not null;uniqueIndex:idx_process_item_consume,priority:1"`
	Process       Process
	ConsumeItemID uint  `gorm:"not null;uniqueIndex:idx_process_item_consume,priority:2;index:idx_process_item_consume_reverse,priority:1"`
	ConsumeItem   Item  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Quantity      int64 `gorm:"not null;default:0"`
}

type ProcessQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewProcessQuery(ctx context.Context, db *gorm.DB) *ProcessQuery {
	return &ProcessQuery{ctx: ctx, db: db}
}

func (q *ProcessQuery) Create(process *Process) error {
	return q.db.WithContext(q.ctx).Create(process).Error
}

func (q *ProcessQuery) UpdateDraft(id uint, updates map[string]any) error {
	return q.db.WithContext(q.ctx).Model(&Process{}).
		Where("id = ? AND status = ?", id, DraftStatusDraft).
		Updates(updates).Error
}

func (q *ProcessQuery) DeleteDraft(id uint) error {
	return q.db.WithContext(q.ctx).
		Where("id = ? AND status = ?", id, DraftStatusDraft).
		Delete(&Process{}).Error
}

func (q *ProcessQuery) SubmitDraft(id uint) error {
	return q.db.WithContext(q.ctx).Model(&Process{}).
		Where("id = ? AND status = ?", id, DraftStatusDraft).
		Update("status", DraftStatusSubmitted).Error
}

func (q *ProcessQuery) Get(id uint, withItems bool) (Process, error) {
	var process Process
	db := q.db.WithContext(q.ctx).Preload("Item")
	if withItems {
		db = db.Preload("Items.ConsumeItem")
	}
	err := db.First(&process, id).Error
	return process, err
}

func (q *ProcessQuery) List(pageSize int, ownerUserID int64, itemID uint, status int32, namePrefix string, itemNamePrefix string, sinceTime *time.Time, cursorUpdatedAt *time.Time, cursorID uint) ([]Process, bool, error) {
	var processes []Process
	db := q.db.WithContext(q.ctx).Model(&Process{})
	if itemNamePrefix != "" {
		db = db.Joins("JOIN items ON items.id = processes.item_id AND items.deleted_at IS NULL").
			Where("items.name LIKE ?", itemNamePrefix+"%")
	}
	if ownerUserID > 0 {
		db = db.Where("processes.owner_user_id = ?", ownerUserID)
	}
	if itemID > 0 {
		db = db.Where("processes.item_id = ?", itemID)
	}
	if status > 0 {
		db = db.Where("processes.status = ?", status)
	}
	if namePrefix != "" {
		db = db.Where("processes.name LIKE ?", namePrefix+"%")
	}
	if sinceTime != nil {
		db = db.Where("processes.updated_at > ?", *sinceTime)
	}
	if cursorUpdatedAt != nil && cursorID > 0 {
		db = db.Where("(processes.updated_at < ? OR (processes.updated_at = ? AND processes.id < ?))", *cursorUpdatedAt, *cursorUpdatedAt, cursorID)
	}
	err := db.Preload("Item").
		Order("processes.updated_at DESC, processes.id DESC").
		Limit(pageSize + 1).
		Find(&processes).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := len(processes) > pageSize
	if hasMore {
		processes = processes[:pageSize]
	}
	return processes, hasMore, nil
}
