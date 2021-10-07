package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	flagUseEnv := flag.String("env", "", "Provides the path to .env files")

	flag.Parse()

	if *flagUseEnv != "" {
		fmt.Printf("Loading .env : %v", flagUseEnv)
		err := godotenv.Load(*flagUseEnv)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		fmt.Println("✅")
	}

	a := App{}
	a.Initialize(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"))
	a.Run(":" + os.Getenv("API_PORT"))
}
