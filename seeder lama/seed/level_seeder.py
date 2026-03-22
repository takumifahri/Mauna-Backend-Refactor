from typing import List, Dict, Any
from sqlalchemy.orm import Session
from datetime import datetime
import os

try:
    from src.models.kamus import Kamus
    from src.models.level import Level
    from src.models.sublevel import SubLevel
    from src.database.seeder import BaseSeeder
except ImportError:
    from ...models.kamus import Kamus
    from ...models.level import Level
    from ...models.sublevel import SubLevel
    from ..seeder import BaseSeeder

class CompleteSeeder(BaseSeeder):
    """Seed Kamus, Levels, dan SubLevels"""
    
    def __init__(self):
        super().__init__()
    
    def run(self):
        """Run complete seeding process tanpa soal"""
        try:
            print("Starting complete seeding (Kamus, Levels, SubLevels)...")
            
            self.seed_kamus()
            self.seed_levels()
            self.seed_sublevels()
            
            print("Complete seeding finished successfully!")
            print("Note: Run SoalSeeder separately for question seeding")
            
        except Exception as e:
            self.db.rollback()
            raise Exception(f"Complete seeding failed: {e}")

    def seed_kamus(self):
        """Seed Kamus data dengan kategori yang dapat dijadikan soal"""
        
        kamus_data = [
            # Numbers (0-9) 
            {"word_text": "0", "definition": "Angka nol ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/0.webm", "image_url_ref": "kamus/0.png", "category": "NUMBERS"},
            {"word_text": "1", "definition": "Angka satu ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/01.webm", "image_url_ref": "kamus/1.png", "category": "NUMBERS"},
            {"word_text": "2", "definition": "Angka dua ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/02.webm", "image_url_ref": "kamus/2.png", "category": "NUMBERS"},
            {"word_text": "3", "definition": "Angka tiga ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/03.webm", "image_url_ref": "kamus/3.png", "category": "NUMBERS"},
            {"word_text": "4", "definition": "Angka empat ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/04.webm", "image_url_ref": "kamus/4.png", "category": "NUMBERS"},
            {"word_text": "5", "definition": "Angka lima ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/05.webm", "image_url_ref": "kamus/5.png", "category": "NUMBERS"},
            {"word_text": "6", "definition": "Angka enam ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/06.webm", "image_url_ref": "kamus/6.png", "category": "NUMBERS"},
            {"word_text": "7", "definition": "Angka tujuh ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/07.webm", "image_url_ref": "kamus/7.png", "category": "NUMBERS"},
            {"word_text": "8", "definition": "Angka delapan ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/08.webm", "image_url_ref": "kamus/8.png", "category": "NUMBERS"},
            {"word_text": "9", "definition": "Angka sembilan ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/09.webm", "image_url_ref": "kamus/9.png", "category": "NUMBERS"},
            
            # Alphabet 
            {"word_text": "A", "definition": "Huruf pertama dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm", "image_url_ref": "kamus/A.png", "category": "ALPHABET"},
            {"word_text": "B", "definition": "Huruf kedua dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/B.webm", "image_url_ref": "kamus/B.png", "category": "ALPHABET"},
            {"word_text": "C", "definition": "Huruf ketiga dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/C.webm", "image_url_ref": "kamus/C.png", "category": "ALPHABET"},
            {"word_text": "D", "definition": "Huruf keempat dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/D.webm", "image_url_ref": "kamus/D.png", "category": "ALPHABET"},
            {"word_text": "E", "definition": "Huruf kelima dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/E.webm", "image_url_ref": "kamus/E.png", "category": "ALPHABET"},
            {"word_text": "F", "definition": "Huruf keenam dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/F.webm", "image_url_ref": "kamus/F.png", "category": "ALPHABET"},
            {"word_text": "G", "definition": "Huruf ketujuh dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/G.webm", "image_url_ref": "kamus/G.png", "category": "ALPHABET"},
            {"word_text": "H", "definition": "Huruf kedelapan dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/H.webm", "image_url_ref": "kamus/H.png", "category": "ALPHABET"},
            {"word_text": "I", "definition": "Huruf kesembilan dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/I.webm", "image_url_ref": "kamus/I.png", "category": "ALPHABET"},
            {"word_text": "J", "definition": "Huruf kesepuluh dalam alfabet", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/J.webm", "image_url_ref": "kamus/J.png", "category": "ALPHABET"},
            {"word_text": "K", "definition": "Huruf kesebelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/K.webm", "image_url_ref": "kamus/K.png", "category": "ALPHABET"},
            {"word_text": "L", "definition": "Huruf keduabelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/L.webm", "image_url_ref": "kamus/L.png", "category": "ALPHABET"},
            {"word_text": "M", "definition": "Huruf ketigabelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/M.webm", "image_url_ref": "kamus/M.png", "category": "ALPHABET"},
            {"word_text": "N", "definition": "Huruf keempatbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/N.webm", "image_url_ref": "kamus/N.png", "category": "ALPHABET"},
            {"word_text": "O", "definition": "Huruf kelimabelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/O.webm", "image_url_ref": "kamus/O.png", "category": "ALPHABET"},
            {"word_text": "P", "definition": "Huruf keenambelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm", "image_url_ref": "kamus/P.png", "category": "ALPHABET"},
            {"word_text": "Q", "definition": "Huruf ketujuhbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Q.webm", "image_url_ref": "kamus/Q.png", "category": "ALPHABET"},
            {"word_text": "R", "definition": "Huruf kedelapanbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/R.webm", "image_url_ref": "kamus/R.png", "category": "ALPHABET"},
            {"word_text": "S", "definition": "Huruf kesembilanbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/S.webm", "image_url_ref": "kamus/S.png", "category": "ALPHABET"},
            {"word_text": "T", "definition": "Huruf keduapuluh dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/T.webm", "image_url_ref": "kamus/T.png", "category": "ALPHABET"},
            {"word_text": "U", "definition": "Huruf keduapuluh satu dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/U.webm", "image_url_ref": "kamus/U.png", "category": "ALPHABET"},
            {"word_text": "V", "definition": "Huruf keduapuluh dua dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/V.webm", "image_url_ref": "kamus/V.png", "category": "ALPHABET"},
            {"word_text": "W", "definition": "Huruf keduapuluh tiga dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/W.webm", "image_url_ref": "kamus/W.png", "category": "ALPHABET"},
            {"word_text": "X", "definition": "Huruf keduapuluh empat dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/X.webm", "image_url_ref": "kamus/X.png", "category": "ALPHABET"},
            {"word_text": "Y", "definition": "Huruf keduapuluh lima dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm", "image_url_ref": "kamus/Y.png", "category": "ALPHABET"},
            {"word_text": "Z", "definition": "Huruf keduapuluh enam dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Z.webm", "image_url_ref": "kamus/Z.png", "category": "ALPHABET"},

            #level2
            # hewan
            {"word_text": "Anjing", "definition": "Hewan peliharaan yang setia", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm", "category": "KOSAKATA"},
            {"word_text": "Bebek", "definition": "Hewan ternak yang suka berenang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", "category": "KOSAKATA"},
            {"word_text": "Kambing", "definition": "Hewan ternak penghasil susu dan daging", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", "category": "KOSAKATA"},
            {"word_text": "Gajah", "definition": "Hewan besar dengan belalai panjang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", "category": "KOSAKATA"},
            {"word_text": "Monyet", "definition": "Hewan cerdas yang suka memanjat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm", "category": "KOSAKATA"},
            {"word_text": "Singa", "definition": "Hewan buas disebut raja hutan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm", "category": "KOSAKATA"},
            {"word_text": "Ular", "definition": "Hewan melata tanpa kaki", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ular.webm", "category": "KOSAKATA"},
            {"word_text": "Semut", "definition": "Hewan kecil yang suka bekerja sama", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Semut.webm", "category": "KOSAKATA"},
            {"word_text": "Kupu-kupu", "definition": "Hewan kecil bersayap indah", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/KupuKupu.webm", "category": "KOSAKATA"},
            {"word_text": "Lebah", "definition": "Hewan penghasil madu", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm", "category": "KOSAKATA"},
            {"word_text": "Kelinci", "definition": "Hewan dengan telinga panjang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm", "category": "KOSAKATA"},
            {"word_text": "Ikan", "definition": "Hewan yang hidup di air", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "category": "KOSAKATA"},
            {"word_text": "Burung", "definition": "Hewan yang bisa terbang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "category": "KOSAKATA"},
            {"word_text": "Kucing", "definition": "Hewan peliharaan yang suka mengeong", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "category": "KOSAKATA"},
            
            #keluarga
            {"word_text": "Ayah", "definition": "Orang tua laki-laki", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", "category": "KOSAKATA"},
            {"word_text": "Kakak", "definition": "Saudara yang lebih tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm", "category": "KOSAKATA"},
            {"word_text": "Kakek", "definition": "Ayah dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm", "category": "KOSAKATA"},
            {"word_text": "Nenek", "definition": "Ibu dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm", "category": "KOSAKATA"},
            {"word_text": "Paman", "definition": "Saudara laki-laki dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm", "category": "KOSAKATA"},
            {"word_text": "Bibi", "definition": "Saudara perempuan dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bibi.webm", "category": "KOSAKATA"},
            {"word_text": "Teman", "definition": "Orang yang akrab dengan kita", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", "category": "KOSAKATA"},
            {"word_text": "Guru", "definition": "Orang yang mengajar di sekolah", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "category": "KOSAKATA"},
            {"word_text": "Sekolah", "definition": "Tempat belajar anak-anak", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sekolah.webm", "category": "KOSAKATA"},
            {"word_text": "Kura-kura", "definition": "Hewan bercangkang keras", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kura-kura.webm", "category": "KOSAKATA"},
            {"word_text": "Katak", "definition": "Hewan kecil yang pandai melompat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Katak.webm", "category": "KOSAKATA"},
            
            #kosakata level 3
            #geometri
            {"word_text": "Lingkaran", "definition": "Bentuk bulat tanpa sudut", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lingkaran.webm", "category": "KOSAKATA"},
            {"word_text": "Segitiga", "definition": "Bentuk dengan tiga sisi dan tiga sudut", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Segitiga.webm", "category": "KOSAKATA"},
            {"word_text": "Persegi", "definition": "Bentuk dengan empat sisi sama panjang dan empat sudut siku-siku", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Persegi.webm", "category": "KOSAKATA"},
            {"word_text": "Garis", "definition": "Bentuk lurus tanpa lebar dan tebal", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Garis.webm", "category": "KOSAKATA"},
            {"word_text": "Kotak", "definition": "Bentuk tiga dimensi dengan enam sisi persegi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kotak.webm", "category": "KOSAKATA"},
            {"word_text": "Bola", "definition": "Bentuk bulat tiga dimensi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bola.webm", "category": "KOSAKATA"},
            #waktu
            {"word_text": "Hari", "definition": "Satuan waktu selama 24 jam", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hari.webm", "category": "KOSAKATA"},
            {"word_text": "Minggu", "definition": "Satuan waktu selama 7 hari", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minggu.webm", "category": "KOSAKATA"},
            {"word_text": "Bulan", "definition": "Satuan waktu selama sekitar 30 hari", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bulan.webm", "category": "KOSAKATA"},
            {"word_text": "Tahun", "definition": "Satuan waktu selama 12 bulan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tahun.webm", "category": "KOSAKATA"},
            {"word_text": "Pagi", "definition": "Waktu setelah matahari terbit hingga siang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pagi.webm", "category": "KOSAKATA"},
            {"word_text": "Siang", "definition": "Waktu setelah pagi hingga sore", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Siang.webm", "category": "KOSAKATA"},
            {"word_text": "Sore", "definition": "Waktu setelah siang hingga matahari terbenam", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sore.webm", "category": "KOSAKATA"},
            {"word_text": "Malam", "definition": "Waktu setelah matahari terbenam hingga pagi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Malam.webm", "category": "KOSAKATA"},
            {"word_text": "Detik", "definition": "Satuan waktu terkecil", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Detik.webm", "category": "KOSAKATA"},
            {"word_text": "Menit", "definition": "Satuan waktu selama 60 detik", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Menit.webm", "category": "KOSAKATA"},
            {"word_text": "Jam", "definition": "Satuan waktu selama 60 menit", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jam.webm", "category": "KOSAKATA"},

            #aktivitas
            # Level 4
            {"word_text": "Bangun", "definition": "Kegiatan setelah tidur", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", "category": "KOSAKATA"},
            {"word_text": "Mandi", "definition": "Membersihkan badan dengan air", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "category": "KOSAKATA"},
            {"word_text": "Sarap", "definition": "Makan pagi sebelum beraktivitas", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarap.webm", "category": "KOSAKATA"},
            {"word_text": "Ajar", "definition": "Kegiatan untuk mendapatkan ilmu", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ajar.webm", "category": "KOSAKATA"},
            {"word_text": "Tulis", "definition": "Membuat huruf atau kata dengan alat tulis", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm", "category": "KOSAKATA"},
            {"word_text": "Baca", "definition": "Melihat dan memahami tulisan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm", "category": "KOSAKATA"},
            {"word_text": "Main", "definition": "Melakukan kegiatan yang menyenangkan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "category": "KOSAKATA"},
            {"word_text": "Tidur", "definition": "Beristirahat dengan memejamkan mata", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "category": "KOSAKATA"},
            {"word_text": "Duduk", "definition": "Posisi badan di atas kursi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", "category": "KOSAKATA"},
            {"word_text": "Nonton", "definition": "Melihat acara di televisi atau layar", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm", "category": "KOSAKATA"},
            {"word_text": "Lari", "definition": "Bergerak cepat dengan kaki", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm", "category": "KOSAKATA"},
            {"word_text": "Lompat", "definition": "Berpindah dengan melompatkan tubuh", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm", "category": "KOSAKATA"},
            {"word_text": "Nyanyi", "definition": "Mengeluarkan suara dengan nada lagu", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nyanyi.webm", "category": "KOSAKATA"},
            {"word_text": "Cuci", "definition": "Membersihkan dengan air dan sabun", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Cuci.webm", "category": "KOSAKATA"},
            {"word_text": "Tangan", "definition": "Anggota badan untuk memegang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tangan.webm", "category": "KOSAKATA"},
            {"word_text": "Sapu", "definition": "Alat untuk membersihkan lantai", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm", "category": "KOSAKATA"},
            {"word_text": "Sikat", "definition": "Alat untuk menggosok dan membersihkan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sikat.webm", "category": "KOSAKATA"},
            {"word_text": "Gigi", "definition": "Bagian mulut untuk mengunyah makanan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gigi.webm", "category": "KOSAKATA"},
            {"word_text": "Lapar", "definition": "Keadaan ingin makan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lapar.webm", "category": "KOSAKATA"},
            {"word_text": "Senang", "definition": "Perasaan gembira", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "category": "KOSAKATA"},
            {"word_text": "Sedih", "definition": "Perasaan tidak bahagia", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm", "category": "KOSAKATA"},
            {"word_text": "Marah", "definition": "Perasaan kesal yang kuat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Marah.webm", "category": "KOSAKATA"},
            {"word_text": "Takut", "definition": "Perasaan tidak berani atau khawatir", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Takut.webm", "category": "KOSAKATA"},
            {"word_text": "Jalan", "definition": "Bergerak dengan kaki perlahan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jalan.webm", "category": "KOSAKATA"},
            {"word_text": "Hujan", "definition": "Air yang jatuh dari langit", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hujan.webm", "category": "KOSAKATA"},
            {"word_text": "Taman", "definition": "Tempat dengan tanaman dan bunga", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Taman.webm", "category": "KOSAKATA"},
            {"word_text": "Mimpi", "definition": "Gambaran dalam tidur", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mimpi.webm", "category": "KOSAKATA"},

            # Imbuhan Awalan
            {"word_text": "Ber-", "definition": "Imbuhan awalan untuk kata kerja", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Ber.webm", "category": "IMBUHAN"},
            {"word_text": "Ter-", "definition": "Imbuhan awalan untuk kata sifat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Ter.webm", "category": "IMBUHAN"},
            {"word_text": "Me-", "definition": "Imbuhan awalan untuk kata kerja aktif", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Me.webm", "category": "IMBUHAN"},
            {"word_text": "Di-", "definition": "Imbuhan awalan untuk kata kerja pasif", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Di.webm", "category": "IMBUHAN"},
            # Imbuhan Akhiran & Partikel
            {"word_text": "-kan", "definition": "Imbuhan akhiran untuk membentuk kata kerja", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Kan.webm", "category": "IMBUHAN"},
            {"word_text": "-i", "definition": "Imbuhan akhiran untuk membentuk kata kerja", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-I.webm", "category": "IMBUHAN"},
            {"word_text": "-an", "definition": "Imbuhan akhiran untuk membentuk kata benda", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-An.webm", "category": "IMBUHAN"},
            {"word_text": "-wan", "definition": "Imbuhan akhiran untuk membentuk kata benda pelaku", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Wan.webm", "category": "IMBUHAN"},
            {"word_text": "-wati", "definition": "Imbuhan akhiran untuk membentuk kata benda pelaku perempuan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Wati.webm", "category": "IMBUHAN"},
            {"word_text": "-man", "definition": "Imbuhan akhiran untuk membentuk kata benda pelaku", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Man.webm", "category": "IMBUHAN"},
            {"word_text": "-ti", "definition": "Imbuhan akhiran untuk membentuk kata benda", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Ti.webm", "category": "IMBUHAN"},
            {"word_text": "-nya", "definition": "Imbuhan akhiran untuk kepemilikan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Nya.webm", "category": "IMBUHAN"},
            # Imbuhan Partikel
            {"word_text": "-pun", "definition": "Partikel penegas", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Pun.webm", "category": "IMBUHAN"},
            {"word_text": "-lah", "definition": "Partikel penegas", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Lah.webm", "category": "IMBUHAN"},
            {"word_text": "-kah", "definition": "Partikel penanya", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Kah.webm", "category": "IMBUHAN"},

        ]
        
        created_count = 0
        for data in kamus_data:
            existing = self.db.query(Kamus).filter(Kamus.word_text == data["word_text"]).first()
            if not existing:
                kamus = Kamus(**data)
                self.db.add(kamus)
                created_count += 1
                print(f"  Created kamus: {data['word_text']}")
            else:
                print(f"  Kamus already exists: {data['word_text']}")
        
        self.db.commit()
        print(f"Kamus seeding completed. Created {created_count} entries.")

    def seed_levels(self):
        """Seed Level data"""
        print("Seeding Levels...")
        
        levels_data = [
            {"name": "Level 1", "description": "Pengenalan Abjad", "tujuan": "Mengenal dan memahami bahasa isyarat alphabet A-Z"},
            {"name": "Level 2", "description": "Pengenalan Kata", "tujuan": "Mengenal dan memahami bahasa isyarat kosakata dasar sehari-hari"},
            {"name": "Level 3", "description": "Angka dan Matematika dasar", "tujuan": "Mengenal dan memahami bahasa isyarat matematika dasar (penjumlahan dan pengurangan)"},
            {"name": "Level 4", "description": "Pengenalan Kata", "tujuan": "Mengenal dan memahami bahasa isyarat kosakata lanjutan"},
        ]
        
        created_count = 0
        for data in levels_data:
            existing = self.db.query(Level).filter(Level.name == data["name"]).first()
            if not existing:
                level = Level(**data)
                self.db.add(level)
                created_count += 1
                print(f"  Created level: {data['name']}")
            else:
                print(f"  Level already exists: {data['name']}")
        
        self.db.commit()
        print(f"Level seeding completed. Created {created_count} levels.")

    def seed_sublevels(self):
        """Seed SubLevel data"""
        print("Seeding SubLevels...")
        
        #1 level 10 sublevels
        sublevels_data = [
            # ======================
            # LEVEL 1 - Huruf dan Angka (Diperbaiki: A–Z lengkap)
            # ======================
            {"name": "Sublevel 1.1", "description": "Huruf A-C", "tujuan": "Belajar isyarat huruf A, B, C", "level_id": 1},
            {"name": "Sublevel 1.2", "description": "Huruf D-F", "tujuan": "Belajar isyarat huruf D, E, F", "level_id": 1},
            {"name": "Sublevel 1.3", "description": "Huruf G-I", "tujuan": "Belajar isyarat huruf G, H, I", "level_id": 1},
            {"name": "Sublevel 1.4", "description": "Huruf J-L", "tujuan": "Belajar isyarat huruf J, K, L", "level_id": 1},
            {"name": "Sublevel 1.5", "description": "Huruf M-O", "tujuan": "Belajar isyarat huruf M, N, O", "level_id": 1},
            {"name": "Sublevel 1.6", "description": "Huruf P-R", "tujuan": "Belajar isyarat huruf P, Q, R", "level_id": 1},
            {"name": "Sublevel 1.7", "description": "Huruf S-U", "tujuan": "Belajar isyarat huruf S, T, U", "level_id": 1},
            {"name": "Sublevel 1.8", "description": "Huruf V-X", "tujuan": "Belajar isyarat huruf V, W, X", "level_id": 1},
            {"name": "Sublevel 1.9", "description": "Huruf Y-Z", "tujuan": "Belajar isyarat huruf Y dan Z", "level_id": 1},
            {"name": "Sublevel 1.10", "description": "Latihan Huruf", "tujuan": "Latihan mengenal isyarat huruf A sampai Z", "level_id": 1},

            # ======================
            # LEVEL 2 - Hewan dan Keluarga (sederhana untuk SD Kelas 1)
            # ======================
            {"name": "Sublevel 2.1", "description": "Hewan Rumah", "tujuan": "Belajar isyarat kucing, anjing, burung, dan ikan", "level_id": 2},
            {"name": "Sublevel 2.2", "description": "Hewan Ternak", "tujuan": "Belajar isyarat ayam, sapi, kambing, dan bebek", "level_id": 2},
            {"name": "Sublevel 2.3", "description": "Hewan Liar", "tujuan": "Belajar isyarat gajah, monyet, singa, dan ular", "level_id": 2},
            {"name": "Sublevel 2.4", "description": "Hewan Kecil", "tujuan": "Belajar isyarat semut, kupu-kupu, lebah, dan kelinci", "level_id": 2},
            {"name": "Sublevel 2.5", "description": "Keluarga Inti", "tujuan": "Belajar isyarat ayah, ibu, kakak, dan adik", "level_id": 2},
            {"name": "Sublevel 2.6", "description": "Keluarga Besar", "tujuan": "Belajar isyarat kakek, nenek, paman, dan bibi", "level_id": 2},
            {"name": "Sublevel 2.7", "description": "Teman dan Guru", "tujuan": "Belajar isyarat teman, guru, dan sekolah", "level_id": 2},
            {"name": "Sublevel 2.8", "description": "Hewan Air", "tujuan": "Belajar isyarat ikan, paus, kura-kura, dan katak", "level_id": 2},
            {"name": "Sublevel 2.9", "description": "Hewan Udara", "tujuan": "Belajar isyarat burung, kupu-kupu, dan lebah", "level_id": 2},
            {"name": "Sublevel 2.10", "description": "Latihan Hewan dan Keluarga", "tujuan": "Latihan mengenal isyarat hewan dan keluarga", "level_id": 2},

            # ======================
            # LEVEL 3 - Matematika Dasar (sederhana)
            # ======================
            {"name": "Sublevel 3.1", "description": "Angka 0-2", "tujuan": "Belajar isyarat angka 0 sampai 2", "level_id": 3},
            {"name": "Sublevel 3.2", "description": "Angka 3-5", "tujuan": "Belajar isyarat angka 3 sampai 5", "level_id": 3},
            {"name": "Sublevel 3.3", "description": "Angka 6-7", "tujuan": "Belajar isyarat angka 6 sampai 7", "level_id": 3},
            {"name": "Sublevel 3.4", "description": "Angka 8-9", "tujuan": "Belajar isyarat angka 8 sampai 9", "level_id": 3},
            {"name": "Sublevel 3.5", "description": "Angka 0-9", "tujuan": "Review isyarat angka 0 sampai 9", "level_id": 3},
            {"name": "Sublevel 3.6", "description": "Penjumlahan", "tujuan": "Belajar isyarat penjumlahan angka 1 digit (0-9)", "level_id": 3},
            {"name": "Sublevel 3.7", "description": "Pengurangan", "tujuan": "Belajar isyarat pengurangan angka 1 digit (0-9)", "level_id": 3},
            {"name": "Sublevel 3.8", "description": "Bentuk Geometri", "tujuan": "Belajar isyarat bentuk geometri dasar (lingkaran, segitiga, kotak, dan garis)", "level_id": 3},
            {"name": "Sublevel 3.9", "description": "Waktu", "tujuan": "Belajar isyarat waktu dasar (hari, minggu, bulan, tahun, pagi, siang, sore, malam)", "level_id": 3},
            {"name": "Sublevel 3.10", "description": "Matematika Dasar", "tujuan": "Latihan mengenal isyarat matematika dasar", "level_id": 3},

            # ======================
            # LEVEL 4 - Aktivitas Sehari-hari (sederhana)
            # ======================
            {"name": "Sublevel 4.1", "description": "Pagi Hari", "tujuan": "Belajar isyarat bangun, mandi, dan sarapan", "level_id": 4},
            {"name": "Sublevel 4.2", "description": "Sekolah", "tujuan": "Belajar isyarat belajar, menulis, membaca, dan bermain", "level_id": 4},
            {"name": "Sublevel 4.3", "description": "Di Rumah", "tujuan": "Belajar isyarat makan, tidur, duduk, dan nonton", "level_id": 4},
            {"name": "Sublevel 4.4", "description": "Bermain", "tujuan": "Belajar isyarat main bola, lari, lompat, dan nyanyi", "level_id": 4},
            {"name": "Sublevel 4.5", "description": "Bersih-bersih", "tujuan": "Belajar isyarat cuci tangan, sapu, dan sikat gigi", "level_id": 4},
            {"name": "Sublevel 4.6", "description": "Makan dan Minum", "tujuan": "Belajar isyarat makan, minum, dan lapar", "level_id": 4},
            {"name": "Sublevel 4.7", "description": "Emosi", "tujuan": "Belajar isyarat senang, sedih, marah, dan takut", "level_id": 4},
            {"name": "Sublevel 4.8", "description": "Kegiatan Luar", "tujuan": "Belajar isyarat jalan, hujan, panas, dan taman", "level_id": 4},
            {"name": "Sublevel 4.9", "description": "Waktu Istirahat", "tujuan": "Belajar isyarat tidur, mimpi, dan bangun", "level_id": 4},
            {"name": "Sublevel 4.10", "description": "Latihan Aktivitas", "tujuan": "Latihan mengenal isyarat aktivitas sehari-hari", "level_id": 4},
        ]
        
        created_count = 0
        for data in sublevels_data:
            existing = self.db.query(SubLevel).filter(
                SubLevel.name == data["name"],
                SubLevel.level_id == data["level_id"]
            ).first()
            
            if not existing:
                sublevel = SubLevel(**data)
                self.db.add(sublevel)
                created_count += 1
                print(f"  Created sublevel: {data['name']} (Level ID: {data['level_id']})")
            else:
                print(f"  SubLevel already exists: {data['name']}")
        
        self.db.commit()
        print(f"SubLevel seeding completed. Created {created_count} sublevels.")
