package main

import "github.com/matejeliash/mepm/internal/tui"

func Resolve(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	tm := tui.NewTuiManager()

	for true {
		tm.PrintAllRecords()
		tm.SelectAction()

	}

	// dbPath := "./test.db"

	// pm, err := passmanager.NewPassManager(dbPath)
	// Resolve(err)

	// var password string

	// fmt.Println("enter password:")
	// fmt.Scan(&password)

	// if !pm.DB.HasMasterTable() {
	// 	err = pm.InitPassManagerDB(password)
	// 	Resolve(err)
	// }

	// err = pm.FetchMasterTable()
	// Resolve(err)

	// for !pm.IsPasswordCorrect(password) {
	// 	fmt.Println("provided password is incorrect, enter again:")
	// 	fmt.Scan(&password)
	// }

	// pm.Password = password

	// err = pm.InsertRecord(password, "username", "info")
	// records, err := pm.GetRecords()
	// Resolve(err)
	// fmt.Println(records)
	// decryptedPassword, err := pm.DecryptPassword(records[0])
	// Resolve(err)
	// fmt.Println(decryptedPassword)

}
