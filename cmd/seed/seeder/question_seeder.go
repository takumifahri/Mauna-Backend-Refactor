package seeder

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// SoalSeeder seeds soal (questions) with auto foreign key assignment
type SoalSeeder struct {
	db *sqlx.DB
}

// soalData represents question seed data
type soalData struct {
	Question        string
	Answer          string
	DictionaryID    int
	SublevelID      int
	VideoURL        *string
	ImageURL        *string
	SoalType        string
	PointGamifikasi int
}

// NewSoalSeeder creates a new SoalSeeder
func NewSoalSeeder(db *sqlx.DB) *SoalSeeder {
	return &SoalSeeder{
		db: db,
	}
}

// Name returns the seeder name
func (s *SoalSeeder) Name() string {
	return "SoalSeeder"
}

// Run executes the soal seeding process
func (s *SoalSeeder) Run() error {
	PrintInfo("Starting Soal seeding...")

	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	soalList := getSoalData()
	successCount := 0
	duplicateCount := 0

	query := `INSERT INTO soal (question, answer, dictionary_id, sublevel_id, video_url, image_url, categories, point_gamifikasi, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`

	for _, soal := range soalList {
		// Check for duplicates
		var count int
		err := tx.Get(&count, "SELECT COUNT(*) FROM soal WHERE question = $1", soal.Question)
		if err != nil {
			PrintError(fmt.Sprintf("Failed to check duplicate soal: %v", err))
			continue
		}

		if count > 0 {
			duplicateCount++
			continue
		}

		_, err = tx.Exec(query,
			soal.Question,
			soal.Answer,
			soal.DictionaryID,
			soal.SublevelID,
			soal.VideoURL,
			soal.ImageURL,
			soal.SoalType,
			soal.PointGamifikasi,
		)

		if err != nil {
			PrintWarning(fmt.Sprintf("Failed to insert soal '%s': %v", soal.Question, err))
			continue
		}

		successCount++
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	PrintSuccess(fmt.Sprintf("Soal seeding completed - Inserted: %d, Duplicates: %d", successCount, duplicateCount))
	return nil
}

// getSoalData returns all soal seed data
func getSoalData() []soalData {
	soalData := []soalData{
		// Level 1: Alphabet A-Z (Sublevel 1.1: A-C)
		{Question: "Tunjukkan isyarat huruf 'A'", Answer: "Kepalkan tangan dengan ibu jari di samping", DictionaryID: 11, SublevelID: 1, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah isyarat huruf?", Answer: "A", DictionaryID: 11, SublevelID: 1, ImageURL: ptrStr("kamus/A.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk huruf 'A' adalah...", Answer: "kamus/A.png", DictionaryID: 11, SublevelID: 1, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Bagaimana cara membuat isyarat huruf 'B'?", Answer: "Tangan terbuka dengan jari rapat, ibu jari menekuk", DictionaryID: 12, SublevelID: 1, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/B.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Isyarat ini adalah huruf?", Answer: "B", DictionaryID: 12, SublevelID: 1, ImageURL: ptrStr("kamus/B.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Huruf 'B' ditunjukkan dengan isyarat...", Answer: "kamus/B.png", DictionaryID: 12, SublevelID: 1, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat huruf 'A' kembali", Answer: "Kepalkan tangan dengan ibu jari di samping", DictionaryID: 11, SublevelID: 1, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat huruf 'B' sekali lagi", Answer: "Tangan terbuka dengan jari rapat, ibu jari menekuk", DictionaryID: 12, SublevelID: 1, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/B.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang terlihat pada gambar?", Answer: "A", DictionaryID: 11, SublevelID: 1, ImageURL: ptrStr("kamus/A.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Gambar berikut menunjukkan huruf?", Answer: "B", DictionaryID: 12, SublevelID: 1, ImageURL: ptrStr("kamus/B.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.2: D-F)
		{Question: "Tunjukkan isyarat untuk huruf 'D'", Answer: "Jari telunjuk tegak, jari lain menyentuh ibu jari", DictionaryID: 14, SublevelID: 2, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/D.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah isyarat untuk huruf?", Answer: "D", DictionaryID: 14, SublevelID: 2, ImageURL: ptrStr("kamus/D.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'D'", Answer: "kamus/D.png", DictionaryID: 14, SublevelID: 2, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'E' dengan benar", Answer: "Jari ditekuk ke dalam menyentuh ibu jari yang tertekuk", DictionaryID: 15, SublevelID: 2, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/E.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "E", DictionaryID: 15, SublevelID: 2, ImageURL: ptrStr("kamus/E.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'F'", Answer: "Jari telunjuk dan ibu jari bertemu, jari lain tegak", DictionaryID: 16, SublevelID: 2, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/F.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf dari gambar isyarat ini", Answer: "F", DictionaryID: 16, SublevelID: 2, ImageURL: ptrStr("kamus/F.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat yang mewakili huruf 'F'?", Answer: "kamus/F.png", DictionaryID: 16, SublevelID: 2, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Coba tunjukkan isyarat huruf 'E' sekali lagi", Answer: "Jari ditekuk ke dalam menyentuh ibu jari yang tertekuk", DictionaryID: 15, SublevelID: 2, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/E.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Pilih gambar yang menunjukkan isyarat 'D'", Answer: "kamus/D.png", DictionaryID: 14, SublevelID: 2, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.3: G-I)
		{Question: "Isyaratkan huruf 'G'", Answer: "Jari telunjuk dan ibu jari sejajar lurus", DictionaryID: 17, SublevelID: 3, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/G.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah isyarat untuk huruf?", Answer: "G", DictionaryID: 17, SublevelID: 3, ImageURL: ptrStr("kamus/G.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'G'", Answer: "kamus/G.png", DictionaryID: 17, SublevelID: 3, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat huruf 'H'", Answer: "Jari telunjuk dan tengah lurus, ibu jari di antara", DictionaryID: 18, SublevelID: 3, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/H.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "H", DictionaryID: 18, SublevelID: 3, ImageURL: ptrStr("kamus/H.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'H'", Answer: "kamus/H.png", DictionaryID: 18, SublevelID: 3, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'I'", Answer: "Jari kelingking tegak", DictionaryID: 19, SublevelID: 3, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/I.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf dari gambar isyarat ini", Answer: "I", DictionaryID: 19, SublevelID: 3, ImageURL: ptrStr("kamus/I.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat huruf 'I' kembali", Answer: "Jari kelingking tegak", DictionaryID: 19, SublevelID: 3, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/I.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang tepat untuk huruf 'H'", Answer: "kamus/H.png", DictionaryID: 18, SublevelID: 3, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.4: J-L)
		{Question: "Lakukan isyarat huruf 'K'", Answer: "Jari kelingking diayun ke bawah membentuk 'K'", DictionaryID: 21, SublevelID: 4, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/K.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar isyarat ini adalah huruf?", Answer: "J", DictionaryID: 20, SublevelID: 4, ImageURL: ptrStr("kamus/J.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'J'", Answer: "kamus/J.png", DictionaryID: 20, SublevelID: 4, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat huruf 'K'", Answer: "Jari telunjuk dan tengah tegak, ibu jari di antara", DictionaryID: 21, SublevelID: 4, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/K.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "K", DictionaryID: 21, SublevelID: 4, ImageURL: ptrStr("kamus/K.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'K'", Answer: "kamus/K.png", DictionaryID: 21, SublevelID: 4, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'L'", Answer: "Jari telunjuk tegak dan ibu jari lurus, membentuk L", DictionaryID: 22, SublevelID: 4, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/L.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf dari gambar isyarat ini", Answer: "L", DictionaryID: 22, SublevelID: 4, ImageURL: ptrStr("kamus/L.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat yang mewakili huruf 'L'?", Answer: "kamus/L.png", DictionaryID: 22, SublevelID: 4, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Coba tunjukkan isyarat huruf 'L' sekali lagi", Answer: "Jari telunjuk tegak dan ibu jari lurus, membentuk L", DictionaryID: 22, SublevelID: 4, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/L.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.5: M-O)
		{Question: "Tunjukkan isyarat huruf 'M'", Answer: "Tiga jari (telunjuk, tengah, manis) di atas ibu jari", DictionaryID: 23, SublevelID: 5, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/M.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar isyarat ini adalah huruf?", Answer: "M", DictionaryID: 23, SublevelID: 5, ImageURL: ptrStr("kamus/M.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'M'", Answer: "kamus/M.png", DictionaryID: 23, SublevelID: 5, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'N' dengan benar", Answer: "Dua jari (telunjuk, tengah) di atas ibu jari", DictionaryID: 24, SublevelID: 5, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/N.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "N", DictionaryID: 24, SublevelID: 5, ImageURL: ptrStr("kamus/N.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'N'", Answer: "kamus/N.png", DictionaryID: 24, SublevelID: 5, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'O'", Answer: "Jari membentuk lingkaran", DictionaryID: 25, SublevelID: 5, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/O.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf dari gambar isyarat ini", Answer: "O", DictionaryID: 25, SublevelID: 5, ImageURL: ptrStr("kamus/O.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat yang mewakili huruf 'O'?", Answer: "kamus/O.png", DictionaryID: 25, SublevelID: 5, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat huruf 'N' kembali", Answer: "Dua jari (telunjuk, tengah) di atas ibu jari", DictionaryID: 24, SublevelID: 5, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/N.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.6: P-R)
		{Question: "Tunjukkan isyarat huruf 'P'", Answer: "Sama dengan 'K', tapi menghadap ke bawah/depan", DictionaryID: 26, SublevelID: 6, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar isyarat ini adalah huruf?", Answer: "P", DictionaryID: 26, SublevelID: 6, ImageURL: ptrStr("kamus/P.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'P'", Answer: "kamus/P.png", DictionaryID: 26, SublevelID: 6, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'Q' dengan benar", Answer: "Jari telunjuk dan ibu jari menunjuk ke bawah", DictionaryID: 27, SublevelID: 6, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Q.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "Q", DictionaryID: 27, SublevelID: 6, ImageURL: ptrStr("kamus/Q.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'Q'", Answer: "kamus/Q.png", DictionaryID: 27, SublevelID: 6, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'R'", Answer: "Jari tengah menyilang di atas jari telunjuk", DictionaryID: 28, SublevelID: 6, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/R.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf dari gambar isyarat ini", Answer: "R", DictionaryID: 28, SublevelID: 6, ImageURL: ptrStr("kamus/R.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat yang mewakili huruf 'R'?", Answer: "kamus/R.png", DictionaryID: 28, SublevelID: 6, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Coba tunjukkan isyarat huruf 'P' sekali lagi", Answer: "Sama dengan 'K', tapi menghadap ke bawah/depan", DictionaryID: 26, SublevelID: 6, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.7: S-U)
		{Question: "Tunjukkan isyarat huruf 'S'", Answer: "Kepalan tangan, ibu jari di depan jari lain", DictionaryID: 29, SublevelID: 7, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/S.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar isyarat ini adalah huruf?", Answer: "S", DictionaryID: 29, SublevelID: 7, ImageURL: ptrStr("kamus/S.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'S'", Answer: "kamus/S.png", DictionaryID: 29, SublevelID: 7, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'T' dengan benar", Answer: "Ibu jari masuk ke antara jari telunjuk dan tengah", DictionaryID: 30, SublevelID: 7, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/T.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "T", DictionaryID: 30, SublevelID: 7, ImageURL: ptrStr("kamus/T.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'T'", Answer: "kamus/T.png", DictionaryID: 30, SublevelID: 7, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'U'", Answer: "Jari telunjuk dan tengah tegak dan rapat", DictionaryID: 31, SublevelID: 7, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/U.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf dari gambar isyarat ini", Answer: "U", DictionaryID: 31, SublevelID: 7, ImageURL: ptrStr("kamus/U.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat yang mewakili huruf 'U'?", Answer: "kamus/U.png", DictionaryID: 31, SublevelID: 7, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Coba tunjukkan isyarat huruf 'S' kembali", Answer: "Kepalan tangan, ibu jari di depan jari lain", DictionaryID: 29, SublevelID: 7, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/S.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.8: V-X)
		{Question: "Tunjukkan isyarat huruf 'V'", Answer: "Jari telunjuk dan tengah tegak membentuk V", DictionaryID: 32, SublevelID: 8, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/V.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar isyarat ini adalah huruf?", Answer: "V", DictionaryID: 32, SublevelID: 8, ImageURL: ptrStr("kamus/V.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'V'", Answer: "kamus/V.png", DictionaryID: 32, SublevelID: 8, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'W' dengan benar", Answer: "Tiga jari (telunjuk, tengah, manis) tegak", DictionaryID: 33, SublevelID: 8, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/W.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "W", DictionaryID: 33, SublevelID: 8, ImageURL: ptrStr("kamus/W.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'W'", Answer: "kamus/W.png", DictionaryID: 33, SublevelID: 8, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'X'", Answer: "Jari telunjuk ditekuk seperti kait", DictionaryID: 34, SublevelID: 8, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/X.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf dari gambar isyarat ini", Answer: "X", DictionaryID: 34, SublevelID: 8, ImageURL: ptrStr("kamus/X.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat yang mewakili huruf 'X'?", Answer: "kamus/X.png", DictionaryID: 34, SublevelID: 8, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Coba tunjukkan isyarat huruf 'V' kembali", Answer: "Jari telunjuk dan tengah tegak membentuk V", DictionaryID: 32, SublevelID: 8, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/V.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.9: Y-Z)
		{Question: "Tunjukkan isyarat huruf 'Y'", Answer: "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", DictionaryID: 35, SublevelID: 9, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar isyarat ini adalah huruf?", Answer: "Y", DictionaryID: 35, SublevelID: 9, ImageURL: ptrStr("kamus/Y.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'Y'", Answer: "kamus/Y.png", DictionaryID: 35, SublevelID: 9, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'Y' dengan benar", Answer: "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", DictionaryID: 35, SublevelID: 9, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", Answer: "Z", DictionaryID: 36, SublevelID: 9, ImageURL: ptrStr("kamus/Z.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'Z'", Answer: "kamus/Z.png", DictionaryID: 36, SublevelID: 9, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'Y' sekali lagi", Answer: "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", DictionaryID: 35, SublevelID: 9, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Identifikasi huruf pada gambar isyarat ini", Answer: "Y", DictionaryID: 35, SublevelID: 9, ImageURL: ptrStr("kamus/Y.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang tepat untuk 'Z'", Answer: "kamus/Z.png", DictionaryID: 36, SublevelID: 9, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Coba tunjukkan isyarat huruf 'Y' di depan kamera", Answer: "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", DictionaryID: 35, SublevelID: 9, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 1: Alphabet A-Z (Sublevel 1.10: Review A-Z)
		{Question: "Tunjukkan isyarat huruf 'A' (Review)", Answer: "Kepalkan tangan dengan ibu jari di samping", DictionaryID: 11, SublevelID: 10, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah isyarat huruf?", Answer: "F", DictionaryID: 16, SublevelID: 10, ImageURL: ptrStr("kamus/F.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'K'", Answer: "kamus/K.png", DictionaryID: 21, SublevelID: 10, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyaratkan huruf 'P'", Answer: "Sama dengan 'K', tapi menghadap ke bawah/depan", DictionaryID: 26, SublevelID: 10, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Huruf apa yang terlihat pada gambar?", Answer: "U", DictionaryID: 31, SublevelID: 10, ImageURL: ptrStr("kamus/U.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tentukan isyarat yang tepat untuk huruf 'Z'", Answer: "kamus/Z.png", DictionaryID: 36, SublevelID: 10, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat huruf 'D'", Answer: "Jari telunjuk tegak, jari lain menyentuh ibu jari", DictionaryID: 14, SublevelID: 10, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/D.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar berikut menunjukkan huruf?", Answer: "M", DictionaryID: 23, SublevelID: 10, ImageURL: ptrStr("kamus/M.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat yang benar untuk huruf 'R'", Answer: "kamus/R.png", DictionaryID: 28, SublevelID: 10, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Coba tunjukkan isyarat huruf 'W'", Answer: "Tiga jari (telunjuk, tengah, manis) tegak", DictionaryID: 33, SublevelID: 10, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/W.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.1: hewan rumah)
		{Question: "Tunjukkan isyarat untuk 'Kucing'", Answer: "kucing", DictionaryID: 50, SublevelID: 11, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan hewan?", Answer: "Kucing", DictionaryID: 50, SublevelID: 11, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Kucing' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", DictionaryID: 50, SublevelID: 11, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Burung'", Answer: "burung", DictionaryID: 49, SublevelID: 11, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan pada gambar ini adalah?", Answer: "Anjing", DictionaryID: 37, SublevelID: 11, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Anjing'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm", DictionaryID: 37, SublevelID: 11, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Burung'", Answer: "burung", DictionaryID: 49, SublevelID: 11, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan ini adalah?", Answer: "Burung", DictionaryID: 49, SublevelID: 11, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Burung' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", DictionaryID: 49, SublevelID: 11, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Ikan'", Answer: "ikan", DictionaryID: 48, SublevelID: 11, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.2: hewan ternak)
		{Question: "Tunjukkan isyarat untuk 'Bebek'", Answer: "bebek", DictionaryID: 38, SublevelID: 12, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan hewan?", Answer: "Bebek", DictionaryID: 38, SublevelID: 12, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Bebek' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", DictionaryID: 38, SublevelID: 12, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Kambing'", Answer: "kambing", DictionaryID: 39, SublevelID: 12, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan pada gambar ini adalah?", Answer: "Kambing", DictionaryID: 39, SublevelID: 12, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Kambing'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", DictionaryID: 39, SublevelID: 12, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Bebek' kembali", Answer: "bebek", DictionaryID: 38, SublevelID: 12, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan apa ini?", Answer: "Bebek", DictionaryID: 38, SublevelID: 12, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Kambing'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", DictionaryID: 39, SublevelID: 12, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Kambing'", Answer: "kambing", DictionaryID: 39, SublevelID: 12, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.3: hewan liar)
		{Question: "Tunjukkan isyarat untuk 'Gajah'", Answer: "gajah", DictionaryID: 40, SublevelID: 13, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan hewan?", Answer: "Gajah", DictionaryID: 40, SublevelID: 13, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Gajah' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", DictionaryID: 40, SublevelID: 13, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Singa'", Answer: "singa", DictionaryID: 42, SublevelID: 13, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan pada gambar ini adalah?", Answer: "Singa", DictionaryID: 42, SublevelID: 13, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Singa'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm", DictionaryID: 42, SublevelID: 13, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Monyet'", Answer: "monyet", DictionaryID: 41, SublevelID: 13, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan apa ini?", Answer: "Monyet", DictionaryID: 41, SublevelID: 13, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Gajah'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", DictionaryID: 40, SublevelID: 13, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Monyet'", Answer: "monyet", DictionaryID: 41, SublevelID: 13, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.4: hewan kecil)
		{Question: "Tunjukkan isyarat untuk 'Kelinci'", Answer: "kelinci", DictionaryID: 47, SublevelID: 14, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan hewan?", Answer: "Kelinci", DictionaryID: 47, SublevelID: 14, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Kelinci' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm", DictionaryID: 47, SublevelID: 14, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Semut'", Answer: "semut", DictionaryID: 44, SublevelID: 14, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Semut.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan pada gambar ini adalah?", Answer: "Semut", DictionaryID: 44, SublevelID: 14, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Semut.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Kupu-kupu'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/KupuKupu.webm", DictionaryID: 45, SublevelID: 14, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Lebah'", Answer: "lebah", DictionaryID: 46, SublevelID: 14, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan apa ini?", Answer: "Lebah", DictionaryID: 46, SublevelID: 14, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Lebah'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm", DictionaryID: 46, SublevelID: 14, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Kelinci'", Answer: "kelinci", DictionaryID: 47, SublevelID: 14, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.5: keluarga inti)
		{Question: "Tunjukkan isyarat untuk 'Ayah'", Answer: "ayah", DictionaryID: 51, SublevelID: 15, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan?", Answer: "Ayah", DictionaryID: 51, SublevelID: 15, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Ayah' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", DictionaryID: 51, SublevelID: 15, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Ibu'", Answer: "ibu", DictionaryID: 52, SublevelID: 15, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Siapa pada gambar ini?", Answer: "Ibu", DictionaryID: 52, SublevelID: 15, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Ibu'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm", DictionaryID: 52, SublevelID: 15, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Kakak'", Answer: "kakak", DictionaryID: 52, SublevelID: 15, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Siapa ini?", Answer: "Kakak", DictionaryID: 52, SublevelID: 15, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Kakak' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm", DictionaryID: 52, SublevelID: 15, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat untuk 'Ayah'", Answer: "ayah", DictionaryID: 51, SublevelID: 15, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.6: keluarga besar)
		{Question: "Tunjukkan isyarat untuk 'Kakek'", Answer: "kakek", DictionaryID: 53, SublevelID: 16, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan?", Answer: "Kakek", DictionaryID: 53, SublevelID: 16, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Kakek' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm", DictionaryID: 53, SublevelID: 16, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Nenek'", Answer: "nenek", DictionaryID: 54, SublevelID: 16, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Siapa pada gambar ini?", Answer: "Nenek", DictionaryID: 54, SublevelID: 16, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Nenek'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm", DictionaryID: 54, SublevelID: 16, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Paman'", Answer: "paman", DictionaryID: 55, SublevelID: 16, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Siapa ini?", Answer: "Paman", DictionaryID: 55, SublevelID: 16, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Paman' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm", DictionaryID: 55, SublevelID: 16, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Bibi'", Answer: "bibi", DictionaryID: 56, SublevelID: 16, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bibi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.7: teman dan guru)
		{Question: "Tunjukkan isyarat untuk 'Teman'", Answer: "teman", DictionaryID: 57, SublevelID: 17, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan?", Answer: "Teman", DictionaryID: 57, SublevelID: 17, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Teman' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", DictionaryID: 57, SublevelID: 17, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Guru'", Answer: "guru", DictionaryID: 58, SublevelID: 17, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Siapa pada gambar ini?", Answer: "Guru", DictionaryID: 58, SublevelID: 17, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Guru'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", DictionaryID: 58, SublevelID: 17, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Teman' kembali", Answer: "teman", DictionaryID: 57, SublevelID: 17, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Siapa ini?", Answer: "Teman", DictionaryID: 57, SublevelID: 17, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Guru'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", DictionaryID: 58, SublevelID: 17, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Guru'", Answer: "guru", DictionaryID: 58, SublevelID: 17, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.8: hewan air)
		{Question: "Tunjukkan isyarat untuk 'Ikan'", Answer: "ikan", DictionaryID: 48, SublevelID: 18, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan hewan?", Answer: "Ikan", DictionaryID: 48, SublevelID: 18, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Ikan' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", DictionaryID: 48, SublevelID: 18, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Hewan air lainnya adalah ikan. Tunjukkan isyarat 'Ikan'", Answer: "ikan", DictionaryID: 48, SublevelID: 18, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan pada gambar ini hidup di?", Answer: "Air", DictionaryID: 48, SublevelID: 18, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk hewan air 'Ikan'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", DictionaryID: 48, SublevelID: 18, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Ikan' sekali lagi", Answer: "ikan", DictionaryID: 48, SublevelID: 18, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan air apa ini?", Answer: "Ikan", DictionaryID: 48, SublevelID: 18, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Ikan'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", DictionaryID: 48, SublevelID: 18, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Ikan'", Answer: "ikan", DictionaryID: 48, SublevelID: 18, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.9: hewan udara)
		{Question: "Tunjukkan isyarat untuk 'Burung'", Answer: "burung", DictionaryID: 49, SublevelID: 19, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan hewan?", Answer: "Burung", DictionaryID: 49, SublevelID: 19, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Burung' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", DictionaryID: 49, SublevelID: 19, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Hewan yang terbang adalah burung. Tunjukkan isyarat 'Burung'", Answer: "burung", DictionaryID: 49, SublevelID: 19, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan pada gambar ini bisa?", Answer: "Terbang", DictionaryID: 49, SublevelID: 19, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk hewan udara 'Burung'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", DictionaryID: 49, SublevelID: 19, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Burung' sekali lagi", Answer: "burung", DictionaryID: 49, SublevelID: 19, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan udara apa ini?", Answer: "Burung", DictionaryID: 49, SublevelID: 19, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Burung'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", DictionaryID: 49, SublevelID: 19, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Burung'", Answer: "burung", DictionaryID: 49, SublevelID: 19, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 2: Basic Words (Sublevel 2.10: latihan campuran)
		{Question: "Tunjukkan isyarat untuk 'Kucing'", Answer: "kucing", DictionaryID: 50, SublevelID: 20, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan?", Answer: "Anjing", DictionaryID: 37, SublevelID: 20, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Ayah' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", DictionaryID: 51, SublevelID: 20, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Ibu'", Answer: "ibu", DictionaryID: 52, SublevelID: 20, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Siapa pada gambar ini?", Answer: "Guru", DictionaryID: 58, SublevelID: 20, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Kucing'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", DictionaryID: 50, SublevelID: 20, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Guru'", Answer: "guru", DictionaryID: 58, SublevelID: 20, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Hewan apa ini?", Answer: "Kucing", DictionaryID: 50, SublevelID: 20, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Ibu'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm", DictionaryID: 52, SublevelID: 20, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Burung'", Answer: "burung", DictionaryID: 49, SublevelID: 20, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 3: Numbers and Math (Sublevel 3.1: Numbers 0-2)
		{Question: "Tunjukkan isyarat angka '0'", Answer: "0", DictionaryID: 1, SublevelID: 21, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/0.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Angka berapa ini?", Answer: "0", DictionaryID: 1, SublevelID: 21, ImageURL: ptrStr("kamus/0.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat angka '0' adalah...", Answer: "kamus/0.png", DictionaryID: 1, SublevelID: 21, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "1 - 1 = ?", Answer: "0", DictionaryID: 1, SublevelID: 21, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Praktikkan isyarat angka '1'", Answer: "1", DictionaryID: 2, SublevelID: 21, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/01.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan angka?", Answer: "1", DictionaryID: 2, SublevelID: 21, ImageURL: ptrStr("kamus/1.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat angka '1'", Answer: "kamus/1.png", DictionaryID: 2, SublevelID: 21, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "0 + 1 = ?", Answer: "1", DictionaryID: 2, SublevelID: 21, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Tunjukkan isyarat angka '2'", Answer: "2", DictionaryID: 3, SublevelID: 21, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/02.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Angka pada gambar?", Answer: "2", DictionaryID: 3, SublevelID: 21, ImageURL: ptrStr("kamus/2.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},

		// Level 3: Numbers and Math (Sublevel 3.2: Numbers 3-5)
		{Question: "Tunjukkan isyarat angka '3'", Answer: "3", DictionaryID: 4, SublevelID: 22, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/03.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Angka berapa ini?", Answer: "3", DictionaryID: 4, SublevelID: 22, ImageURL: ptrStr("kamus/3.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat angka '3' adalah...", Answer: "kamus/3.png", DictionaryID: 4, SublevelID: 22, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "2 + 1 = ?", Answer: "3", DictionaryID: 4, SublevelID: 22, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Praktikkan isyarat angka '4'", Answer: "4", DictionaryID: 5, SublevelID: 22, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/04.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan angka?", Answer: "4", DictionaryID: 5, SublevelID: 22, ImageURL: ptrStr("kamus/4.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat angka '4'", Answer: "kamus/4.png", DictionaryID: 5, SublevelID: 22, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "3 + 1 = ?", Answer: "4", DictionaryID: 5, SublevelID: 22, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Tunjukkan isyarat angka '5'", Answer: "5", DictionaryID: 6, SublevelID: 22, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/05.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Angka pada gambar?", Answer: "5", DictionaryID: 6, SublevelID: 22, ImageURL: ptrStr("kamus/5.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},

		// Level 3: Numbers and Math (Sublevel 3.3: Numbers 6-7)
		{Question: "Tunjukkan isyarat angka '6'", Answer: "6", DictionaryID: 7, SublevelID: 23, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/06.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Angka berapa ini?", Answer: "6", DictionaryID: 7, SublevelID: 23, ImageURL: ptrStr("kamus/6.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat angka '6' adalah...", Answer: "kamus/6.png", DictionaryID: 7, SublevelID: 23, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "5 + 1 = ?", Answer: "6", DictionaryID: 7, SublevelID: 23, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Praktikkan isyarat angka '7'", Answer: "7", DictionaryID: 8, SublevelID: 23, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/07.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan angka?", Answer: "7", DictionaryID: 8, SublevelID: 23, ImageURL: ptrStr("kamus/7.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat angka '7'", Answer: "kamus/7.png", DictionaryID: 8, SublevelID: 23, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "6 + 1 = ?", Answer: "7", DictionaryID: 8, SublevelID: 23, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "7 - 1 = ?", Answer: "6", DictionaryID: 7, SublevelID: 23, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Tunjukkan kembali angka '6'", Answer: "6", DictionaryID: 7, SublevelID: 23, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/06.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 3: Numbers and Math (Sublevel 3.4: Numbers 8-9)
		{Question: "Tunjukkan isyarat angka '8'", Answer: "8", DictionaryID: 9, SublevelID: 24, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/08.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Angka berapa ini?", Answer: "8", DictionaryID: 9, SublevelID: 24, ImageURL: ptrStr("kamus/8.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat angka '8' adalah...", Answer: "kamus/8.png", DictionaryID: 9, SublevelID: 24, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "7 + 1 = ?", Answer: "8", DictionaryID: 9, SublevelID: 24, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Praktikkan isyarat angka '9'", Answer: "9", DictionaryID: 10, SublevelID: 24, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/09.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan angka?", Answer: "9", DictionaryID: 10, SublevelID: 24, ImageURL: ptrStr("kamus/9.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat angka '9'", Answer: "kamus/9.png", DictionaryID: 10, SublevelID: 24, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "8 + 1 = ?", Answer: "9", DictionaryID: 10, SublevelID: 24, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "9 - 1 = ?", Answer: "8", DictionaryID: 9, SublevelID: 24, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Tunjukkan kembali angka '8'", Answer: "8", DictionaryID: 9, SublevelID: 24, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/08.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 3: Numbers and Math (Sublevel 3.5: Review 0-9)
		{Question: "Tunjukkan isyarat angka '5'", Answer: "5", DictionaryID: 6, SublevelID: 25, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/05.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Angka berapa ini?", Answer: "3", DictionaryID: 4, SublevelID: 25, ImageURL: ptrStr("kamus/3.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat angka '9' adalah...", Answer: "kamus/9.png", DictionaryID: 10, SublevelID: 25, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "4 + 3 = ?", Answer: "7", DictionaryID: 8, SublevelID: 25, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Praktikkan isyarat angka '2'", Answer: "2", DictionaryID: 3, SublevelID: 25, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/02.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan angka?", Answer: "6", DictionaryID: 7, SublevelID: 25, ImageURL: ptrStr("kamus/6.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat angka '1'", Answer: "kamus/1.png", DictionaryID: 2, SublevelID: 25, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "9 - 5 = ?", Answer: "4", DictionaryID: 5, SublevelID: 25, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Tunjukkan isyarat angka '0'", Answer: "0", DictionaryID: 1, SublevelID: 25, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/0.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "8 - 0 = ?", Answer: "8", DictionaryID: 9, SublevelID: 25, SoalType: "MATEMATIKA", PointGamifikasi: 15},

		// Level 3: Numbers and Math (Sublevel 3.6: Penjumlahan)
		{Question: "1 + 1 = ?", Answer: "2", DictionaryID: 3, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "2 + 3 = ?", Answer: "5", DictionaryID: 6, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Berapa hasil 3 + 2?", Answer: "5", DictionaryID: 6, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "4 + 1 = ?", Answer: "5", DictionaryID: 6, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "5 + 3 = ?", Answer: "8", DictionaryID: 9, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Berapa 2 + 2?", Answer: "4", DictionaryID: 5, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "1 + 2 = ?", Answer: "3", DictionaryID: 4, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "0 + 5 = ?", Answer: "5", DictionaryID: 6, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "3 + 3 = ?", Answer: "6", DictionaryID: 7, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "1 + 8 = ?", Answer: "9", DictionaryID: 10, SublevelID: 26, SoalType: "MATEMATIKA", PointGamifikasi: 15},

		// Level 3: Numbers and Math (Sublevel 3.7: Pengurangan)
		{Question: "5 - 2 = ?", Answer: "3", DictionaryID: 4, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "8 - 3 = ?", Answer: "5", DictionaryID: 6, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Berapa hasil 9 - 4?", Answer: "5", DictionaryID: 6, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "7 - 2 = ?", Answer: "5", DictionaryID: 6, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "6 - 1 = ?", Answer: "5", DictionaryID: 6, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Berapa 4 - 2?", Answer: "2", DictionaryID: 3, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "9 - 5 = ?", Answer: "4", DictionaryID: 5, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "3 - 3 = ?", Answer: "0", DictionaryID: 1, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "8 - 1 = ?", Answer: "7", DictionaryID: 8, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "2 - 2 = ?", Answer: "0", DictionaryID: 1, SublevelID: 27, SoalType: "MATEMATIKA", PointGamifikasi: 15},

		// Level 3: Numbers and Math (Sublevel 3.8: Geometri)
		{Question: "Bentuk ini adalah?", Answer: "Lingkaran", DictionaryID: 62, SublevelID: 28, ImageURL: ptrStr("kamus/Lingkaran.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Lingkaran'", Answer: "Lingkaran", DictionaryID: 62, SublevelID: 28, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lingkaran.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Isyarat 'Segitiga' adalah...", Answer: "kamus/Segitiga.png", DictionaryID: 63, SublevelID: 28, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Persegi'", Answer: "Persegi", DictionaryID: 64, SublevelID: 28, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Persegi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Bentuk apa ini?", Answer: "Segitiga", DictionaryID: 63, SublevelID: 28, ImageURL: ptrStr("kamus/Segitiga.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Gambar ini adalah?", Answer: "Persegi", DictionaryID: 64, SublevelID: 28, ImageURL: ptrStr("kamus/Persegi.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat 'Lingkaran'", Answer: "kamus/Lingkaran.png", DictionaryID: 62, SublevelID: 28, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Segitiga'", Answer: "Segitiga", DictionaryID: 63, SublevelID: 28, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Segitiga.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Bentuk dengan 3 sisi adalah?", Answer: "Segitiga", DictionaryID: 63, SublevelID: 28, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Bentuk dengan 4 sisi sama adalah?", Answer: "Persegi", DictionaryID: 64, SublevelID: 28, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},

		// Level 3: Numbers and Math (Sublevel 3.9: Waktu)
		{Question: "Tunjukkan isyarat 'Hari'", Answer: "Hari", DictionaryID: 68, SublevelID: 29, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hari.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Kapan matahari terbit?", Answer: "Pagi", DictionaryID: 72, SublevelID: 29, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Pagi'", Answer: "Pagi", DictionaryID: 72, SublevelID: 29, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pagi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Waktu untuk makan siang adalah?", Answer: "Siang", DictionaryID: 73, SublevelID: 29, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Isyarat 'Sore' adalah...", Answer: "kamus/Sore.png", DictionaryID: 74, SublevelID: 29, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Malam'", Answer: "Malam", DictionaryID: 75, SublevelID: 29, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Malam.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "7 hari membentuk satu?", Answer: "Minggu", DictionaryID: 69, SublevelID: 29, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Minggu'", Answer: "Minggu", DictionaryID: 69, SublevelID: 29, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minggu.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Waktu tidur adalah?", Answer: "Malam", DictionaryID: 75, SublevelID: 29, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "12 bulan membentuk satu?", Answer: "Tahun", DictionaryID: 71, SublevelID: 29, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},

		// Level 3: Numbers and Math (Sublevel 3.10: Latihan Campuran)
		{Question: "5 + 3 = ?", Answer: "8", DictionaryID: 9, SublevelID: 30, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Bentuk ini adalah?", Answer: "Lingkaran", DictionaryID: 62, SublevelID: 30, ImageURL: ptrStr("kamus/Lingkaran.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "9 - 4 = ?", Answer: "5", DictionaryID: 6, SublevelID: 30, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Tunjukkan isyarat 'Segitiga'", Answer: "Segitiga", DictionaryID: 63, SublevelID: 30, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Segitiga.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "6 + 2 = ?", Answer: "8", DictionaryID: 9, SublevelID: 30, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Pilih isyarat angka '7'", Answer: "kamus/7.png", DictionaryID: 8, SublevelID: 30, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Waktu untuk tidur adalah?", Answer: "Malam", DictionaryID: 75, SublevelID: 30, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "10 - 3 = ?", Answer: "7", DictionaryID: 8, SublevelID: 30, SoalType: "MATEMATIKA", PointGamifikasi: 15},
		{Question: "Praktikkan isyarat angka '4'", Answer: "4", DictionaryID: 5, SublevelID: 30, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/04.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Bentuk dengan 4 sisi sama adalah?", Answer: "Persegi", DictionaryID: 64, SublevelID: 30, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},

		// Level 4: Basic Activities (Sublevel 4.1: Pagi Hari)
		{Question: "Tunjukkan isyarat untuk 'Bangun'", Answer: "bangun", DictionaryID: 79, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Bangun", DictionaryID: 79, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Mandi' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", DictionaryID: 80, SublevelID: 31, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Mandi'", Answer: "mandi", DictionaryID: 80, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Mandi", DictionaryID: 80, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Bangun'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", DictionaryID: 79, SublevelID: 31, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Sarapan'", Answer: "sarapan", DictionaryID: 81, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarapan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Di pagi hari, setelah bangun kita?", Answer: "Mandi", DictionaryID: 80, SublevelID: 31, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Sarapan", DictionaryID: 81, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarapan.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Bangun'", Answer: "bangun", DictionaryID: 79, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 4: Basic Activities (Sublevel 4.2: Sekolah)
		{Question: "Tunjukkan isyarat untuk 'Main'", Answer: "main", DictionaryID: 85, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Ajar", DictionaryID: 82, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ajar.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Tulis' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm", DictionaryID: 83, SublevelID: 32, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Baca'", Answer: "Baca", DictionaryID: 84, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Tulis", DictionaryID: 83, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Main'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", DictionaryID: 85, SublevelID: 32, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Duduk'", Answer: "duduk", DictionaryID: 87, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Di sekolah, siswa biasanya?", Answer: "Ajar", DictionaryID: 82, SublevelID: 32, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Baca", DictionaryID: 84, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Nonton'", Answer: "nonton", DictionaryID: 88, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 4: Basic Activities (Sublevel 4.3: Di Rumah)
		{Question: "Tunjukkan isyarat untuk 'Tidur'", Answer: "tidur", DictionaryID: 86, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Tidur", DictionaryID: 86, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Duduk' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", DictionaryID: 87, SublevelID: 33, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Nonton'", Answer: "nonton", DictionaryID: 88, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Duduk", DictionaryID: 87, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Makan'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", DictionaryID: 97, SublevelID: 33, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Makan'", Answer: "makan", DictionaryID: 97, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Di rumah, kita biasanya?", Answer: "Tidur", DictionaryID: 86, SublevelID: 33, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Nonton", DictionaryID: 88, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Duduk'", Answer: "duduk", DictionaryID: 87, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Level 4: Basic Activities (Sublevel 4.4package seeder
		// Level 4: Basic Activities (Sublevel 4.1-4.10)
		// Sublevel 31: Pagi Hari (Morning Activities) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Bangun'", Answer: "bangun", DictionaryID: 79, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Bangun", DictionaryID: 79, SublevelID: 31, ImageURL: ptrStr("kamus/bangun.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Mandi' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", DictionaryID: 80, SublevelID: 31, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Mandi'", Answer: "mandi", DictionaryID: 80, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Mandi", DictionaryID: 80, SublevelID: 31, ImageURL: ptrStr("kamus/mandi.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Bangun'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", DictionaryID: 79, SublevelID: 31, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Sarapan'", Answer: "sarapan", DictionaryID: 81, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarapan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Di pagi hari, setelah bangun kita?", Answer: "Mandi", DictionaryID: 80, SublevelID: 31, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Sarapan", DictionaryID: 81, SublevelID: 31, ImageURL: ptrStr("kamus/sarapan.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Bangun'", Answer: "bangun", DictionaryID: 79, SublevelID: 31, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Sublevel 32: Sekolah (School) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Main'", Answer: "main", DictionaryID: 85, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Ajar", DictionaryID: 82, SublevelID: 32, ImageURL: ptrStr("kamus/ajar.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Tulis' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm", DictionaryID: 83, SublevelID: 32, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Baca'", Answer: "Baca", DictionaryID: 84, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Tulis", DictionaryID: 83, SublevelID: 32, ImageURL: ptrStr("kamus/tulis.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Main'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", DictionaryID: 85, SublevelID: 32, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Duduk'", Answer: "duduk", DictionaryID: 87, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Di sekolah, siswa biasanya?", Answer: "Ajar", DictionaryID: 82, SublevelID: 32, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Baca", DictionaryID: 84, SublevelID: 32, ImageURL: ptrStr("kamus/baca.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Nonton'", Answer: "nonton", DictionaryID: 88, SublevelID: 32, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Sublevel 33: Di Rumah (At Home) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Tidur'", Answer: "tidur", DictionaryID: 86, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Tidur", DictionaryID: 86, SublevelID: 33, ImageURL: ptrStr("kamus/tidur.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Duduk' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", DictionaryID: 87, SublevelID: 33, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Nonton'", Answer: "nonton", DictionaryID: 88, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Duduk", DictionaryID: 87, SublevelID: 33, ImageURL: ptrStr("kamus/duduk.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Makan'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", DictionaryID: 97, SublevelID: 33, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Makan'", Answer: "makan", DictionaryID: 97, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Di rumah, kita biasanya?", Answer: "Tidur", DictionaryID: 86, SublevelID: 33, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Nonton", DictionaryID: 88, SublevelID: 33, ImageURL: ptrStr("kamus/nonton.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Duduk'", Answer: "duduk", DictionaryID: 87, SublevelID: 33, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Sublevel 34: Bermain (Playing) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Lari'", Answer: "lari", DictionaryID: 89, SublevelID: 34, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Lari", DictionaryID: 89, SublevelID: 34, ImageURL: ptrStr("kamus/lari.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Lompat' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm", DictionaryID: 90, SublevelID: 34, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Nyanyi'", Answer: "nyanyi", DictionaryID: 91, SublevelID: 34, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nyanyi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Lompat", DictionaryID: 90, SublevelID: 34, ImageURL: ptrStr("kamus/lompat.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Main'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", DictionaryID: 85, SublevelID: 34, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Main'", Answer: "main", DictionaryID: 85, SublevelID: 34, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Saat bermain, anak-anak suka?", Answer: "Lari", DictionaryID: 89, SublevelID: 34, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Nyanyi", DictionaryID: 91, SublevelID: 34, ImageURL: ptrStr("kamus/nyanyi.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Lompat'", Answer: "lompat", DictionaryID: 90, SublevelID: 34, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Sublevel 35: Bersih-bersih (Cleaning) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Cuci'", Answer: "cuci", DictionaryID: 92, SublevelID: 35, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Cuci.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Cuci", DictionaryID: 92, SublevelID: 35, ImageURL: ptrStr("kamus/cuci.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Sapu' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm", DictionaryID: 94, SublevelID: 35, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Mandi'", Answer: "mandi", DictionaryID: 80, SublevelID: 35, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Sapu", DictionaryID: 94, SublevelID: 35, ImageURL: ptrStr("kamus/sapu.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Mandi'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", DictionaryID: 80, SublevelID: 35, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Untuk membersihkan badan, kita?", Answer: "Mandi", DictionaryID: 80, SublevelID: 35, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Untuk membersihkan lantai, kita?", Answer: "Sapu", DictionaryID: 94, SublevelID: 35, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Gambar ini adalah aktivitas?", Answer: "Mandi", DictionaryID: 80, SublevelID: 35, ImageURL: ptrStr("kamus/mandi.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Sapu'", Answer: "sapu", DictionaryID: 94, SublevelID: 35, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Sublevel 36: Makan Minum (Eating & Drinking) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Makan'", Answer: "makan", DictionaryID: 97, SublevelID: 36, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Makan", DictionaryID: 97, SublevelID: 36, ImageURL: ptrStr("kamus/makan.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Minum' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minum.webm", DictionaryID: 97, SublevelID: 36, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Minum'", Answer: "minum", DictionaryID: 97, SublevelID: 36, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minum.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas pada gambar ini adalah?", Answer: "Minum", DictionaryID: 97, SublevelID: 36, ImageURL: ptrStr("kamus/minum.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Saat haus, kita?", Answer: "Minum", DictionaryID: 97, SublevelID: 36, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Saat lapar, kita?", Answer: "Makan", DictionaryID: 97, SublevelID: 36, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Makan'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", DictionaryID: 97, SublevelID: 36, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Minum'", Answer: "minum", DictionaryID: 97, SublevelID: 36, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minum.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Makan'", Answer: "makan", DictionaryID: 97, SublevelID: 36, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Sublevel 37: Emosi (Emotions) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Senang'", Answer: "senang", DictionaryID: 98, SublevelID: 37, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan perasaan?", Answer: "Senang", DictionaryID: 98, SublevelID: 37, ImageURL: ptrStr("kamus/senang.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Sedih' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm", DictionaryID: 99, SublevelID: 37, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Sedih'", Answer: "sedih", DictionaryID: 99, SublevelID: 37, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Perasaan pada gambar ini adalah?", Answer: "Sedih", DictionaryID: 99, SublevelID: 37, ImageURL: ptrStr("kamus/sedih.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Saat ulang tahun, kita merasa?", Answer: "Senang", DictionaryID: 98, SublevelID: 37, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Saat kehilangan mainan, kita merasa?", Answer: "Sedih", DictionaryID: 99, SublevelID: 37, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Senang'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", DictionaryID: 98, SublevelID: 37, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lawan dari senang adalah?", Answer: "Sedih", DictionaryID: 99, SublevelID: 37, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Senang'", Answer: "senang", DictionaryID: 98, SublevelID: 37, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},

		// Sublevel 38: Kegiatan Luar (Outdoor Activities) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Pergi'", Answer: "pergi", DictionaryID: 102, SublevelID: 38, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pergi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Jalan", DictionaryID: 102, SublevelID: 38, ImageURL: ptrStr("kamus/jalan.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Saat keluar rumah, kita biasanya?", Answer: "Jalan", DictionaryID: 102, SublevelID: 38, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Jalan'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jalan.webm", DictionaryID: 102, SublevelID: 38, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Untuk pergi ke sekolah, kita?", Answer: "Jalan", DictionaryID: 102, SublevelID: 38, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Pergi'", Answer: "pergi", DictionaryID: 102, SublevelID: 38, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pergi.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan?", Answer: "Pergi", DictionaryID: 102, SublevelID: 38, ImageURL: ptrStr("kamus/pergi.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat 'Datang' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Datang.webm", DictionaryID: 102, SublevelID: 38, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Lari'", Answer: "lari", DictionaryID: 89, SublevelID: 38, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Lawan dari pergi adalah?", Answer: "Datang", DictionaryID: 102, SublevelID: 38, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},

		// Sublevel 39: Waktu Istirahat (Rest Time) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Tidur'", Answer: "tidur", DictionaryID: 86, SublevelID: 39, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Tidur", DictionaryID: 86, SublevelID: 39, ImageURL: ptrStr("kamus/tidur.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Saat malam hari, kita?", Answer: "Tidur", DictionaryID: 86, SublevelID: 39, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Saat lelah, kita perlu?", Answer: "Tidur", DictionaryID: 86, SublevelID: 39, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Tidur'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", DictionaryID: 86, SublevelID: 39, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Aktivitas yang dilakukan di tempat tidur adalah?", Answer: "Tidur", DictionaryID: 86, SublevelID: 39, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Setelah tidur, kita akan?", Answer: "Bangun", DictionaryID: 79, SublevelID: 39, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Tidur' kembali", Answer: "tidur", DictionaryID: 86, SublevelID: 39, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Tidur yang cukup membuat badan?", Answer: "Sehat", DictionaryID: 86, SublevelID: 39, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Waktu tidur yang baik adalah?", Answer: "Malam", DictionaryID: 75, SublevelID: 39, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},

		// Sublevel 40: Latihan (Practice/Review) - 10 Soal
		{Question: "Tunjukkan isyarat untuk 'Main'", Answer: "main", DictionaryID: 85, SublevelID: 40, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Gambar ini menunjukkan aktivitas?", Answer: "Main", DictionaryID: 85, SublevelID: 40, ImageURL: ptrStr("kamus/main.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Isyarat untuk 'Makan' adalah...", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", DictionaryID: 97, SublevelID: 40, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Praktikkan isyarat 'Tidur'", Answer: "tidur", DictionaryID: 86, SublevelID: 40, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Perasaan pada gambar ini adalah?", Answer: "Senang", DictionaryID: 98, SublevelID: 40, ImageURL: ptrStr("kamus/senang.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Pilih isyarat untuk 'Main'", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", DictionaryID: 85, SublevelID: 40, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Tunjukkan isyarat 'Senang'", Answer: "senang", DictionaryID: 98, SublevelID: 40, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
		{Question: "Aktivitas apa ini?", Answer: "Ajar", DictionaryID: 82, SublevelID: 40, ImageURL: ptrStr("kamus/ajar.png"), SoalType: "TEBAK_GAMBAR", PointGamifikasi: 10},
		{Question: "Manakah isyarat untuk 'Tidur'?", Answer: "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", DictionaryID: 86, SublevelID: 40, SoalType: "PILIHAN_GANDA", PointGamifikasi: 10},
		{Question: "Lakukan isyarat 'Lari'", Answer: "lari", DictionaryID: 89, SublevelID: 40, VideoURL: ptrStr("http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm"), SoalType: "OPEN_CAMERA", PointGamifikasi: 10},
	}

	return soalData
}
