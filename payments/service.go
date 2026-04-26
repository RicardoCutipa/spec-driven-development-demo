package payments

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Money int64

type CreatePaymentRequest struct {
	AmountCents Money  `json:"amount_cents"`
	Currency    string `json:"currency"`
	CustomerID  string `json:"customer_id"`
}

type Payment struct {
	ID          string `json:"id"`
	AmountCents Money  `json:"amount_cents"`
	Currency    string `json:"currency"`
	CustomerID  string `json:"customer_id"`
	Status      string `json:"status"`
}

type IdempotencyRecord struct {
	RequestHash string
	Payment     Payment
}

type IdempotencyStore interface {
	Get(ctx context.Context, key string) (IdempotencyRecord, bool, error)
	Put(ctx context.Context, key string, record IdempotencyRecord) error
}

type Repository interface {
	Create(ctx context.Context, payment Payment) (Payment, error)
}

type IDFunc func() string

type Service struct {
	idem IdempotencyStore
	repo Repository
	idFn IDFunc
}

func NewService(idem IdempotencyStore, repo Repository, idFn IDFunc) *Service {
	return &Service{idem: idem, repo: repo, idFn: idFn}
}

type CreateResult struct {
	StatusCode int
	Payment    Payment
	Err        error
}

func (s *Service) CreatePayment(ctx context.Context, idempotencyKey string, rawBody []byte) CreateResult {
	if idempotencyKey == "" || len(idempotencyKey) > 128 {
		return CreateResult{StatusCode: 400, Err: errors.New("invalid idempotency key")}
	}

	var req CreatePaymentRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		return CreateResult{StatusCode: 400, Err: errors.New("invalid json")}
	}

	if req.AmountCents <= 0 || req.Currency != "USD" || req.CustomerID == "" {
		return CreateResult{StatusCode: 400, Err: errors.New("invalid request")}
	}

	requestHash := hashPayload(rawBody)

	if rec, ok, err := s.idem.Get(ctx, idempotencyKey); err != nil {
		return CreateResult{StatusCode: 500, Err: err}
	} else if ok {
		if rec.RequestHash != requestHash {
			return CreateResult{StatusCode: 409, Err: errors.New("idempotency key reused with different payload")}
		}
		return CreateResult{StatusCode: 200, Payment: rec.Payment}
	}

	payment := Payment{
		ID:          s.idFn(),
		AmountCents: req.AmountCents,
		Currency:    req.Currency,
		CustomerID:  req.CustomerID,
		Status:      "authorized",
	}

	created, err := s.repo.Create(ctx, payment)
	if err != nil {
		return CreateResult{StatusCode: 500, Err: err}
	}

	if err := s.idem.Put(ctx, idempotencyKey, IdempotencyRecord{RequestHash: requestHash, Payment: created}); err != nil {
		return CreateResult{StatusCode: 500, Err: err}
	}

	return CreateResult{StatusCode: 201, Payment: created}
}

func hashPayload(rawBody []byte) string {
	sum := sha256.Sum256(rawBody)
	return hex.EncodeToString(sum[:])
}

