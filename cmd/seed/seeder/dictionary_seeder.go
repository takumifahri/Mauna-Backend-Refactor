package seeder

import (
    "database/sql"
    "fmt"
    "github.com/jmoiron/sqlx"
    "REFACTORING_MAUNA/internal/domain/constant"
)

type KamusSeeder struct {
    db *sqlx.DB
}

func NewKamusSeeder(db *sqlx.DB) *KamusSeeder {
    return &KamusSeeder{db: db}
}

func (s *KamusSeeder) Name() string {
    return "KamusSeeder"
}

type kamusData struct {
    WordText     string
    Definition   string
    VideoURL     *string
    ImageURLRef  *string
    Category     constant.DictionaryCategory
}

func (s *KamusSeeder) Run() error {
    PrintInfo("🌱 Seeding Kamus (Dictionary)...")

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

    kamusDataList := s.getKamusData()

    createdCount := 0
    for _, data := range kamusDataList {
        var existingID sql.NullInt64
        err := tx.QueryRow("SELECT id FROM kamus WHERE word_text = $1", data.WordText).Scan(&existingID)

        if err == nil && existingID.Valid {
            PrintWarning(fmt.Sprintf("Kamus already exists: %s", data.WordText))
            continue
        }

        if err != nil && err != sql.ErrNoRows {
            PrintError(fmt.Sprintf("Database error: %v", err))
            return err
        }

        _, err = tx.Exec(
            "INSERT INTO kamus (word_text, definition, video_url, image_url_ref, category, created_at) VALUES ($1, $2, $3, $4, $5, NOW())",
            data.WordText,
            data.Definition,
            data.VideoURL,
            data.ImageURLRef,
            data.Category,
        )

        if err != nil {
            PrintError(fmt.Sprintf("Failed to create kamus: %v", err))
            return err
        }

        PrintSuccess(fmt.Sprintf("Created kamus: %s", data.WordText))
        createdCount++
    }

    if err := tx.Commit(); err != nil {
        PrintError(fmt.Sprintf("Failed to commit transaction: %v", err))
        return err
    }

    PrintSuccess(fmt.Sprintf("Kamus seeding completed. Created %d entries.", createdCount))
    return nil
}

