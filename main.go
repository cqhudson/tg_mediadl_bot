package main

import (
    "log"
    "os" // needed to load env variables loaded from godotenv

    "github.com/joho/godotenv"
)

func main() {
    
    // First load the .env file
    log.Print("Loading .env file")
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    log.Print(".env file successfully loaded")

    // Fetch Telegram API token from .env
    telegramToken := os.Getenv("TG_API_KEY")
    log.Printf("Your Telegram API token is %s", telegramToken) // This line is for local debugging and should be disabled when running in prod
}
