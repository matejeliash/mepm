package db

import (
	"github.com/glebarez/sqlite"
	"github.com/matejeliash/mepm/internal/models"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func CreateDb(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	} else {
		return &Database{DB: db}, nil
	}
}

func (d *Database) HasMasterTable() bool {
	return d.DB.Migrator().HasTable(&models.MasterTable{})
}

func (d *Database) CreateMasterTable(tableStruct *models.MasterTable) error {
	err := d.DB.AutoMigrate(&models.MasterTable{})
	if err != nil {
		return err
	}
	err = d.DB.Create(tableStruct).Error
	return err
}

func (d *Database) GetMasterTable() (*models.MasterTable, error) {

	var table models.MasterTable
	err := d.DB.First(&table).Error
	return &table, err
}
