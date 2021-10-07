package main

import "os"

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"))
	a.Run(":" + os.Getenv("API_PORT"))
}
