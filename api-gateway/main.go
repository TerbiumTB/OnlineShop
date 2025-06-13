package main

import (
	"log"
	"os"
)

func main() {
	lg := log.New(os.Stdout, "API Gateway ", log.LstdFlags)

	lg.Println("Starting API Gateway")
}
