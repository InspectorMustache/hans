package hans

import (
	"math/rand"
	"regexp"
	"testing"
)

// getRandString generates a random string of length n
func getRandString(n int) string {
	pool := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÜÖÄabcdefghijklmnopqrstuvwxyzüöä")
	s := make([]rune, n)

	for pos := range s {
		s[pos] = pool[rand.Intn(len(pool))]
	}

	return string(s)
}

func TestValidate(t *testing.T) {
	// TODO: need some sort of test here
	return
}

func TestRxPop(t *testing.T) {
	first_cl := charLine(getRandString(rand.Intn(25)))
	sec_cl := charLine(getRandString(rand.Intn(25)))
	cl := first_cl + sec_cl
	rx, err := regexp.Compile(string(first_cl))
	if err != nil {
		t.Error("Error compiling regexp from", first_cl)
	}

	s, err := cl.rxPop(rx)
	if !(charLine(s) == first_cl &&
		cl == sec_cl &&
		err == nil) {
		t.Errorf("%s popped off of %s returned %s",
			first_cl, first_cl + sec_cl, cl)
	}
}

func TestStringToCharLine(t *testing.T) {
	s := []string{"one", "two", "three", "four"}
	cl := []charLine{"one", "two", "three", "four"}

	if len(s) != len(cl) {
		t.Errorf("Error casting %s to %s", s, cl)
	}

	for n, v := range s {
		if cl[n] != charLine(v) {
			t.Errorf("Error casting %s to %s", s, cl)
		}
	}
}

func TestGetCharDict(t *testing.T) {
	cd, err := GetCharDict()
	// TODO: check some keys here for their correct values
	if err != nil {
		t.Error("Error getting character dictionary")
	}
}
