package main

import (
	"log"

	"github.com/Sekarfo/P_blog/app"
)

/*
import (
<<<<<<< Updated upstream
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
    if err == nil {
        log.SetOutput(file)
    } else {
        log.Info("Failed to log to file, using default stderr")
    }
=======
	"log"
	"personal_blog/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
>>>>>>> Stashed changes
}
*/

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}

}
