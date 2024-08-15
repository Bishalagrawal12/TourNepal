package repository

import (
	"time"

	"github.com/Atul-Ranjan12/tourism/internal/models"
)

type DatabaseRepo interface {
	InsertNewUser(reg models.UserRegistration) error
	CheckTable() error
	Authenticate(email, testPassword string) (int, string, error)
	FindUserByID(id int) (models.User, error)
	UserExists(email string) (bool, error)
	GetVerificationCode(user models.User) (int, error)
	IncrementVerification(user models.User) error
	AddMerchantAddress(address models.MerchantAddress) error
	AddMerchantDocuments(docs models.MerchantDocument) (int, error)
	AddMerchant(mer models.MerchantData) error
	GetAddressIDFromUser(userID int) (int, error)
	GetMerchantIDFromUserID(userID int) (int, error)

	// Bus Basic Fucntions
	AddBusToDatabase(bus models.AddBusData) error

	// Activity Basic Functions
	AddActivityToDatabase(activity models.AddActivityData) error
	GetAllActivity(merchantID int) ([]models.AddActivityData, error)
	GetActivityByID(activityID int) (models.AddActivityData, error)
	UpdateActivityInfo(activityID int, i models.AddActivityData) error
	DeleteActivityByID(activityID int) error

	GetAllBus(merchantID int) ([]models.AddBusData, error)
	GetBusByID(busID int) (models.AddBusData, error)
	UpdateBusInfo(busID int, i models.AddBusData) error
	DeleteBusByID(busID int) error
	MakeBusReservation(busRes models.BusReservationData) error

	// Bus Reservation Functions
	GetAllBusReservations(showNew bool, mid int) ([]models.BusReservationData, error)
	GetReservationByID(id int) (models.BusReservationData, error)
	ProcessReservation(table string, id int) error
	UpdateBusReservation(res models.BusReservationData, id int) error
	DeleteBusReservation(id int) error
	AddNewHotelRoom(hotel models.HotelRoom) error
	GetAllHotelRooms(merchantID int) ([]models.HotelRoom, error)
	GetRoomByID(id int) (models.HotelRoom, error)
	DeleteRoomByID(id int) error
	UpdateRoom(hotel models.HotelRoom, id int) error

	// Hotel Room Reservation Functions
	MakeHotelReservation(res models.HotelRoomReservation) error
	GetAllHotelReservations(showNew bool, mid int) ([]models.HotelRoomReservation, error)
	GetHotelReseravtionByID(id int) (models.HotelRoomReservation, error)
	UpdateHotelReservation(res models.HotelRoomReservation, id int) error

	//Recreational Activity Reservation Functions
	MakeActivityReservation(res models.ActivityReservation) error
	GetAllActivityReservations(showNew bool, mid int) ([]models.ActivityReservation, error)
	GetActivityReseravtionByID(id int) (models.ActivityReservation, error)
	UpdateActivityReservation(res models.ActivityReservation, id int) error

	// In general Delete Reservation
	DeleteReservation(tableName string, id int) error

	// Uptility Fucntion to add to the reveiws
	InsertReview(r models.Review) error

	// Functions to get all the item reveiws:
	GetItemReviews(catID int, itemID int) ([]models.ItemReview, error)

	//Month wise retrieval of reservations
	GetActivityReservationByMonth(month int, mid int) (models.ReservationCalendar, error)
	GetHotelReservationByMonth(month int, mid int) (models.ReservationCalendar, error)
	GetBusReservationByMonth(month int, mid int) (models.ReservationCalendar, error)

	//DayWise retrieval of Reservation Data
	GetBusReservationByDay(day time.Time, mid int) ([]models.BusReservationData,error)
	GetHotelReservationByDay(day time.Time, mid int) ([]models.HotelRoomReservation,error)
	GetActivityReservationByDay(day time.Time, mid int) ([]models.ActivityReservation,error)


	//dashboard functions for reservation count
	GetTotalReservationCountHotel(mid int) (int, error)
	GetProcessedReservationCountHotel(mid int)(int, error)

	GetTotalReservationCountBus(mid int) (int, error)
	GetProcessedReservationCountBus(mid int)(int, error)

	GetTotalReservationCountActivity(mid int) (int, error)
	GetProcessedReservationCountActivity(mid int)(int, error)

	// Get All Assets in the system
	GetAllHotelRoomsInSystem() ([]models.HotelRoom, error)
	GetAllBusInSystem() ([]models.AddBusData, error)
	GetAllActivityInSystem() ([]models.AddActivityData, error)

	// Funcitons for the front en
	GetTopHotels(topn int) ([]models.HotelRoom, error)
	GetTopBus(topn int) ([]models.AddBusData, error)
	GetTopActivity(topn int) ([]models.AddActivityData, error)
}

