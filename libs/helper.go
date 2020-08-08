package libs

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func FileUrl(path string) string {
	return "file:///" + strings.TrimLeft(strings.ReplaceAll(path, "\\", "/"), "/")
}

func RandInt(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return r.Intn(max-min) + min
}
