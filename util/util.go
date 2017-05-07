package util

import (
	"strconv"
	"time"
	"math/rand"
	"crypto/md5"
	"io"
	"fmt"
	"encoding/base64"
)

func ToAddr(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

func Random(array []string) string {
	rand.Seed(time.Now().Unix())
	return array[rand.Intn(len(array))]
}

func Md5(key string) string {
	h := md5.New()
	io.WriteString(h, key)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RandStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Base64Encode(value string) string {
	return base64.URLEncoding.EncodeToString([]byte(value))
}

func Base64Decode(key string) string {
	bytes, _ := base64.URLEncoding.DecodeString(key)
	return string(bytes)
}
