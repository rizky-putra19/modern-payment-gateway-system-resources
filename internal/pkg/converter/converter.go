package converter

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"github.com/skip2/go-qrcode"
)

// ToString converts any value to string
func ToString(v interface{}) string {
	result := ""
	if v == nil {
		return ""
	}
	switch v := v.(type) {
	case string:
		result = v
	case int:
		result = strconv.Itoa(v)
	case int32:
		result = strconv.Itoa(int(v))
	case int64:
		result = strconv.FormatInt(v, 10)
	case bool:
		result = strconv.FormatBool(v)
	case float32:
		result = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		result = strconv.FormatFloat(v, 'f', -1, 64)
	case []uint8:
		result = string(v)
	default:
		resultJSON, err := json.Marshal(v)
		if err == nil {
			result = string(resultJSON)
		} else {
			log.Printf("failed to convert data: [%#v] with error: [%v]\n", v, err)
		}
	}

	return result
}

func ToBase64Img(text string) string {
	byteImg, err := qrcode.Encode(text, qrcode.High, 256)
	if err != nil {
		slog.Errorw("failed to create image with error", err.Error())
	}
	base64Img := base64.StdEncoding.EncodeToString(byteImg)

	return fmt.Sprintf("data:image/png;base64,%v", base64Img)
}

func ToMD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// ToBool convert any value to boolean
func ToBool(v interface{}) bool {
	var result bool
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, _ = strconv.ParseBool(str)
	case int:
		result = v != 0
	default:
		// do nothing
	}

	return result
}

// ToInt converts any value to int
func ToInt(v interface{}) int {
	result := 0
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		result, _ = strconv.Atoi(str)
	case int:
		result = v
	case int32:
		result = int(v)
	case int64:
		result = int(v)
	case float32:
		result = int(v)
	case float64:
		result = int(v)
	case []byte:
		result, _ = strconv.Atoi(string(v))
	default:
		result = 0
	}

	return result
}

// ToInt64 converts any value to int64
func ToInt64(v interface{}) int64 {
	result := int64(0)
	switch v := v.(type) {
	case string:
		str := strings.TrimSpace(v)
		x, _ := strconv.Atoi(str)
		result = int64(x)
	case int:
		result = int64(v)
	case int32:
		result = int64(v)
	case int64:
		result = v
	case float32:
		result = int64(v)
	case float64:
		result = int64(v)
	case []byte:
		x, _ := strconv.Atoi(string(v))
		result = int64(x)
	default:
		result = 0
	}

	return result
}

// ToArrayOfInt convert any value to []int
func ToArrayOfInt(v interface{}) []int {
	var result []int
	switch v := v.(type) {
	case string:
		_ = json.Unmarshal([]byte(v), &result)
	case []string:
		b := v
		for _, vv := range b {
			result = append(result, ToInt(vv))
		}
	case [][]byte:
		b := v
		for _, vv := range b {
			result = append(result, ToInt(vv))
		}
	default:
		// do nothing
	}

	return result
}

// ToArrayOfString convert any value to []string
func ToArrayOfString(v interface{}) []string {
	var result []string
	switch v := v.(type) {
	case string:
		_ = json.Unmarshal([]byte(v), &result)
	case [][]byte:
		b := v
		for _, vv := range b {
			result = append(result, string(vv))
		}
	default:
		// do nothing
	}

	return result
}

func FromStringToIntAmount(val string) int {
	var result int
	str := strings.TrimSpace(val)
	num := strings.Split(str, ".")
	result, _ = strconv.Atoi(num[0])

	return result
}

func FormattedCompletionRate(averageDuration time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d", int(averageDuration.Hours()), int(averageDuration.Minutes())%60, int(averageDuration.Seconds())%60)
}
