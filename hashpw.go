package main

import (
	"fmt"
	"os"

	"hotsauceshop/lib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hashpw <password>")
		os.Exit(1)
	}
	pw := os.Args[1]

	if len(pw) < 8 {
		fmt.Println("Password must be at least 8 characters long")
		os.Exit(1)
	}

	hashedPw, err := lib.HashPassword(pw)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(hashedPw)
}
