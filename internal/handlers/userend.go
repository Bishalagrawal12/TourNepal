package handlers

import (
	"log"
	"net/http"

	"github.com/Atul-Ranjan12/tourism/internal/models"
	"github.com/Atul-Ranjan12/tourism/internal/render"
)

func (m *Repository) ShowHome(w http.ResponseWriter, r *http.Request) {
	allHotels, err := m.DB.GetTopHotels(4)
	if err != nil {
		log.Println(err)
		return
	}
	allBus, err := m.DB.GetTopBus(5)
	if err != nil {
		log.Println(err)
		return
	}
	allActivity, err := m.DB.GetTopActivity(5)
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]interface{})
	log.Println("Length of allHotels is: ", len(allHotels))

	// Putting the values on the data variable
	data["hotels"] = allHotels
	data["bus"] = allBus
	data["activity"] = allActivity

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
