package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Atul-Ranjan12/tourism/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// This package contains all the functionalities and queries of the database
// Currently this application is using the postgres mailing server

// InsertNewUser adds a new registration to the database
func (m *PostgresDBRepo) InsertNewUser(reg models.UserRegistration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (first_name, last_name, email, 
			password, age, gender, access_level, phone_number, mail_verification_code, created_at, updated_at)
		VALUES  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := m.DB.ExecContext(ctx, query,
		reg.FirstName,
		reg.LastName,
		reg.Email,
		reg.HashedPassword,
		reg.Age,
		reg.Gender,
		1,
		reg.PhoneNumber,
		reg.VerificationCode,
		reg.CreatedAt,
		reg.UpdatedAt,
	)
	if err != nil {
		log.Println("Error eexecuting the query")
		return err
	}

	return nil
}

func (m *PostgresDBRepo) UserExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT COUNT(*) FROM users WHERE email = $1
	`
	var numRows int
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&numRows)
	if err != nil {
		log.Println("ERROR: Error executing query to check if the user exists")
		return false, err
	}

	if numRows == 0 {
		return false, nil
	}
	return true, nil
}

func (m *PostgresDBRepo) CheckTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select * from merchants
	`

	_, err := m.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	log.Println("Succesful query execution")
	return nil
}

// Function to Authenticate the User for a sign in
func (m *PostgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password from users WHERE email = $1", email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	// Compare password with my password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println("Invalid Password")
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		log.Println("Unexpected Error")
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// Function to find the user by ID
func (m *PostgresDBRepo) FindUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, phone_number, user_is_verified
		FROM users
		WHERE id = $1
	`
	var user models.User

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.IsVerified,
	)
	if err != nil {
		log.Println("Probable error in query execution: FindUserBID")
		return user, err
	}

	return user, nil
}

func (m *PostgresDBRepo) IncrementVerification(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("Increment Verification Called :: Incrementing to: ", user.IsVerified+1)

	query := `
		UPDATE users
		SET user_is_verified = $1
		WHERE id = $2
	`

	_, err := m.DB.ExecContext(ctx, query, user.IsVerified+1, user.ID)
	if err != nil {
		return err
	}
	return nil
}

// Funciton to get the verification code from the user
func (m *PostgresDBRepo) GetVerificationCode(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var verificationCode int

	query := `
		SELECT mail_verification_code 
		FROM users
		WHERE id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, user.ID)
	err := row.Scan(&verificationCode)
	if err != nil {
		return 0, err
	}

	return verificationCode, err
}

// Function to add the merchant address to the database
func (m *PostgresDBRepo) AddMerchantAddress(address models.MerchantAddress) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO merchant_address (city, state, country, address_line_1, address_line_2, created_at, updated_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := m.DB.QueryContext(ctx, query,
		address.City,
		address.State,
		address.Country,
		address.AddressLine1,
		address.AddressLine2,
		address.CreatedAt,
		address.UpdatedAt,
		address.UserID,
	)
	if err != nil {
		log.Println("Error executing query: Inserting into merchant_address")
		return err
	}

	return nil
}

func (m *PostgresDBRepo) AddMerchantDocuments(docs models.MerchantDocument) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO merchant_documents (document_link, document_id, image, created_at, updated_at, user_id)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	var newID int

	err := m.DB.QueryRowContext(ctx, query,
		docs.DocumentLink,
		docs.DocumentID,
		docs.ImageFile,
		docs.CreatedAt,
		docs.UpdatedAt,
		docs.UserID,
	).Scan(&newID)

	if err != nil {
		log.Println("Error executing query : Inserting into merchant_documents")
		return 0, err
	}
	return newID, nil
}

