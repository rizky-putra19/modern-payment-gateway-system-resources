package helper

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"golang.org/x/crypto/bcrypt"
)

func TransformToExternalBankCode(code string) string {
	return strings.ReplaceAll(code, "IDR_", "")
}
func TransformToInternalBankCode(code string) string {
	return fmt.Sprintf("IDR_%v", code)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CompareTwoStrings(stringOne, stringTwo string) float32 {
	// removeSpaces(&stringOne, &stringTwo)

	if value := returnEarlyIfPossible(stringOne, stringTwo); value >= 0 {
		return value
	}

	firstBigrams := make(map[string]int)
	for i := 0; i < len(stringOne)-1; i++ {
		a := fmt.Sprintf("%c", stringOne[i])
		b := fmt.Sprintf("%c", stringOne[i+1])

		bigram := a + b

		var count int

		if value, ok := firstBigrams[bigram]; ok {
			count = value + 1
		} else {
			count = 1
		}

		firstBigrams[bigram] = count
	}

	var intersectionSize float32
	intersectionSize = 0

	for i := 0; i < len(stringTwo)-1; i++ {
		a := fmt.Sprintf("%c", stringTwo[i])
		b := fmt.Sprintf("%c", stringTwo[i+1])

		bigram := a + b

		var count int

		if value, ok := firstBigrams[bigram]; ok {
			count = value
		} else {
			count = 0
		}

		if count > 0 {
			firstBigrams[bigram] = count - 1
			intersectionSize = intersectionSize + 1
		}
	}

	return (2.0 * intersectionSize) / (float32(len(stringOne)) + float32(len(stringTwo)) - 2)
}

func returnEarlyIfPossible(stringOne, stringTwo string) float32 {
	// if both are empty strings
	if len(stringOne) == 0 && len(stringTwo) == 0 {
		return 1
	}

	// if only one is empty string
	if len(stringOne) == 0 || len(stringTwo) == 0 {
		return 0
	}

	// identical
	if stringOne == stringTwo {
		return 1
	}

	// both are 1-letter strings
	if len(stringOne) == 1 && len(stringTwo) == 1 {
		return 0
	}

	// if either is a 1-letter string
	if len(stringOne) < 2 || len(stringTwo) < 2 {
		return 0
	}

	return -1
}

func EncodeToMd5(msg string) string {
	hash := md5.Sum([]byte(msg))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func Contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

func GetInputNumber(phoneNumber string) string {
	prefix := phoneNumber[:4]
	return prefix
}

// Check and add .00 if necessary, then change back to float64
func FormatFloat64(num float64) float64 {
	// Check whether num has a decimal part or not
	if num == float64(int(num)) {
		// Add .00
		formattedStr := strconv.FormatFloat(num, 'f', 2, 64)
		// Change back to float64
		formattedFloat, err := strconv.ParseFloat(formattedStr, 64)
		if err != nil {
			fmt.Println("Error converting string to float64:", err)
			return num
		}
		return formattedFloat
	}
	return num
}

func SplitString(input string) []string {
	// Remove the leading and trailing brackets
	input = strings.TrimPrefix(input, "[")
	input = strings.TrimSuffix(input, "]")

	// Split the string by commas
	items := strings.Split(input, ",")

	// Trim spaces from each item and store them in the slice
	var result []string
	for _, item := range items {
		trimmedItem := strings.TrimSpace(item)
		result = append(result, trimmedItem)
	}

	return result
}

// GenerateRandomString generates a random alphanumeric string of a given length
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomPinNumericString(length int) string {
	const charset = "" +
		"0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func StringToSignatureSymmetric(text string, clientSecret string) string {
	h := hmac.New(sha512.New, []byte(clientSecret))
	h.Write([]byte(text))
	d := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(d)
}

func GenerateTime(hours int) string {
	// Create a location for GMT+7
	location := time.FixedZone("GMT+7", 7*3600)

	// Get the current time in the GMT+7 timezone
	now := time.Now().In(location)

	// Set the time to 00:00:00
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)

	times := startOfDay.Format("2006-01-02 15:04:05")

	// Add specified hours to the start of the day
	if hours != 0 {
		newTime := startOfDay.Add(time.Duration(hours) * time.Hour)
		times = newTime.Format("2006-01-02 15:04:05")
	}

	return times
}

func FormattedUsingPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

func GenerateTiers(payload []entity.MerchantPaychannel) []dto.TiersDto {
	var tiers []dto.TiersDto
	var respTier []dto.TiersDto
	tierId := 1

	// Generate the initial list of tiers
	for i := 'A'; i <= 'J'; i++ {
		tier := dto.TiersDto{
			Id:    tierId,
			Tiers: fmt.Sprintf("TIER-%c", i),
		}
		tiers = append(tiers, tier)
		tierId++
	}

	// Create a map to track existing segments in the payload
	existingSegments := make(map[string]bool)
	for _, load := range payload {
		existingSegments[load.Segment] = true
	}

	// Filter out tiers that already exist in the payload
	for _, tier := range tiers {
		if !existingSegments[tier.Tiers] {
			respTier = append(respTier, tier)
		}
	}

	return respTier
}

func FloatPtr(f float64) *float64 {
	return &f
}

func StringPtr(s string) *string {
	return &s
}

func ExtractDate(datetimeStr string) (string, error) {
	// Format dari string datetime yang akan di-parse
	const layout = "2006-01-02 15:04:05"

	// Parse string ke tipe time.Time
	datetime, err := time.Parse(layout, datetimeStr)
	if err != nil {
		return "", err
	}

	// Format time.Time menjadi string hanya tanggal
	dateStr := datetime.Format("02/01/2006")
	return dateStr, nil
}

// IsValidEmail checks if the provided email address is valid
func IsValidEmail(email string) bool {
	// Define the regular expression pattern for a valid email address
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	emailRegex := regexp.MustCompile(emailRegexPattern)

	// Check if the email matches the pattern
	return emailRegex.MatchString(email)
}

func HashString(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		slog.Infof("failed to hash string")
		return "", err
	}

	return string(hash), nil
}

func CheckingFirstAndLastStr(str string) (first string, last string, err error) {
	words := strings.Fields(str)

	if len(words) > 0 {
		firstWord := words[0]
		if len(words) > 1 {
			restOfWords := strings.Join(words[1:], " ")
			return firstWord, restOfWords, nil
		} else {
			return firstWord, firstWord, nil
		}
	}

	return "", "", errors.New("there is no string")
}
