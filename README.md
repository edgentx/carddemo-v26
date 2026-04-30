# carddemo-v26

**CardDemo COBOL Modernization — DDD Hexagonal Architecture (Go / chi)**

A reference implementation by [Edgent LLC](https://www.edgentllc.com) of the IBM-published [CardDemo](https://github.com/aws-samples/aws-mainframe-modernization-carddemo) sample — a credit-card account-management application originally written in COBOL/JCL/CICS — re-platformed onto a modern stack:

- **Domain-Driven Design** with hexagonal (ports-and-adapters) layout
- **Go 1.21** + [chi](https://github.com/go-chi/chi) HTTP router
- **MongoDB** for aggregate persistence (one repo adapter per bounded context)
- **godog** BDD test suite covering domain rules and HTTP handlers

The same patterns we use here are how we'd approach federal mainframe-modernization work: preserve the COBOL/JCL business-logic semantics (debt aging, demand-letter generation, batch settlement, transaction posting, etc.) while moving the system to cloud-native services and storage.

## Live demo

A working demo of the modernized 3270 terminal interface is hosted at **[dev.card.vforce360.ai](https://dev.card.vforce360.ai)** — green-screen terminal emulator backed by this Go API.

## Architecture

```
┌─────────────────────────────────────────────────────┐
│  Inbound (HTTP / chi)                               │
│  cmd/server  →  src/app/{context}/handler           │
└──────────────────────────┬──────────────────────────┘
                           │ commands & queries
                           ▼
┌─────────────────────────────────────────────────────┐
│  Application layer  (src/app)                       │
│   - command handlers, query handlers                │
│   - port interfaces (repository, services)          │
└──────────────────────────┬──────────────────────────┘
                           │ pure domain calls
                           ▼
┌─────────────────────────────────────────────────────┐
│  Domain layer  (src/domain)                         │
│   - aggregates: Account, Card, CardPolicy,          │
│     Transaction, BatchSettlement, ExportJob,        │
│     Report, UserProfile                             │
│   - invariants, domain events                       │
└──────────────────────────┬──────────────────────────┘
                           │ adapter implementations
                           ▼
┌─────────────────────────────────────────────────────┐
│  Infrastructure  (src/infrastructure)               │
│   - MongoDB repositories per aggregate              │
│   - shared persistence helpers                      │
└─────────────────────────────────────────────────────┘
```

## Bounded contexts

Each context owns its aggregates, commands, queries, and persistence:

| Context | Responsibility |
|---------|---------------|
| `account` | Account lifecycle, status transitions |
| `userprofile` | Customer / user identity records |
| `card` | Card issuance, lost-card reporting, status |
| `cardpolicy` | Spending limits, fraud thresholds, policy updates |
| `transaction` | Transaction posting, authorization, reversal |
| `batchsettlement` | End-of-day batch open / reconcile / close |
| `report` | Statement and operational report generation |
| `exportjob` | Async export workflows (initiate / complete) |
| `shared` | Cross-cutting types and routing |

## Layout

```
cmd/server/         # HTTP entry point (chi router on :8080)
src/
  domain/           # pure aggregates and invariants
  app/              # commands, queries, port interfaces
  infrastructure/   # MongoDB adapters per context
mocks/              # in-memory implementations of ports for tests
tests/              # godog BDD suite + mock repositories
Makefile
go.mod
```

## Build / test / run

```bash
make build      # go build ./...
make test       # go test ./... -v -cover
make lint       # golangci-lint run

go run ./cmd/server      # listens on :8080, /health returns 200
```

## Origin

This implementation is derived from the open-source [aws-mainframe-modernization-carddemo-sample](https://github.com/aws-samples/aws-mainframe-modernization-carddemo) — the COBOL artifact AWS and IBM publish for federal and enterprise teams evaluating mainframe-replatforming approaches.

The Go port here keeps the original business behavior (account, card, transaction, batch settlement, reporting flows) and demonstrates how the same semantics map cleanly onto DDD aggregates with clear ports, testable in isolation from the mainframe scheduler and storage formats it originally relied on.

## About Edgent

[Edgent LLC](https://www.edgentllc.com) is a Service-Disabled Veteran-Owned Small Business based in Bastrop, TX, focused on federal training, AI-augmented engineering, and platform modernization. UEI `UZFLEDGS3PC5`.
