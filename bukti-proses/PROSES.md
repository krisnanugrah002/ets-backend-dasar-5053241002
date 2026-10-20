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
Sesuai dengan permintaan Change Request B yang mengharuskan hasil `bentrok_dengan` diurutkan secara ID, saya menggunakan fungsi `sort.Slice`.Seluruh objek peminjaman yang beririsan ke dalam sebuah *slice* sementara bernama `pemKandidat`. Setelah iterasi `map` selesai dan bentrok terkonfirmasi karena kapasitas terlampaui, *slice* `pemKandidat` tersebut diurutkan berdasarkan `ID` naik, lalu disalin ke *slice* `bentrokDengan` bertipe `[]int` sebagai hasil akhir. Hal ini memastikan keluaran array JSON selalu konsisten.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: 404 saat pengujian API di postman
- Cara menemukan: Melakukan check pada URL ternyata ada . yang semestinya tidak ada di URL.
- Penyebab dan solusi: Penyebabnya ada . di akhir URL sehingga membuat 404. Solusinya menghapus . tersebut.

### Bug 2
- Gejala: Terjadi masalah pada pengujian postman
- Cara menemukan: menggunakan postman
- Penyebab dan solusi: Memperbaiki request

## 5. Change request
**Change Request B (Daftar Bentrok):**
Saat terjadi bentrok pada `POST /peminjaman`, sistem kini mengembalikan HTTP 409 beserta field `bentrok_dengan` yang memuat array ID peminjaman yang menghalangi rentang slot, diurutkan secara naik[cite: 6].
**Keputusan Desain:** 
Karena penyimpanan menggunakan tipe data `map` di Go yang urutan iterasinya selalu acak, saya memutuskan untuk mengumpulkan seluruh objek peminjaman yang beririsan ke dalam sebuah *slice* sementara bernama `pemKandidat`[cite: 9]. Setelah iterasi `map` selesai, saya memakai fungsi `sort.Slice` untuk mengurutkan objek di dalam *slice* tersebut berdasarkan `ID` naik[cite: 9]. Setelah urut, barulah ID-nya diekstrak ke dalam array `[]int` sebagai hasil akhir[cite: 9]. Ini memastikan keluaran JSON selalu konsisten sesuai permintaan CR[cite: 6].
