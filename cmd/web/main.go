package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Atul-Ranjan12/tourism/internal/config"
	"github.com/Atul-Ranjan12/tourism/internal/driver"
	"github.com/Atul-Ranjan12/tourism/internal/handlers"
	"github.com/Atul-Ranjan12/tourism/internal/helpers"
	"github.com/Atul-Ranjan12/tourism/internal/models"
	"github.com/Atul-Ranjan12/tourism/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":5000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	os.Setenv("TZ", "UTC")
	db, err := run()
	if err != nil {
		log.Println("An unexpected error occured while initializing the database, quiting with error: ", err)
	}
	defer db.SQL.Close()
	// Close the Mailing Channel
	defer close(app.MailChan)
	fmt.Println("Starting Mailing Channel...")
	listenForMail()

	fmt.Println("Starting Server on localhost port", portNumber)

	// Set up server configurations:
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	serverError := server.ListenAndServe()
	log.Fatal(serverError)
}

// run handles major initialization processes for the app
func run() (*driver.DB, error) {
	// Register Premitive Types for the Session (Variables the session needs to store)
	gob.Register(models.User{})

	// Set up Channel to send the mail
	mailChan := make(chan models.ConfirmationMailData)
	app.MailChan = mailChan

	// Set up configuration variable to not be in production mode
	app.InProduction = false

	// Set up Error and information log
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	// Initialize the session for the application
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	// Add the session to the application config
	app.Session = session

	// TODO: Handle connection to the database in this section
	log.Println("Connecting to database....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=tourism user=postgres password=krishna")
	if err != nil {
		log.Fatal("Cannot connect to the database, dying..")
	}
	log.Println("Connected to the database")

	// Handle creation of template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = true

	// Create Handlers and Repository
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// Create Renderer and Helpers
	render.NewRenderer(&app)
	helpers.NewHelper(&app)

	return db, nil
}
