package util

import (
	"strconv"
)

var count int

func Temp_next() string {
	s := "x_" + strconv.Itoa(count)
	count++
	return s
}
