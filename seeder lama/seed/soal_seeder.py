from typing import List, Dict, Any
from sqlalchemy.orm import Session
from datetime import datetime
import random
import enum
import os

try:
    from src.models.kamus import Kamus
    from src.models.level import Level
    from src.models.sublevel import SubLevel
    from src.models.soal import Soal
    from src.database.seeder import BaseSeeder
except ImportError:
    from ...models.kamus import Kamus
    from ...models.level import Level
    from ...models.sublevel import SubLevel
    from ...models.soal import Soal
    from ..seeder import BaseSeeder

    kamus_data = [
        # Numbers (0-9) 
        {"id": 1, "word_text": "0", "definition": "Angka nol ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/0.webm", "image_url_ref": "kamus/0.png", "category": "NUMBERS"},
        {"id": 2, "word_text": "1", "definition": "Angka satu ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/01.webm", "image_url_ref": "kamus/1.png", "category": "NUMBERS"},
        {"id": 3, "word_text": "2", "definition": "Angka dua ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/02.webm", "image_url_ref": "kamus/2.png", "category": "NUMBERS"},
        {"id": 4, "word_text": "3", "definition": "Angka tiga ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/03.webm", "image_url_ref": "kamus/3.png", "category": "NUMBERS"},
        {"id": 5, "word_text": "4", "definition": "Angka empat ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/04.webm", "image_url_ref": "kamus/4.png", "category": "NUMBERS"},
        {"id": 6, "word_text": "5", "definition": "Angka lima ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/05.webm", "image_url_ref": "kamus/5.png", "category": "NUMBERS"},
        {"id": 7, "word_text": "6", "definition": "Angka enam ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/06.webm", "image_url_ref": "kamus/6.png", "category": "NUMBERS"},
        {"id": 8, "word_text": "7", "definition": "Angka tujuh ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/07.webm", "image_url_ref": "kamus/7.png", "category": "NUMBERS"},
        {"id": 9, "word_text": "8", "definition": "Angka delapan ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/08.webm", "image_url_ref": "kamus/8.png", "category": "NUMBERS"},
        {"id": 10, "word_text": "9", "definition": "Angka sembilan ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/09.webm", "image_url_ref": "kamus/9.png", "category": "NUMBERS"},
        
        # Alphabet 
        {"id": 11, "word_text": "A", "definition": "Huruf pertama dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm", "image_url_ref": "kamus/A.png", "category": "ALPHABET"},
        {"id": 12, "word_text": "B", "definition": "Huruf kedua dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/B.webm", "image_url_ref": "kamus/B.png", "category": "ALPHABET"},
        {"id": 13, "word_text": "C", "definition": "Huruf ketiga dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/C.webm", "image_url_ref": "kamus/C.png", "category": "ALPHABET"},
        {"id": 14, "word_text": "D", "definition": "Huruf keempat dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/D.webm", "image_url_ref": "kamus/D.png", "category": "ALPHABET"},
        {"id": 15, "word_text": "E", "definition": "Huruf kelima dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/E.webm", "image_url_ref": "kamus/E.png", "category": "ALPHABET"},
        {"id": 16, "word_text": "F", "definition": "Huruf keenam dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/F.webm", "image_url_ref": "kamus/F.png", "category": "ALPHABET"},
        {"id": 17, "word_text": "G", "definition": "Huruf ketujuh dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/G.webm", "image_url_ref": "kamus/G.png", "category": "ALPHABET"},
        {"id": 18, "word_text": "H", "definition": "Huruf kedelapan dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/H.webm", "image_url_ref": "kamus/H.png", "category": "ALPHABET"},
        {"id": 19, "word_text": "I", "definition": "Huruf kesembilan dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/I.webm", "image_url_ref": "kamus/I.png", "category": "ALPHABET"},
        {"id": 20, "word_text": "J", "definition": "Huruf kesepuluh dalam alfabet", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/J.webm", "image_url_ref": "kamus/J.png", "category": "ALPHABET"},
        {"id": 21, "word_text": "K", "definition": "Huruf kesebelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/K.webm", "image_url_ref": "kamus/K.png", "category": "ALPHABET"},
        {"id": 22, "word_text": "L", "definition": "Huruf keduabelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/L.webm", "image_url_ref": "kamus/L.png", "category": "ALPHABET"},
        {"id": 23, "word_text": "M", "definition": "Huruf ketigabelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/M.webm", "image_url_ref": "kamus/M.png", "category": "ALPHABET"},
        {"id": 24, "word_text": "N", "definition": "Huruf keempatbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/N.webm", "image_url_ref": "kamus/N.png", "category": "ALPHABET"},
        {"id": 25, "word_text": "O", "definition": "Huruf kelimabelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/O.webm", "image_url_ref": "kamus/O.png", "category": "ALPHABET"},
        {"id": 26, "word_text": "P", "definition": "Huruf keenambelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm", "image_url_ref": "kamus/P.png", "category": "ALPHABET"},
        {"id": 27, "word_text": "Q", "definition": "Huruf ketujuhbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Q.webm", "image_url_ref": "kamus/Q.png", "category": "ALPHABET"},
        {"id": 28, "word_text": "R", "definition": "Huruf kedelapanbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/R.webm", "image_url_ref": "kamus/R.png", "category": "ALPHABET"},
        {"id": 29, "word_text": "S", "definition": "Huruf kesembilanbelas dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/S.webm", "image_url_ref": "kamus/S.png", "category": "ALPHABET"},
        {"id": 30, "word_text": "T", "definition": "Huruf keduapuluh dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/T.webm", "image_url_ref": "kamus/T.png", "category": "ALPHABET"},
        {"id": 31, "word_text": "U", "definition": "Huruf keduapuluh satu dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/U.webm", "image_url_ref": "kamus/U.png", "category": "ALPHABET"},
        {"id": 32, "word_text": "V", "definition": "Huruf keduapuluh dua dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/V.webm", "image_url_ref": "kamus/V.png", "category": "ALPHABET"},
        {"id": 33, "word_text": "W", "definition": "Huruf keduapuluh tiga dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/W.webm", "image_url_ref": "kamus/W.png", "category": "ALPHABET"},
        {"id": 34, "word_text": "X", "definition": "Huruf keduapuluh empat dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/X.webm", "image_url_ref": "kamus/X.png", "category": "ALPHABET"},
        {"id": 35, "word_text": "Y", "definition": "Huruf keduapuluh lima dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm", "image_url_ref": "kamus/Y.png", "category": "ALPHABET"},
        {"id": 36, "word_text": "Z", "definition": "Huruf keduapuluh enam dalam alfabet ", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Z.webm", "image_url_ref": "kamus/Z.png", "category": "ALPHABET"},

        #level2
        # hewan
        {"id": 37, "word_text": "Anjing", "definition": "Hewan peliharaan yang setia", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm", "category": "KOSAKATA"},
        {"id": 38, "word_text": "Bebek", "definition": "Hewan ternak yang suka berenang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", "category": "KOSAKATA"},
        {"id": 39, "word_text": "Kambing", "definition": "Hewan ternak penghasil susu dan daging", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", "category": "KOSAKATA"},
        {"id": 40, "word_text": "Gajah", "definition": "Hewan besar dengan belalai panjang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", "category": "KOSAKATA"},
        {"id": 41, "word_text": "Monyet", "definition": "Hewan cerdas yang suka memanjat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm", "category": "KOSAKATA"},
        {"id": 42, "word_text": "Singa", "definition": "Hewan buas disebut raja hutan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm", "category": "KOSAKATA"},
        {"id": 43, "word_text": "Ular", "definition": "Hewan melata tanpa kaki", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ular.webm", "category": "KOSAKATA"},
        {"id": 44, "word_text": "Semut", "definition": "Hewan kecil yang suka bekerja sama", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Semut.webm", "category": "KOSAKATA"},
        {"id": 45, "word_text": "Kupu-kupu", "definition": "Hewan kecil bersayap indah", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/KupuKupu.webm", "category": "KOSAKATA"},
        {"id": 46, "word_text": "Lebah", "definition": "Hewan penghasil madu", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm", "category": "KOSAKATA"},
        {"id": 47, "word_text": "Kelinci", "definition": "Hewan dengan telinga panjang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm", "category": "KOSAKATA"},
        {"id": 48, "word_text": "Ikan", "definition": "Hewan yang hidup di air", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "category": "KOSAKATA"},
        {"id": 49, "word_text": "Burung", "definition": "Hewan yang bisa terbang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "category": "KOSAKATA"},
        {"id": 50, "word_text": "Kucing", "definition": "Hewan peliharaan yang suka mengeong", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "category": "KOSAKATA"},
        
        #keluarga
        {"id": 51, "word_text": "Ayah", "definition": "Orang tua laki-laki", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", "category": "KOSAKATA"},
        {"id": 52, "word_text": "Kakak", "definition": "Saudara yang lebih tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm", "category": "KOSAKATA"},
        {"id": 53, "word_text": "Kakek", "definition": "Ayah dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm", "category": "KOSAKATA"},
        {"id": 54, "word_text": "Nenek", "definition": "Ibu dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm", "category": "KOSAKATA"},
        {"id": 55, "word_text": "Paman", "definition": "Saudara laki-laki dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm", "category": "KOSAKATA"},
        {"id": 56, "word_text": "Bibi", "definition": "Saudara perempuan dari orang tua", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bibi.webm", "category": "KOSAKATA"},
        {"id": 57, "word_text": "Teman", "definition": "Orang yang akrab dengan kita", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", "category": "KOSAKATA"},
        {"id": 58, "word_text": "Guru", "definition": "Orang yang mengajar di sekolah", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "category": "KOSAKATA"},
        {"id": 59, "word_text": "Sekolah", "definition": "Tempat Ajar anak-anak", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sekolah.webm", "category": "KOSAKATA"},
        {"id": 60, "word_text": "Kura-kura", "definition": "Hewan bercangkang keras", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kura-kura.webm", "category": "KOSAKATA"},
        {"id": 61, "word_text": "Katak", "definition": "Hewan kecil yang pandai melompat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Katak.webm", "category": "KOSAKATA"},
        
        #kosakata level 3
        #geometri
        {"id": 62, "word_text": "Lingkaran", "definition": "Bentuk bulat tanpa sudut", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lingkaran.webm", "category": "KOSAKATA"},
        {"id": 63, "word_text": "Segitiga", "definition": "Bentuk dengan tiga sisi dan tiga sudut", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Segitiga.webm", "category": "KOSAKATA"},
        {"id": 64, "word_text": "Persegi", "definition": "Bentuk dengan empat sisi sama panjang dan empat sudut siku-siku", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Persegi.webm", "category": "KOSAKATA"},
        {"id": 65, "word_text": "Garis", "definition": "Bentuk lurus tanpa lebar dan tebal", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Garis.webm", "category": "KOSAKATA"},
        {"id": 66, "word_text": "Kotak", "definition": "Bentuk tiga dimensi dengan enam sisi persegi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kotak.webm", "category": "KOSAKATA"},
        {"id": 67, "word_text": "Bola", "definition": "Bentuk bulat tiga dimensi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bola.webm", "category": "KOSAKATA"},
        #waktu
        {"id": 68, "word_text": "Hari", "definition": "Satuan waktu selama 24 jam", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hari.webm", "category": "KOSAKATA"},
        {"id": 69, "word_text": "Minggu", "definition": "Satuan waktu selama 7 hari", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minggu.webm", "category": "KOSAKATA"},
        {"id": 70, "word_text": "Bulan", "definition": "Satuan waktu selama sekitar 30 hari", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bulan.webm", "category": "KOSAKATA"},
        {"id": 71, "word_text": "Tahun", "definition": "Satuan waktu selama 12 bulan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tahun.webm", "category": "KOSAKATA"},
        {"id": 72, "word_text": "Pagi", "definition": "Waktu setelah matahari terbit hingga siang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pagi.webm", "category": "KOSAKATA"},
        {"id": 73, "word_text": "Siang", "definition": "Waktu setelah pagi hingga sore", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Siang.webm", "category": "KOSAKATA"},
        {"id": 74, "word_text": "Sore", "definition": "Waktu setelah siang hingga matahari terbenam", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sore.webm", "category": "KOSAKATA"},
        {"id": 75, "word_text": "Malam", "definition": "Waktu setelah matahari terbenam hingga pagi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Malam.webm", "category": "KOSAKATA"},
        {"id": 76, "word_text": "Detik", "definition": "Satuan waktu terkecil", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Detik.webm", "category": "KOSAKATA"},
        {"id": 77, "word_text": "Menit", "definition": "Satuan waktu selama 60 detik", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Menit.webm", "category": "KOSAKATA"},
        {"id": 78, "word_text": "Jam", "definition": "Satuan waktu selama 60 menit", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jam.webm", "category": "KOSAKATA"},

        #aktivitas
        # Level 4
        {"id": 79, "word_text": "Bangun", "definition": "Kegiatan setelah tidur", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", "category": "KOSAKATA"},
        {"id": 80, "word_text": "Mandi", "definition": "Membersihkan badan dengan air", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "category": "KOSAKATA"},
        {"id": 81, "word_text": "Sarapan", "definition": "Makan pagi sebelum beraktivitas", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarapan.webm", "category": "KOSAKATA"},
        {"id": 82, "word_text": "Ajar", "definition": "Kegiatan untuk mendapatkan ilmu", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ajar.webm", "category": "KOSAKATA"},
        {"id": 83, "word_text": "Tulis", "definition": "Membuat huruf atau kata dengan alat tulis", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm", "category": "KOSAKATA"},
        {"id": 84, "word_text": "Baca", "definition": "Melihat dan memahami tulisan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm", "category": "KOSAKATA"},
        {"id": 85, "word_text": "Main", "definition": "Melakukan kegiatan yang menyenangkan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "category": "KOSAKATA"},
        {"id": 86, "word_text": "Tidur", "definition": "Beristirahat dengan memejamkan mata", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "category": "KOSAKATA"},
        {"id": 87, "word_text": "Duduk", "definition": "Posisi badan di atas kursi", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", "category": "KOSAKATA"},
        {"id": 88, "word_text": "Nonton", "definition": "Melihat acara di televisi atau layar", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm", "category": "KOSAKATA"},
        {"id": 89, "word_text": "Lari", "definition": "Bergerak cepat dengan kaki", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm", "category": "KOSAKATA"},
        {"id": 90, "word_text": "Lompat", "definition": "Berpindah dengan melompatkan tubuh", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm", "category": "KOSAKATA"},
        {"id": 91, "word_text": "Nyanyi", "definition": "Mengeluarkan suara dengan nada lagu", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nyanyi.webm", "category": "KOSAKATA"},
        {"id": 92, "word_text": "Cuci", "definition": "Membersihkan dengan air dan sabun", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Cuci.webm", "category": "KOSAKATA"},
        {"id": 93, "word_text": "Tangan", "definition": "Anggota badan untuk memegang", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tangan.webm", "category": "KOSAKATA"},
        {"id": 94, "word_text": "Sapu", "definition": "Alat untuk membersihkan lantai", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm", "category": "KOSAKATA"},
        {"id": 95, "word_text": "Sikat", "definition": "Alat untuk menggosok dan membersihkan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sikat.webm", "category": "KOSAKATA"},
        {"id": 96, "word_text": "Gigi", "definition": "Bagian mulut untuk mengunyah makanan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gigi.webm", "category": "KOSAKATA"},
        {"id": 97, "word_text": "Lapar", "definition": "Keadaan ingin makan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lapar.webm", "category": "KOSAKATA"},
        {"id": 98, "word_text": "Senang", "definition": "Perasaan gembira", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "category": "KOSAKATA"},
        {"id": 99, "word_text": "Sedih", "definition": "Perasaan tidak bahagia", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm", "category": "KOSAKATA"},
        {"id": 100, "word_text": "Marah", "definition": "Perasaan kesal yang kuat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Marah.webm", "category": "KOSAKATA"},
        {"id": 101, "word_text": "Takut", "definition": "Perasaan tidak berani atau khawatir", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Takut.webm", "category": "KOSAKATA"},
        {"id": 102, "word_text": "Jalan", "definition": "Bergerak dengan kaki perlahan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jalan.webm", "category": "KOSAKATA"},
        {"id": 103, "word_text": "Hujan", "definition": "Air yang jatuh dari langit", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hujan.webm", "category": "KOSAKATA"},
        {"id": 104, "word_text": "Taman", "definition": "Tempat dengan tanaman dan bunga", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Taman.webm", "category": "KOSAKATA"},
        {"id": 105, "word_text": "Mimpi", "definition": "Gambaran dalam tidur", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mimpi.webm", "category": "KOSAKATA"},

        # Imbuhan Awalan
        {"id": 106, "word_text": "Ber-", "definition": "Imbuhan awalan untuk kata kerja", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Ber.webm", "category": "IMBUHAN"},
        {"id": 107, "word_text": "Ter-", "definition": "Imbuhan awalan untuk kata sifat", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Ter.webm", "category": "IMBUHAN"},
        {"id": 108, "word_text": "Me-", "definition": "Imbuhan awalan untuk kata kerja aktif", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Me.webm", "category": "IMBUHAN"},
        {"id": 109, "word_text": "Di-", "definition": "Imbuhan awalan untuk kata kerja pasif", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Awalan-Di.webm", "category": "IMBUHAN"},
        # Imbuhan Akhiran & Partikel
        {"id": 110, "word_text": "-kan", "definition": "Imbuhan akhiran untuk membentuk kata kerja", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Kan.webm", "category": "IMBUHAN"},
        {"id": 111, "word_text": "-i", "definition": "Imbuhan akhiran untuk membentuk kata kerja", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-I.webm", "category": "IMBUHAN"},
        {"id": 112, "word_text": "-an", "definition": "Imbuhan akhiran untuk membentuk kata benda", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-An.webm", "category": "IMBUHAN"},
        {"id": 113, "word_text": "-wan", "definition": "Imbuhan akhiran untuk membentuk kata benda pelaku", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Wan.webm", "category": "IMBUHAN"},
        {"id": 114, "word_text": "-wati", "definition": "Imbuhan akhiran untuk membentuk kata benda pelaku perempuan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Wati.webm", "category": "IMBUHAN"},
        {"id": 115, "word_text": "-man", "definition": "Imbuhan akhiran untuk membentuk kata benda pelaku", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Man.webm", "category": "IMBUHAN"},
        {"id": 116, "word_text": "-ti", "definition": "Imbuhan akhiran untuk membentuk kata benda", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Ti.webm", "category": "IMBUHAN"},
        {"id": 117, "word_text": "-nya", "definition": "Imbuhan akhiran untuk kepemilikan", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Akhiran-Nya.webm", "category": "IMBUHAN"},
        # Imbuhan Partikel
        {"id": 118, "word_text": "-pun", "definition": "Partikel penegas", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Pun.webm", "category": "IMBUHAN"},
        {"id": 119, "word_text": "-lah", "definition": "Partikel penegas", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Lah.webm", "category": "IMBUHAN"},
        {"id": 120, "word_text": "-kah", "definition": "Partikel penanya", "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/imbuhan/Partikel-Kah.webm", "category": "IMBUHAN"},

    ]

class SoalSeeder(BaseSeeder):
    """Seed Soal dengan assignment foreign keys otomatis"""

    def __init__(self):
        super().__init__()

    def run(self):
        """Jalankan proses seeding soal"""
        try:
            print("Starting Soal seeding with auto foreign key assignment...")
            
            self.seed_soal()
            
            print("Soal seeding finished successfully!")
            
        except Exception as e:
            self.db.rollback()
            raise Exception(f"Soal seeding failed: {e}")

    def seed_soal(self):
        print("Seeding Soal...")
        
        soal_data = []
        
        # Level 1: Alphabet A-Z (Sublevel 1.1: A-C)
        sublevel_id = 1
        soal_data.extend([
            {"question": "Tunjukkan isyarat huruf 'A'", "answer": "Kepalkan tangan dengan ibu jari di samping", "dictionary_id": 11, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini adalah isyarat huruf?", "answer": "A", "dictionary_id": 11, "sublevel_id": sublevel_id, "image_url": "kamus/A.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk huruf 'A' adalah...", "answer": "kamus/A.png", "dictionary_id": 11, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Bagaimana cara membuat isyarat huruf 'B'?", "answer": "Tangan terbuka dengan jari rapat, ibu jari menekuk", "dictionary_id": 12, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/B.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Isyarat ini adalah huruf?", "answer": "B", "dictionary_id": 12, "sublevel_id": sublevel_id, "image_url": "kamus/B.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Huruf 'B' ditunjukkan dengan isyarat...", "answer": "kamus/B.png", "dictionary_id": 12, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat huruf 'A' kembali", "answer": "Kepalkan tangan dengan ibu jari di samping", "dictionary_id": 11, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Tunjukkan isyarat huruf 'B' sekali lagi", "answer": "Tangan terbuka dengan jari rapat, ibu jari menekuk", "dictionary_id": 12, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/B.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang terlihat pada gambar?", "answer": "A", "dictionary_id": 11, "sublevel_id": sublevel_id, "image_url": "kamus/A.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Gambar berikut menunjukkan huruf?", "answer": "B", "dictionary_id": 12, "sublevel_id": sublevel_id, "image_url": "kamus/B.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
        ])

        sublevel_id = 2
        # Level 1: Alphabet A-Z (Sublevel 1.2: D-F) - 10 Soal
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk huruf 'D'", "answer": "Jari telunjuk tegak, jari lain menyentuh ibu jari", "dictionary_id": 14, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/D.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini adalah isyarat untuk huruf?", "answer": "D", "dictionary_id": 14, "sublevel_id": sublevel_id, "image_url": "kamus/D.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'D'", "answer": "kamus/D.png", "dictionary_id": 14, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'E' dengan benar", "answer": "Jari ditekuk ke dalam menyentuh ibu jari yang tertekuk", "dictionary_id": 15, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/E.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "E", "dictionary_id": 15, "sublevel_id": sublevel_id, "image_url": "kamus/E.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Lakukan isyarat huruf 'F'", "answer": "Jari telunjuk dan ibu jari bertemu, jari lain tegak", "dictionary_id": 16, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/F.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf dari gambar isyarat ini", "answer": "F", "dictionary_id": 16, "sublevel_id": sublevel_id, "image_url": "kamus/F.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat yang mewakili huruf 'F'?", "answer": "kamus/F.png", "dictionary_id": 16, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Coba tunjukkan isyarat huruf 'E' sekali lagi", "answer": "Jari ditekuk ke dalam menyentuh ibu jari yang tertekuk", "dictionary_id": 15, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/E.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Pilih gambar yang menunjukkan isyarat 'D'", "answer": "kamus/D.png", "dictionary_id": 14, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
        ])

        sublevel_id = 3
        # Level 1: Alphabet A-Z (Sublevel 1.3: G-I) - 10 Soal
        soal_data.extend([
            {"question": "Isyaratkan huruf 'G'", "answer": "Jari telunjuk dan ibu jari sejajar lurus", "dictionary_id": 17, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/G.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini adalah isyarat untuk huruf?", "answer": "G", "dictionary_id": 17, "sublevel_id": sublevel_id, "image_url": "kamus/G.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'G'", "answer": "kamus/G.png", "dictionary_id": 17, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat huruf 'H'", "answer": "Jari telunjuk dan tengah lurus, ibu jari di antara", "dictionary_id": 18, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/H.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "H", "dictionary_id": 18, "sublevel_id": sublevel_id, "image_url": "kamus/H.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'H'", "answer": "kamus/H.png", "dictionary_id": 18, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat huruf 'I'", "answer": "Jari kelingking tegak", "dictionary_id": 19, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/I.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf dari gambar isyarat ini", "answer": "I", "dictionary_id": 19, "sublevel_id": sublevel_id, "image_url": "kamus/I.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Praktikkan isyarat huruf 'I' kembali", "answer": "Jari kelingking tegak", "dictionary_id": 19, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/I.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Pilih isyarat yang tepat untuk huruf 'H'", "answer": "kamus/H.png", "dictionary_id": 18, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
        ])

        sublevel_id = 4
        # Level 1: Alphabet A-Z (Sublevel 1.4: J-L) - 10 Soal
        soal_data.extend([
            {"question": "Lakukan isyarat huruf 'K'", "answer": "Jari kelingking diayun ke bawah membentuk 'K'", "dictionary_id": 21, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/K.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar isyarat ini adalah huruf?", "answer": "J", "dictionary_id": 20, "sublevel_id": sublevel_id, "image_url": "kamus/J.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'J'", "answer": "kamus/J.png", "dictionary_id": 20, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat huruf 'K'", "answer": "Jari telunjuk dan tengah tegak, ibu jari di antara", "dictionary_id": 21, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/K.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "K", "dictionary_id": 21, "sublevel_id": sublevel_id, "image_url": "kamus/K.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'K'", "answer": "kamus/K.png", "dictionary_id": 21, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'L'", "answer": "Jari telunjuk tegak dan ibu jari lurus, membentuk L", "dictionary_id": 22, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/L.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf dari gambar isyarat ini", "answer": "L", "dictionary_id": 22, "sublevel_id": sublevel_id, "image_url": "kamus/L.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat yang mewakili huruf 'L'?", "answer": "kamus/L.png", "dictionary_id": 22, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Coba tunjukkan isyarat huruf 'L' sekali lagi", "answer": "Jari telunjuk tegak dan ibu jari lurus, membentuk L", "dictionary_id": 22, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/L.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        sublevel_id = 5
        # Level 1: Alphabet A-Z (Sublevel 1.5: M-O) - 10 Soal
        soal_data.extend([
            {"question": "Tunjukkan isyarat huruf 'M'", "answer": "Tiga jari (telunjuk, tengah, manis) di atas ibu jari", "dictionary_id": 23, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/M.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar isyarat ini adalah huruf?", "answer": "M", "dictionary_id": 23, "sublevel_id": sublevel_id, "image_url": "kamus/M.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'M'", "answer": "kamus/M.png", "dictionary_id": 23, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'N' dengan benar", "answer": "Dua jari (telunjuk, tengah) di atas ibu jari", "dictionary_id": 24, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/N.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "N", "dictionary_id": 24, "sublevel_id": sublevel_id, "image_url": "kamus/N.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'N'", "answer": "kamus/N.png", "dictionary_id": 24, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat huruf 'O'", "answer": "Jari membentuk lingkaran", "dictionary_id": 25, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/O.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf dari gambar isyarat ini", "answer": "O", "dictionary_id": 25, "sublevel_id": sublevel_id, "image_url": "kamus/O.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat yang mewakili huruf 'O'?", "answer": "kamus/O.png", "dictionary_id": 25, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat huruf 'N' kembali", "answer": "Dua jari (telunjuk, tengah) di atas ibu jari", "dictionary_id": 24, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/N.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        sublevel_id = 6
        # Level 1: Alphabet A-Z (Sublevel 1.6: P-R) - 10 Soal
        soal_data.extend([
            {"question": "Tunjukkan isyarat huruf 'P'", "answer": "Sama dengan 'K', tapi menghadap ke bawah/depan", "dictionary_id": 26, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar isyarat ini adalah huruf?", "answer": "P", "dictionary_id": 26, "sublevel_id": sublevel_id, "image_url": "kamus/P.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'P'", "answer": "kamus/P.png", "dictionary_id": 26, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'Q' dengan benar", "answer": "Jari telunjuk dan ibu jari menunjuk ke bawah", "dictionary_id": 27, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Q.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "Q", "dictionary_id": 27, "sublevel_id": sublevel_id, "image_url": "kamus/Q.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'Q'", "answer": "kamus/Q.png", "dictionary_id": 27, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat huruf 'R'", "answer": "Jari tengah menyilang di atas jari telunjuk", "dictionary_id": 28, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/R.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf dari gambar isyarat ini", "answer": "R", "dictionary_id": 28, "sublevel_id": sublevel_id, "image_url": "kamus/R.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat yang mewakili huruf 'R'?", "answer": "kamus/R.png", "dictionary_id": 28, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Coba tunjukkan isyarat huruf 'P' sekali lagi", "answer": "Sama dengan 'K', tapi menghadap ke bawah/depan", "dictionary_id": 26, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        sublevel_id = 7
        # Level 1: Alphabet A-Z (Sublevel 1.7: S-U) - 10 Soal
        soal_data.extend([
            {"question": "Tunjukkan isyarat huruf 'S'", "answer": "Kepalan tangan, ibu jari di depan jari lain", "dictionary_id": 29, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/S.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar isyarat ini adalah huruf?", "answer": "S", "dictionary_id": 29, "sublevel_id": sublevel_id, "image_url": "kamus/S.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'S'", "answer": "kamus/S.png", "dictionary_id": 29, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'T' dengan benar", "answer": "Ibu jari masuk ke antara jari telunjuk dan tengah", "dictionary_id": 30, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/T.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "T", "dictionary_id": 30, "sublevel_id": sublevel_id, "image_url": "kamus/T.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'T'", "answer": "kamus/T.png", "dictionary_id": 30, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat huruf 'U'", "answer": "Jari telunjuk dan tengah tegak dan rapat", "dictionary_id": 31, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/U.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf dari gambar isyarat ini", "answer": "U", "dictionary_id": 31, "sublevel_id": sublevel_id, "image_url": "kamus/U.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat yang mewakili huruf 'U'?", "answer": "kamus/U.png", "dictionary_id": 31, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Coba tunjukkan isyarat huruf 'S' kembali", "answer": "Kepalan tangan, ibu jari di depan jari lain", "dictionary_id": 29, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/S.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        sublevel_id = 8
        # Level 1: Alphabet A-Z (Sublevel 1.8: V-X) - 10 Soal
        soal_data.extend([
            {"question": "Tunjukkan isyarat huruf 'V'", "answer": "Jari telunjuk dan tengah tegak membentuk V", "dictionary_id": 32, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/V.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar isyarat ini adalah huruf?", "answer": "V", "dictionary_id": 32, "sublevel_id": sublevel_id, "image_url": "kamus/V.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'V'", "answer": "kamus/V.png", "dictionary_id": 32, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'W' dengan benar", "answer": "Tiga jari (telunjuk, tengah, manis) tegak", "dictionary_id": 33, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/W.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "W", "dictionary_id": 33, "sublevel_id": sublevel_id, "image_url": "kamus/W.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'W'", "answer": "kamus/W.png", "dictionary_id": 33, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat huruf 'X'", "answer": "Jari telunjuk ditekuk seperti kait", "dictionary_id": 34, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/X.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf dari gambar isyarat ini", "answer": "X", "dictionary_id": 34, "sublevel_id": sublevel_id, "image_url": "kamus/X.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat yang mewakili huruf 'X'?", "answer": "kamus/X.png", "dictionary_id": 34, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Coba tunjukkan isyarat huruf 'V' kembali", "answer": "Jari telunjuk dan tengah tegak membentuk V", "dictionary_id": 32, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/V.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        sublevel_id = 9
        # Level 1: Alphabet A-Z (Sublevel 1.9: Y-Z) - 10 Soal
        soal_data.extend([
            {"question": "Tunjukkan isyarat huruf 'Y'", "answer": "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", "dictionary_id": 35, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar isyarat ini adalah huruf?", "answer": "Y", "dictionary_id": 35, "sublevel_id": sublevel_id, "image_url": "kamus/Y.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'Y'", "answer": "kamus/Y.png", "dictionary_id": 35, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'Y' dengan benar", "answer": "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", "dictionary_id": 35, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang ditunjukkan oleh isyarat pada gambar?", "answer": "Z", "dictionary_id": 36, "sublevel_id": sublevel_id, "image_url": "kamus/Z.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'Z'", "answer": "kamus/Z.png", "dictionary_id": 36, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat huruf 'Y' sekali lagi", "answer": "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", "dictionary_id": 35, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Identifikasi huruf pada gambar isyarat ini", "answer": "Y", "dictionary_id": 35, "sublevel_id": sublevel_id, "image_url": "kamus/Y.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang tepat untuk 'Z'", "answer": "kamus/Z.png", "dictionary_id": 36, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Coba tunjukkan isyarat huruf 'Y' di depan kamera", "answer": "Jari kelingking dan ibu jari terbuka, jari lain ditekuk", "dictionary_id": 35, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Y.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        sublevel_id = 10
        # Level 1: Alphabet A-Z (Sublevel 1.10: Review A-Z) - 10 Soal
        soal_data.extend([
            {"question": "Tunjukkan isyarat huruf 'A' (Review)", "answer": "Kepalkan tangan dengan ibu jari di samping", "dictionary_id": 11, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/A.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini adalah isyarat huruf?", "answer": "F", "dictionary_id": 16, "sublevel_id": sublevel_id, "image_url": "kamus/F.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'K'", "answer": "kamus/K.png", "dictionary_id": 21, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyaratkan huruf 'P'", "answer": "Sama dengan 'K', tapi menghadap ke bawah/depan", "dictionary_id": 26, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/P.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Huruf apa yang terlihat pada gambar?", "answer": "U", "dictionary_id": 31, "sublevel_id": sublevel_id, "image_url": "kamus/U.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tentukan isyarat yang tepat untuk huruf 'Z'", "answer": "kamus/Z.png", "dictionary_id": 36, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat huruf 'D'", "answer": "Jari telunjuk tegak, jari lain menyentuh ibu jari", "dictionary_id": 14, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/D.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar berikut menunjukkan huruf?", "answer": "M", "dictionary_id": 23, "sublevel_id": sublevel_id, "image_url": "kamus/M.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat yang benar untuk huruf 'R'", "answer": "kamus/R.png", "dictionary_id": 28, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Coba tunjukkan isyarat huruf 'W'", "answer": "Tiga jari (telunjuk, tengah, manis) tegak", "dictionary_id": 33, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/W.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.1: hewan rumah)
        sublevel_id = 11
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Kucing'", "answer": "kucing", "dictionary_id": 50, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan hewan?", "answer": "Kucing", "dictionary_id": 50, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Kucing' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "dictionary_id": 50, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Burung'", "answer": "burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan pada gambar ini adalah?", "answer": "Anjing", "dictionary_id": 37, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Anjing'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm", "dictionary_id": 37, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Burung'", "answer": "burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan ini adalah?", "answer": "Burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Burung' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "dictionary_id": 49, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Ikan'", "answer": "ikan", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.2: hewan ternak)
        sublevel_id = 12
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Bebek'", "answer": "bebek", "dictionary_id": 38, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan hewan?", "answer": "Bebek", "dictionary_id": 38, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Bebek' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", "dictionary_id": 38, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Kambing'", "answer": "kambing", "dictionary_id": 39, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan pada gambar ini adalah?", "answer": "Kambing", "dictionary_id": 39, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Kambing'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", "dictionary_id": 39, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Bebek' kembali", "answer": "bebek", "dictionary_id": 38, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan apa ini?", "answer": "Bebek", "dictionary_id": 38, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bebek.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Kambing'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", "dictionary_id": 39, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Kambing'", "answer": "kambing", "dictionary_id": 39, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kambing.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.3: hewan liar)
        sublevel_id = 13
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Gajah'", "answer": "gajah", "dictionary_id": 40, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan hewan?", "answer": "Gajah", "dictionary_id": 40, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Gajah' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", "dictionary_id": 40, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Singa'", "answer": "singa", "dictionary_id": 42, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan pada gambar ini adalah?", "answer": "Singa", "dictionary_id": 42, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Singa'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Singa.webm", "dictionary_id": 42, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Monyet'", "answer": "monyet", "dictionary_id": 41, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan apa ini?", "answer": "Monyet", "dictionary_id": 41, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Gajah'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Gajah.webm", "dictionary_id": 40, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Monyet'", "answer": "monyet", "dictionary_id": 41, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Monyet.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.4: hewan kecil)
        sublevel_id = 14
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Kelinci'", "answer": "kelinci", "dictionary_id": 47, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan hewan?", "answer": "Kelinci", "dictionary_id": 47, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Kelinci' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm", "dictionary_id": 47, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Semut'", "answer": "semut", "dictionary_id": 44, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Semut.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan pada gambar ini adalah?", "answer": "Semut", "dictionary_id": 44, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Semut.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Kupu-kupu'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/KupuKupu.webm", "dictionary_id": 45, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Lebah'", "answer": "lebah", "dictionary_id": 46, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan apa ini?", "answer": "Lebah", "dictionary_id": 46, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Lebah'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lebah.webm", "dictionary_id": 46, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Kelinci'", "answer": "kelinci", "dictionary_id": 47, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kelinci.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.5: keluarga inti)
        sublevel_id = 15
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Ayah'", "answer": "ayah", "dictionary_id": 51, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan?", "answer": "Ayah", "dictionary_id": 51, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Ayah' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", "dictionary_id": 51, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Ibu'", "answer": "ibu", "dictionary_id": 52, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Siapa pada gambar ini?", "answer": "Ibu", "dictionary_id": 52, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Ibu'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm", "dictionary_id": 52, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Kakak'", "answer": "kakak", "dictionary_id": 52, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Siapa ini?", "answer": "Kakak", "dictionary_id": 52, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Kakak' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakak.webm", "dictionary_id": 52, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat untuk 'Ayah'", "answer": "ayah", "dictionary_id": 51, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.6: keluarga besar)
        sublevel_id = 16
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Kakek'", "answer": "kakek", "dictionary_id": 53, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan?", "answer": "Kakek", "dictionary_id": 53, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Kakek' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kakek.webm", "dictionary_id": 53, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Nenek'", "answer": "nenek", "dictionary_id": 54, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Siapa pada gambar ini?", "answer": "Nenek", "dictionary_id": 54, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Nenek'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nenek.webm", "dictionary_id": 54, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Paman'", "answer": "paman", "dictionary_id": 55, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Siapa ini?", "answer": "Paman", "dictionary_id": 55, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Paman' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Paman.webm", "dictionary_id": 55, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Bibi'", "answer": "bibi", "dictionary_id": 56, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bibi.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.7: teman dan guru)
        sublevel_id = 17
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Teman'", "answer": "teman", "dictionary_id": 57, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan?", "answer": "Teman", "dictionary_id": 57, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Teman' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", "dictionary_id": 57, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Guru'", "answer": "guru", "dictionary_id": 58, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Siapa pada gambar ini?", "answer": "Guru", "dictionary_id": 58, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Guru'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "dictionary_id": 58, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Teman' kembali", "answer": "teman", "dictionary_id": 57, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Siapa ini?", "answer": "Teman", "dictionary_id": 57, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Teman.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Guru'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "dictionary_id": 58, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Guru'", "answer": "guru", "dictionary_id": 58, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.8: hewan air)
        sublevel_id = 18
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Ikan'", "answer": "ikan", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan hewan?", "answer": "Ikan", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Ikan' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "dictionary_id": 48, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Hewan air lainnya adalah ikan. Tunjukkan isyarat 'Ikan'", "answer": "ikan", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan pada gambar ini hidup di?", "answer": "Air", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk hewan air 'Ikan'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "dictionary_id": 48, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Ikan' sekali lagi", "answer": "ikan", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan air apa ini?", "answer": "Ikan", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Ikan'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "dictionary_id": 48, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Ikan'", "answer": "ikan", "dictionary_id": 48, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ikan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.9: hewan udara)
        sublevel_id = 19
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Burung'", "answer": "burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan hewan?", "answer": "Burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Burung' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "dictionary_id": 49, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Hewan yang terbang adalah burung. Tunjukkan isyarat 'Burung'", "answer": "burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan pada gambar ini bisa?", "answer": "Terbang", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk hewan udara 'Burung'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "dictionary_id": 49, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Burung' sekali lagi", "answer": "burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan udara apa ini?", "answer": "Burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Burung'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "dictionary_id": 49, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Burung'", "answer": "burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 2: Basic Words (Sublevel 2.10: latihan campuran)
        sublevel_id = 20
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Kucing'", "answer": "kucing", "dictionary_id": 50, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan?", "answer": "Anjing", "dictionary_id": 37, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Anjing.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Ayah' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ayah.webm", "dictionary_id": 51, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Ibu'", "answer": "ibu", "dictionary_id": 52, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Siapa pada gambar ini?", "answer": "Guru", "dictionary_id": 58, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Kucing'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "dictionary_id": 50, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Guru'", "answer": "guru", "dictionary_id": 58, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Guru.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Hewan apa ini?", "answer": "Kucing", "dictionary_id": 50, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Kucing.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Ibu'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ibu.webm", "dictionary_id": 52, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Burung'", "answer": "burung", "dictionary_id": 49, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Burung.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.1: Numbers 0-2)
        sublevel_id = 21
        soal_data.extend([
            {"question": "Tunjukkan isyarat angka '0'", "answer": "0", "dictionary_id": 1, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/0.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Angka berapa ini?", "answer": "0", "dictionary_id": 1, "sublevel_id": sublevel_id, "image_url": "kamus/0.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat angka '0' adalah...", "answer": "kamus/0.png", "dictionary_id": 1, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "1 - 1 = ?", "answer": "0", "dictionary_id": 1, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Praktikkan isyarat angka '1'", "answer": "1", "dictionary_id": 2, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/01.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan angka?", "answer": "1", "dictionary_id": 2, "sublevel_id": sublevel_id, "image_url": "kamus/1.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat angka '1'", "answer": "kamus/1.png", "dictionary_id": 2, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "0 + 1 = ?", "answer": "1", "dictionary_id": 2, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Tunjukkan isyarat angka '2'", "answer": "2", "dictionary_id": 3, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/02.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Angka pada gambar?", "answer": "2", "dictionary_id": 3, "sublevel_id": sublevel_id, "image_url": "kamus/2.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
        ])

        # Level 3: Numbers and Math (Sublevel 3.2: Numbers 3-5)
        sublevel_id = 22
        soal_data.extend([
            {"question": "Tunjukkan isyarat angka '3'", "answer": "3", "dictionary_id": 4, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/03.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Angka berapa ini?", "answer": "3", "dictionary_id": 4, "sublevel_id": sublevel_id, "image_url": "kamus/3.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat angka '3' adalah...", "answer": "kamus/3.png", "dictionary_id": 4, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "2 + 1 = ?", "answer": "3", "dictionary_id": 4, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Praktikkan isyarat angka '4'", "answer": "4", "dictionary_id": 5, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/04.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan angka?", "answer": "4", "dictionary_id": 5, "sublevel_id": sublevel_id, "image_url": "kamus/4.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat angka '4'", "answer": "kamus/4.png", "dictionary_id": 5, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "3 + 1 = ?", "answer": "4", "dictionary_id": 5, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Tunjukkan isyarat angka '5'", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/05.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Angka pada gambar?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "image_url": "kamus/5.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
        ])

        # Level 3: Numbers and Math (Sublevel 3.3: Numbers 6-7)
        sublevel_id = 23
        soal_data.extend([
            {"question": "Tunjukkan isyarat angka '6'", "answer": "6", "dictionary_id": 7, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/06.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Angka berapa ini?", "answer": "6", "dictionary_id": 7, "sublevel_id": sublevel_id, "image_url": "kamus/6.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat angka '6' adalah...", "answer": "kamus/6.png", "dictionary_id": 7, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "5 + 1 = ?", "answer": "6", "dictionary_id": 7, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Praktikkan isyarat angka '7'", "answer": "7", "dictionary_id": 8, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/07.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan angka?", "answer": "7", "dictionary_id": 8, "sublevel_id": sublevel_id, "image_url": "kamus/7.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat angka '7'", "answer": "kamus/7.png", "dictionary_id": 8, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "6 + 1 = ?", "answer": "7", "dictionary_id": 8, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "7 - 1 = ?", "answer": "6", "dictionary_id": 7, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Tunjukkan kembali angka '6'", "answer": "6", "dictionary_id": 7, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/06.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.4: Numbers 8-9)
        sublevel_id = 24
        soal_data.extend([
            {"question": "Tunjukkan isyarat angka '8'", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/08.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Angka berapa ini?", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "image_url": "kamus/8.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat angka '8' adalah...", "answer": "kamus/8.png", "dictionary_id": 9, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "7 + 1 = ?", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Praktikkan isyarat angka '9'", "answer": "9", "dictionary_id": 10, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/09.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan angka?", "answer": "9", "dictionary_id": 10, "sublevel_id": sublevel_id, "image_url": "kamus/9.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat angka '9'", "answer": "kamus/9.png", "dictionary_id": 10, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "8 + 1 = ?", "answer": "9", "dictionary_id": 10, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "9 - 1 = ?", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Tunjukkan kembali angka '8'", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/08.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.5: Review 0-9)
        sublevel_id = 25
        soal_data.extend([
            {"question": "Tunjukkan isyarat angka '5'", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/05.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Angka berapa ini?", "answer": "3", "dictionary_id": 4, "sublevel_id": sublevel_id, "image_url": "kamus/3.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat angka '9' adalah...", "answer": "kamus/9.png", "dictionary_id": 10, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "4 + 3 = ?", "answer": "7", "dictionary_id": 8, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Praktikkan isyarat angka '2'", "answer": "2", "dictionary_id": 3, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/02.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan angka?", "answer": "6", "dictionary_id": 7, "sublevel_id": sublevel_id, "image_url": "kamus/6.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat angka '1'", "answer": "kamus/1.png", "dictionary_id": 2, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "9 - 5 = ?", "answer": "4", "dictionary_id": 5, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Tunjukkan isyarat angka '0'", "answer": "0", "dictionary_id": 1, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/0.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "8 - 0 = ?", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.6: Penjumlahan)
        sublevel_id = 26
        soal_data.extend([
            {"question": "1 + 1 = ?", "answer": "2", "dictionary_id": 3, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "2 + 3 = ?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Berapa hasil 3 + 2?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "4 + 1 = ?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "5 + 3 = ?", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Berapa 2 + 2?", "answer": "4", "dictionary_id": 5, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "1 + 2 = ?", "answer": "3", "dictionary_id": 4, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "0 + 5 = ?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "3 + 3 = ?", "answer": "6", "dictionary_id": 7, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "1 + 8 = ?", "answer": "9", "dictionary_id": 10, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.7: Pengurangan)
        sublevel_id = 27
        soal_data.extend([
            {"question": "5 - 2 = ?", "answer": "3", "dictionary_id": 4, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "8 - 3 = ?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Berapa hasil 9 - 4?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "7 - 2 = ?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "6 - 1 = ?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Berapa 4 - 2?", "answer": "2", "dictionary_id": 3, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "9 - 5 = ?", "answer": "4", "dictionary_id": 5, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "3 - 3 = ?", "answer": "0", "dictionary_id": 1, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "8 - 1 = ?", "answer": "7", "dictionary_id": 8, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "2 - 2 = ?", "answer": "0", "dictionary_id": 1, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.8: Geometri)
        sublevel_id = 28
        soal_data.extend([
            {"question": "Bentuk ini adalah?", "answer": "Lingkaran", "dictionary_id": 62, "sublevel_id": sublevel_id, "image_url": "kamus/Lingkaran.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Tunjukkan isyarat 'Lingkaran'", "answer": "Lingkaran", "dictionary_id": 62, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lingkaran.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Isyarat 'Segitiga' adalah...", "answer": "kamus/Segitiga.png", "dictionary_id": 63, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Persegi'", "answer": "Persegi", "dictionary_id": 64, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Persegi.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Bentuk apa ini?", "answer": "Segitiga", "dictionary_id": 63, "sublevel_id": sublevel_id, "image_url": "kamus/Segitiga.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Gambar ini adalah?", "answer": "Persegi", "dictionary_id": 64, "sublevel_id": sublevel_id, "image_url": "kamus/Persegi.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat 'Lingkaran'", "answer": "kamus/Lingkaran.png", "dictionary_id": 62, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Segitiga'", "answer": "Segitiga", "dictionary_id": 63, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Segitiga.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Bentuk dengan 3 sisi adalah?", "answer": "Segitiga", "dictionary_id": 63, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Bentuk dengan 4 sisi sama adalah?", "answer": "Persegi", "dictionary_id": 64, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.9: Waktu)
        sublevel_id = 29
        soal_data.extend([
            {"question": "Tunjukkan isyarat 'Hari'", "answer": "Hari", "dictionary_id": 68, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Hari.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Kapan matahari terbit?", "answer": "Pagi", "dictionary_id": 72, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Pagi'", "answer": "Pagi", "dictionary_id": 72, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pagi.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Waktu untuk makan siang adalah?", "answer": "Siang", "dictionary_id": 73, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Isyarat 'Sore' adalah...", "answer": "kamus/Sore.png", "dictionary_id": 74, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Malam'", "answer": "Malam", "dictionary_id": 75, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Malam.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "7 hari membentuk satu?", "answer": "Minggu", "dictionary_id": 69, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Minggu'", "answer": "Minggu", "dictionary_id": 69, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minggu.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Waktu tidur adalah?", "answer": "Malam", "dictionary_id": 75, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "12 bulan membentuk satu?", "answer": "Tahun", "dictionary_id": 71, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
        ])

        # Level 3: Numbers and Math (Sublevel 3.10: Latihan Campuran)
        sublevel_id = 30
        soal_data.extend([
            {"question": "5 + 3 = ?", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Bentuk ini adalah?", "answer": "Lingkaran", "dictionary_id": 62, "sublevel_id": sublevel_id, "image_url": "kamus/Lingkaran.png", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "9 - 4 = ?", "answer": "5", "dictionary_id": 6, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Tunjukkan isyarat 'Segitiga'", "answer": "Segitiga", "dictionary_id": 63, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Segitiga.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "6 + 2 = ?", "answer": "8", "dictionary_id": 9, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Pilih isyarat angka '7'", "answer": "kamus/7.png", "dictionary_id": 8, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Waktu untuk tidur adalah?", "answer": "Malam", "dictionary_id": 75, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "10 - 3 = ?", "answer": "7", "dictionary_id": 8, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.MATEMATIKA},
            {"question": "Praktikkan isyarat angka '4'", "answer": "4", "dictionary_id": 5, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/angka/04.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Bentuk dengan 4 sisi sama adalah?", "answer": "Persegi", "dictionary_id": 64, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
        ])

        # Level 4: Basic Activities (Sublevel 4.1: Pagi Hari)
        sublevel_id = 31
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Bangun'", "answer": "bangun", "dictionary_id": 79, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Bangun", "dictionary_id": 79, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Mandi' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "dictionary_id": 80, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Mandi'", "answer": "mandi", "dictionary_id": 80, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Aktivitas pada gambar ini adalah?", "answer": "Mandi", "dictionary_id": 80, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Bangun'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", "dictionary_id": 79, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Sarapan'", "answer": "sarapan", "dictionary_id": 81, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarapan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Di pagi hari, setelah bangun kita?", "answer": "Mandi", "dictionary_id": 80, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Gambar ini adalah aktivitas?", "answer": "Sarapan", "dictionary_id": 81, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sarapan.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Lakukan isyarat 'Bangun'", "answer": "bangun", "dictionary_id": 79, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Bangun.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 4: Basic Activities (Sublevel 4.2: Sekolah)
        sublevel_id = 32
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Main'", "answer": "main", "dictionary_id": 85, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Ajar", "dictionary_id": 82, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ajar.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Tulis' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm", "dictionary_id": 83, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Baca'", "answer": "Baca", "dictionary_id": 84, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Aktivitas pada gambar ini adalah?", "answer": "Tulis", "dictionary_id": 83, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tulis.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Main'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "dictionary_id": 85, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Duduk'", "answer": "duduk", "dictionary_id": 87, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Di sekolah, siswa biasanya?", "answer": "Ajar", "dictionary_id": 82, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Gambar ini adalah aktivitas?", "answer": "Baca", "dictionary_id": 84, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Baca.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Lakukan isyarat 'Nonton'", "answer": "nonton", "dictionary_id": 88, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 4: Basic Activities (Sublevel 4.3: Di Rumah)
        sublevel_id = 33
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Tidur'", "answer": "tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Duduk' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", "dictionary_id": 87, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Nonton'", "answer": "nonton", "dictionary_id": 88, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Aktivitas pada gambar ini adalah?", "answer": "Duduk", "dictionary_id": 87, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Makan'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", "dictionary_id": 97, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Makan'", "answer": "makan", "dictionary_id": 97, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Di rumah, kita biasanya?", "answer": "Tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Gambar ini adalah aktivitas?", "answer": "Nonton", "dictionary_id": 88, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nonton.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Lakukan isyarat 'Duduk'", "answer": "duduk", "dictionary_id": 87, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Duduk.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 4: Basic Activities (Sublevel 4.4: Bermain)
        sublevel_id = 34
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Lari'", "answer": "lari", "dictionary_id": 89, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Lari", "dictionary_id": 89, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Lompat' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm", "dictionary_id": 90, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Nyanyi'", "answer": "nyanyi", "dictionary_id": 91, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nyanyi.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Aktivitas pada gambar ini adalah?", "answer": "Lompat", "dictionary_id": 90, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Main'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "dictionary_id": 85, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Main'", "answer": "main", "dictionary_id": 85, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Saat bermain, anak-anak suka?", "answer": "Lari", "dictionary_id": 89, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Gambar ini adalah aktivitas?", "answer": "Nyanyi", "dictionary_id": 91, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Nyanyi.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Lakukan isyarat 'Lompat'", "answer": "lompat", "dictionary_id": 90, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lompat.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 4: Basic Activities (Sublevel 4.5: Bersih-bersih)
        sublevel_id = 35
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Cuci'", "answer": "cuci", "dictionary_id": 92, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Cuci.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Cuci", "dictionary_id": 92, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Cuci.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Sapu' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm", "dictionary_id": 94, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Mandi'", "answer": "mandi", "dictionary_id": 80, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Aktivitas pada gambar ini adalah?", "answer": "Sapu", "dictionary_id": 94, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Mandi'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "dictionary_id": 80, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Untuk membersihkan badan, kita?", "answer": "Mandi", "dictionary_id": 80, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Untuk membersihkan lantai, kita?", "answer": "Sapu", "dictionary_id": 94, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Gambar ini adalah aktivitas?", "answer": "Mandi", "dictionary_id": 80, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Mandi.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Lakukan isyarat 'Sapu'", "answer": "sapu", "dictionary_id": 94, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sapu.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 4: Basic Activities (Sublevel 4.6: Makan Minum)
        sublevel_id = 36
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Makan'", "answer": "makan", "dictionary_id": 97, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Makan", "dictionary_id": 97, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Minum' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minum.webm", "dictionary_id": 97, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Minum'", "answer": "minum", "dictionary_id": 97, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minum.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Aktivitas pada gambar ini adalah?", "answer": "Minum", "dictionary_id": 97, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minum.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Saat haus, kita?", "answer": "Minum", "dictionary_id": 97, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Saat lapar, kita?", "answer": "Makan", "dictionary_id": 97, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Pilih isyarat untuk 'Makan'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", "dictionary_id": 97, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Minum'", "answer": "minum", "dictionary_id": 97, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Minum.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Lakukan isyarat 'Makan'", "answer": "makan", "dictionary_id": 97, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 4: Basic Activities (Sublevel 4.7: Emosi)
        sublevel_id = 37
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Senang'", "answer": "senang", "dictionary_id": 98, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan perasaan?", "answer": "Senang", "dictionary_id": 98, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Sedih' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm", "dictionary_id": 99, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Sedih'", "answer": "sedih", "dictionary_id": 99, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Perasaan pada gambar ini adalah?", "answer": "Sedih", "dictionary_id": 99, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Sedih.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Saat ulang tahun, kita merasa?", "answer": "Senang", "dictionary_id": 98, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Saat kehilangan mainan, kita merasa?", "answer": "Sedih", "dictionary_id": 99, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Pilih isyarat untuk 'Senang'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "dictionary_id": 98, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lawan dari senang adalah?", "answer": "Sedih", "dictionary_id": 99, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Senang'", "answer": "senang", "dictionary_id": 98, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])

        # Level 4: Basic Activities (Sublevel 4.8: Kegiatan Luar)
        sublevel_id = 38
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Datang'", "answer": "datang", "dictionary_id": 102, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Datang.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Jalan", "dictionary_id": 102, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jalan.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Saat keluar rumah, kita biasanya?", "answer": "Jalan", "dictionary_id": 102, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Pilih isyarat untuk 'Jalan'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Jalan.webm", "dictionary_id": 102, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Untuk pergi ke sekolah, kita?", "answer": "Jalan", "dictionary_id": 102, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Pergi'", "answer": "pergi", "dictionary_id": 102, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pergi.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan?", "answer": "Pergi", "dictionary_id": 102, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Pergi.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat 'Datang' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Datang.webm", "dictionary_id": 102, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Lari'", "answer": "lari", "dictionary_id": 89, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Lawan dari pergi adalah?", "answer": "Datang", "dictionary_id": 102, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
        ])

        # Level 4: Basic Activities (Sublevel 4.9: Waktu Istirahat)
        sublevel_id = 39
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Tidur'", "answer": "tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Saat malam hari, kita?", "answer": "Tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Saat lelah, kita perlu?", "answer": "Tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Pilih isyarat untuk 'Tidur'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "dictionary_id": 86, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Aktivitas yang dilakukan di tempat tidur adalah?", "answer": "Tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Setelah tidur, kita akan?", "answer": "Bangun", "dictionary_id": 79, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Tidur' kembali", "answer": "tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Tidur yang cukup membuat badan?", "answer": "Sehat", "dictionary_id": 86, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Waktu tidur yang baik adalah?", "answer": "Malam", "dictionary_id": 75, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
        ])

        # Level 4: Basic Activities (Sublevel 4.10: Latihan)
        sublevel_id = 40
        soal_data.extend([
            {"question": "Tunjukkan isyarat untuk 'Main'", "answer": "main", "dictionary_id": 85, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Gambar ini menunjukkan aktivitas?", "answer": "Main", "dictionary_id": 85, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Isyarat untuk 'Makan' adalah...", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Makan.webm", "dictionary_id": 97, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Praktikkan isyarat 'Tidur'", "answer": "tidur", "dictionary_id": 86, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Perasaan pada gambar ini adalah?", "answer": "Senang", "dictionary_id": 98, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Pilih isyarat untuk 'Main'", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Main.webm", "dictionary_id": 85, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Tunjukkan isyarat 'Senang'", "answer": "senang", "dictionary_id": 98, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Senang.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
            {"question": "Aktivitas apa ini?", "answer": "Ajar", "dictionary_id": 82, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Ajar.webm", "categories": Soal.TypeSoalEnum.TEBAK_GAMBAR},
            {"question": "Manakah isyarat untuk 'Tidur'?", "answer": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Tidur.webm", "dictionary_id": 86, "sublevel_id": sublevel_id, "categories": Soal.TypeSoalEnum.PILIHAN_GANDA},
            {"question": "Lakukan isyarat 'Lari'", "answer": "lari", "dictionary_id": 89, "sublevel_id": sublevel_id, "video_url": "http://pmpk.kemdikbud.go.id/sibi/SIBI/katadasar/Lari.webm", "categories": Soal.TypeSoalEnum.OPEN_CAMERA},
        ])


        for data in soal_data:
            existing = self.db.query(Soal).filter(Soal.question == data["question"]).first()
            if existing:
                print(f"  Soal already exists: {data['question'][:50]}...")
                continue

            soal = Soal(**data)
            self.db.add(soal)

        self.db.commit()
