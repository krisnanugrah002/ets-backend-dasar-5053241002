package main

import "sync"

// Store menyimpan data di memori dengan aman untuk goroutine
type Store struct {
	mu         sync.Mutex
	sumberDaya map[int]SumberDaya
	peminjaman map[int]Peminjaman
	nextSDID   int
	nextPemID  int
}

func NewStore() *Store {
	return &Store{
		sumberDaya: make(map[int]SumberDaya),
		peminjaman: make(map[int]Peminjaman),
		nextSDID:   1,
		nextPemID:  1,
	}
}
