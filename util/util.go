package util

import "strconv"

func ToAddr(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}
