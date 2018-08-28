package hans

import (
	"math/rand"
	"regexp"
	"testing"
)

const testHtml = `<html>
<head>
<meta charset="UTF-8"/>
</head>
<body>
<pre>									
	一	1	一	一	1		*	0		M	*
	丁	2	吕	一	1		亅	1		MN	一
	丂	2	一	丂	2		*	0		MVS	一
	七	2	一	七	2		*	0		JU	一
	丄	2	一	丄	2		*	0		LM	一
	丅	2	一	丅	2		*	0		ML	一
	丆	2	吕	一	2		丿	0		MH	一
	万	3	一	万	3		*	0		MS	一
	丈	3	一	丈	3		*	0		JK	一
	三	3	回	二	2		一	1		MMM	一
	上	3	一	上	3		*	0		YM	一
	下	3	吕	一	1		卜	2		MY	一
	丌	3	一	丌	3		*	0		ML	一
	不	4	一	不	4		*	0		MF	一
	与	4	一	与	4		*	0		YSM	一
	丏	4	一	丏	4		*	0		MLVS	一
	丐	4	一	丐	4		*	0		MYVS	一
	丑	4	一	丑	4		*	0		NG	一
	丒	4	吕	刃	3		一	1		SKM	一
	专	4	一	专	4		*	0		QNI	一
	且	5	一	且	5		*	0		BM	一
	丕	5	吕	不	4		一	1		MFM	一
	世	4	一	世	4		*	0		PT	一
	丗	4	一	丗	4		*	0		TJ	一
	丘	5	一	丘	5		*	0		OM	一
	丙	5	吕	一	1		内	4		MOB	一
	业	5	一	业	5		*	0		TC	一
	丛	5	咒	人	2		一	1		OOM	一
	东	5	一	东	5		*	0		KD	一
	丝	5	一	丝	5		*	0		VVM	一
	丞	6	吕	氶	5		一	1		NEM	一
	丟	6	吕	王	4		厶	2		MGI	一
	丠	6	吕	北	5		一	1		LPM	一
	両	6	一	両	6		*	0		MUB	一
	丢	6	吕	壬	4		厶	2		HGI	一
	丣	7	一	丣	7		*	0		MLLS	一
	两	7	一	两	7		*	0		MOOB	一
	严	7	一	严	7		*	0		MTCH	一
	並	8	吕	丷	2		亚	6		TTC	一
	丧	8	一	丧	8		*	0		GCV	一
</pre>
</body>
</html>`

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
			first_cl, first_cl+sec_cl, cl)
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

	if (len(sSlice) != len(clSlice)) || !identical {
		t.Errorf("Function stringsToCharLines took %s and returned %s",
			sSlice, clSlice)
	}
}

// TestExtractTableLines
func TestExtractTableLines(t *testing.T) {
	clSlice, err := extractTableLines([]byte(testHtml))
	if err != nil {
		t.Error(err.Error())
	}

	// get number of lines in testHtml
	rxPre := regexp.MustCompile(`(?s)<pre>\s*\n(.+?)\s*\n</pre>`)
	rxNL := regexp.MustCompile(`\n`)
	s := rxPre.FindStringSubmatch(testHtml)[1]
	n := len(rxNL.FindAllStringIndex(s, -1)) + 1

	if len(clSlice) != n {
		t.Errorf(`extractTableLines extracted %d lines when it should have extracted %d.
Extracted lines:
%s

HTML:
%s`, len(clSlice), n, clSlice, testHtml)
	}

}

func TestGetCharDict(t *testing.T) {
	cd, err := GetCharDict()

	if err != nil {
		t.Error("Error getting character dictionary")
	}

	strokesTable := []struct {
		in string
		out int64
	} {
		{"我", 7},
		{"赢", 20},
		{"籠", 22},
		{"三", 3},
		{"忽", 8},
	}

	for _, testItem := range strokesTable {
		if cd.dict[testItem.in].keyStrokes != testItem.out {
			t.Errorf(`Error checking CharDict entries. keyStroke value of %s
			should be %d but is %d`, testItem.in, testItem.out, cd.dict[testItem.in].keyStrokes)
			t.Error(cd.dict[testItem.in])
		}
	}
}
