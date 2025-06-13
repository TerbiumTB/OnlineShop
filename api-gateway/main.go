package main

import (
	"log"
	"os"
)

func main() {
	lg := log.New(os.Stdout, "File storage ", log.LstdFlags)

	lg.Println("Starting api gateway")
}
