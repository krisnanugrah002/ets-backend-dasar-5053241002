package main

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type Resource struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Stock int    `json:"stock"`
}

type Peminjaman struct {
	ID         int    `json:"id"`
	ResourceID int    `json:"sumber_daya_id"`
	Borrower   string `json:"peminjam"`
	Date       string `json:"tanggal"`
	StartSlot  int    `json:"slot_mulai"`
	EndSlot    int    `json:"slot_selesai"`
	Quantity   int    `json:"jumlah"`
	Status     string `json:"status"`
	Reason     string `json:"alasan"`
}

type CampusStore struct {
	mu               sync.Mutex
	daftarSumberDaya []Resource
	daftarPeminjaman []Peminjaman
	nextResourceID   int
	nextBookingID    int
}

func NewCampusStore() *CampusStore {
	return &CampusStore{
		resources:      make([]Resource, 0),
		bookings:       make([]Peminjaman, 0),
		nextResourceID: 1,
		nextBookingID:  1,
	}
}

func (s *CampusStore) AddResource(name, resType string, stock int) (Resource, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// [DAUR ULANG file.go - Pertemuan 04 & 06] Trim spasi dengan strings.TrimSpace
	name = strings.TrimSpace(name)
	resType = strings.TrimSpace(resType)

	if name == "" || stock < 1 || (resType != "ruang" && resType != "alat") {
		return Resource{}, ErrInvalidInput
	}

	res := Resource{
		ID:    s.nextResourceID,
		Name:  name,
		Type:  resType,
		Stock: stock,
	}
	s.nextResourceID++
	s.resources = append(s.resources, res)
	return res, nil
}

func (s *CampusStore) ListResources() []Resource {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Resource, len(s.resources))
	copy(out, s.resources)
	return out
}
