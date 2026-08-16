package application

import (
	"context"
	"time"

	"sort"

	"github.com/google/uuid"
	"github.com/wyw14/cry005/internal/domain"
)

type Repository interface {
	ListAll(context.Context) ([]domain.Item, error)
	ListVisible(context.Context, string) ([]domain.Item, error)
	DoOnce(context.Context, string, func() (domain.Item, error)) (domain.Item, error)
	AppendEvent(context.Context, domain.Event) error
	Replay(context.Context, int64) ([]domain.Event, error)
	SaveAfter(context.Context, domain.Item, time.Duration) error
	Count(context.Context) int
	Snapshot(context.Context) (domain.Snapshot, error)
	ReplaceSnapshot(context.Context, domain.Snapshot, string) error
}

type Service struct{ repo Repository }

func New(repo Repository) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context, actorScope string) ([]domain.Item, error) {
	items, err := s.repo.ListVisible(ctx, actorScope)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Item, 0, len(items))
	for _, item := range items {
		if domain.CanAccess(actorScope, item.Scope) {
			out = append(out, item)
		}
	}
	return out, nil

}

func (s *Service) Create(ctx context.Context, scope, actor, idempotencyKey, payload string) (domain.Item, error) {
	scopeKey := domain.IdempotencyScope(scope, actor, "create", idempotencyKey)

	return s.repo.DoOnce(ctx, scopeKey, func() (domain.Item, error) {
		time.Sleep(2 * time.Millisecond)
		item := domain.Item{ID: uuid.NewString(), Scope: scope, OwnerID: actor, State: "queued", Payload: payload, Version: 1, CreatedAt: time.Now().UTC()}
		if err := s.repo.SaveAfter(ctx, item, 0); err != nil {
			return domain.Item{}, err
		}
		return item, nil
	})
}

func (s *Service) RecordEvent(ctx context.Context, itemID, state string, id int64) error {
	return s.repo.AppendEvent(ctx, domain.Event{ID: id, ItemID: itemID, State: state, CreatedAt: time.Now().UTC()})
}

func (s *Service) Replay(ctx context.Context, after int64) ([]domain.Event, error) {
	events, err := s.repo.Replay(ctx, after)
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].ID < events[j].ID })
	return events, nil

}

func (s *Service) RunCancelable(ctx context.Context, item domain.Item, delay time.Duration) error {
	background := context.Background()
	return s.repo.SaveAfter(background, item, delay)

}

func (s *Service) ApplyAtomic(ctx context.Context, targetState, failAt string) error {
	before, err := s.repo.Snapshot(ctx)
	if err != nil {
		return err
	}
	next := domain.ApplyWorkflow(before, targetState)
	if err := s.repo.ReplaceSnapshot(ctx, next, failAt); err != nil {
		return err
	}
	return nil

}
