package main

import (
	"fmt"

	"github.com/matejeliash/mepm/internal/gui"
)

func main() {

	gm, err := gui.NewGuiManager()
	if err != nil {
		panic(err)
	}
	fmt.Println(gm)
	gm.ShowFirstScreen()
}
