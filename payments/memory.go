package payments

import (
	"context"
	"sync"
)

type MemoryIdempotencyStore struct {
	mu   sync.Mutex
	data map[string]IdempotencyRecord
}

func NewMemoryIdempotencyStore() *MemoryIdempotencyStore {
	return &MemoryIdempotencyStore{data: map[string]IdempotencyRecord{}}
}

func (s *MemoryIdempotencyStore) Get(ctx context.Context, key string) (IdempotencyRecord, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.data[key]
	return rec, ok, nil
}

func (s *MemoryIdempotencyStore) Put(ctx context.Context, key string, record IdempotencyRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = record
	return nil
}

type MemoryRepo struct {
	mu       sync.Mutex
	Payments []Payment
}

func (r *MemoryRepo) Create(ctx context.Context, payment Payment) (Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Payments = append(r.Payments, payment)
	return payment, nil
}

