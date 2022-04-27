package session_repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/go-redis/redis/v8"
)

type repository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Create(ctx context.Context, refreshToken string, session *domains.Session) error {
	stringify, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("sessions:%s", refreshToken)
	if err := r.client.Set(ctx, key, string(stringify), 0).Err(); err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByRefreshToken(ctx context.Context, refreshToken string) (*domains.Session, error) {
	key := fmt.Sprintf("sessions:%s", refreshToken)
	sessionRaw, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var session *domains.Session
	if err := json.Unmarshal([]byte(sessionRaw), &session); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *repository) Remove(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf("sessions:%s", refreshToken)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
