package passmanager

import (
	"encoding/base64"

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

	passSaltBase64 := base64.StdEncoding.EncodeToString(passSalt)
	keySaltBase64 := base64.StdEncoding.EncodeToString(keySalt)

	mt.KeySalt = keySaltBase64
	mt.PasswordSalt = passSaltBase64

	hash := secutils.HashPassword([]byte(password), passSalt, params)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	mt.PasswordHash = hashBase64

	err = pm.DB.CreateMasterTable(mt)
	return err

}
