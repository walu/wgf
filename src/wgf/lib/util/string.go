package util

import (
	"math/rand"
	"strconv"
	"time"
)

func Uuid() string {
	now := time.Now()
	unixtimestamp := now.Unix()
	rand.Seed(unixtimestamp)

	pre := strconv.FormatInt(unixtimestamp, 36)
	suf := strconv.FormatInt(rand.Int63(), 36)
	return pre + suf
}
