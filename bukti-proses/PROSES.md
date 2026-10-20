# PROSES

## 1. Persona dan masalahnya
Persona adalah Admin Himpunan Mahasiswa.
Masalah yang diselesaikan adalah ruang himpunan tidak boleh langsung terpakai. Seluruh pengajuan ruangan harus disetujui oleh admin, apabila terdapat penolakan, admin akan menulis alasan penolakan pengajuan.

## 2. Rancangan struktur data
- **SumberDaya:** Menggunakan `struct` berisi `ID`, `Nama`, `Jenis`, dan `Stok`.
- **Peminjaman:** Menggunakan `struct` dengan tambahan field `Status` dan `Alasan` khusus untuk Persona C.
- **Penyimpanan:** Disimpan di memori menggunakan `map[int]SumberDaya` dan `map[int]Peminjaman`.
- **Concurrency:** Data dibungkus dalam struct `Store` yang memiliki `sync.Mutex` agar aman saat diakses dari goroutine yang berjalan bersamaan.

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: 404 saat pengujian API di postman
- Cara menemukan: Melakukan check pada URL ternyata ada . yang semestinya tidak ada di URL.
- Penyebab dan solusi: Penyebabnya ada . di akhir URL sehingga membuat 404. Solusinya menghapus . tersebut.

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
