# PROJECT1 - Backend Mauna REFACTOR using (Go Native)

Backend ini dibangun menggunakan bahasa **Go** dengan pendekatan **Clean Architecture**. Proyek ini difokuskan pada performa maksimal dengan meminimalisir penggunaan *library* pihak ketiga (Native) dan menggunakan **Raw SQL** melalui `sqlx` untuk kendali penuh atas optimasi database.

---

## 🛠️ Tech Stack
* **Language**: Go 1.25+
* **Database**: PostgreSQL
* **Driver**: `jmoiron/sqlx` (Raw SQL)
* **Migrations**: `golang-migrate`
* **OS Environment**: Arch Linux / Fedora
* **Workflow**: Makefile

---

## IMPORTANT NOTES GOLANG!
Capital letter (PascalCase)  = EXPORTED (Public)   ✅ Bisa diakses dari package lain  
lowercase camelCase          = UNEXPORTED (Private) ❌ Hanya bisa di package yang sama

---
## 📂 Struktur Proyek
Sesuai dengan standar *Go Project Layout*, struktur folder dipisahkan untuk menjaga skalabilitas:
```
Mauna-Backend-Refactor/
├── .env                         # Environment variables local (git ignored)
├── .env.example                 # Template environment variables
├── .gitignore                   # Git ignore rules
├── docker-compose.observability.yml # Jaeger, Prometheus, Loki, Alloy, Grafana
├── Makefile                     # Build, migration, seed, observability commands
├── README.md                    # Project documentation
├── go.mod                       # Go module definition & dependencies
├── go.sum                       # Dependency checksums
│
├── cmd/                         # Application entry points
│   ├── app/
│   │   └── main.go              # Main API application bootstrap
│   └── seed/
│       ├── main.go              # Database seed entry point
│       └── seeder/              # Seeder implementations
│           ├── badge_seeder.go
│           ├── base.go
│           ├── dictionary_seeder.go
│           ├── level_seeder.go
│           ├── question_seeder.go
│           ├── shop_seeder.go
│           ├── sublevel_seeder.go
│           ├── user_badge_seeder.go
│           └── user_seeder.go
│
├── config/
│   └── config.go                # Static config helpers
│
├── docs/                        # Generated/auxiliary API docs package
│
├── internal/                    # Private application code
│   ├── delivery/http/
│   │   ├── debug_log.go         # Shared HTTP error logging helper
│   │   ├── handler.go           # Root and health handlers
│   │   ├── response.go          # Shared HTTP response helpers
│   │   ├── route.go             # HTTP dependency wiring and route registry
│   │   ├── swagger.go           # Swagger UI and OpenAPI spec handlers
│   │   ├── swagger/
│   │   │   └── openapi.json     # OpenAPI 3 specification
│   │   ├── handler/             # HTTP handlers grouped by feature
│   │   │   ├── auth/
│   │   │   │   ├── change_password.go
│   │   │   │   ├── cookie.go
│   │   │   │   ├── debug_log.go
│   │   │   │   ├── handler.go
│   │   │   │   ├── login.go
│   │   │   │   ├── logout.go
│   │   │   │   ├── refresh_token.go
│   │   │   │   ├── register.go
│   │   │   │   └── response.go
│   │   │   └── user/
│   │   │       └── profile/
│   │   │           └── geProfile.go
│   │   ├── middleware/
│   │   │   ├── auth.go          # JWT auth middleware
│   │   │   └── debug_log.go     # Middleware error logging helper
│   │   └── routes/
│   │       └── auth_routes.go
│   │
│   ├── domain/                  # Enterprise/domain rules and contracts
│   │   ├── errors.go            # Domain and wrapped internal errors
│   │   ├── repository.go        # Repository interfaces
│   │   └── entities/            # Domain entities
│   │       ├── badge.go
│   │       ├── daily_task.go
│   │       ├── dictionary.go
│   │       ├── level.go
│   │       ├── progress.go
│   │       ├── question.go
│   │       ├── shop_item.go
│   │       ├── token_blacklist.go # Token revocation
│   │       └── user.go
│   │
│   ├── dto/                     # Request/response DTOs
│   │   ├── auth_dto.go
│   │   └── common_dto.go
│   │
│   ├── repository/              # Data access implementations
│   │   └── auth_repository.go
│   │
│   ├── service/                 # Usecase implementations / business logic
│   │   └── auth_service.go
│   │
│   ├── usecase/                 # Application ports / usecase contracts
│   │   └── auth_usecase.go
│   └── utils/                   # Internal utilities
│
├── migration/                   # Database schema versioning
│   └── 000001-013_*.up/down.sql # SQL migrations
│
├── model/                       # Pre-trained ML models
│   ├── mauna_alphabet_label_map.npy
│   ├── mauna_alphabet_model.pkl
│   ├── mauna_number_label_map.npy
│   └── mauna_number_model.pkl
│
├── observability/               # Local observability stack config
│   ├── alloy/
│   │   └── config.alloy         # Log shipping to Loki
│   ├── grafana/provisioning/datasources/
│   │   └── datasources.yml      # Grafana datasource provisioning
│   └── prometheus/
│       └── prometheus.yml       # Prometheus scrape config
│
└── pkg/                         # Reusable infrastructure packages
    ├── database/
    │   └── connection.go        # PostgreSQL setup via sqlx
    │
    ├── logger/
    │   └── logger.go            # Structured JSON logger
    │
    ├── observability/
    │   ├── http.go              # HTTP metrics, tracing, request logging
    │   └── otel.go              # OpenTelemetry tracer provider
    │
    ├── security/
    │   ├── encryption.go        # AES-256-GCM encryption/decryption
    │   ├── hash.go              # SHA256, SHA512, MD5 helpers
    │   ├── jwt.go               # JWT TokenManager implementation
    │   └── password.go          # Argon2id password hashing
    │
    ├── errors/
    └── validation/
```
---

