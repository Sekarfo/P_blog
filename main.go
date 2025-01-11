package main

import (
    "net/http"
    "os"

    log "github.com/sirupsen/logrus"

    "personal_blog/config"
    "personal_blog/routes"
)

func init() {
    // Set log format to JSON
    log.SetFormatter(&log.JSONFormatter{})

    // Create log file
    file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if (err == nil) {
        log.SetOutput(file)
    } else {
        log.Info("Failed to log to file, using default stderr")
    }
}

func main() {
    db := config.InitDB()
    config.AutoMigrateDB(db)
    router := routes.SetupRoutes(db)

    log.Info("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}