package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	return &RedisStore{
		Client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

// HoldSeat uses Redis SETNX to ensure only one user can hold a given
// seat for a show at a time. Returns true if the hold was acquired.
func (r *RedisStore) HoldSeat(ctx context.Context, showID, seatID, userID string) (bool, error) {
	key := "hold:" + showID + ":" + seatID
	// EX 300 = 5 minute TTL, matches the booking hold window
	return r.Client.SetNX(ctx, key, userID, 5*time.Minute).Result()
}

// VerifyHold checks that a hold on a seat still belongs to expectedUserID.
func (r *RedisStore) VerifyHold(ctx context.Context, showID, seatID, expectedUserID string) bool {
	key := "hold:" + showID + ":" + seatID
	val, err := r.Client.Get(ctx, key).Result()
	return err == nil && val == expectedUserID
}

// ReleaseHold removes a hold, e.g. after a successful booking or explicit cancel.
func (r *RedisStore) ReleaseHold(ctx context.Context, showID, seatID string) error {
	key := "hold:" + showID + ":" + seatID
	return r.Client.Del(ctx, key).Err()
}
