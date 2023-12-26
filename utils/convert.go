package utils

import (
	"fmt"
	"time"
)

// ----------------------------------------------
// exported funtions
// ----------------------------------------------

func ConvertUnixToIso(unixTime int64) string {
	t := time.UnixMilli(unixTime)
	fmt.Println(t)
	isoTime := t.Format("2006-01-02T15:04:05.000Z")
	return isoTime
}