func (s *KamusSeeder) getKamusData() []kamusData {
    return []kamusData{
        // Numbers (0-9)
        {WordText: "0", Definition: "Angka nol", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/0.webm"), ImageURLRef: ptrStr("kamus/0.png"), Category: constant.CategoryNumbers},
        {WordText: "1", Definition: "Angka satu", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/01.webm"), ImageURLRef: ptrStr("kamus/1.png"), Category: constant.CategoryNumbers},
        {WordText: "2", Definition: "Angka dua", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/02.webm"), ImageURLRef: ptrStr("kamus/2.png"), Category: constant.CategoryNumbers},
        {WordText: "3", Definition: "Angka tiga", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/03.webm"), ImageURLRef: ptrStr("kamus/3.png"), Category: constant.CategoryNumbers},
        {WordText: "4", Definition: "Angka empat", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/04.webm"), ImageURLRef: ptrStr("kamus/4.png"), Category: constant.CategoryNumbers},
        {WordText: "5", Definition: "Angka lima", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/05.webm"), ImageURLRef: ptrStr("kamus/5.png"), Category: constant.CategoryNumbers},
        {WordText: "6", Definition: "Angka enam", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/06.webm"), ImageURLRef: ptrStr("kamus/6.png"), Category: constant.CategoryNumbers},
        {WordText: "7", Definition: "Angka tujuh", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/07.webm"), ImageURLRef: ptrStr("kamus/7.png"), Category: constant.CategoryNumbers},
        {WordText: "8", Definition: "Angka delapan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/08.webm"), ImageURLRef: ptrStr("kamus/8.png"), Category: constant.CategoryNumbers},
        {WordText: "9", Definition: "Angka sembilan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/09.webm"), ImageURLRef: ptrStr("kamus/9.png"), Category: constant.CategoryNumbers},

        // Alphabet A-Z
        {WordText: "A", Definition: "Huruf pertama dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm"), ImageURLRef: ptrStr("kamus/A.png"), Category: constant.CategoryAlphabet},
        {WordText: "B", Definition: "Huruf kedua dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/B.webm"), ImageURLRef: ptrStr("kamus/B.png"), Category: constant.CategoryAlphabet},
        {WordText: "C", Definition: "Huruf ketiga dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/C.webm"), ImageURLRef: ptrStr("kamus/C.png"), Category: constant.CategoryAlphabet},
        {WordText: "D", Definition: "Huruf keempat dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/D.webm"), ImageURLRef: ptrStr("kamus/D.png"), Category: constant.CategoryAlphabet},
        {WordText: "E", Definition: "Huruf kelima dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/E.webm"), ImageURLRef: ptrStr("kamus/E.png"), Category: constant.CategoryAlphabet},
        {WordText: "F", Definition: "Huruf keenam dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/F.webm"), ImageURLRef: ptrStr("kamus/F.png"), Category: constant.CategoryAlphabet},
        {WordText: "G", Definition: "Huruf ketujuh dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/G.webm"), ImageURLRef: ptrStr("kamus/G.png"), Category: constant.CategoryAlphabet},
        {WordText: "H", Definition: "Huruf kedelapan dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/H.webm"), ImageURLRef: ptrStr("kamus/H.png"), Category: constant.CategoryAlphabet},
        {WordText: "I", Definition: "Huruf kesembilan dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/I.webm"), ImageURLRef: ptrStr("kamus/I.png"), Category: constant.CategoryAlphabet},
        {WordText: "J", Definition: "Huruf kesepuluh dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/J.webm"), ImageURLRef: ptrStr("kamus/J.png"), Category: constant.CategoryAlphabet},
        {WordText: "K", Definition: "Huruf kesebelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/K.webm"), ImageURLRef: ptrStr("kamus/K.png"), Category: constant.CategoryAlphabet},
        {WordText: "L", Definition: "Huruf keduabelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/L.webm"), ImageURLRef: ptrStr("kamus/L.png"), Category: constant.CategoryAlphabet},
        {WordText: "M", Definition: "Huruf ketigabelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/M.webm"), ImageURLRef: ptrStr("kamus/M.png"), Category: constant.CategoryAlphabet},
        {WordText: "N", Definition: "Huruf keempatbelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/N.webm"), ImageURLRef: ptrStr("kamus/N.png"), Category: constant.CategoryAlphabet},
        {WordText: "O", Definition: "Huruf kelimabelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/O.webm"), ImageURLRef: ptrStr("kamus/O.png"), Category: constant.CategoryAlphabet},
        {WordText: "P", Definition: "Huruf keenambelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm"), ImageURLRef: ptrStr("kamus/P.png"), Category: constant.CategoryAlphabet},
        {WordText: "Q", Definition: "Huruf ketujuhbelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Q.webm"), ImageURLRef: ptrStr("kamus/Q.png"), Category: constant.CategoryAlphabet},
        {WordText: "R", Definition: "Huruf kedelapanbelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/R.webm"), ImageURLRef: ptrStr("kamus/R.png"), Category: constant.CategoryAlphabet},
        {WordText: "S", Definition: "Huruf kesembilanbelas dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/S.webm"), ImageURLRef: ptrStr("kamus/S.png"), Category: constant.CategoryAlphabet},
        {WordText: "T", Definition: "Huruf keduapuluh dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/T.webm"), ImageURLRef: ptrStr("kamus/T.png"), Category: constant.CategoryAlphabet},
        {WordText: "U", Definition: "Huruf keduapuluh satu dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/U.webm"), ImageURLRef: ptrStr("kamus/U.png"), Category: constant.CategoryAlphabet},
        {WordText: "V", Definition: "Huruf keduapuluh dua dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/V.webm"), ImageURLRef: ptrStr("kamus/V.png"), Category: constant.CategoryAlphabet},
        {WordText: "W", Definition: "Huruf keduapuluh tiga dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/W.webm"), ImageURLRef: ptrStr("kamus/W.png"), Category: constant.CategoryAlphabet},
        {WordText: "X", Definition: "Huruf keduapuluh empat dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/X.webm"), ImageURLRef: ptrStr("kamus/X.png"), Category: constant.CategoryAlphabet},
        {WordText: "Y", Definition: "Huruf keduapuluh lima dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm"), ImageURLRef: ptrStr("kamus/Y.png"), Category: constant.CategoryAlphabet},
        {WordText: "Z", Definition: "Huruf keduapuluh enam dalam alfabet", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Z.webm"), ImageURLRef: ptrStr("kamus/Z.png"), Category: constant.CategoryAlphabet},

        // Hewan
        {WordText: "Anjing", Definition: "Hewan peliharaan yang setia", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm"), Category: constant.CategoryKosakata},
        {WordText: "Bebek", Definition: "Hewan ternak yang suka berenang", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kambing", Definition: "Hewan ternak penghasil susu dan daging", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm"), Category: constant.CategoryKosakata},
        {WordText: "Gajah", Definition: "Hewan besar dengan belalai panjang", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm"), Category: constant.CategoryKosakata},
        {WordText: "Monyet", Definition: "Hewan cerdas yang suka memanjat", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm"), Category: constant.CategoryKosakata},
        {WordText: "Singa", Definition: "Hewan buas disebut raja hutan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm"), Category: constant.CategoryKosakata},
        {WordText: "Ular", Definition: "Hewan melata tanpa kaki", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ular.webm"), Category: constant.CategoryKosakata},
        {WordText: "Semut", Definition: "Hewan kecil yang suka bekerja sama", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Semut.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kupu-kupu", Definition: "Hewan kecil bersayap indah", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/KupuKupu.webm"), Category: constant.CategoryKosakata},
        {WordText: "Lebah", Definition: "Hewan penghasil madu", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kelinci", Definition: "Hewan dengan telinga panjang", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm"), Category: constant.CategoryKosakata},
        {WordText: "Ikan", Definition: "Hewan yang hidup di air", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), Category: constant.CategoryKosakata},
        {WordText: "Burung", Definition: "Hewan yang bisa terbang", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kucing", Definition: "Hewan peliharaan yang suka mengeong", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm"), Category: constant.CategoryKosakata},

        // Keluarga
        {WordText: "Ayah", Definition: "Orang tua laki-laki", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kakak", Definition: "Saudara yang lebih tua", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kakek", Definition: "Ayah dari orang tua", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm"), Category: constant.CategoryKosakata},
        {WordText: "Nenek", Definition: "Ibu dari orang tua", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm"), Category: constant.CategoryKosakata},
        {WordText: "Paman", Definition: "Saudara laki-laki dari orang tua", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm"), Category: constant.CategoryKosakata},
        {WordText: "Bibi", Definition: "Saudara perempuan dari orang tua", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bibi.webm"), Category: constant.CategoryKosakata},
        {WordText: "Teman", Definition: "Orang yang akrab dengan kita", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm"), Category: constant.CategoryKosakata},
        {WordText: "Guru", Definition: "Orang yang mengajar di sekolah", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm"), Category: constant.CategoryKosakata},
        {WordText: "Sekolah", Definition: "Tempat belajar anak-anak", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sekolah.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kura-kura", Definition: "Hewan bercangkang keras", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kura-kura.webm"), Category: constant.CategoryKosakata},
        {WordText: "Katak", Definition: "Hewan kecil yang pandai melompat", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Katak.webm"), Category: constant.CategoryKosakata},

        // Geometri
        {WordText: "Lingkaran", Definition: "Bentuk bulat tanpa sudut", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lingkaran.webm"), Category: constant.CategoryKosakata},
        {WordText: "Segitiga", Definition: "Bentuk dengan tiga sisi dan tiga sudut", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Segitiga.webm"), Category: constant.CategoryKosakata},
        {WordText: "Persegi", Definition: "Bentuk dengan empat sisi sama panjang dan empat sudut siku-siku", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Persegi.webm"), Category: constant.CategoryKosakata},
        {WordText: "Garis", Definition: "Bentuk lurus tanpa lebar dan tebal", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Garis.webm"), Category: constant.CategoryKosakata},
        {WordText: "Kotak", Definition: "Bentuk tiga dimensi dengan enam sisi persegi", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kotak.webm"), Category: constant.CategoryKosakata},
        {WordText: "Bola", Definition: "Bentuk bulat tiga dimensi", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bola.webm"), Category: constant.CategoryKosakata},

        // Waktu
        {WordText: "Hari", Definition: "Satuan waktu selama 24 jam", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hari.webm"), Category: constant.CategoryKosakata},
        {WordText: "Minggu", Definition: "Satuan waktu selama 7 hari", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minggu.webm"), Category: constant.CategoryKosakata},
        {WordText: "Bulan", Definition: "Satuan waktu selama sekitar 30 hari", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bulan.webm"), Category: constant.CategoryKosakata},
        {WordText: "Tahun", Definition: "Satuan waktu selama 12 bulan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tahun.webm"), Category: constant.CategoryKosakata},
        {WordText: "Pagi", Definition: "Waktu setelah matahari terbit hingga siang", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pagi.webm"), Category: constant.CategoryKosakata},
        {WordText: "Siang", Definition: "Waktu setelah pagi hingga sore", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Siang.webm"), Category: constant.CategoryKosakata},
        {WordText: "Sore", Definition: "Waktu setelah siang hingga matahari terbenam", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sore.webm"), Category: constant.CategoryKosakata},
        {WordText: "Malam", Definition: "Waktu setelah matahari terbenam hingga pagi", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Malam.webm"), Category: constant.CategoryKosakata},
        {WordText: "Detik", Definition: "Satuan waktu terkecil", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Detik.webm"), Category: constant.CategoryKosakata},
        {WordText: "Menit", Definition: "Satuan waktu selama 60 detik", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Menit.webm"), Category: constant.CategoryKosakata},
        {WordText: "Jam", Definition: "Satuan waktu selama 60 menit", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jam.webm"), Category: constant.CategoryKosakata},

        // Aktivitas
        {WordText: "Bangun", Definition: "Kegiatan setelah tidur", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm"), Category: constant.CategoryKosakata},
        {WordText: "Mandi", Definition: "Membersihkan badan dengan air", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm"), Category: constant.CategoryKosakata},
        {WordText: "Sarap", Definition: "Makan pagi sebelum beraktivitas", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarap.webm"), Category: constant.CategoryKosakata},
        {WordText: "Ajar", Definition: "Kegiatan untuk mendapatkan ilmu", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ajar.webm"), Category: constant.CategoryKosakata},
        {WordText: "Tulis", Definition: "Membuat huruf atau kata dengan alat tulis", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm"), Category: constant.CategoryKosakata},
        {WordText: "Baca", Definition: "Melihat dan memahami tulisan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm"), Category: constant.CategoryKosakata},
        {WordText: "Main", Definition: "Melakukan kegiatan yang menyenangkan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm"), Category: constant.CategoryKosakata},
        {WordText: "Tidur", Definition: "Beristirahat dengan memejamkan mata", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm"), Category: constant.CategoryKosakata},
        {WordText: "Duduk", Definition: "Posisi badan di atas kursi", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm"), Category: constant.CategoryKosakata},
        {WordText: "Nonton", Definition: "Melihat acara di televisi atau layar", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm"), Category: constant.CategoryKosakata},
        {WordText: "Lari", Definition: "Bergerak cepat dengan kaki", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm"), Category: constant.CategoryKosakata},
        {WordText: "Lompat", Definition: "Berpindah dengan melompatkan tubuh", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm"), Category: constant.CategoryKosakata},
        {WordText: "Nyanyi", Definition: "Mengeluarkan suara dengan nada lagu", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nyanyi.webm"), Category: constant.CategoryKosakata},
        {WordText: "Cuci", Definition: "Membersihkan dengan air dan sabun", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Cuci.webm"), Category: constant.CategoryKosakata},
        {WordText: "Tangan", Definition: "Anggota badan untuk memegang", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tangan.webm"), Category: constant.CategoryKosakata},
        {WordText: "Sapu", Definition: "Alat untuk membersihkan lantai", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm"), Category: constant.CategoryKosakata},
        {WordText: "Sikat", Definition: "Alat untuk menggosok dan membersihkan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sikat.webm"), Category: constant.CategoryKosakata},
        {WordText: "Gigi", Definition: "Bagian mulut untuk mengunyah makanan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gigi.webm"), Category: constant.CategoryKosakata},
        {WordText: "Lapar", Definition: "Keadaan ingin makan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lapar.webm"), Category: constant.CategoryKosakata},
        {WordText: "Senang", Definition: "Perasaan gembira", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm"), Category: constant.CategoryKosakata},
        {WordText: "Sedih", Definition: "Perasaan tidak bahagia", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm"), Category: constant.CategoryKosakata},
        {WordText: "Marah", Definition: "Perasaan kesal yang kuat", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Marah.webm"), Category: constant.CategoryKosakata},
        {WordText: "Takut", Definition: "Perasaan tidak berani atau khawatir", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Takut.webm"), Category: constant.CategoryKosakata},
        {WordText: "Jalan", Definition: "Bergerak dengan kaki perlahan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jalan.webm"), Category: constant.CategoryKosakata},
        {WordText: "Hujan", Definition: "Air yang jatuh dari langit", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hujan.webm"), Category: constant.CategoryKosakata},
        {WordText: "Taman", Definition: "Tempat dengan tanaman dan bunga", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Taman.webm"), Category: constant.CategoryKosakata},
        {WordText: "Mimpi", Definition: "Gambaran dalam tidur", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mimpi.webm"), Category: constant.CategoryKosakata},

        // Imbuhan
        {WordText: "Ber-", Definition: "Imbuhan awalan untuk kata kerja", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Ber.webm"), Category: constant.CategoryImbuhan},
        {WordText: "Ter-", Definition: "Imbuhan awalan untuk kata sifat", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Ter.webm"), Category: constant.CategoryImbuhan},
        {WordText: "Me-", Definition: "Imbuhan awalan untuk kata kerja aktif", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Me.webm"), Category: constant.CategoryImbuhan},
        {WordText: "Di-", Definition: "Imbuhan awalan untuk kata kerja pasif", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Di.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-kan", Definition: "Imbuhan akhiran untuk membentuk kata kerja", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Kan.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-i", Definition: "Imbuhan akhiran untuk membentuk kata kerja", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-I.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-an", Definition: "Imbuhan akhiran untuk membentuk kata benda", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-An.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-wan", Definition: "Imbuhan akhiran untuk membentuk kata benda pelaku", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Wan.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-wati", Definition: "Imbuhan akhiran untuk membentuk kata benda pelaku perempuan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Wati.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-man", Definition: "Imbuhan akhiran untuk membentuk kata benda pelaku", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Man.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-ti", Definition: "Imbuhan akhiran untuk membentuk kata benda", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Ti.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-nya", Definition: "Imbuhan akhiran untuk kepemilikan", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Nya.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-pun", Definition: "Partikel penegas", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Pun.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-lah", Definition: "Partikel penegas", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Lah.webm"), Category: constant.CategoryImbuhan},
        {WordText: "-kah", Definition: "Partikel penanya", VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Kah.webm"), Category: constant.CategoryImbuhan},
    }
}
