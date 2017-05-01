package util

import (
	"strconv"
	"time"
	"math/rand"
)

func ToAddr(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

func Random(array []string) string {
	rand.Seed(time.Now().Unix())
	return array[rand.Intn(len(array))]
}