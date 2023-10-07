package testing

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	date := time.Now().Local().Format("2006-01-02 15:04:05")

	fmt.Println(date)
}
