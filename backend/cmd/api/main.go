package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// ctx := context.Background()

	fmt.Printf("Main.go")
}
