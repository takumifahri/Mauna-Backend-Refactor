package seeder

import (
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
)

type SublevelSeeder struct {
    db *sqlx.DB
}

func NewSublevelSeeder(db *sqlx.DB) *SublevelSeeder {
    return &SublevelSeeder{db: db}
}

func (s *SublevelSeeder) Name() string {
    return "SublevelSeeder"
}

type sublevelData struct {
    Name        string
    Description string
    Tujuan      string
    LevelID     int64
}

func (s *SublevelSeeder) Run() error {
    PrintInfo("🎯 Seeding Sublevels...")

    tx, err := s.db.Beginx()
    if err != nil {
        PrintError(fmt.Sprintf("Failed to start transaction: %v", err))
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    sublevelsData := s.getSublevelsData()

    createdCount := 0
    for _, data := range sublevelsData {
        var existingID sql.NullInt64
        err := tx.QueryRow(
            "SELECT id FROM sublevel WHERE name = $1 AND level_id = $2",
            data.Name, data.LevelID,
        ).Scan(&existingID)

        if err == nil && existingID.Valid {
            PrintWarning(fmt.Sprintf("Sublevel already exists: %s", data.Name))
            continue
        }

        if err != nil && err != sql.ErrNoRows {
            PrintError(fmt.Sprintf("Database error: %v", err))
            return err
        }

        _, err = tx.Exec(
            "INSERT INTO sublevel (name, description, tujuan, level_id, created_at) VALUES ($1, $2, $3, $4, NOW())",
            data.Name,
            data.Description,
            data.Tujuan,
            data.LevelID,
        )

        if err != nil {
            PrintError(fmt.Sprintf("Failed to create sublevel: %v", err))
            return err
        }

        PrintSuccess(fmt.Sprintf("Created sublevel: %s (Level %d)", data.Name, data.LevelID))
        createdCount++
    }

    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Sublevel seeding completed. Created %d sublevels.", createdCount))
    return nil
}

func (s *SublevelSeeder) getSublevelsData() []sublevelData {
    return []sublevelData{
        // LEVEL 1
        {Name: "Sublevel 1.1", Description: "Huruf A-C", Tujuan: "Belajar isyarat huruf A, B, C", LevelID: 1},
        {Name: "Sublevel 1.2", Description: "Huruf D-F", Tujuan: "Belajar isyarat huruf D, E, F", LevelID: 1},
        {Name: "Sublevel 1.3", Description: "Huruf G-I", Tujuan: "Belajar isyarat huruf G, H, I", LevelID: 1},
        {Name: "Sublevel 1.4", Description: "Huruf J-L", Tujuan: "Belajar isyarat huruf J, K, L", LevelID: 1},
        {Name: "Sublevel 1.5", Description: "Huruf M-O", Tujuan: "Belajar isyarat huruf M, N, O", LevelID: 1},
        {Name: "Sublevel 1.6", Description: "Huruf P-R", Tujuan: "Belajar isyarat huruf P, Q, R", LevelID: 1},
        {Name: "Sublevel 1.7", Description: "Huruf S-U", Tujuan: "Belajar isyarat huruf S, T, U", LevelID: 1},
        {Name: "Sublevel 1.8", Description: "Huruf V-X", Tujuan: "Belajar isyarat huruf V, W, X", LevelID: 1},
        {Name: "Sublevel 1.9", Description: "Huruf Y-Z", Tujuan: "Belajar isyarat huruf Y dan Z", LevelID: 1},
        {Name: "Sublevel 1.10", Description: "Latihan Huruf", Tujuan: "Latihan mengenal isyarat huruf A sampai Z", LevelID: 1},

        // LEVEL 2
        {Name: "Sublevel 2.1", Description: "Hewan Rumah", Tujuan: "Belajar isyarat kucing, anjing, burung, dan ikan", LevelID: 2},
        {Name: "Sublevel 2.2", Description: "Hewan Ternak", Tujuan: "Belajar isyarat ayam, sapi, kambing, dan bebek", LevelID: 2},
        {Name: "Sublevel 2.3", Description: "Hewan Liar", Tujuan: "Belajar isyarat gajah, monyet, singa, dan ular", LevelID: 2},
        {Name: "Sublevel 2.4", Description: "Hewan Kecil", Tujuan: "Belajar isyarat semut, kupu-kupu, lebah, dan kelinci", LevelID: 2},
        {Name: "Sublevel 2.5", Description: "Keluarga Inti", Tujuan: "Belajar isyarat ayah, ibu, kakak, dan adik", LevelID: 2},
        {Name: "Sublevel 2.6", Description: "Keluarga Besar", Tujuan: "Belajar isyarat kakek, nenek, paman, dan bibi", LevelID: 2},
        {Name: "Sublevel 2.7", Description: "Teman dan Guru", Tujuan: "Belajar isyarat teman, guru, dan sekolah", LevelID: 2},
        {Name: "Sublevel 2.8", Description: "Hewan Air", Tujuan: "Belajar isyarat ikan, paus, kura-kura, dan katak", LevelID: 2},
        {Name: "Sublevel 2.9", Description: "Hewan Udara", Tujuan: "Belajar isyarat burung, kupu-kupu, dan lebah", LevelID: 2},
        {Name: "Sublevel 2.10", Description: "Latihan Hewan dan Keluarga", Tujuan: "Latihan mengenal isyarat hewan dan keluarga", LevelID: 2},

        // LEVEL 3
        {Name: "Sublevel 3.1", Description: "Angka 0-2", Tujuan: "Belajar isyarat angka 0 sampai 2", LevelID: 3},
        {Name: "Sublevel 3.2", Description: "Angka 3-5", Tujuan: "Belajar isyarat angka 3 sampai 5", LevelID: 3},
        {Name: "Sublevel 3.3", Description: "Angka 6-7", Tujuan: "Belajar isyarat angka 6 sampai 7", LevelID: 3},
        {Name: "Sublevel 3.4", Description: "Angka 8-9", Tujuan: "Belajar isyarat angka 8 sampai 9", LevelID: 3},
        {Name: "Sublevel 3.5", Description: "Angka 0-9", Tujuan: "Review isyarat angka 0 sampai 9", LevelID: 3},
        {Name: "Sublevel 3.6", Description: "Penjumlahan", Tujuan: "Belajar isyarat penjumlahan angka 1 digit (0-9)", LevelID: 3},
        {Name: "Sublevel 3.7", Description: "Pengurangan", Tujuan: "Belajar isyarat pengurangan angka 1 digit (0-9)", LevelID: 3},
        {Name: "Sublevel 3.8", Description: "Bentuk Geometri", Tujuan: "Belajar isyarat bentuk geometri dasar (lingkaran, segitiga, kotak, dan garis)", LevelID: 3},
        {Name: "Sublevel 3.9", Description: "Waktu", Tujuan: "Belajar isyarat waktu dasar (hari, minggu, bulan, tahun, pagi, siang, sore, malam)", LevelID: 3},
        {Name: "Sublevel 3.10", Description: "Matematika Dasar", Tujuan: "Latihan mengenal isyarat matematika dasar", LevelID: 3},

        // LEVEL 4
        {Name: "Sublevel 4.1", Description: "Pagi Hari", Tujuan: "Belajar isyarat bangun, mandi, dan sarapan", LevelID: 4},
        {Name: "Sublevel 4.2", Description: "Sekolah", Tujuan: "Belajar isyarat belajar, menulis, membaca, dan bermain", LevelID: 4},
        {Name: "Sublevel 4.3", Description: "Di Rumah", Tujuan: "Belajar isyarat makan, tidur, duduk, dan nonton", LevelID: 4},
        {Name: "Sublevel 4.4", Description: "Bermain", Tujuan: "Belajar isyarat main bola, lari, lompat, dan nyanyi", LevelID: 4},
        {Name: "Sublevel 4.5", Description: "Bersih-bersih", Tujuan: "Belajar isyarat cuci tangan, sapu, dan sikat gigi", LevelID: 4},
        {Name: "Sublevel 4.6", Description: "Makan dan Minum", Tujuan: "Belajar isyarat makan, minum, dan lapar", LevelID: 4},
        {Name: "Sublevel 4.7", Description: "Emosi", Tujuan: "Belajar isyarat senang, sedih, marah, dan takut", LevelID: 4},
        {Name: "Sublevel 4.8", Description: "Kegiatan Luar", Tujuan: "Belajar isyarat jalan, hujan, panas, dan taman", LevelID: 4},
        {Name: "Sublevel 4.9", Description: "Waktu Istirahat", Tujuan: "Belajar isyarat tidur, mimpi, dan bangun", LevelID: 4},
        {Name: "Sublevel 4.10", Description: "Latihan Aktivitas", Tujuan: "Latihan mengenal isyarat aktivitas sehari-hari", LevelID: 4},
    }
}