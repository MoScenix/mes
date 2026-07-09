package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Name        string `gorm:"type:varchar(100);not null"`
	Unit        string `gorm:"type:varchar(20);not null"`
	Description string `gorm:"type:varchar(255);not null;default:''"`

	TotalCount       int64 `gorm:"not null;default:0"`
	InStockCount     int64 `gorm:"not null;default:0"`
	ReservedCount    int64 `gorm:"not null;default:0"`
	OutStockCount    int64 `gorm:"not null;default:0"`
	PendingCount     int64 `gorm:"not null;default:0"`
	QualifiedCount   int64 `gorm:"not null;default:0"`
	UnqualifiedCount int64 `gorm:"not null;default:0"`
	AvailableCount   int64 `gorm:"not null;default:0"`
}

type ItemQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewItemQuery(ctx context.Context, db *gorm.DB) *ItemQuery {
	return &ItemQuery{ctx: ctx, db: db}
}

func (q *ItemQuery) Create(item *Item) error {
	return q.db.WithContext(q.ctx).Create(item).Error
}

func (q *ItemQuery) Update(id uint, updates map[string]any) error {
	return q.db.WithContext(q.ctx).Model(&Item{}).Where("id = ?", id).Updates(updates).Error
}

func (q *ItemQuery) Get(id uint) (Item, error) {
	var item Item
	err := q.db.WithContext(q.ctx).First(&item, id).Error
	return item, err
}

func (q *ItemQuery) List(pageSize int, namePrefix string, cursorUpdatedAt *time.Time, cursorID uint) ([]Item, bool, error) {
	var items []Item
	db := q.db.WithContext(q.ctx).Model(&Item{})
	if namePrefix != "" {
		db = db.Where("name LIKE ?", namePrefix+"%")
	}
	if cursorUpdatedAt != nil && cursorID > 0 {
		db = db.Where("(updated_at < ? OR (updated_at = ? AND id < ?))", *cursorUpdatedAt, *cursorUpdatedAt, cursorID)
	}
	err := db.Order("updated_at DESC, id DESC").Limit(pageSize + 1).Find(&items).Error
	if err != nil {
		return nil, false, err
	}
	hasMore := len(items) > pageSize
	if hasMore {
		items = items[:pageSize]
	}
	return items, hasMore, nil
}
