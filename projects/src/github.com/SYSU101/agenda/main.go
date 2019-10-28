package main

import (
	"fmt"

	"github.com/SYSU101/agenda/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