## 🚀 Cara Menjalankan
Pastikan kamu sudah menginstal `migrate` di sistem kamu (via `pacman -S migrate` atau `go install`).

1. **Setup Environment**:
   Sesuaikan konfigurasi di file `.env`.

2. **Jalankan Migrasi Database**:
   ```bash
   make migrate-up
   ```
3. **Jalankan Aplikasi**:
   ```bash
   make run
   ```

4. **Buka Swagger UI**:
   ```text
   http://localhost:8081/swagger/
   ```

OpenAPI JSON tersedia di:

```text
http://localhost:8081/swagger/openapi.json
```

## Observability Lokal
Project ini sudah menyiapkan structured logging, Prometheus metrics, Loki logs, Grafana, dan OpenTelemetry tracing ke Jaeger.

1. Jalankan stack observability:
   ```bash
   make observability-up
   ```

2. Set environment observability di `.env`:
   ```env
   LOG_FILE=logs/mauna-api.log
   OTEL_TRACES_EXPORTER=otlp
   OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
   ```

3. Jalankan aplikasi:
   ```bash
   make run
   ```

4. Buka tools observability:
   ```text
   Grafana:    http://localhost:3000
   Prometheus: http://localhost:9090
   Jaeger:     http://localhost:16686
   Loki:       http://localhost:3100
   ```

Grafana login default untuk lokal adalah `admin` / `admin`. Datasource Prometheus, Loki, dan Jaeger sudah diprovision otomatis.

Endpoint metrics aplikasi tersedia di:

```text
http://localhost:8081/metrics
```

Contoh query Prometheus:

```promql
sum(rate(mauna_http_requests_total[5m])) by (method, path, status)
histogram_quantile(0.95, sum(rate(mauna_http_request_duration_seconds_bucket[5m])) by (le, path))
```

Contoh query Loki di Grafana Explore:

```logql
{app="mauna-backend"} | json
{app="mauna-backend"} | json | trace_id != ""
```

Log request membawa `trace_id` dan `span_id`, jadi error di Loki bisa dicari kembali di Jaeger.

## Auth JWT dan Cookies
Login akan mengembalikan token di JSON response dan juga menulis cookie:

```text
access_token  - HttpOnly, 24 jam
refresh_token - HttpOnly, 7 hari
```

Endpoint protected seperti `POST /api/auth/change-password` membaca JWT dari cookie `access_token`.
Sebagai fallback untuk client non-browser, header ini juga didukung:

```http
Authorization: Bearer <access_token>
```

Refresh token bisa dikirim dari cookie `refresh_token` atau body JSON:

```json
{ "refresh_token": "..." }
```

Untuk local development gunakan:

```env
COOKIE_SECURE=false
```

Untuk HTTPS production gunakan:

```env
COOKIE_SECURE=true
JWT_SECRET_KEY=secret-yang-panjang-dan-random
```

   ## 📜 Makefile Commands
Gunakan perintah `make` untuk mempercepat alur kerja DevOps di terminal Arch Linux:
* `make run` - Menjalankan aplikasi secara lokal dari `cmd/api/main.go`.
* `make build` - Kompilasi aplikasi menjadi file binary executable di folder `bin/`.
* `make test` - Menjalankan unit testing untuk seluruh modul.
* `make migrate-create name=...` - Membuat file migrasi baru (Up & Down SQL) di folder `migrations/`.
* `make migrate-up` - Menerapkan semua perubahan skema ke database PostgreSQL.
* `make migrate-down` - Membatalkan satu langkah migrasi terakhir (Rollback).
* `make observability-up` - Menjalankan Jaeger lokal untuk tracing.
* `make observability-down` - Mematikan Jaeger lokal.
* `make observability-logs` - Melihat log Jaeger lokal.
* `make tidy` - Merapikan `go.mod` dan melakukan standarisasi format kode (`go fmt`).

---

## 🛡️ Keamanan (Security Focus)
Sesuai dengan spesifikasi sistem, keamanan diimplementasikan pada beberapa lapisan:
* **Password Hashing**: Menggunakan `bcrypt` untuk enkripsi satu arah pada kredensial user.
* **SQL Injection Protection**: Memanfaatkan *parameterized queries* bawaan `sqlx` (menggunakan placeholder `$1, $2`).
* **JWT Authentication**: Proteksi rute API menggunakan *middleware* untuk memastikan akses hanya diberikan kepada user yang valid.
* **Environment Protection**: Data sensitif seperti kredensial database disimpan dalam file `.env` (di-ignore oleh Git).

---

## 🏗️ Alur Data (Data Flow)
Proyek ini mengikuti aturan dependensi **Clean Architecture**:
`Delivery (HTTP)` -> `Service (Business Logic)` -> `Repository (Database)` -> `Domain (Contract & Entity)`

Setiap layer hanya berkomunikasi melalui *interface* yang didefinisikan di dalam folder `domain` untuk menjaga kode tetap *testable* dan modular.



---
*Developed by takumifahri | Developed on MSI Cyborg 15 A12VF (RTX 4060)*
