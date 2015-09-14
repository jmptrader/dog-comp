package util

import (
	"strconv"
)

var count int

func Next() string {
	s := "x_" + strconv.Itoa(count)
	count++
	return s
}
