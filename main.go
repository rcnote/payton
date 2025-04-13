package main

import (
	"fmt"
	cmd "payton/command"
	"payton/config"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[Start Server Err!!!] ", err)
		}
	}()
	config.Init()
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
