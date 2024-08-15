package forms

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

// Defining a custom form struct
type Form struct {
	url.Values
	Errors errors
}

// Initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(key string) bool {
	x := f.Get(key)
	if x == "" {
		f.Errors.Add(key, "This field cannot be empty")
		return false
	}
	return true
}

// Checks all of the requuired fields in a generalized way
func (f *Form) Required(keys ...string) {
	for _, key := range keys {
		value := f.Get(key)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(key, "This field cannot be blank")
		}
	}
}

// Checks for minimum length
func (f *Form) MinLength(key string, length int, r *http.Request) bool {
	value := f.Get(key)
	if len(value) < length {
		f.Errors.Add(key, fmt.Sprintf("Minimum number of length shoulld be: %d", length))
		return false
	}
	return true
}

// Checks if Email is valid
func (f *Form) IsEmailValid(key string) {
	if !govalidator.IsEmail(f.Get(key)) {
		f.Errors.Add(key, "Please enter a valid Email address")
	}
}

// Returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Funtion to hash the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Validate if the function is valid
// func (f *Form) IsPasswordValid(password, confirmPassword string) (string, error) {
// 	// hash, _ := bcrypt.GenerateFromPassword([]byte(f.Get(password)), 12)
// 	hash, _ := HashPassword(password)

// 	err := bcrypt.CompareHashAndPassword(hash, []byte(f.Get(confirmPassword)))
// 	if err == bcrypt.ErrMismatchedHashAndPassword {
// 		log.Println("Password Mismatch")
// 		f.Errors.Add(confirmPassword, "Password does not match")
// 		return "", err
// 	} else if err != nil {
// 		log.Println("Unexpected Error")
// 		f.Errors.Add(confirmPassword, "Unexpected Error")
// 		return "", err
// 	}

// 	return string(hash), nil
// }

// Check if password and confirmPassword are the same
func (f *Form) IsPasswordValid(password, confirmPassword string) {
	if f.Get(password) != f.Get(confirmPassword) {
		f.Errors.Add(password, "Password and confirm password must be the same")
		f.Errors.Add(confirmPassword, "Password and confirm password must be the same")
	}
}

// Check if user has checked the box
func (f *Form) HasUserAccepted(key string) {
	if f.Get(key) == "agreed" {
		return
	} else {
		f.Errors.Add(key, "Have to agree to the terms of service")
	}
}

// Helper function to check if the email exist
func (f *Form) FormValidateUser(key string, value bool) {
	if value {
		f.Errors.Add(key, "User already Exists, unable to create account")
		return
	} else {
		return
	}
}

// Function to add the verification error for the code
func (f *Form) AddVerificationError() {
	f.Errors.Add("verification_code", "Verification Code is not the same")
}

// Function to check if the image file is less than 200 kb
func IsValidFileSize(handler *multipart.FileHeader, x int64) bool {
	return handler.Size <= x*1024
}

// Funciton to validate if the field is an integer, if yes return the value
func (f *Form) ConvertToInt(key string) int {
	value := f.Get(key)
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		f.Errors.Add(key, "The Value you entered is not an integer")
		return 0
	} else {
		return valueInt
	}
}
