package main

import (
	"os"
	"fmt"
	
	"import/geolite"
	"main/locdata"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "import" {
			fmt.Println("Running import 1.")
			geolite.Import1()
			os.Exit(0)
		}
		
		if os.Args[1] == "import2" {
			fmt.Println("Running import 2.")
			geolite.Import2()
			os.Exit(0)
		}
	}
	
	locdata.Init()
}