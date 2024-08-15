package models

import "time"

type UserRegistration struct {
	FirstName        string
	LastName         string
	Age              string
	Gender           string
	Email            string
	HashedPassword   string
	PhoneNumber      string
	VerificationCode int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type User struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	IsVerified int
}

// Structure for the mail data
type ConfirmationMailData struct {
	To      string
	From    string
	Subject string
	Content string
}

type MerchantAddress struct {
	City         string
	State        string
	Country      string
	AddressLine1 string
	AddressLine2 string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserID       int
}

type MerchantDocument struct {
	DocumentID   string
	DocumentLink string
	ImageFile    []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserID       int
}

type MerchantData struct {
	UserID     int
	AddressID  int
	DocumentID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Type to add the bus form
type AddBusData struct {
	BusID       int
	MerchantID  int
	BusName     string
	BusModel    string
	BusAddress  string
	BusStart    string
	BusEnd      string
	BusNumSeats int
	BusNumPlate string
	BusPAN      string
	Price       int
	Image       []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Type to add the activity form
type AddActivityData struct {
	ActivityID          int
	MerchantID          int
	ActivityName        string
	ActivityDescription string
	ActivityPrice       int
	ActivityDuration    int
	MaxGroupSize        int
	AgeRestriction      int
	PhoneNumber         string
	Email               string
	Location            string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Image               []byte
}

// Bus Reservation Model
type BusReservationData struct {
	ReservationID   int
	BusID           int
	FirstName       string
	LastName        string
	ReservationDate time.Time
	NumPassengers   int
	From            string
	Stop            string
	PhoneNumber     string
	Email           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Bus             AddBusData
}

// Model for the Hotel/ Hotel Room
type HotelRoom struct {
	HotelID              int
	MerchantID           int
	HotelName            string
	HotelRoomName        string
	HotelType            string
	HotelAddress         string
	HotelPAN             string
	HotelNumRooms        int
	HotelPhone1          string
	HotelPhone2          string
	HotelRoomDescription string
	Price                int
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Image                []byte
}

// Editable Startdate, enddate, numPeople, phone, email
// Structure to make the Hotel Room Reservation
type HotelRoomReservation struct {
	ReservationID int
	HotelID       int
	FirstName     string
	LastName      string
	ResDateStart  time.Time
	ResDateEnd    time.Time
	NumPeople     int
	PhoneNumber   string
	Email         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          HotelRoom
}

//structure to make reservation for an activity

type ActivityReservation struct {
	ReservationID int
	ActivityID    int
	FirstName     string
	LastName      string
	ResDate       time.Time
	NumPeople     int
	PhoneNumber   string
	Email         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Activity      AddActivityData
}

// Utility Model to Post a Item Review
type Review struct {
	ReviewID   int
	CategoryID int
	ItemID     int
	FirstName  string
	LastName   string
	Phone      string
	Email      string
	Stars      int
	Review     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Define a type for the reveiw structure to be displayed in the webpage
type ItemReview struct {
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Stars     int
	Review    string
}

// Define a type to describe all the reviews for a particular merchant item
type CatItemReview struct {
	ItemName string
	Review   []ItemReview
}

type CalenderDay struct {
	Day    int
	NumRes int
}

func (d *CalenderDay) UpdateNumRes(x int) {
	d.NumRes = x
}

type ReservationCalendar struct {
	Month        string
	Reservations map[int]int
}

type ReservationsCount struct{
	TotalHotelRes int
	ProcessedHotelRes float64
	TotalBusRes int
	ProcessedBusRes float64
	TotalActivityRes int
	ProcessedActivityRes float64
}
