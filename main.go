package main

import (
	"fmt"
	"os"

	"github.com/jld3103/segelboot/cmd"
)

func main() {
	err := cmd.NewRootCmd().Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
