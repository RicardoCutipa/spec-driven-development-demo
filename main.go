package main

import (
	"context"
	"fmt"

	"github.com/RicardoCutipa/spec-driven-development-demo/payments"
)

func main() {
	ctx := context.Background()

	idem := payments.NewMemoryIdempotencyStore()
	repo := &payments.MemoryRepo{}
	i := 0
	nextID := func() string {
		i++
		return fmt.Sprintf("pay_%d", i)
	}

	svc := payments.NewService(idem, repo, nextID)

	body := []byte(`{"amount_cents":1200,"currency":"USD","customer_id":"cus_1"}`)
	r1 := svc.CreatePayment(ctx, "idem-key-1", body)
	r2 := svc.CreatePayment(ctx, "idem-key-1", body)

	fmt.Printf("First call:  status=%d id=%s\n", r1.StatusCode, r1.Payment.ID)
	fmt.Printf("Replay call: status=%d id=%s\n", r2.StatusCode, r2.Payment.ID)
}
