package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/wyw14/cry005/internal/domain"
)

type onceResult struct {
	item domain.Item
	err  error
	done chan struct{}
}

type Store struct {
	mu       sync.RWMutex
	items    map[string]domain.Item
	once     map[string]*onceResult
	events   []domain.Event
	snapshot domain.Snapshot
}

func New() *Store {
	return &Store{items: map[string]domain.Item{}, once: map[string]*onceResult{}, snapshot: domain.Snapshot{ID: "workflow", State: "queued", Secondary: "ready", Version: 1}}
}

func (s *Store) ListAll(ctx context.Context) ([]domain.Item, error) {
	if err := domain.CheckContext(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Item, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *Store) ListVisible(ctx context.Context, scope string) ([]domain.Item, error) {
	if err := domain.CheckContext(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Item, 0)
	for _, item := range s.items {
		if item.Scope == scope {
			out = append(out, item)
		}
	}
	return out, nil

}

func (s *Store) DoOnce(ctx context.Context, key string, fn func() (domain.Item, error)) (domain.Item, error) {
	s.mu.Lock()
	if existing, ok := s.once[key]; ok {
		done := existing.done
		s.mu.Unlock()
		select {
		case <-done:
			return existing.item, existing.err
		case <-ctx.Done():
			return domain.Item{}, ctx.Err()
		}
	}
	entry := &onceResult{done: make(chan struct{})}
	s.once[key] = entry
	s.mu.Unlock()
	item, err := fn()
	s.mu.Lock()
	entry.item, entry.err = item, err
	close(entry.done)
	s.mu.Unlock()
	return item, err

}

func (s *Store) AppendEvent(ctx context.Context, event domain.Event) error {
	if err := domain.CheckContext(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
	return nil
}

func (s *Store) Replay(ctx context.Context, after int64) ([]domain.Event, error) {
	if err := domain.CheckContext(ctx); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Event, 0)
	for _, event := range s.events {
		if event.ID > after {
			out = append(out, event)
		}
	}

	return out, nil
}

func (s *Store) SaveAfter(ctx context.Context, item domain.Item, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}

	if err := domain.CheckContext(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item.ID] = item
	return nil
}

func (s *Store) Count(ctx context.Context) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

func (s *Store) Snapshot(ctx context.Context) (domain.Snapshot, error) {
	if err := domain.CheckContext(ctx); err != nil {
		return domain.Snapshot{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.snapshot, nil
}

func (s *Store) ReplaceSnapshot(ctx context.Context, next domain.Snapshot, failAt string) error {
	if err := domain.CheckContext(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshot.State = next.State
	s.snapshot.Version = next.Version
	if failAt == "member" {
		return errors.New("simulated secondary write failure")
	}
	s.snapshot.Secondary = next.Secondary
	return nil

}
