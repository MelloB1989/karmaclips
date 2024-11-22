package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

func GenerateOTP() string {
	otp := rand.Intn(900000) + 100000
	return strconv.Itoa(otp)
}

func GenerateID() string {
	id, _ := gonanoid.Generate("qwertyuiopasdfghjklzxcvbnm1234567890_-", 10)
	return id
}

// VerifyPhoneNumber checks if the phone number is in the format +917569236628
func VerifyPhoneNumber(phone string) bool {
	// Define the regular expression pattern
	phonePattern := `^\+91\d{10}$`
	// Compile the regular expression
	re := regexp.MustCompile(phonePattern)
	// Check if the phone number matches the pattern
	return re.MatchString(phone)
}

// Ensure that single-digit day and month are zero-padded
func FormatDate(input string) string {
	// Split the date into day, month, and year
	parts := strings.Split(input, "/")
	if len(parts) != 3 {
		return input // If the input format is wrong, return as is.
	}

	// Pad single-digit day and month with leading zero if necessary
	if len(parts[0]) == 1 {
		parts[0] = "0" + parts[0]
	}
	if len(parts[1]) == 1 {
		parts[1] = "0" + parts[1]
	}

	// Reconstruct the date string in the desired format "DD/MM/YYYY"
	return strings.Join(parts, "/")
}

// Function to parse date with single-digit handling
func ParseDate(input string) (time.Time, error) {
	// Format the input date to ensure it's in "DD/MM/YYYY"
	formattedDate := FormatDate(input)

	// Parse the formatted date
	date, err := time.Parse("02/01/2006", formattedDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %v", err)
	}
	return date, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GetNow() string {
	return time.Now().String()
}
