package utils

import (
	"fmt"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"time"
)

//тут вспомогательные функции

// округление  вверх
func RoundUp(num, multiple int) int {
	if num%multiple == 0 {
		return num
	}
	return num + (multiple - num%multiple)
}

// перевод минут в время формата 00:00
func MinuteToTime(min int) string {
	return fmt.Sprintf("%02d:%02d", min/60, min%60)
}

// время в минутах от начала дня
func MinuteFromStartDay(time time.Time) int {
	hours, minute, _ := time.Clock()
	return hours*60 + minute
}

// перевод time.Time в дату
func TimeToDate(t time.Time) string {
	return t.Format("2006 January	02")
}

//маскировка

func Encrypt(data string, key []byte) (string, error) {
	encrypted, err := jwe.Encrypt(
		[]byte(data),
		jwe.WithKey(jwa.A256GCMKW, key),
		jwe.WithContentEncryption(jwa.A256GCM),
	)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
