package db

import (
	"fmt"

	"math/rand"

	"strings"

	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {

	rand.Seed(time.Now().UnixNano())

}

// randomInt generates a random integer between min and max

func randomInt(min, max int64) int64 {

	return min + rand.Int63n(max-min+1)

}

// randomString generates a random string of length n

func randomString(n int) string {

	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {

		c := alphabet[rand.Intn(k)]

		sb.WriteByte(c)

	}

	return sb.String()

}

// randomEmail generates a random email

func randomEmail() string {

	return fmt.Sprintf("%s@%s.com", randomString(6), randomString(4))

}

// randomBool generates a random boolean

func randomBool() bool {

	return rand.Intn(2) == 1

}