// Funciton adds user and document references to the merchants table
func (m *PostgresDBRepo) AddMerchant(mer models.MerchantData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO merchants (user_id, address_id, document_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := m.DB.QueryContext(ctx, query,
		mer.UserID,
		mer.AddressID,
		mer.DocumentID,
		mer.CreatedAt,
		mer.UpdatedAt,
	)
	if err != nil {
		log.Println("Error executing query : Inserting into merchants tablel")
		return err
	}
	return nil
}

func (m *PostgresDBRepo) GetAddressIDFromUser(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id FROM merchant_address WHERE user_id = $1
	`
	var addressID int
	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(&addressID)
	if err != nil {
		return 0, err
	}
	return addressID, err
}

// Function to get the userid form the merchatn id
func (m *PostgresDBRepo) GetMerchantIDFromUserID(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id FROM merchants WHERE user_id = $1`

	var merchantID int
	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(&merchantID)
	if err != nil {
		return 0, err
	}

	return merchantID, nil
}

// Add activity to database
func (m *PostgresDBRepo) AddActivityToDatabase(activity models.AddActivityData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO activity (activity_name, activity_description, activity_price, activity_duration, max_size, min_age, phone_num, email, location,
						merchant_id, created_at, updated_at, image)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,$12, $13)
	`
	_, err := m.DB.QueryContext(ctx, query,
		activity.ActivityName,
		activity.ActivityDescription,
		activity.ActivityPrice,
		activity.ActivityDuration,
		activity.MaxGroupSize,
		activity.AgeRestriction,
		activity.PhoneNumber,
		activity.Email,
		activity.Location,
		activity.MerchantID,
		activity.CreatedAt,
		activity.UpdatedAt,
		activity.Image,
	)
	if err != nil {
		log.Println("Error executing query: ", err)
		return err
	}

	return nil
}

//get All activity from the database

func (m *PostgresDBRepo) GetAllActivity(merchantID int) ([]models.AddActivityData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, activity_name, activity_description, activity_price, activity_duration, max_size, min_age, phone_num, email, location,
						merchant_id, created_at, updated_at
		FROM activity
		WHERE merchant_id = $1
	`
	var activities []models.AddActivityData

	rows, err := m.DB.QueryContext(ctx, query, merchantID)
	if err != nil {
		log.Println("Could not execute query: GetAllActivity ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i models.AddActivityData
		err := rows.Scan(
			&i.ActivityID,
			&i.ActivityName,
			&i.ActivityDescription,
			&i.ActivityPrice,
			&i.ActivityDuration,
			&i.MaxGroupSize,
			&i.AgeRestriction,
			&i.PhoneNumber,
			&i.Email,
			&i.Location,
			&i.MerchantID,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning the rows into variables")
			return activities, err
		}

		activities = append(activities, i)
	}
	if err = rows.Err(); err != nil {
		return activities, err
	}
	return activities, nil
}

// Fucntion to get Activity details by ID
func (m *PostgresDBRepo) GetActivityByID(activityID int) (models.AddActivityData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, activity_name, activity_description, activity_price, activity_duration, max_size, min_age, phone_num, email, location,
						merchant_id, created_at, updated_at, image
		FROM activity
		WHERE id = $1
	`
	var i models.AddActivityData

	row := m.DB.QueryRowContext(ctx, query, activityID)
	err := row.Scan(
		&i.ActivityID,
		&i.ActivityName,
		&i.ActivityDescription,
		&i.ActivityPrice,
		&i.ActivityDuration,
		&i.MaxGroupSize,
		&i.AgeRestriction,
		&i.PhoneNumber,
		&i.Email,
		&i.Location,
		&i.MerchantID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Image,
	)
	if err != nil {
		return i, err
	}
	return i, nil
}

// Update the Activity Details in the page
func (m *PostgresDBRepo) UpdateActivityInfo(activityID int, i models.AddActivityData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE activity
		SET activity_name=$1, activity_description=$2, activity_price=$3, activity_duration=$4, max_size=$5, min_age=$6, phone_num=$7, email=$8, location=$9,
			merchant_id=$10, created_at=$11, updated_at=$12
		WHERE id = $13
	`

	_, err := m.DB.QueryContext(ctx, query,
		i.ActivityName,
		i.ActivityDescription,
		i.ActivityPrice,
		i.ActivityDuration,
		i.MaxGroupSize,
		i.AgeRestriction,
		i.PhoneNumber,
		i.Email,
		i.Location,
		i.MerchantID,
		i.CreatedAt,
		i.UpdatedAt,
		activityID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Function to delete activity by id
func (m *PostgresDBRepo) DeleteActivityByID(activityID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM activity WHERE id=$1`

	_, err := m.DB.ExecContext(ctx, query, activityID)
	if err != nil {
		return err
	}
	return nil
}

// Add the bus details to the server
func (m *PostgresDBRepo) AddBusToDatabase(bus models.AddBusData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO bus (bus_name, bus_source, bus_destination, bus_model, bus_no_plate, num_seats, office_pan, office_address, price,
						merchant_id, created_at, updated_at, image)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
`
	_, err := m.DB.QueryContext(ctx, query,
		bus.BusName,
		bus.BusStart,
		bus.BusEnd,
		bus.BusModel,
		bus.BusNumPlate,
		bus.BusNumSeats,
		bus.BusPAN,
		bus.BusAddress,
		bus.Price,
		bus.MerchantID,
		bus.CreatedAt,
		bus.UpdatedAt,
		bus.Image,
	)
	if err != nil {
		log.Println("Error executing query: ", err)
		return err
	}

	return nil
}

func (m *PostgresDBRepo) GetAllBus(merchantID int) ([]models.AddBusData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, bus_name, bus_source, bus_destination, bus_model, bus_no_plate, num_seats, office_pan, office_address, 
			   merchant_id, created_at, updated_at
		FROM bus
		WHERE merchant_id = $1
	`
	var busses []models.AddBusData

	rows, err := m.DB.QueryContext(ctx, query, merchantID)
	if err != nil {
		log.Println("Could not execute query: GetAllBus ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i models.AddBusData
		err := rows.Scan(
			&i.BusID,
			&i.BusName,
			&i.BusStart,
			&i.BusEnd,
			&i.BusModel,
			&i.BusNumPlate,
			&i.BusNumSeats,
			&i.BusPAN,
			&i.BusAddress,
			&i.MerchantID,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning the rows into variables")
			return busses, err
		}

		busses = append(busses, i)
	}
	if err = rows.Err(); err != nil {
		return busses, err
	}
	return busses, nil
}

// Fucntion to get bus details by ID
func (m *PostgresDBRepo) GetBusByID(busID int) (models.AddBusData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, bus_name, bus_source, bus_destination, bus_model, bus_no_plate, num_seats, office_pan, office_address, 
			   merchant_id, created_at, updated_at, price, image
		FROM bus
		WHERE id = $1
	`
	var i models.AddBusData

	row := m.DB.QueryRowContext(ctx, query, busID)
	err := row.Scan(
		&i.BusID,
		&i.BusName,
		&i.BusStart,
		&i.BusEnd,
		&i.BusModel,
		&i.BusNumPlate,
		&i.BusNumSeats,
		&i.BusPAN,
		&i.BusAddress,
		&i.MerchantID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Price,
		&i.Image,
	)
	if err != nil {
		return i, err
	}
	return i, nil
}

// Update the Bus Details in the page
func (m *PostgresDBRepo) UpdateBusInfo(busID int, i models.AddBusData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE bus
		SET bus_name = $1, bus_source = $2, bus_destination = $3, bus_model = $4,
			bus_no_plate = $5, num_seats = $6, office_pan = $7, office_address = $8,
			merchant_id = $9, created_at = $10, updated_at = $11, price = $12
		WHERE id = $13
	`

	_, err := m.DB.QueryContext(ctx, query,
		i.BusName,
		i.BusStart,
		i.BusEnd,
		i.BusModel,
		i.BusNumPlate,
		i.BusNumSeats,
		i.BusPAN,
		i.BusAddress,
		i.MerchantID,
		i.CreatedAt,
		i.UpdatedAt,
		i.Price,
		busID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Function to delete bus by id
func (m *PostgresDBRepo) DeleteBusByID(busID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM bus WHERE id=$1`

	_, err := m.DB.ExecContext(ctx, query, busID)
	if err != nil {
		return err
	}
	return nil
}

// Function to make a new Bus Reservation
func (m *PostgresDBRepo) MakeBusReservation(busRes models.BusReservationData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO bus_reservations (bus_id, first_name, last_name, reservation_date, num_passangers, 
			start, stop, phone_number, email, processed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := m.DB.QueryContext(ctx, query,
		busRes.BusID,
		busRes.FirstName,
		busRes.LastName,
		busRes.ReservationDate,
		busRes.NumPassengers,
		busRes.From,
		busRes.Stop,
		busRes.PhoneNumber,
		busRes.Email,
		0,
		busRes.CreatedAt,
		busRes.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// Function to get all the bus Reservations from the database
func (m *PostgresDBRepo) GetAllBusReservations(showNew bool, mid int) ([]models.BusReservationData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var busRes []models.BusReservationData
	var query string
	var processed int

	if showNew {
		processed = 0
	} else {
		processed = 1
	}

	query = fmt.Sprintf(`
			SELECT  br.id, bus_id, first_name, last_name, reservation_date, num_passangers, 
					start, stop, phone_number, email, bus_name, bus_no_plate
			FROM bus_reservations br
			LEFT JOIN bus b ON (br.bus_id = b.id)
			WHERE br.processed = %d AND merchant_id = $1
			ORDER BY br.reservation_date ASC
		`, processed)

	rows, err := m.DB.QueryContext(ctx, query, mid)
	if err != nil {
		log.Println("Cannot execute this query select from bus reservations table")
		return busRes, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.BusReservationData
		err := rows.Scan(
			&i.ReservationID,
			&i.BusID,
			&i.FirstName,
			&i.LastName,
			&i.ReservationDate,
			&i.NumPassengers,
			&i.From,
			&i.Stop,
			&i.PhoneNumber,
			&i.Email,
			&i.Bus.BusName,
			&i.Bus.BusNumPlate,
		)
		if err != nil {
			log.Println("Error scanning the rows into the variables")
			return busRes, err
		}

		busRes = append(busRes, i)
	}

	if err = rows.Err(); err != nil {
		return busRes, err
	}
	return busRes, nil
}

// Get One reeservation information from ID
func (m *PostgresDBRepo) GetReservationByID(id int) (models.BusReservationData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var i models.BusReservationData

	query := `
			SELECT  br.id, bus_id, first_name, last_name, reservation_date, num_passangers, 
					start, stop, phone_number, email, bus_name, bus_no_plate
			FROM bus_reservations br
			LEFT JOIN bus b ON (br.bus_id = b.id)
			WHERE br.id = $1
			ORDER BY br.reservation_date ASC
		`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.ReservationID,
		&i.BusID,
		&i.FirstName,
		&i.LastName,
		&i.ReservationDate,
		&i.NumPassengers,
		&i.From,
		&i.Stop,
		&i.PhoneNumber,
		&i.Email,
		&i.Bus.BusName,
		&i.Bus.BusNumPlate,
	)
	if err != nil {
		return i, err
	}

	return i, nil
}

// Function to process a reseravtion
func (m *PostgresDBRepo) ProcessReservation(table string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		UPDATE %s
		SET processed = 1
		WHERE id = $1
	`, table)

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// Function to update bus reservation
func (m *PostgresDBRepo) UpdateBusReservation(res models.BusReservationData, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE bus_reservations
		SET start = $1, stop = $2, phone_number = $3, email = $4
		WHERE id = $5
	`
	_, err := m.DB.ExecContext(ctx, query, res.From, res.Stop, res.PhoneNumber, res.Email, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostgresDBRepo) DeleteBusReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE from bus_reservations WHERE id = $1
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// Function to add a new Hotel to the database
func (m *PostgresDBRepo) AddNewHotelRoom(hotel models.HotelRoom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO hotel_room (hotel_name, hotel_room_name, hotel_type, hotel_address, hotel_pan, hotel_num_room, hotel_phone_1, hotel_phone_2, merchant_id, hotel_description, price, created_at, updated_at, image)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := m.DB.ExecContext(ctx, query,
		hotel.HotelName,
		hotel.HotelRoomName,
		hotel.HotelType,
		hotel.HotelAddress,
		hotel.HotelPAN,
		hotel.HotelNumRooms,
		hotel.HotelPhone1,
		hotel.HotelPhone2,
		hotel.MerchantID,
		hotel.HotelRoomDescription,
		hotel.Price,
		hotel.CreatedAt,
		hotel.UpdatedAt,
		hotel.Image,
	)
	if err != nil {
		return err
	}
	return nil
}

// Funciton to get all the hotel reservations:
func (m *PostgresDBRepo) GetAllHotelRooms(merchantID int) ([]models.HotelRoom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, hotel_name, hotel_room_name, hotel_type, hotel_address, hotel_pan, hotel_num_room, hotel_phone_1, hotel_phone_2, merchant_id, hotel_description, created_at, updated_at
	FROM hotel_room 
	WHERE merchant_id = $1
	`
	var rooms []models.HotelRoom

	rows, err := m.DB.QueryContext(ctx, query, merchantID)
	if err != nil {
		log.Println("Could not execute this query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i models.HotelRoom
		err := rows.Scan(
			&i.HotelID,
			&i.HotelName,
			&i.HotelRoomName,
			&i.HotelType,
			&i.HotelAddress,
			&i.HotelPAN,
			&i.HotelNumRooms,
			&i.HotelPhone1,
			&i.HotelPhone2,
			&i.MerchantID,
			&i.HotelRoomDescription,
			&i.CreatedAt,
			&i.UpdatedAt,
		)

		if err != nil {
			log.Println("Error scanning the rows into variables")
			return rooms, err
		}

		rooms = append(rooms, i)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, err
}

// Function to get a room by ID
func (m *PostgresDBRepo) GetRoomByID(id int) (models.HotelRoom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, hotel_name, hotel_room_name, hotel_type, 
	hotel_address, hotel_pan, hotel_num_room, hotel_phone_1, 
	hotel_phone_2, merchant_id, hotel_description, 
	created_at, updated_at, price, image
	FROM hotel_room 
	WHERE id = $1 
	`
	var i models.HotelRoom

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.HotelID,
		&i.HotelName,
		&i.HotelRoomName,
		&i.HotelType,
		&i.HotelAddress,
		&i.HotelPAN,
		&i.HotelNumRooms,
		&i.HotelPhone1,
		&i.HotelPhone2,
		&i.MerchantID,
		&i.HotelRoomDescription,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Price,
		&i.Image,
	)
	if err != nil {
		return i, err
	}
	return i, err
}

// Function to Delete a Bus
func (m *PostgresDBRepo) DeleteRoomByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	DELETE from hotel_room WHERE id = $1
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return err
}

// Function to Update a Bus
func (m *PostgresDBRepo) UpdateRoom(hotel models.HotelRoom, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE hotel_room
	SET hotel_name = $1, hotel_room_name = $2, hotel_type = $3, hotel_address = $4, hotel_pan = $5, hotel_num_room = $6, hotel_phone_1 = $7, 
		hotel_phone_2 = $8, merchant_id = $9, hotel_description = $10, created_at = $11, updated_at = $12, price = $13
	WHERE id = $14
	`
	_, err := m.DB.ExecContext(ctx, query,
		hotel.HotelName,
		hotel.HotelRoomName,
		hotel.HotelType,
		hotel.HotelAddress,
		hotel.HotelPAN,
		hotel.HotelNumRooms,
		hotel.HotelPhone1,
		hotel.HotelPhone2,
		hotel.MerchantID,
		hotel.HotelRoomDescription,
		hotel.CreatedAt,
		hotel.UpdatedAt,
		hotel.Price,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

// Adds a reservation to the database
func (m *PostgresDBRepo) MakeHotelReservation(res models.HotelRoomReservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO hotel_reservations (first_name, last_name, hotel_id, reservation_date_start, 
		reservation_date_end, num_people, phone_number, email, processed, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := m.DB.QueryContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.HotelID,
		res.ResDateStart,
		res.ResDateEnd,
		res.NumPeople,
		res.PhoneNumber,
		res.Email,
		0,
		res.CreatedAt,
		res.UpdatedAt,
	)
	if err != nil {
		return nil
	}
	return nil
}

// Gets all the Hotel Reservations From the Database
func (m *PostgresDBRepo) GetAllHotelReservations(showNew bool, mid int) ([]models.HotelRoomReservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res []models.HotelRoomReservation
	var query string
	var processed int

	if showNew {
		processed = 0
	} else {
		processed = 1
	}

	query = fmt.Sprintf(`
		SELECT r.id, first_name, last_name, hotel_id, reservation_date_start, 
		reservation_date_end, num_people, phone_number, email, hotel_name, hotel_room_name
		FROM hotel_reservations r
		LEFT JOIN hotel_room h ON (r.hotel_id = h.id)
		WHERE r.processed = %d AND merchant_id = $1
		ORDER BY r.reservation_date_end ASC
	`, processed)

	rows, err := m.DB.QueryContext(ctx, query, mid)
	if err != nil {
		log.Println("Cannot execute the select query: ")
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.HotelRoomReservation
		err := rows.Scan(
			&i.ReservationID,
			&i.FirstName,
			&i.LastName,
			&i.HotelID,
			&i.ResDateStart,
			&i.ResDateEnd,
			&i.NumPeople,
			&i.PhoneNumber,
			&i.Email,
			&i.Room.HotelName,
			&i.Room.HotelRoomName,
		)
		if err != nil {
			log.Println("Error scanning the rows into variables: GetAllHotelResrevations")
			return res, err
		}
		res = append(res, i)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

// Get One Hotel Reservation By ID
func (m *PostgresDBRepo) GetHotelReseravtionByID(id int) (models.HotelRoomReservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var i models.HotelRoomReservation

	query := `
		SELECT r.id, first_name, last_name, hotel_id, reservation_date_start, 
		reservation_date_end, num_people, phone_number, email, hotel_name, hotel_room_name
		FROM hotel_reservations r
		LEFT JOIN hotel_room h ON (r.hotel_id = h.id)
		WHERE r.id = $1
		ORDER BY r.reservation_date_end ASC
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.ReservationID,
		&i.FirstName,
		&i.LastName,
		&i.HotelID,
		&i.ResDateStart,
		&i.ResDateEnd,
		&i.NumPeople,
		&i.PhoneNumber,
		&i.Email,
		&i.Room.HotelName,
		&i.Room.HotelRoomName,
	)
	if err != nil {
		return i, err
	}

	return i, nil
}

// Funciton to Update a Reservation :: Hotel
func (m *PostgresDBRepo) UpdateHotelReservation(res models.HotelRoomReservation, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE hotel_reservations
		SET num_people = $1, 
			phone_number = $2, email = $3
		WHERE id = $4
	`
	_, err := m.DB.ExecContext(ctx, query, res.NumPeople, res.PhoneNumber, res.Email, id)
	if err != nil {
		return err
	}
	return nil
}

// Adds an activity reservation to the database
func (m *PostgresDBRepo) MakeActivityReservation(res models.ActivityReservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO activity_reservations (first_name, last_name, activity_id, reservation_date, num_people, phone_number, email, processed, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := m.DB.QueryContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.ActivityID,
		res.ResDate,
		res.NumPeople,
		res.PhoneNumber,
		res.Email,
		0,
		res.CreatedAt,
		res.UpdatedAt,
	)
	if err != nil {
		return nil
	}
	return nil
}

// Gets all the Activity Reservations From the Database
func (m *PostgresDBRepo) GetAllActivityReservations(showNew bool, mid int) ([]models.ActivityReservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res []models.ActivityReservation
	var query string
	var processed int

	if showNew {
		processed = 0
	} else {
		processed = 1
	}

	query = fmt.Sprintf(`
		SELECT ar.id, first_name, last_name, activity_id, reservation_date, num_people, phone_number, ar.email, activity_name, location
		FROM activity_reservations ar
		LEFT JOIN activity a ON (ar.activity_id = a.id)
		WHERE ar.processed = %d AND merchant_id = $1
		ORDER BY ar.reservation_date ASC
	`, processed)

	rows, err := m.DB.QueryContext(ctx, query, mid)
	if err != nil {
		log.Println("Cannot execute the select query: ")
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.ActivityReservation
		err := rows.Scan(
			&i.ReservationID,
			&i.FirstName,
			&i.LastName,
			&i.ActivityID,
			&i.ResDate,
			&i.NumPeople,
			&i.PhoneNumber,
			&i.Email,
			&i.Activity.ActivityName,
			&i.Activity.Location,
		)
		if err != nil {
			log.Println("Error scanning the rows into variables: GetAllActivityResrevations")
			return res, err
		}
		res = append(res, i)
	}

	if err = rows.Err(); err != nil {
		return res, err
	}
	return res, nil
}

// Get One Activity Reservation By ID
func (m *PostgresDBRepo) GetActivityReseravtionByID(id int) (models.ActivityReservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var i models.ActivityReservation

	query := `
		SELECT a.id, first_name, last_name, activity_id, reservation_date, num_people, phone_number, ar.email, activity_name, location
		FROM activity_reservations ar
		LEFT JOIN activity a ON (ar.activity_id = a.id)
		WHERE ar.id = $1
		ORDER BY ar.reservation_date ASC
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&i.ReservationID,
		&i.FirstName,
		&i.LastName,
		&i.ActivityID,
		&i.ResDate,
		&i.NumPeople,
		&i.PhoneNumber,
		&i.Email,
		&i.Activity.ActivityName,
		&i.Activity.Location,
	)
	if err != nil {
		return i, err
	}

	return i, nil
}

// Funciton to Update a Reservation :: Recreational Activity
func (m *PostgresDBRepo) UpdateActivityReservation(res models.ActivityReservation, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE activity_reservations
		SET num_people = $1, 
			phone_number = $2, email = $3
		WHERE id = $4
	`
	_, err := m.DB.ExecContext(ctx, query, res.NumPeople, res.PhoneNumber, res.Email, id)
	if err != nil {
		return err
	}
	return nil
}

// Delete a Reseravtion In General
func (m *PostgresDBRepo) DeleteReservation(tableName string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, tableName)

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// 1. Create a function to post the reviews
func (m *PostgresDBRepo) InsertReview(r models.Review) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO reviews (review_cat, review_item, first_name, last_name, email, phone, stars, review, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := m.DB.ExecContext(ctx, query,
		r.CategoryID,
		r.ItemID,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		r.Stars,
		r.Review,
		r.CreatedAt,
		r.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// Function to get the particular Item Review from the database:
func (m *PostgresDBRepo) GetItemReviews(catID int, itemID int) ([]models.ItemReview, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var itemReview []models.ItemReview

	query := `
		SELECT first_name, last_name, phone, email, stars, review 
		FROM reviews 
		WHERE review_cat = $1 AND review_item = $2
	`
	rows, err := m.DB.QueryContext(ctx, query, catID, itemID)
	if err != nil {
		log.Println("Could not execute query")
		return itemReview, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.ItemReview
		err := rows.Scan(
			&i.FirstName,
			&i.LastName,
			&i.Phone,
			&i.Email,
			&i.Stars,
			&i.Review,
		)
		if err != nil {
			log.Println("Error scanning the rows into variables")
			return itemReview, err
		}

		itemReview = append(itemReview, i)
	}
	if err = rows.Err(); err != nil {
		return itemReview, err
	}
	return itemReview, nil
}



//function to get the number of Activity reservations in each day of the month

func (m *PostgresDBRepo) GetActivityReservationByMonth(month int, mid int) (models.ReservationCalendar, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.ReservationCalendar

	now := time.Now()
	targetMonth := time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, now.Location())

	// Query to get the reservations for the given month and activity_id
	query := `
        SELECT DATE_TRUNC('day', activity_reservations.reservation_date) AS day,
		COUNT(activity_reservations.id) AS num_reservations
		FROM activity_reservations
		FULL OUTER JOIN activity ON activity_reservations.activity_id = activity.id
	    WHERE DATE_TRUNC('month', activity_reservations.reservation_date) = DATE_TRUNC('month', $2::DATE)
		AND activity.merchant_id = $1
		GROUP BY DATE_TRUNC('day', activity_reservations.reservation_date)
		ORDER BY day ASC
    `

	rows, err := m.DB.QueryContext(ctx, query, mid, targetMonth)
	if err != nil {
		log.Println("Cannot execute the select query: ", err)
		return res, err
	}
	defer rows.Close()

	// Create a map to hold the number of reservations for each day of the month
	reservationMap := make(map[int]int)
	for i := 1; i <= 31; i++ {
		reservationMap[i] = 0
	}

	// Loop through the rows of the result set and update the reservationMap
	for rows.Next() {
		var day time.Time
		var numReservations int
		err := rows.Scan(&day, &numReservations)
		if err != nil {
			log.Println("Error scanning row: ", err)
			return res, err
		}
		reservationMap[day.Day()] = numReservations
	}

	// Create a ReservationCalendar object for the given month using the reservationMap
	var rc models.ReservationCalendar
	rc.Month = time.Month(month).String()
	rc.Reservations = reservationMap

	// Append the ReservationCalendar object to the result slice
	res = rc

	return res, nil

}

// function to get the number of Hotel reservations in each day of the month
func (m *PostgresDBRepo) GetHotelReservationByMonth(month int, mid int) (models.ReservationCalendar, error) {
	// Convert the month integer to a string with leading zeros

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.ReservationCalendar

	now := time.Now()
	targetMonth := time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, now.Location())

	// Query to get the reservations for the given month and activity_id
	query := `
        SELECT DATE_TRUNC('day', hotel_reservations.reservation_date_start) AS day,
		COUNT(hotel_reservations.id) AS num_reservations
		FROM hotel_reservations
		FULL OUTER JOIN hotel_room ON hotel_reservations.hotel_id = hotel_room.id
	    WHERE DATE_TRUNC('month', hotel_reservations.reservation_date_start) = DATE_TRUNC('month', $2::DATE)
		AND hotel_room.merchant_id = $1
		GROUP BY DATE_TRUNC('day', hotel_reservations.reservation_date_start)
		ORDER BY day ASC
    `

	rows, err := m.DB.QueryContext(ctx, query, mid, targetMonth)
	if err != nil {
		log.Println("Cannot execute the select query: ", err)
		return res, err
	}
	defer rows.Close()

	// Create a map to hold the number of reservations for each day of the month
	reservationMap := make(map[int]int)
	for i := 1; i <= 31; i++ {
		reservationMap[i] = 0
	}

	// Loop through the rows of the result set and update the reservationMap
	for rows.Next() {
		var day time.Time
		var numReservations int
		err := rows.Scan(&day, &numReservations)
		if err != nil {
			log.Println("Error scanning row: ", err)
			return res, err
		}
		reservationMap[day.Day()] = numReservations
	}

	// Create a ReservationCalendar object for the given month using the reservationMap
	var rc models.ReservationCalendar
	rc.Month = time.Month(month).String()
	rc.Reservations = reservationMap

	// Append the ReservationCalendar object to the result slice
	res = rc

	return res, nil

}

// function to get the number of Bus reservations in each day of the month
func (m *PostgresDBRepo) GetBusReservationByMonth(month int, mid int) (models.ReservationCalendar, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.ReservationCalendar

	now := time.Now()
	targetMonth := time.Date(now.Year(), time.Month(month), 1, 0, 0, 0, 0, now.Location())

	// Query to get the reservations for the given month and activity_id
	query := `
        SELECT DATE_TRUNC('day', bus_reservations.reservation_date) AS day,
		COUNT(bus_reservations.id) AS num_reservations
		FROM bus_reservations
		FULL OUTER JOIN bus ON bus_reservations.bus_id = bus.id
	    WHERE DATE_TRUNC('month', bus_reservations.reservation_date) = DATE_TRUNC('month', $2::DATE)
		AND bus.merchant_id = $1
		GROUP BY DATE_TRUNC('day', bus_reservations.reservation_date)
		ORDER BY day ASC
    `

	rows, err := m.DB.QueryContext(ctx, query, mid, targetMonth)
	if err != nil {
		log.Println("Cannot execute the select query: ", err)
		return res, err
	}
	defer rows.Close()

	// Create a map to hold the number of reservations for each day of the month
	reservationMap := make(map[int]int)
	for i := 1; i <= 31; i++ {
		reservationMap[i] = 0
	}

	// Loop through the rows of the result set and update the reservationMap
	for rows.Next() {
		var day time.Time
		var numReservations int
		err := rows.Scan(&day, &numReservations)
		if err != nil {
			log.Println("Error scanning row: ", err)
			return res, err
		}
		reservationMap[day.Day()] = numReservations
	}

	// Create a ReservationCalendar object for the given month using the reservationMap
	var rc models.ReservationCalendar
	rc.Month = time.Month(month).String()
	rc.Reservations = reservationMap

	// Append the ReservationCalendar object to the result slice
	res = rc

	return res, nil

}

// function to get the Details of of Bus reservations in each day of the month
func (m *PostgresDBRepo) GetBusReservationByDay(date time.Time, mid int) ([]models.BusReservationData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res []models.BusReservationData

	// Query to get the reservations for the given month and reservation_id
	query := `
        SELECT bus_reservations.id
		FROM bus_reservations
		FULL OUTER JOIN bus ON bus_reservations.bus_id = bus.id
	    WHERE bus_reservations.reservation_date::date = $2::DATE
		AND bus.merchant_id = $1
		GROUP BY bus_reservations.id,DATE_TRUNC('day', bus_reservations.reservation_date)
		ORDER BY bus_reservations.reservation_date ASC
    `

	rows, err := m.DB.QueryContext(ctx, query, mid, date)
	if err != nil {
		log.Println("Cannot execute the select query: ", err)
		return res, err
	}
	defer rows.Close()

	// Loop through the rows of the result set and update the res
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println("Error scanning row: ", err)
			return res, err
		}
		busReservation, err := m.GetReservationByID(id)
		if err != nil {
			log.Println("Error getting reservation by ID: ", err)
			return res, err
		}
		res = append(res, busReservation)

	}

	return res, nil
}

// function to get the Details of of Hotel reservations in each day of the month
func (m *PostgresDBRepo) GetHotelReservationByDay(date time.Time, mid int) ([]models.HotelRoomReservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res []models.HotelRoomReservation

	// Query to get the reservations for the given month and reservation ID
	query := `
        SELECT hotel_reservations.id
		FROM hotel_reservations
		FULL OUTER JOIN hotel_room ON hotel_reservations.hotel_id = hotel_room.id
	    WHERE hotel_reservations.reservation_date_start::date = $2::DATE
		AND hotel_room.merchant_id = $1
		GROUP BY hotel_reservations.id,DATE_TRUNC('day', hotel_reservations.reservation_date_start)
		ORDER BY hotel_reservations.reservation_date_start ASC
    `

	rows, err := m.DB.QueryContext(ctx, query, mid, date)
	if err != nil {
		log.Println("Cannot execute the select query: ", err)
		return res, err
	}
	defer rows.Close()

	// Loop through the rows of the result set and update the reservation
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println("Error scanning row: ", err)
			return res, err
		}
		hotelReservation, err := m.GetHotelReseravtionByID(id)
		if err != nil {
			log.Println("Error getting reservation by ID: ", err)
			return res, err
		}
		res = append(res, hotelReservation)

	}

	return res, nil
}

// function to get the Details of of Activity reservations in each day of the month
func (m *PostgresDBRepo) GetActivityReservationByDay(date time.Time, mid int) ([]models.ActivityReservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res []models.ActivityReservation

	// Query to get the reservations for the given month and activity_reservation_id
	query := `
        SELECT activity_reservations.id
		FROM activity_reservations
		FULL OUTER JOIN activity ON activity_reservations.activity_id = activity.id
	    WHERE activity_reservations.reservation_date::date = $2::DATE
		AND activity.merchant_id = $1
		GROUP BY activity_reservations.id,DATE_TRUNC('day', activity_reservations.reservation_date)
		ORDER BY activity_reservations.reservation_date ASC
    `

	rows, err := m.DB.QueryContext(ctx, query, mid, date)
	if err != nil {
		log.Println("Cannot execute the select query: ", err)
		return res, err
	}
	defer rows.Close()

	// Loop through the rows of the result set and update the reservation
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Println("Error scanning row: ", err)
			return res, err
		}
		activityReservation, err := m.GetActivityReseravtionByID(id)
		if err != nil {
			log.Println("Error getting reservation by ID: ", err)
			return res, err
		}
		res = append(res, activityReservation)

	}

	return res, nil
}

// fucntion to get total reservation count for hotels

func (m *PostgresDBRepo) GetTotalReservationCountHotel(mid int) (int,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res int
	// Query to get the reservations count 
	query := `
        SELECT COUNT(hr.id)
		FROM hotel_reservations hr
		LEFT JOIN hotel_room h ON (hr.hotel_id = h.id)
		WHERE merchant_id=$1
    `

	row := m.DB.QueryRowContext(ctx, query,mid)
	err := row.Scan(&res)
	if err != nil {
		log.Println("ERROR: Error executing query to check if the user exists")
		return 0, err
	}
	return res, nil
}

// fucntion to get processed reservation count for hotels

func (m *PostgresDBRepo) GetProcessedReservationCountHotel(mid int) (int,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res int
	// Query to get the reservations count 
	query := `
        SELECT COUNT(hr.id) 
		FROM hotel_reservations hr
		LEFT JOIN hotel_room h ON (hr.hotel_id = h.id)
		WHERE processed =1 AND merchant_id=$1
    `
	row := m.DB.QueryRowContext(ctx, query,mid)
	err := row.Scan(&res)
	if err != nil {
		log.Println("ERROR: Error executing query to check if the user exists")
		return 0, err
	}
	return res, nil
}

// fucntion to get total reservation count for bus

func (m *PostgresDBRepo) GetTotalReservationCountBus(mid int) (int,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res int
	// Query to get the reservations count 
	query := `
        SELECT COUNT(br.id)
		FROM bus_reservations br
		LEFT JOIN bus b ON (br.bus_id = b.id)
		WHERE merchant_id=$1
    `

	row := m.DB.QueryRowContext(ctx, query,mid)
	err := row.Scan(&res)
	if err != nil {
		log.Println("ERROR: Error executing query to check if the user exists")
		return 0, err
	}
	return res, nil
}

// fucntion to get processed reservation count for bus

func (m *PostgresDBRepo) GetProcessedReservationCountBus(mid int) (int,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res int
	// Query to get the reservations count 
	query := `
        SELECT COUNT(br.id)
		FROM bus_reservations br
		LEFT JOIN bus b ON (br.bus_id = b.id)
		WHERE processed =1 AND merchant_id=$1
    `
	row := m.DB.QueryRowContext(ctx, query,mid)
	err := row.Scan(&res)
	if err != nil {
		log.Println("ERROR: Error executing query to check if the user exists")
		return 0, err
	}
	return res, nil
}

// fucntion to get total reservation count for activity

func (m *PostgresDBRepo) GetTotalReservationCountActivity(mid int) (int,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res int
	// Query to get the reservations count 
	query := `
        SELECT COUNT(ar.id) 
		FROM activity_reservations ar
		LEFT JOIN activity a ON (ar.activity_id = a.id)
		WHERE merchant_id=$1
    `

	row := m.DB.QueryRowContext(ctx, query,mid)
	err := row.Scan(&res)
	if err != nil {
		log.Println("ERROR: Error executing query to check if the user exists")
		return 0, err
	}
	return res, nil
}

// fucntion to get processed reservation count for activity

func (m *PostgresDBRepo) GetProcessedReservationCountActivity(mid int) (int,error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res int
	// Query to get the reservations count 
	query := `
        SELECT COUNT(ar.id)
		FROM activity_reservations ar
		LEFT JOIN activity a ON (ar.activity_id = a.id)
		WHERE processed =1 AND merchant_id=$1
    `
	row := m.DB.QueryRowContext(ctx, query,mid)
	err := row.Scan(&res)
	if err != nil {
		log.Println("ERROR: Error executing query to check if the user exists")
		return 0, err
	}
	return res, nil
}
