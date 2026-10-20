# PROSES

## 1. Persona dan masalahnya
Pengguna adalah Admin Himpunan Mahasiswa.
Masalah yang diselesaikan adalah ruang himpunan tidak boleh langsung terpakai

## 2. Rancangan struktur data
- `SumberDaya`: struct menyimpan ID, nama, jenis, dan stok.
- `Peminjaman`: struct menyimpan ID, ID sumber daya, peminjam, tanggal, rentang slot, jumlah, status (default `"diajukan"`), dan alasan (default `""`).
- Penyimpanan *in-memory* memakai `map[int]SumberDaya` dan `map[int]Peminjaman`.
- Menggunakan `sync.Mutex` untuk memastikan keamanan data saat *concurrency*.

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
