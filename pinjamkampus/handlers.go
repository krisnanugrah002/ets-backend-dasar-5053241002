package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type App struct {
	store *Store
}

func NewHandler() http.Handler {
	app := &App{store: NewStore()}
	mux := http.NewServeMux()

	// B1
	mux.HandleFunc("POST /sumber-daya", app.handleCreateSumberDaya)
	mux.HandleFunc("GET /sumber-daya", app.handleGetSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", app.handleGetSumberDayaByID)

	// B2, B3, B6
	mux.HandleFunc("POST /peminjaman", app.handleCreatePeminjaman)
	mux.HandleFunc("GET /peminjaman", app.handleGetPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", app.handleGetPeminjamanByID)
	mux.HandleFunc("PUT /peminjaman/{id}", app.handleUpdatePeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", app.handleDeletePeminjaman)

	// Persona C
	mux.HandleFunc("POST /peminjaman/{id}/setujui", app.handleSetujui)
	mux.HandleFunc("POST /peminjaman/{id}/tolak", app.handleTolak)
	mux.HandleFunc("POST /peminjaman/{id}/selesaikan", app.handleSelesaikan)

	return mux
}

// Utils
func sendError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func sendJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Handler Sumber Daya
func (a *App) handleCreateSumberDaya(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Nama  string `json:"nama"`
		Jenis string `json:"jenis"`
		Stok  int    `json:"stok"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}

	req.Nama = strings.TrimSpace(req.Nama)
	if req.Nama == "" || len(req.Nama) > 50 {
		sendError(w, http.StatusBadRequest, "Nama tidak valid")
		return
	}
	if req.Jenis != "ruang" && req.Jenis != "alat" {
		sendError(w, http.StatusBadRequest, "Jenis tidak valid")
		return
	}
	if req.Stok < 1 {
		sendError(w, http.StatusBadRequest, "Stok tidak valid")
		return
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	sd := SumberDaya{
		ID:    a.store.nextSDID,
		Nama:  req.Nama,
		Jenis: req.Jenis,
		Stok:  req.Stok,
	}
	a.store.sumberDaya[sd.ID] = sd
	a.store.nextSDID++

	sendJSON(w, http.StatusCreated, sd)
}

func (a *App) handleGetSumberDaya(w http.ResponseWriter, r *http.Request) {
	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	res := make([]SumberDaya, 0)
	for i := 1; i < a.store.nextSDID; i++ {
		if sd, ok := a.store.sumberDaya[i]; ok {
			res = append(res, sd)
		}
	}
	sendJSON(w, http.StatusOK, res)
}

func (a *App) handleGetSumberDayaByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	sd, ok := a.store.sumberDaya[id]
	if !ok {
		sendError(w, http.StatusNotFound, "Sumber daya tidak ditemukan")
		return
	}
	sendJSON(w, http.StatusOK, sd)
}

// --- Validasi & Bentrok ---
func validasiPeminjaman(req Peminjaman, store *Store) (int, string) {
	req.Peminjam = strings.TrimSpace(req.Peminjam)
	if req.Peminjam == "" || len(req.Peminjam) > 50 {
		return http.StatusBadRequest, "Peminjam tidak valid"
	}
	if _, err := time.Parse("2006-01-02", req.Tanggal); err != nil {
		return http.StatusBadRequest, "Tanggal tidak valid"
	}
	if req.SlotMulai < 1 || req.SlotMulai > 8 || req.SlotSelesai < 1 || req.SlotSelesai > 8 || req.SlotMulai > req.SlotSelesai {
		return http.StatusBadRequest, "Slot tidak valid"
	}
	sd, ok := store.sumberDaya[req.SumberDayaID]
	if !ok {
		return http.StatusBadRequest, "Sumber daya tidak ditemukan"
	}
	if req.Jumlah < 1 || req.Jumlah > sd.Stok {
		return http.StatusBadRequest, "Jumlah tidak valid"
	}
	return 200, ""
}

func cekBentrok(req Peminjaman, store *Store, ignoreID int) bool {
	terpakai := make([]int, 9)
	for _, p := range store.peminjaman {
		if p.SumberDayaID == req.SumberDayaID && p.Tanggal == req.Tanggal && p.ID != ignoreID {
			// Persona C: menahan slot jika diajukan atau disetujui
			if p.Status == "diajukan" || p.Status == "disetujui" {
				for i := p.SlotMulai; i <= p.SlotSelesai; i++ {
					terpakai[i] += p.Jumlah
				}
			}
		}
	}
	sd := store.sumberDaya[req.SumberDayaID]
	for i := req.SlotMulai; i <= req.SlotSelesai; i++ {
		if terpakai[i]+req.Jumlah > sd.Stok {
			return true
		}
	}
	return false
}

// Handler Peminjaman
func (a *App) handleCreatePeminjaman(w http.ResponseWriter, r *http.Request) {
	var req Peminjaman
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	if code, msg := validasiPeminjaman(req, a.store); code != 200 {
		sendError(w, code, msg)
		return
	}
	if cekBentrok(req, a.store, 0) {
		sendError(w, http.StatusConflict, "Slot bentrok")
		return
	}

	req.ID = a.store.nextPemID
	req.Status = "diajukan" // Khas Persona C
	req.Alasan = ""
	req.Peminjam = strings.TrimSpace(req.Peminjam)

	a.store.peminjaman[req.ID] = req
	a.store.nextPemID++

	sendJSON(w, http.StatusCreated, req)
}

func (a *App) handleGetPeminjaman(w http.ResponseWriter, r *http.Request) {
	qSD := r.URL.Query().Get("sumber_daya_id")
	qTgl := r.URL.Query().Get("tanggal")
	qStatus := r.URL.Query().Get("status")

	var filterSD int
	var err error
	if qSD != "" {
		filterSD, err = strconv.Atoi(qSD)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Filter sumber_daya_id tidak valid")
			return
		}
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	res := make([]Peminjaman, 0)
	for i := 1; i < a.store.nextPemID; i++ {
		if p, ok := a.store.peminjaman[i]; ok {
			match := true
			if qSD != "" && p.SumberDayaID != filterSD {
				match = false
			}
			if qTgl != "" && p.Tanggal != qTgl {
				match = false
			}
			if qStatus != "" && p.Status != qStatus {
				match = false
			}
			if match {
				res = append(res, p)
			}
		}
	}
	sendJSON(w, http.StatusOK, res)
}

func (a *App) handleGetPeminjamanByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	p, ok := a.store.peminjaman[id]
	if !ok {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}
	sendJSON(w, http.StatusOK, p)
}

func (a *App) handleUpdatePeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	oldP, ok := a.store.peminjaman[id]
	if !ok {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}

	var req Peminjaman
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}

	if code, msg := validasiPeminjaman(req, a.store); code != 200 {
		sendError(w, code, msg)
		return
	}

	req.Status = oldP.Status
	req.Alasan = oldP.Alasan

	if cekBentrok(req, a.store, id) {
		sendError(w, http.StatusConflict, "Slot bentrok")
		return
	}

	req.ID = id
	req.Peminjam = strings.TrimSpace(req.Peminjam)
	a.store.peminjaman[id] = req

	sendJSON(w, http.StatusOK, req)
}

func (a *App) handleDeletePeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()

	if _, ok := a.store.peminjaman[id]; !ok {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}
	delete(a.store.peminjaman, id)
	w.WriteHeader(http.StatusNoContent)
}

// Handler Persona C
func (a *App) handleSetujui(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}
	a.store.mu.Lock()
	defer a.store.mu.Unlock()
	p, ok := a.store.peminjaman[id]
	if !ok {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}
	if p.Status != "diajukan" {
		sendError(w, http.StatusConflict, "Status bukan diajukan")
		return
	}
	p.Status = "disetujui"
	a.store.peminjaman[id] = p
	sendJSON(w, http.StatusOK, p)
}

func (a *App) handleTolak(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}
	var req struct {
		Alasan string `json:"alasan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Body tidak valid")
		return
	}
	req.Alasan = strings.TrimSpace(req.Alasan)
	if req.Alasan == "" || len(req.Alasan) > 100 {
		sendError(w, http.StatusBadRequest, "Alasan wajib diisi, max 100 char")
		return
	}

	a.store.mu.Lock()
	defer a.store.mu.Unlock()
	p, ok := a.store.peminjaman[id]
	if !ok {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}
	if p.Status != "diajukan" {
		sendError(w, http.StatusConflict, "Status bukan diajukan")
		return
	}
	p.Status = "ditolak"
	p.Alasan = req.Alasan
	a.store.peminjaman[id] = p
	sendJSON(w, http.StatusOK, p)
}

func (a *App) handleSelesaikan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}
	a.store.mu.Lock()
	defer a.store.mu.Unlock()
	p, ok := a.store.peminjaman[id]
	if !ok {
		sendError(w, http.StatusNotFound, "Peminjaman tidak ada")
		return
	}
	if p.Status != "disetujui" {
		sendError(w, http.StatusConflict, "Status bukan disetujui")
		return
	}
	p.Status = "selesai"
	a.store.peminjaman[id] = p
	sendJSON(w, http.StatusOK, p)
}
