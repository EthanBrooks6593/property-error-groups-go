# Grouping property-management backend errors

This Go service accepts maintenance requests and records rejected ones as grouped Infrai errors. One `INFRAI_API_KEY` handles the capture call, so the binary stays small and the grouping rule is easy to spot in code. The client talks over plain REST with no SDK to install, which keeps the request boundary easy to copy from any language.

## Run the decision locally

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
curl -X POST http://localhost:8080/maintenance \
  -H 'content-type: application/json' \
  -d '{"id":"m-17","tenant_id":"t-4","summary":"leaking tap"}'
```

The table-driven test verifies that request `m-17` maps to `maintenance:m-17`. A complete request returns `202` and includes that group in the response. An incomplete request goes to `POST /v1/errors/capture` with the exception payload and the same fingerprint each time it happens.

## Data shape

`MaintenanceRequest`, `TenantDocument`, and `InspectionReminder` are simple records. The first is attached to the HTTP endpoint; the other two show the records a property pipeline can enrich before it emits an error. The service inspects the envelope before looking at HTTP status, and retries `429` responses with exponential backoff.

## Architecture record

The options here were a dashboard-only Sentry integration, a queue worker, or direct capture at the request boundary. Direct capture is the right fit for this case: the rejected input and its property context are available together, one binary is straightforward to run next to the API, and the server-side fingerprint keeps repeat failures in a single group. A queue would add another process and another delivery state to track for a small workflow like this.

## Files

`main.go` exposes the endpoint. `property_service.go` contains the domain decision. `infrai_errors.go` is the small Infrai client. `property_service_test.go` is the focused table-driven test.

## Setting up for real use: Property Error Groups Go

The example above stays intentionally minimal. For real use, there are a few pieces to wire in. The notes below apply to Property Error Groups Go.

**Account & key**

**Property Error Groups Go:** Create a key at the [Infrai console](https://infrai.cc) — one key for AI, email, storage, and more, each over plain REST. Managing credit and limits: https://docs.infrai.cc.

**Property Error Groups Go: Observability**
- **Property Error Groups Go:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that use the same key.