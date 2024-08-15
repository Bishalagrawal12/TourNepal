package handlers

import (
	"time"

	"github.com/Atul-Ranjan12/tourism/internal/models"
)

// This funciton gets the number of invalid days before the first day of the week
func getInvalidDays(date time.Time) int {
	daysOfWeek := []int{6, 0, 1, 2, 3, 4, 5}
	dayOfWeekIndex := int(date.Weekday())
	return daysOfWeek[dayOfWeekIndex]
}

// Function to add the invalid days in the month
func addInvalidDaysStart(date time.Time) []int {
	numDays := getInvalidDays(date)

	days := []int{26, 27, 28, 29, 30, 31}
	returnDays := []int{}

	for i := len(days) - numDays; i < len(days); i++ {
		returnDays = append(returnDays, days[i])
	}

	return returnDays
}

// Funciton to add the invalid days in the end of the month
func addInvalidEnd(date time.Time) []int {
	numDays := getInvalidDays(date)

	days := []int{1, 2, 3, 4, 5, 6, 7}
	returnDays := []int{}

	counter := 1
	for i := numDays; i < len(days)-1; i++ {
		returnDays = append(returnDays, counter)
		counter++
	}

	return returnDays
}

// Funciton to add the active days into the array
func addActiveDays(date time.Time) [][]models.CalenderDay {
	// Find the number of days in the month
	lastOfMonth := date.AddDate(0, 1, -1)
	numDays := lastOfMonth.Day()

	var calDay [][]models.CalenderDay

	invalidDays := addInvalidDaysStart(date)
	var i int = 1

	var nextWeeekStart int = 1
	var firstWeek []models.CalenderDay
	for j := len(invalidDays); j < 7; j++ {
		day := models.CalenderDay{
			Day: i,
		}
		firstWeek = append(firstWeek, day)
		i++
		nextWeeekStart++
	}

	calDay = append(calDay, firstWeek)

	// Addding the other weeks into the calender
	var counter int = 1
	var week []models.CalenderDay

	var lastWeekStart int
	for j := nextWeeekStart; j < numDays; j++ {
		week = append(week, models.CalenderDay{Day: j})
		if counter >= 7 {
			calDay = append(calDay, week)
			week = []models.CalenderDay{}
			counter = 1
			lastWeekStart = j
			continue
		}
		counter++
	}

	// Add the last week to the array
	var lastWeek []models.CalenderDay
	for j := 1; j <= 7-len(addInvalidEnd(lastOfMonth)); j++ {
		lastWeek = append(lastWeek, models.CalenderDay{Day: lastWeekStart + 1})
		lastWeekStart++
	}

	calDay = append(calDay, lastWeek)

	return calDay
}
