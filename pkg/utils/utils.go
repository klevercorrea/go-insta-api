package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Davincible/goinsta/v3"
	"golang.org/x/time/rate"
)

func RandomTime(minSeconds, maxSeconds int) time.Duration {
	return time.Duration(minSeconds+rand.Intn(maxSeconds-minSeconds)) * time.Second
}

func Countdown(sleepDuration int) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Printf("Sleeping for %d seconds\n", sleepDuration)
	for i := sleepDuration; i > 0; i-- {
		<-ticker.C
		fmt.Printf("\r%d", i)
	}
	fmt.Printf("\rDone sleeping for %d seconds\n", sleepDuration)
}

func ShortSleep() {
	sleepDuration := time.Duration(10+rand.Intn(30)) * time.Second
	time.Sleep(sleepDuration)
}

func LongSleep() {
	sleepDuration := time.Duration(300+rand.Intn(600)) * time.Second
	fmt.Printf("Sleep for: %v\n", sleepDuration)
	time.Sleep(sleepDuration)
}

func RateLimiter() *rate.Limiter {
	limiter := rate.NewLimiter(2, 1)
	return limiter
}

func Contains(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func MustMarshalJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// ImportFromBytes imports instagram configuration from an array of bytes.
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportFromBytes(inputBytes []byte) (*goinsta.Instagram, error) {
	return goinsta.ImportReader(bytes.NewReader(inputBytes))
}

// ImportFromBase64String imports instagram configuration from a base64 encoded string.
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportFromBase64String(base64String string) (*goinsta.Instagram, error) {
	sDec, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	return ImportFromBytes(sDec)
}

// ExportAsBytes exports selected *Instagram object as []byte
func ExportAsBytes(insta *goinsta.Instagram) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := insta.ExportIO(buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// ExportAsBase64String exports selected *Instagram object as base64 string
func ExportAsBase64String(insta *goinsta.Instagram) (string, error) {
	bytes, err := insta.ExportAsBytes()
	if err != nil {
		return "", err
	}

	sEnc := base64.StdEncoding.EncodeToString(bytes)
	return sEnc, nil
}
