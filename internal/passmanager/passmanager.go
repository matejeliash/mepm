package passmanager

import (
	"fmt"

	"github.com/matejeliash/mepm/internal/db"
	"github.com/matejeliash/mepm/internal/models"
	"github.com/matejeliash/mepm/internal/secutils"
)

type PassManager struct {
	Password    string
	MasterTable models.MasterTable
	DB          *db.Database
}

func NewPassManager(dbPath string) (*PassManager, error) {

	db, err := db.CreateDb(dbPath)
	if err != nil {
		return nil, err
	}
	pm := &PassManager{}
	pm.DB = db
	return pm, nil

}

func (pm *PassManager) FetchMasterTable() error {
	mt, err := pm.DB.GetMasterTable()
	if err != nil {
		return err
	}
	pm.MasterTable = *mt

	return nil

}

func (pm *PassManager) InitPassManagerDB(password string) error {

	// generate salts
	passSalt, err := secutils.GenerateSalt()
	if err != nil {
		return err
	}
	keySalt, err := secutils.GenerateSalt()
	if err != nil {
		return err
	}
	params := secutils.GetDefaultKeyGenParams()

	mt := &models.MasterTable{}
	mt.KeyGenParams = models.KeyGenParams(*params)

	// passSaltBase64 := base64.StdEncoding.EncodeToString(passSalt)
	// keySaltBase64 := base64.StdEncoding.EncodeToString(keySalt)

	// mt.KeySalt = keySaltBase64
	// mt.PasswordSalt = passSaltBase64

	mt.KeySalt = keySalt
	mt.PasswordSalt = passSalt

	hash := secutils.HashPassword([]byte(password), passSalt, params)
	// hashBase64 := base64.StdEncoding.EncodeToString(hash)
	// mt.PasswordHash = hashBase64

	mt.PasswordHash = hash

	err = pm.DB.CreateMasterTable(mt)
	return err

}

func (pm *PassManager) InsertRecord(recordPassword, username, info string) error {

	//passwordSalt, err := base64.StdEncoding.DecodeString(pm.MasterTable.PasswordSalt)

	key := secutils.GenerateKey([]byte(pm.Password), pm.MasterTable.KeySalt, (*secutils.KeyGenParams)(&pm.MasterTable.KeyGenParams))

	encryptedBytes, err := secutils.EncryptAES([]byte(recordPassword), key)
	if err != nil {
		return err
	}

	//encryptedBytesBase64 := base64.StdEncoding.EncodeToString(encryptedBytes)
	record := &models.Record{
		Info:              info,
		Username:          username,
		EncryptedPassword: encryptedBytes,
	}
	err = pm.DB.InsertRecord(record)
	return err
}

func (pm *PassManager) InsertNote(title, text string) error {

	//passwordSalt, err := base64.StdEncoding.DecodeString(pm.MasterTable.PasswordSalt)

	key := secutils.GenerateKey([]byte(pm.Password), pm.MasterTable.KeySalt, (*secutils.KeyGenParams)(&pm.MasterTable.KeyGenParams))

	encryptedBytes, err := secutils.EncryptAES([]byte(text), key)
	if err != nil {
		return err
	}

	note := &models.Note{
		Title:         title,
		EncryptedText: encryptedBytes,
	}
	err = pm.DB.InsertNote(note)
	return err
}

func (pm *PassManager) GetRecords() ([]models.Record, error) {
	records, err := pm.DB.GetRecordsTable()
	if err != nil {
		return nil, err
	}
	return records, nil

}

func (pm *PassManager) GetNotes() ([]models.Note, error) {

	notes, err := pm.DB.GetNotesTable()
	if err != nil {
		return nil, err
	}
	return notes, nil

}

func (pm *PassManager) DecryptPassword(record models.Record) (string, error) {
	fmt.Println(pm.Password)

	key := secutils.GenerateKey([]byte(pm.Password), pm.MasterTable.KeySalt, (*secutils.KeyGenParams)(&pm.MasterTable.KeyGenParams))
	fmt.Println(key)

	decryptedBytes, err := secutils.DecryptAES(record.EncryptedPassword, key)

	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}

func (pm *PassManager) DecryptNote(note models.Note) (string, error) {
	fmt.Println(pm.Password)

	key := secutils.GenerateKey([]byte(pm.Password), pm.MasterTable.KeySalt, (*secutils.KeyGenParams)(&pm.MasterTable.KeyGenParams))
	fmt.Println(key)

	decryptedBytes, err := secutils.DecryptAES(note.EncryptedText, key)

	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}

func (pm *PassManager) IsPasswordCorrect(password string) bool {
	return secutils.ComparePasswordAndHash(
		[]byte(password),
		pm.MasterTable.PasswordSalt,
		pm.MasterTable.PasswordHash,
		(*secutils.KeyGenParams)(&pm.MasterTable.KeyGenParams),
	)
}

func (pm *PassManager) UpdateRecord(id int, info, username, recordPassword string) error {

	key := secutils.GenerateKey([]byte(pm.Password), pm.MasterTable.KeySalt, (*secutils.KeyGenParams)(&pm.MasterTable.KeyGenParams))

	encryptedBytes, err := secutils.EncryptAES([]byte(recordPassword), key)
	if err != nil {
		return err
	}

	record := models.Record{}
	record.ID = uint(id)
	record.Info = info
	record.Username = username
	record.EncryptedPassword = encryptedBytes

	err = pm.DB.UpdateRecord(&record)

	return err

}
