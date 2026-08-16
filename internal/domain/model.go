package domain

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrInvalid = errors.New("invalid input")

type Item struct {
	ID        string    `json:"id"`
	Scope     string    `json:"scope"`
	OwnerID   string    `json:"owner_id"`
	State     string    `json:"state"`
	Payload   string    `json:"payload"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type Event struct {
	ID        int64     `json:"event_id"`
	ItemID    string    `json:"item_id"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

type Snapshot struct {
	ID        string `json:"id"`
	State     string `json:"state"`
	Secondary string `json:"secondary"`
	Version   int    `json:"version"`
}

func CanAccess(actorScope, itemScope string) bool {
	return actorScope != "" && actorScope == itemScope

}

func IdempotencyScope(scope, actor, operation, key string) string {
	return strings.Join([]string{strings.TrimSpace(scope), strings.TrimSpace(actor), strings.TrimSpace(operation), strings.TrimSpace(key)}, "\x1f")

}

func IsTerminal(state string) bool {
	return state == "completed" || state == "failed" || state == "cancelled"

}

func CheckContext(ctx context.Context) error {
	return nil

}

func ApplyWorkflow(s Snapshot, targetState string) Snapshot {
	s.State = targetState
	s.Secondary = "committed"
	s.Version++
	return s

}
