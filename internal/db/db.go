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
	if err != nil {
		return err
	}
	err = d.DB.AutoMigrate(&models.Record{})
	if err != nil {
		return err
	}
	err = d.DB.AutoMigrate(&models.Note{})
	if err != nil {
		return err
	}

	return err

}

func (d *Database) GetMasterTable() (*models.MasterTable, error) {

	var table models.MasterTable
	err := d.DB.First(&table).Error
	return &table, err
}

func (d *Database) GetRecordsTable() ([]models.Record, error) {

	var records []models.Record
	err := d.DB.Find(&records).Error
	return records, err
}

func (d *Database) GetNotesTable() ([]models.Note, error) {

	var notes []models.Note
	err := d.DB.Find(&notes).Error
	return notes, err
}

func (d *Database) InsertRecord(record *models.Record) error {
	err := d.DB.Create(record).Error
	return err

}

func (d *Database) InsertNote(note *models.Note) error {
	err := d.DB.Create(note).Error
	return err

}

func (d *Database) RemoveRecord(record *models.Record) error {

	err := d.DB.Delete(&models.Record{}, record.ID).Error
	return err

}

func (d *Database) RemoveNote(note *models.Note) error {

	err := d.DB.Delete(&models.Note{}, note.ID).Error
	return err

}

func (d *Database) UpdateRecord(record *models.Record) error {
	err := d.DB.Model(&models.Record{}).Where("id = ?", record.ID).Updates(record).Error
	return err

}
