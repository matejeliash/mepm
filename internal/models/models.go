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
	KeySalt      []byte
	PasswordSalt []byte
	PasswordHash []byte
	KeyGenParams
}

type Record struct {
	gorm.Model
	EncryptedPassword []byte
	Info              string
	Username          string
}

type Note struct {
	gorm.Model
	Title         string
	EncryptedText []byte
}
