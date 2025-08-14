package main

import (
	"fmt"

	"github.com/matejeliash/mepm/internal/passmanager"
)

func Resolve(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	dbPath := "./test.db"

	pm, err := passmanager.NewPassManager(dbPath)
	Resolve(err)

	err = pm.InitPassManagerDB("heslo")
	Resolve(err)

	mt, err := pm.DB.GetMasterTable()
	Resolve(err)

	fmt.Printf("%+v", mt)

}
