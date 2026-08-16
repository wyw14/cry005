package application_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/wyw14/cry005/internal/application"
	"github.com/wyw14/cry005/internal/domain"
	"github.com/wyw14/cry005/internal/repository/memory"
	httptransport "github.com/wyw14/cry005/internal/transport/http"
)

func TestRoomMembersStayWithinRoom(t *testing.T) {
	repo := memory.New()
	svc := application.New(repo)
	_, _ = svc.Create(context.Background(), "alpha", "owner-a", "a", "visible")
	_, _ = svc.Create(context.Background(), "beta", "owner-b", "b", "hidden")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
	req.Header.Set("X-Scope", "alpha")
	rr := httptest.NewRecorder()
	httptransport.New(svc).Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var items []domain.Item
	if err := json.Unmarshal(rr.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Scope != "alpha" {
		t.Fatalf("scope leak: %+v", items)
	}
}

func TestConcurrentJoinUsesOneMembership(t *testing.T) {
	repo := memory.New()
	svc := application.New(repo)
	const workers = 24
	ids := make(chan string, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			item, err := svc.Create(context.Background(), "alpha", "actor", "same-key", "payload")
			if err != nil {
				errs <- err
				return
			}
			ids <- item.ID
		}()
	}
	wg.Wait()
	close(ids)
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	first := ""
	for id := range ids {
		if first == "" {
			first = id
		}
		if id != first {
			t.Fatalf("idempotent requests returned different IDs: %s and %s", first, id)
		}
	}
	if got := repo.Count(context.Background()); got != 1 {
		t.Fatalf("stored %d items, want 1", got)
	}
}

func TestEventReplayPreservesOrder(t *testing.T) {
	repo := memory.New()
	svc := application.New(repo)
	_ = svc.RecordEvent(context.Background(), "x", "running", 1)
	_ = svc.RecordEvent(context.Background(), "x", "running", 2)
	_ = svc.RecordEvent(context.Background(), "x", "completed", 3)
	events, err := svc.Replay(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 || events[0].ID != 2 || events[1].ID != 3 || events[1].State != "completed" {
		t.Fatalf("bad replay: %+v", events)
	}
}

func TestCanceledExpiryJobDoesNotCloseRoom(t *testing.T) {
	repo := memory.New()
	svc := application.New(repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := svc.RunCancelable(ctx, domain.Item{ID: "late", Scope: "alpha"}, 20*time.Millisecond)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v, want context.Canceled", err)
	}
	if got := repo.Count(context.Background()); got != 0 {
		t.Fatalf("canceled operation persisted %d items", got)
	}
}

func TestStartGameFailureRollsBack(t *testing.T) {
	repo := memory.New()
	svc := application.New(repo)
	before, _ := repo.Snapshot(context.Background())
	err := svc.ApplyAtomic(context.Background(), "completed", "member")
	if err == nil {
		t.Fatal("expected injected failure")
	}
	after, _ := repo.Snapshot(context.Background())
	if after != before {
		t.Fatalf("partial commit: before=%+v after=%+v", before, after)
	}
}

func TestHTTPCreateValidation(t *testing.T) {
	repo := memory.New()
	svc := application.New(repo)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	httptransport.New(svc).Router().ServeHTTP(rr, req)
	if rr.Code != 422 {
		t.Fatalf("status=%d", rr.Code)
	}
}
