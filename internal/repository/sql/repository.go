package sqlrepo

import (
	"context"

	"errors"

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
	return r.db.WithContext(context.Background())

}

func (r *Repository) ReplaceSnapshot(ctx context.Context, next domain.Snapshot, failAt string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if failAt == "member" {
			return errors.New("simulated secondary write failure")
		}
		return tx.Model(&snapshotRow{}).Where("id = ?", next.ID).Updates(map[string]any{"state": next.State, "secondary": next.Secondary, "version": next.Version}).Error
	})

}
