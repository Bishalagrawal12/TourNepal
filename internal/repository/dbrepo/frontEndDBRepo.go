package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/Atul-Ranjan12/tourism/internal/models"
)

// This file contains all the database funcitons necessary to be displayed in the
// Front end of this website

// 1. Funciton to get the top 5 hotels
func (m *PostgresDBRepo) GetTopHotels(topn int) ([]models.HotelRoom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT hotel_id, COUNT(*) as num_reservations
	FROM hotel_reservations
	GROUP BY hotel_id
	ORDER BY num_reservations DESC
	LIMIT $1;
	`

	var topHotels []models.HotelRoom

	rows, err := m.DB.QueryContext(ctx, query, topn)
	if err != nil {
		return topHotels, err
	}

	for rows.Next() {
		var id int
		var numRes int
		err := rows.Scan(
			&id,
			&numRes,
		)
		if err != nil {
			return topHotels, err
		}
		hotel, err := m.GetRoomByID(id)
		if err != nil {
			log.Println("ERROR: Error getting the hotel room")
			return topHotels, err
		}
		topHotels = append(topHotels, hotel)
	}

	var length int = len(topHotels)

	if length < topn {
		nq := `
		SELECT DISTINCT hr.id 
		FROM hotel_room hr
		WHERE NOT EXISTS (
  			SELECT 1
  			FROM hotel_reservations hres
  			WHERE hr.id = hres.hotel_id
		)
		`
		nr, err := m.DB.QueryContext(ctx, nq)
		if err != nil {
			return topHotels, err
		}
		var counter int = 0
		for nr.Next() {
			var id int
			err := nr.Scan(&id)
			if err != nil {
				return topHotels, err
			}
			hotel, err := m.GetRoomByID(id)
			if err != nil {
				log.Println("Error getting the room by id")
				return topHotels, err
			}
			if counter < topn-length {
				topHotels = append(topHotels, hotel)
				counter++
			} else {
				break
			}
		}
	}

	if err := rows.Err(); err != nil {
		return topHotels, err
	}
	return topHotels, nil
}

// 2. Function to get the top 5 buses
func (m *PostgresDBRepo) GetTopBus(topn int) ([]models.AddBusData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT bus_id, COUNT(*) as num_reservations
	FROM bus_reservations
	GROUP BY bus_id
	ORDER BY num_reservations DESC
	LIMIT $1;
	`

	var topBus []models.AddBusData

	rows, err := m.DB.QueryContext(ctx, query, topn)
	if err != nil {
		return topBus, err
	}

	for rows.Next() {
		var id int
		var numRes int
		err := rows.Scan(
			&id,
			&numRes,
		)
		if err != nil {
			return topBus, err
		}
		bus, err := m.GetBusByID(id)
		if err != nil {
			log.Println("ERROR: Error getting the hotel room")
			return topBus, err
		}
		topBus = append(topBus, bus)
	}

	var length int = len(topBus)
	if length < topn {
		nq := `
		SELECT DISTINCT b.id 
		FROM bus b 
		WHERE NOT EXISTS (
  			SELECT 1
  			FROM bus_reservations br 
  			WHERE b.id = br.bus_id  
		);
		`
		nr, err := m.DB.QueryContext(ctx, nq)
		if err != nil {
			return topBus, err
		}
		var counter int = 0
		for nr.Next() {
			var id int
			err := nr.Scan(&id)
			if err != nil {
				return topBus, err
			}
			bus, err := m.GetBusByID(id)
			if err != nil {
				log.Println("Error getting the room by id")
				return topBus, err
			}
			if counter < topn {
				topBus = append(topBus, bus)
				counter++
			} else {
				break
			}
		}
	}

	if err := rows.Err(); err != nil {
		return topBus, err
	}
	return topBus, nil
}

func (m *PostgresDBRepo) GetTopActivity(topn int) ([]models.AddActivityData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT activity_id, COUNT(*) as num_reservations
	FROM activity_reservations
	GROUP BY activity_id
	ORDER BY num_reservations DESC
	LIMIT $1;
	`

	var topAc []models.AddActivityData

	rows, err := m.DB.QueryContext(ctx, query, topn)
	if err != nil {
		return topAc, err
	}

	for rows.Next() {
		var id int
		var numRes int
		err := rows.Scan(
			&id,
			&numRes,
		)
		if err != nil {
			return topAc, err
		}
		ac, err := m.GetActivityByID(id)
		if err != nil {
			log.Println("ERROR: Error getting the hotel room")
			return topAc, err
		}
		topAc = append(topAc, ac)
	}

	if len(topAc) < topn {
		nq := `
		SELECT DISTINCT ac.id 
		FROM activity ac
		WHERE NOT EXISTS (
  			SELECT 1
  			FROM activity_reservations ar 
  			WHERE ac.id = ar.activity_id 
		);
		`
		nr, err := m.DB.QueryContext(ctx, nq)
		if err != nil {
			return topAc, err
		}
		var counter int = 0
		for nr.Next() {
			var id int
			err := nr.Scan(&id)
			if err != nil {
				return topAc, err
			}
			activity, err := m.GetActivityByID(id)
			if err != nil {
				log.Println("Error getting the room by id")
				return topAc, err
			}
			if counter < topn {
				topAc = append(topAc, activity)
				counter++
			} else {
				break
			}
		}
	}

	if err := rows.Err(); err != nil {
		return topAc, err
	}
	return topAc, nil
}

// 3. Funciton to get the top 5 Recreational activiteis

// 4. Function to get all the bus registered
func (m *PostgresDBRepo) GetAllBusInSystem() ([]models.AddBusData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, bus_name, bus_source, bus_destination, bus_model, bus_no_plate, num_seats, office_pan, office_address, 
			   merchant_id, created_at, updated_at, image
		FROM bus
	`
	var busses []models.AddBusData

	rows, err := m.DB.QueryContext(ctx, query)
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
			&i.Image,
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

// 5. Function to get all the hotels registered
func (m *PostgresDBRepo) GetAllHotelRoomsInSystem() ([]models.HotelRoom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, hotel_name, hotel_room_name, hotel_type, hotel_address, hotel_pan, hotel_num_room, 
	hotel_phone_1, hotel_phone_2, merchant_id, hotel_description, created_at, updated_at, image
	FROM hotel_room 
	`
	var rooms []models.HotelRoom

	rows, err := m.DB.QueryContext(ctx, query)
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
			&i.Image,
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

// 6. Funciton to get all the recreational activities registered
func (m *PostgresDBRepo) GetAllActivityInSystem() ([]models.AddActivityData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, activity_name, activity_description, activity_price, activity_duration, max_size, min_age, phone_num, email, location,
						merchant_id, created_at, updated_at, image
		FROM activity
	`
	var activities []models.AddActivityData

	rows, err := m.DB.QueryContext(ctx, query)
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
			&i.Image,
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
