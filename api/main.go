package main

import (
	"database/sql"
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
		fmt.Printf("Loading .env : %v ", *flagUseEnv)
		err := godotenv.Load(*flagUseEnv)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		fmt.Println("✅")
	}

	checkEnv()

	initDB(os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"))

	a := App{}
	a.Initialize()
	a.Run(":" + os.Getenv("API_PORT"))
}

// initDB creates a global connection pool from identifiers.
// DB is global, because it is the connection pool not the connection itself.
func initDB(user, password, dbname, hostname, port string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, hostname, port)

	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

// checkEnv ensures all necessary env data is present.
func checkEnv() {
	mandatoryEnvs := []string{
		"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
		"POSTGRES_HOST", "POSTGRES_PORT", "API_PORT", "POSTGRES_SCHEMA", "MAX_FEATURE", "LIMIT_FEATURE"}

	for i := 0; i < len(mandatoryEnvs); i++ {
		theEnv, isPresent := os.LookupEnv(mandatoryEnvs[i])
		if !isPresent {
			log.Panicf("Sorry, env %v is mandatory. Please check environment variable or use --env <path> options", mandatoryEnvs[i])
		}
		fmt.Printf("* %s : %s \n", mandatoryEnvs[i], theEnv)
	}
}
