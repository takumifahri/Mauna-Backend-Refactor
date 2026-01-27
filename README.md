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

## 📂 Struktur Proyek
Sesuai dengan standar *Go Project Layout*, struktur folder dipisahkan untuk menjaga skalabilitas:

* **`cmd/api/`**: Tempat file `main.go` berada. Berfungsi sebagai *entry point* untuk inisialisasi database dan *dependency injection*.
* **`internal/domain/`**: Berisi *Entity* (struktur tabel DB) dan *Interface* (kontrak bisnis).
* **`internal/service/`**: Tempat logika bisnis utama (Usecase).
* **`internal/repository/`**: Implementasi akses data langsung ke database menggunakan Raw SQL.
* **`internal/delivery/http/`**: Handler API, Middleware, dan Routing.
* **`migrations/`**: Berisi file `.sql` untuk perubahan skema database (Version Control for Database).
* **`pkg/`**: Library pembantu yang bersifat umum (seperti Helper keamanan atau Utils).

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