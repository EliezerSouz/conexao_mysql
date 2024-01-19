package modulos

import (
	"fmt"
	"os"
)

func LimparXmls() error {

	dirName := "fs-xmls"

	err := os.RemoveAll(dirName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Directory", dirName, "removed successfully")
	}
	return nil
}
