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

// TestPutInSlice 
func TestPutInSlice(t *testing.T) {
	s := make([]string, 2)

	testTable := []struct {
		pos int
		item string
	} {
		{3, "test1"},
		{12, "test2"},
		{16, "test3"},
		{255, "test4"},
		{23401, "test5"},
	}

	for _, testItem := range testTable {
		putInStringSlice(&s, testItem.pos, &testItem.item)
		if s[testItem.pos] != testItem.item {
			t.Errorf(`Tried putting %s at position %d into slice of original
			length but exited with this slice: %s`,
				testItem.item, testItem.pos, s)
		}
	}
}

func TestStringToCharLine(t *testing.T) {
	s := []string{"one", "two", "three", "four"}
	cl := []charLine{"one", "two", "three", "four"}

	for n, v := range s {
		if cl[n] != charLine(v) {
			t.Errorf("Error casting string %s to charLines", v)
		}
	}
}

func TestStringsToCharLines(t *testing.T) {
	sSlice := []string{"one", "two", "three", "four"}
	clSlice := stringsToCharLines(sSlice)

	identical := true

	for n, v := range clSlice {
		if v != charLine(sSlice[n]) {
			identical = false
			break
		}
	}

	if (len(sSlice) !=  len(clSlice)) || ! identical {
		t.Errorf("Function stringsToCharLines took %s and returned %s",
		sSlice, clSlice)
	}
}

// func TestGetCharDict(t *testing.T) {
// 	cd, err := GetCharDict()
// 	// TODO: check some keys here for their correct values
// 	if err != nil {
// 		t.Error("Error getting character dictionary")
// 	}

// 	strokesTable := []struct {
// 		in string
// 		out int64
// 	} {
// 		{"我", 7},
// 		{"赢", 20},
// 		{"籠", 22},
// 		{"三", 3},
// 		{"忽", 8},
// 	}

// 	for _, testItem := range strokesTable {
// 		if cd.dict[testItem.in].keyStrokes != testItem.out {
// 			t.Errorf(`Error checking CharDict entries. keyStroke value of %s
// 			should be %d but is %d`, testItem.in, testItem.out, cd.dict[testItem.in].keyStrokes)
// 			t.Error(cd.dict[testItem.in])
// 		}
// 	}
// }
