package main

// SumberDaya merepresentasikan ruang atau alat
type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"` // "ruang" atau "alat"
	Stok  int    `json:"stok"`
}

// Peminjaman merepresentasikan transaksi peminjaman
type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"` // Persona C: "diajukan", "disetujui", "ditolak", "selesai"
	Alasan       string `json:"alasan"` // Persona C: alasan jika ditolak, default ""
}
