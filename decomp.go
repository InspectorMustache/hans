package hans

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"bytes"
)

var lineRx = regexp.MustCompile(`^\s+(.)\s+(\d+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\S+)$`)
var keyRx = regexp.MustCompile(`^\s+(.)`)
var infoRx = regexp.MustCompile(`\s+(\d+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\S+)$`)
var preOpenRx = regexp.MustCompile(`^<pre>\s*\n(\s+.)`)
var preCloseRx = regexp.MustCompile(`\n\s*</pre>$`)

type charInfo struct {
	// move this to map key
	// keyChar string
	keyStrokes        int64
	compType          string
	firstPart         string
	firstPartStrokes  int64
	firstPartVerify   bool
	secondPart        string
	secondPartStrokes int64
	secondPartVerify  bool
	cangjie           string
	radical           string
	structVerify      bool
	components        []string
}

type charDict struct {
	sha  [32]byte
	dict map[string]charInfo
}

type charLine string

// string returns string representation of charLine.
func (cl *charLine) string() string {
	return string(*cl)
}

func (cl *charLine) validate() error {
	// check if charLine matches the general format for charlines
	// return nil if it does
	if lineRx.MatchString(cl.string()) == false {
		return fmt.Errorf("Line %s doesnt' match regex %s.", *cl,
			lineRx.String())
	} else {
		return nil
	}
}

// rxPop returns the match for the regexp rx and trims it off
func (cl *charLine) rxPop(rx *regexp.Regexp) (string, error) {
	r := rx.FindString(cl.string())
	s := rx.ReplaceAllString(cl.string(), "")
	*cl = charLine(s)
	return r, nil
}

// populate individually parses a slice of charLines, creating entries in
// charDicts dictionary for each element
func (cd *charDict) populate(lines []charLine) {
	cd.dict = make(map[string]charInfo, len(lines))

	for _, l := range lines {
		k, e := l.rxPop(keyRx)
		if e != nil {
			continue
		}

		rest, e := l.rxPop(infoRx)
		if e != nil {
			continue
		}

		ci := new(charInfo)
		ci.populate(rest)
		cd.dict[k] = *ci
	}

}

// populate parses the string line and populates the fields of its parent ccs.
func (ccs *charInfo) populate(line string) {
	// line is assumed to already be stripped of the keychar
	// set all the other fields
	if infoRx.MatchString(line) {
		fieldSlice := infoRx.FindStringSubmatch(line)
		// we can assume there are no errors because the \d regex
		// matched
		ccs.keyStrokes, _ = strconv.ParseInt(fieldSlice[1], 10, 64)
		ccs.compType = fieldSlice[2]
		ccs.firstPart = fieldSlice[3]
		ccs.firstPartStrokes, _ = strconv.ParseInt(fieldSlice[4], 10, 64)
		ccs.secondPart = fieldSlice[6]
		ccs.secondPartStrokes, _ = strconv.ParseInt(fieldSlice[7], 10, 64)
		ccs.cangjie = fieldSlice[9]
		ccs.radical = fieldSlice[10]
		ccs.structVerify = true

		if fieldSlice[5] == "?" {
			ccs.firstPartVerify = false
		} else {
			ccs.firstPartVerify = true
		}

		if fieldSlice[8] == "?" {
			ccs.secondPartVerify = false
		} else {
			ccs.secondPartVerify = true
		}
	} else {
		ccs.structVerify = false
	}
}

// stringsToCharLines converts a slice of strings to a slice of charLines.
func stringsToCharLines(stringS []string) []charLine {
	clSlice := make([]charLine, 0, len(stringS))
	for _, v := range stringS {
		cl := charLine(v)
		clSlice = append(clSlice, cl)
	}

	return clSlice
}

// extractTableLines reads the content of b and parses it as html which is
// then split into a slice of charLines and returned.
func extractTableLines(b []byte) ([]charLine, error) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	lines := make([]charLine, 0, 21800)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "pre" {
			s := strings.Split(cleanNode(n), "\n")
			if len(s) > 0 {
				lines = append(lines, stringsToCharLines(s)...)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return lines, nil
}

// cleanNode returns the html node n as a processable string.
func cleanNode(n *html.Node) string {
	var b strings.Builder
	w := io.Writer(&b)
	html.Render(w, n)
	s := b.String()
	// if preOpenRx doesn't match we are at the opening paragraph
	if !preOpenRx.MatchString(s) {
		return ""
	}
	s = preOpenRx.ReplaceAllString(s, "$1")
	s = preCloseRx.ReplaceAllString(s, "")
	return s
}

// GetCharDict downloads the webpage of the Wikimedia decomposition project and
// turns it into a charDict.
func GetCharDict() (charDict, error) {
	resp, err := http.Get(
		"https://commons.wikimedia.org/wiki/Commons:Chinese_characters_decomposition")

	if err != nil {
		return charDict{}, err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return charDict{}, err
	}

	l, err := extractTableLines(raw)
	if err != nil {
		return charDict{}, err
	}

	cd := charDict{sha: sha256.Sum256(raw)}
	cd.populate(l)

	return cd, nil
}
