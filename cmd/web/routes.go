package main

import (
	"net/http"

	"github.com/Atul-Ranjan12/tourism/internal/config"
	"github.com/Atul-Ranjan12/tourism/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Set up Multiplexer configuration
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	// ---------------------User Side Routes -----------------------

	mux.Get("/", handlers.Repo.ShowHome)

	// ---------------------Admin Side Routes -----------------------

	// Login Page Routes
	mux.Get("/login", handlers.Repo.ShowLogin)
	mux.Post("/login", handlers.Repo.PostLogin)

	// SignUp Page Routes
	mux.Get("/signup", handlers.Repo.ShowSignUp)
	mux.Post("/signup", handlers.Repo.PostSignUp)

	// User logout
	mux.Get("/logout", handlers.Repo.Logout)

	// Reservation Section
	// 1. Bus Reservation
	mux.Get("/make-bus-reservation", handlers.Repo.MakeBusReservation)
	mux.Post("/make-bus-reservation", handlers.Repo.PostMakeBusReservation)

	// 2. Hotel Room Reservation
	mux.Get("/make-hotel-reservation", handlers.Repo.ShowMakeHotelReservation)
	mux.Post("/make-hotel-reservation", handlers.Repo.PostShowMakeHotelReservation)

	// 3.Recreational Activity Reservation
	mux.Get("/make-activity-reservation", handlers.Repo.ShowMakeActivityReservation)
	mux.Post("/make-activity-reservation", handlers.Repo.PostShowMakeActivityReservation)


	// Utility
	// 1. Write Reviews for Items:
	mux.Get("/add-review", handlers.Repo.ShowAddReview)
	mux.Post("/add-review", handlers.Repo.PostShowAddReview)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// TODO: Set up Routes for Logged In Users
	mux.Route("/merchant", func(mux chi.Router) {
		// Check if user is authenticated
		mux.Use(Auth)

		// Verification Code :: Stage 1
		mux.Get("/{src}/verification", handlers.Repo.ShowAdminVerification)
		mux.Post("/{src}/verification", handlers.Repo.PostShowAdminVerification)

		// Address Verification :: Stage 2
		mux.Get("/{src}/verification-address", handlers.Repo.ShowAdminAddress)
		mux.Post("/{src}/verification-address", handlers.Repo.PostShowAdminAddress)

		// Document Verificaiotn :: Stage 3
		mux.Get("/{src}/verification-documents", handlers.Repo.ShowDocumentsVerification)
		mux.Post("/{src}/verification-documents", handlers.Repo.PostShowDocumentsVerification)

		// Merchant Dashboard Items
		// 1. Merchant Add item Seciotn
		mux.Get("/{src}/merchant-add-items", handlers.Repo.AdminAddMerchantItems)

		// 2. Merchant Reservation Calender Section
		mux.Get("/{src}/reservation-calender", handlers.Repo.ShowReservationCalender)

		// 3. Merchant Dashboard Section
		mux.Get("/{src}/dashboard", handlers.Repo.ShowAdminDashboard)

		// 4. Show All Reservations :: Unprocessed
		mux.Get("/{src}/merchant-show-reservations", handlers.Repo.ShowAllReservations)

		// 5. Show All Reservations :: Processed
		mux.Get("/{src}/merchant-show-reservations-processed", handlers.Repo.ShowReservationsProcessed)

		// 6. Show Reviews Section
		mux.Get("/{src}/merchant-reviews", handlers.Repo.ShowReviewsPage)

		// Add Items Section
		// Bus
		// 1. Merchant Add Bus Section:
		mux.Get("/{src}/add-bus", handlers.Repo.AdminAddBus)
		mux.Post("/{src}/add-bus", handlers.Repo.PostAdminAddBus)

		// 2. Merchant SHow and Edit Bus section
		mux.Get("/{src}/add-bus/{id}", handlers.Repo.AdminShowBus)
		mux.Post("/{src}/add-bus/{id}", handlers.Repo.PostAdminUpdateBus)

		// 3. Delete the bus
		mux.Get("/{src}/add-bus/delete/{id}", handlers.Repo.PostAdminDeleteBus)

		// 4. Show One Single Bus Reservation
		mux.Get("/{src}/merchant-show-reservations/{id}", handlers.Repo.ShowOneReservation)
		mux.Post("/{src}/merchant-show-reservations/{id}", handlers.Repo.PostShowOneReservation)

		// 5. Link to Process the Bus Reservation
		mux.Get("/{src}/merchant-show-reservations/{id}/process", handlers.Repo.ProcessBusReservation)

		// 6. Function to delete the bus reservations
		mux.Get("/{src}/delete-reservation/{id}", handlers.Repo.DeleteBusReservation)

		// Hotel Reservations:
		// 1.  Merchant Add Hotel Resrvation Section
		mux.Get("/{src}/add-hotel", handlers.Repo.AdminAddHotel)
		mux.Post("/{src}/add-hotel", handlers.Repo.PostAdminAddHotel)

		// 2. Mercant Show and Edit the Bus Section
		mux.Get("/{src}/add-hotel/{id}", handlers.Repo.AdminShowOneHotel)
		mux.Post("/{src}/add-hotel/{id}", handlers.Repo.PostAdminShowOneHotel)

		// 3. Delete the reservation
		mux.Get("/{src}/add-hotel/delete/{id}", handlers.Repo.DeleteRoom)

		// 4. Show One Single Hotel Reservation
		mux.Get("/{src}/merchant-show-reservations/{id}/hotel", handlers.Repo.ShowOneHotelReservation)
		mux.Post("/{src}/merchant-show-reservations/{id}/hotel", handlers.Repo.PostShowOneHotelReservation)

		// Show One Single Activity Reservation
		mux.Get("/{src}/merchant-show-reservations/{id}/activity", handlers.Repo.ShowOneActivityReservation)
		mux.Post("/{src}/merchant-show-reservations/{id}/activity", handlers.Repo.PostShowOneActivityReservation)

		// 5.  Link to Process the Hotel Reservation
		mux.Get("/{src}/merchant-show-reservations/{id}/hotel/process", handlers.Repo.ProcessHotelReservations)

		// Link to Process the Activity Reservation
		mux.Get("/{src}/merchant-show-reservations/{id}/activity/process", handlers.Repo.ProcessActivityReservations)

		// 6. Function to Delete the Hotel Reservation
		mux.Get("/{src}/delete-reservation/{id}/hotel", handlers.Repo.DeleteHotelReservation)

		// Function to Delete the Activity Reservation
		mux.Get("/{src}/delete-reservation/{id}/activity", handlers.Repo.DeleteActivityReservation)

		// Recreational Activities
		// 1. merchant add recreational route
		mux.Get("/{src}/add-activity", handlers.Repo.AdminAddRecreationalActivity)
		mux.Post("/{src}/add-activity", handlers.Repo.PostAdminAddRecreationalActivity)

		// 2. merchant show and edit recreational route
		mux.Get("/{src}/add-activity/{id}", handlers.Repo.AdminShowRecreationalActivity)
		mux.Post("/{src}/add-activity/{id}", handlers.Repo.PostAdminUpdateRecreationalActivity)

		// 3. merchant delete recreational activity
		mux.Get("/{src}/add-activity/delete/{id}", handlers.Repo.PostAdminDeleteActivity)

		// Show the Reservation per day
		mux.Get("/{src}/show-reservations-per-day", handlers.Repo.ShowReservationsPerDay)

		//test for count of reservation
		mux.Get("/{src}/getActivityReservations", handlers.Repo.GetActivityByMonth)
		mux.Get("/{src}/getHotelReservations", handlers.Repo.GetHotelByMonth)
		mux.Get("/{src}/getBusReservations", handlers.Repo.GetBusByMonth)

		// test for getting the day wise details
		mux.Get("/{src}/getBusDetailsByDate", handlers.Repo.GetBusByDay)
		mux.Get("/{src}/getHotelDetailsByDate", handlers.Repo.GetHotelByDay)
		mux.Get("/{src}/getActivityDetailsByDate", handlers.Repo.GetActivityByDay)

	})

	return mux
}
