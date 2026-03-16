Architecture of the Security Platform

security-platform/
│
├── gateway/ → API Gateway (entry point)
│
├── auth-service/ → Authentication + RBAC
│
├── shared/ → shared libraries
│ ├── logger
│ ├── middleware
│ ├── errors
│ ├── config
│ └── utils
│
├── infra/
│ ├── docker
│ ├── prometheus
│ ├── grafana
│ └── jaeger
│
└── scripts/

---

System Architecture

Client
│
▼
API Gateway (Gin)
│
├── Rate Limiter
├── CSRF protection
├── JWT validation
├── Request logging
├── Circuit breaker
├── Distributed tracing
│
▼
Auth Service
│
├── Login
├── Register
├── Token rotation
├── Refresh tokens
├── RBAC
├── Permission checks
│
▼
Database (Postgres)

- Redis
  ├── Rate limiting
  ├── Session store
  ├── Token blacklist
  ├── Refresh rotation

- Queue (optional)
  ├── audit logs
  └── async email

---

Architecture principle:

cmd/ → entry points

internal/ → application logic

pkg/ → reusable libraries

configs/ → configuration files

---

Authentication
Access Token (JWT)
Refresh Token (Rotating)
Refresh Token Revocation
Token Blacklist
CSRF Token
Password hashing (argon2id)
Session tracking
Device fingerprint

Authorization
RBAC
Roles
Permissions
Role inheritance
Policy middleware

Security Controls
Rate limiting
IP throttling
Brute force protection
Account lock
CSRF protection
Audit logging
Security headers
Request validation

Observability
Structured logging (Zap)
Metrics (Prometheus)
Tracing (OpenTelemetry + Jaeger)
Health checks

Reliability
Circuit breaker
Retry
Timeouts
Queue workers
Graceful shutdown

---

Redis will later store:

refresh tokens

session data

rate limiter state

CSRF tokens

token blacklist

---
