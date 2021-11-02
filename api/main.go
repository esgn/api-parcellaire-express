package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// init checks configuration at module initialization.
func init() {
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
}

// checkEnv ensures all necessary env data is present.
// panic in case of missing env.
// hide
func checkEnv() {
	mandatoryEnvs := []string{
		"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB",
		"POSTGRES_HOST", "POSTGRES_PORT", "API_PORT", "POSTGRES_SCHEMA",
		"MAX_FEATURE", "VIEWER_URL"}

	optionalEnvs := []string{"VIEWER_URL"}

	reIsPassword := regexp.MustCompile(`(?i)password|passwd|key`)

	for _, theEnvName := range mandatoryEnvs {
		theEnv, isPresent := os.LookupEnv(theEnvName)
		if !isPresent {
			log.Fatalf("Sorry, env %v is mandatory. Please check environment variable or use --env <path> options", theEnvName)
		}

		if reIsPassword.MatchString(theEnvName) {
			theEnv = "********"
		}

		fmt.Printf("* %s : %s \n", theEnvName, theEnv)
	}

	for _, theEnvName := range optionalEnvs {
		theEnv, isPresent := os.LookupEnv(theEnvName)

		if !isPresent {
			theEnv = "?"
		} else if reIsPassword.MatchString(theEnvName) {
			theEnv = "********"
		}

		fmt.Printf("# %s : %s \n", theEnvName, theEnv)
	}
}

// initDB creates a global connection pool from identifiers.
// DB is global, because it is the connection pool not the connection itself.
func initDB(user, password, dbname, hostname, port string) *sql.DB {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, hostname, port)

	DB, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return DB
}

func main() {
	// I choosed to initiliazed DB outside the app
	// because of the global nature of the connection pool (check DB.Close() comment)
	DB := initDB(os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"))
	// not really necessary, but maniac decision 💩
	defer DB.Close()

	a := App{}
	a.Initialize(DB)
	a.Run(":" + os.Getenv("API_PORT"))
}
