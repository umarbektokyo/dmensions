package main

import (
	"dmensions/internal/utils"
	"fmt"
)

func main() {
	utils.Splash()
	fmt.Println("Version " + utils.VERSION)
}
