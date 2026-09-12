// internal/repository/session.go
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xanity-07/openmat/internal/model/session"
	"github.com/xanity-07/openmat/internal/server"
)

type SessionRepository struct {
	server *server.Server
}

func NewSessionRepository(server *server.Server) *SessionRepository {
	return &SessionRepository{server: server}
}

func sessionKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

func (r *SessionRepository) Create(ctx context.Context, sessionID string, data session.Session, ttl time.Duration) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	// Note the .Err() call here — a raw *redis.StatusCmd from Set() is
	// never nil on its own, so skipping .Err() makes every call look like
	// a failure even when Redis says OK.
	if err := r.server.Redis.Set(ctx, sessionKey(sessionID), payload, ttl).Err(); err != nil {
		return fmt.Errorf("store session in redis: %w", err)
	}

	return nil
}

func (r *SessionRepository) Get(ctx context.Context, sessionID string) (*session.Session, error) {
	val, err := r.server.Redis.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // key doesn't exist — not logged in, not an error
		}
		return nil, fmt.Errorf("get session from redis: %w", err)
	}

	var data session.Session
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &data, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	if err := r.server.Redis.Del(ctx, sessionKey(sessionID)).Err(); err != nil {
		return fmt.Errorf("delete session from redis: %w", err)
	}
	return nil
}
