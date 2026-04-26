# Spec: Idempotent Create Payment

This spec defines the behavior of an idempotent "create payment" operation. The core idea:
retrying the same request must not create duplicate payments.

## Operation

`POST /payments`

### Required header

`Idempotency-Key`: non-empty string, max length 128.

### Request body (JSON)

- `amount_cents` (integer, required, must be > 0)
- `currency` (string, required, must be `"USD"` for now)
- `customer_id` (string, required, non-empty)

### Behavior

1. If `Idempotency-Key` is missing/empty/too long -> **400**
2. If the JSON is invalid -> **400**
3. If any field is invalid -> **400**
4. If the idempotency key is **new**:
   - create exactly one payment
   - return **201** with the payment
5. If the idempotency key was used before:
   - if the request payload is identical -> return **200** with the same payment as before
   - if the request payload differs -> return **409**

### Invariants

- A successful replay (same key + same payload) returns the same `payment_id`.
- A key cannot be reused for a different payload.

