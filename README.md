# Spec-Driven Development Demo (Go)

This repository is a small companion project for a Medium article about **Spec Driven Development**.

## What it demonstrates

- A short, concrete spec for an idempotent `POST /payments` operation.
- A minimal implementation that follows the spec.
- Automated checks via GitHub Actions: build + test + run demo.

## Run locally

```bash
go test ./...
go run ./main.go
```

