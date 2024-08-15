package config

import (
	"log"
	"text/template"

	"github.com/Atul-Ranjan12/tourism/internal/models"
	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.ConfirmationMailData
}
