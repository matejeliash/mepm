package models

import "gorm.io/gorm"

type KeyGenParams struct {
	Iters       uint32
	MemoryBytes uint32
	Threads     uint8
	KeyLen      uint32
}

type MasterTable struct {
	gorm.Model
	KeySalt      string
	PasswordSalt string
	PasswordHash string
	KeyGenParams
}
