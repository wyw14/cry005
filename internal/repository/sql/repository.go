package sqlrepo

import (
	"context"

	"github.com/wyw14/cry005/internal/domain"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

type snapshotRow struct {
	ID, State, Secondary string
	Version              int
}

func (r *Repository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)

}

func (r *Repository) ReplaceSnapshot(ctx context.Context, next domain.Snapshot, failAt string) error {
	return r.db.WithContext(ctx).Model(&snapshotRow{}).Where("id = ?", next.ID).Updates(map[string]any{"state": next.State, "version": next.Version}).Error

}
