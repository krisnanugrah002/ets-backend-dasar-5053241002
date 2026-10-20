# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## Checkpoint 1: Setup AI dan Project
- **Prompt:** "Cek lampiran file yang saya upload, bantu saya untuk membuat direktori project, dan step untuk setup project"
- **Dipakai:** Pembuatan file `main.go`, `models.go`, `store.go`, dan `handlers.go`
- **Diubah:** -

## Checkpoint 1: Pembuatan Struktur Awal dan B1-B3
- **Prompt:** "Buatkan struktur file dasar Go 1.22 dengan `http.NewServeMux`, in-memory map, dan mutex untuk proyek PinjamKampus Persona C."
- **Dipakai:** Kerangka dasar `main.go`, `models.go`, `store.go`, dan `handlers.go` untuk rute awal B1, B2, dan B3.
- **Diubah:** Menambahkan *field* `Status` (default `"diajukan"`) dan `Alasan` pada struct `Peminjaman` karena AI awalnya hanya membuat struktur default.

## Checkpoint 1: Pengecekan Error
- **Prompt:** "Cek apakah ada yang salah dengan URL API, dikarenakan ketika di uji pada postman keluar 404. URL adalah berikut: http://localhost:8080/sumber-daya.
- **Dipakai:** http://localhost:8080/sumber-daya
- **Diubah:** Mengubah URL dan menghapus . yang membuat error 404.

## Checkpoint 2: Validasi, Bentrok, Concurrency, dan Alur Persetujuan (B4-B7 & Persona C)
- **Prompt:** "Lengkapi fungsi B4-B7 dengan validasi input yang ketat, cek bentrok kapasitas, penerapan mutex untuk concurrency, serta tambahkan handler khusus Persona C yaitu setujui, tolak (dengan alasan), dan selesaikan."
- **Dipakai:** Fungsi pembantu `validasiPeminjaman` dan `cekBentrok`, serta implementasi penuh pada handler `handleUpdatePeminjaman`, `handleDeletePeminjaman`, `handleSetujui`, `handleTolak`, dan `handleSelesaikan` di `handlers.go`.
- **Diubah:** Hasil dari AI awalnya menghitung semua peminjaman dalam pengecekan bentrok. Saya mengubah logika di dalam fungsi `cekBentrok` agar HANYA peminjaman yang berstatus `"diajukan"` dan `"disetujui"` yang ikut dihitung menahan slot, sesuai dengan aturan Persona C. Selain itu, saya memastikan penempatan `defer a.store.mu.Unlock()` diletakkan tepat setelah inisiasi *lock* agar tidak terjadi *deadlock* saat program melakukan *return* lebih awal karena gagal validasi.

## Checkpoint 3: Implementasi Change Request B: Menampilkan Daftar Bentrok (`bentrok_dengan`)
- **Prompt:** "Sudah aku tambahkan file markdown untuk CR-B, coba kamu cek lalu lakukan implementasi tersebut di kode ini. JANGAN UPDATE BUKTI PROSES. Hanya code aja"
- **Dipakai:** Modifikasi fungsi `cekBentrok` dan handler `handleCreatePeminjaman` pada `pinjamkampus/handlers.go`, serta penambahan fungsi uji `TestChangeRequestB` pada `pinjamkampus/handlers_test.go`.
- **Diubah:** 
  1. Mengubah nilai kembalian fungsi `cekBentrok` menjadi `(bool, []int)`.
  2. Menambahkan filter pengumpulan peminjaman yang beririsan rentang slotnya (`p.SlotMulai <= req.SlotSelesai && req.SlotMulai <= p.SlotSelesai`), berstatus menahan slot (`diajukan` dan `disetujui`), pada tanggal dan sumber daya yang sama.
  3. Mengurutkan ID peminjaman yang bentrok secara menaik menggunakan `sort.Slice`.
  4. Mengubah respons error 409 pada `POST /peminjaman` agar menyertakan field `"bentrok_dengan": bentrokIDs` di samping field `"error": "slot sudah terpakai"`.
  5. Memastikan endpoint `PUT` tetap tidak diubah (tidak menyertakan `bentrok_dengan`) sesuai instruksi lembar CR-B.



