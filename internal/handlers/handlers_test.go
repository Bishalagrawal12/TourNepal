package handlers

import (
	"net/http"
	"testing"

	"github.com/Atul-Ranjan12/tourism/internal/config"
	"github.com/Atul-Ranjan12/tourism/internal/repository"
)

func TestRepository_ShowReservationCalender(t *testing.T) {
	type fields struct {
		App *config.AppConfig
		DB  repository.DatabaseRepo
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Repository{
				App: tt.fields.App,
				DB:  tt.fields.DB,
			}
			m.ShowReservationCalender(tt.args.w, tt.args.r)
		})
	}
}
