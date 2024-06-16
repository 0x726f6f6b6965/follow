package mocks

import (
	"context"
	"time"
)

var (
	SetFunc    func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetFunc    func(ctx context.Context, key string) (string, error)
	DeleteFunc func(ctx context.Context, key string) error
)

type MockSotrageCache struct{}

func (m *MockSotrageCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return SetFunc(ctx, key, value, ttl)
}
func (m *MockSotrageCache) Get(ctx context.Context, key string) (string, error) {
	return GetFunc(ctx, key)
}
func (m *MockSotrageCache) Delete(ctx context.Context, key string) error {
	return DeleteFunc(ctx, key)
}
