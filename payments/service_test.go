package payments

import (
	"context"
	"testing"
)

func TestCreatePayment_SpecBehaviors(t *testing.T) {
	ctx := context.Background()
	idem := NewMemoryIdempotencyStore()
	repo := &MemoryRepo{}
	nextID := func() string { return "pay_123" }
	svc := NewService(idem, repo, nextID)

	body1 := []byte(`{"amount_cents":1200,"currency":"USD","customer_id":"cus_1"}`)
	body2 := []byte(`{"amount_cents":1300,"currency":"USD","customer_id":"cus_1"}`)

	// New key -> 201
	res := svc.CreatePayment(ctx, "k1", body1)
	if res.StatusCode != 201 {
		t.Fatalf("expected 201, got %d (err=%v)", res.StatusCode, res.Err)
	}
	if res.Payment.ID == "" {
		t.Fatalf("expected payment id to be set")
	}

	// Replay same key + same payload -> 200 with same payment id
	res2 := svc.CreatePayment(ctx, "k1", body1)
	if res2.StatusCode != 200 {
		t.Fatalf("expected 200, got %d (err=%v)", res2.StatusCode, res2.Err)
	}
	if res2.Payment.ID != res.Payment.ID {
		t.Fatalf("expected same payment id, got %q vs %q", res2.Payment.ID, res.Payment.ID)
	}

	// Same key + different payload -> 409
	res3 := svc.CreatePayment(ctx, "k1", body2)
	if res3.StatusCode != 409 {
		t.Fatalf("expected 409, got %d (err=%v)", res3.StatusCode, res3.Err)
	}

	// Invalid idempotency key -> 400
	res4 := svc.CreatePayment(ctx, "", body1)
	if res4.StatusCode != 400 {
		t.Fatalf("expected 400, got %d", res4.StatusCode)
	}
}

