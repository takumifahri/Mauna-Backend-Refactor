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
├── .env                      # Environment variables (local, git ignored)
├── .env.example              # Template environment variables
├── .gitignore                # Git ignore rules
├── Makefile                  # Build commands (run, build, test, migrate, seed)
├── README.md                 # Project documentation
├── go.mod                    # Go module definition & dependencies
├── go.sum                    # Dependency checksums
│
├── cmd/                      # Entry point aplikasi
│   ├── app/
│   │   └── main.go           # Aplikasi utama + DI container
│   └── seed/
│       ├── main.go           # Seed entry point
│       └── seeder/           # Seeder implementations
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
│   └── config.go             # Load .env & manage configuration
│
├── internal/                 # Core business logic (private)
│   ├── delivery/http/
│   │   ├── handler.go        # Global handlers (health, root)
│   │   ├── route.go          # Route registry & management
│   │   ├── handler/          # HTTP request handlers per feature
│   │   │   ├── admin/        # Admin endpoints (empty)
│   │   │   ├── auth/         # Auth endpoints
│   │   │   │   ├── change_password.go
│   │   │   │   ├── handler.go
│   │   │   │   ├── login.go
│   │   │   │   ├── logout.go
│   │   │   │   ├── refresh_token.go
│   │   │   │   └── register.go
│   │   │   └── user/         # User endpoints (empty)
│   │   ├── middleware/       # HTTP middleware (empty)
│   │   └── routes/           # Route registration
│   │       └── auth_routes.go
│   │
│   ├── domain/               # Business rules & interfaces
│   │   ├── errors.go         # Domain error definitions
│   │   ├── repository.go     # Repository interfaces
│   │   └── entities/         # Data models
│   │       ├── badge.go      # Badge & UserBadge
│   │       ├── daily_task.go # Daily task tracking
│   │       ├── dictionary.go # Kamus (vocabulary)
│   │       ├── level.go      # Level & SubLevel
│   │       ├── progress.go   # Progress status enum
│   │       ├── question.go   # Soal (question types)
│   │       ├── shop_item.go  # ShopItem & Inventory
│   │       ├── token_blacklist.go # Token revocation
│   │       └── user.go       # User role & tier
│   │
│   ├── dto/                  # Data Transfer Objects
│   │   ├── auth_dto.go       # Login, Register, ChangePassword DTO
│   │   └── common_dto.go     # Response wrapper & ErrorResponse
│   │
│   ├── repository/           # Data access layer
│   │   └── auth_repository.go # User CRUD operations
│   │
│   ├── service/              # Business logic layer
│   │   └── auth_service.go   # Auth usecases
│   │
│   ├── routes/               # Route registry (empty)
│   └── utils/                # Utility functions (empty)
│
├── migration/                # Database schema versioning
│   └── 000001-013_*.up/down.sql # 13 migrations
│
├── model/                    # Pre-trained ML models
│   ├── mauna_alphabet_label_map.npy
│   ├── mauna_alphabet_model.pkl
│   ├── mauna_number_label_map.npy
│   └── mauna_number_model.pkl
│
└── pkg/                      # Reusable packages (public)
    ├── database/
    │   └── connection.go      # PostgreSQL setup via sqlx
    │
    ├── security/             # Security utilities
    │   ├── encryption.go      # AES-256-GCM encryption/decryption
    │   ├── hash.go            # SHA256, SHA512, MD5 functions
    │   ├── jwt.go             # JWT token generation & verification
    │   └── password.go        # Argon2id password hashing
    │
    ├── errors/               # Custom errors (empty)
    ├── logger/               # Logging utilities (empty)
    └── validation/           # Input validation (empty)
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
   ## 📜 Makefile Commands
Gunakan perintah `make` untuk mempercepat alur kerja DevOps di terminal Arch Linux:
* `make run` - Menjalankan aplikasi secara lokal dari `cmd/api/main.go`.
* `make build` - Kompilasi aplikasi menjadi file binary executable di folder `bin/`.
* `make test` - Menjalankan unit testing untuk seluruh modul.
* `make migrate-create name=...` - Membuat file migrasi baru (Up & Down SQL) di folder `migrations/`.
* `make migrate-up` - Menerapkan semua perubahan skema ke database PostgreSQL.
* `make migrate-down` - Membatalkan satu langkah migrasi terakhir (Rollback).
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