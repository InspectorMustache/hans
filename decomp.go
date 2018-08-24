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
)


var lineRx, _ = regexp.Compile(`^\s+(.)\s+(\d+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\S+)$`)
var keyRx, _ = regexp.Compile(`^\s+(.)`)
var infoRx, _ = regexp.Compile(`\s+(\d+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\d+)\s+(\*| |\?)\s+(\S+)\s+(\S+)$`)
var preOpenRx, _ = regexp.Compile(`^<pre>\s*\n(\s{8}.)`)
var preCloseRx, _ = regexp.Compile(`\n\s*</pre>$`)

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

// putInStringSlice attempts to put a string into slice at position pos. if that
// position is out of bounds, the slice capacity is extended
func putInStringSlice(slice *[]string, pos int, item *string) {
	if cap(*slice) > pos {
		(*slice)[pos] = *item
	} else {
		nSlice := make([]string, getNewSliceSize(cap(*slice), pos))
		copy(nSlice, *slice)
		*slice = nSlice
		(*slice)[pos] = *item
	}
}

// putInCharLineSlice attempts to put a charLine into slice at position pos. if
// that position is out of bounds, the slice capacity is extended
func putInCharLineSlice(slice *[]charLine, pos int, item *charLine) {
	if cap(*slice) > pos {
		(*slice)[pos] = *item
	} else {
		nSlice := make([]charLine, getNewSliceSize(cap(*slice), pos))
		copy(nSlice, *slice)
		*slice = nSlice
		(*slice)[pos] = *item
	}
}

// getNewSliceSize keeps increasing oldSize by a fifth until it at least reaches minSize 
func getNewSliceSize(oldSize int, minSize int) int {
	for oldSize < minSize + 1 {
		plus := (oldSize / 5) + 1 // increase by at least 1
		oldSize = oldSize + plus
	}
	return oldSize
}

// stringsToCharLines converts a slice of strings to a slice of charLines.
func stringsToCharLines(stringS []string) []charLine {
	clSlice := make([]charLine, len(stringS))
	for k, v := range stringS {
		cl := charLine(v)
		putInCharLineSlice(&clSlice, k, &cl)
	}

	return clSlice
}

// getShaFromReader finishes reading from r and then returns the sha256sum of
// the yielded content.
func getShaFromReader(r io.Reader) ([32]byte, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return [32]byte{}, err
	}
	sha := sha256.Sum256(raw)

	return sha, nil
}

// extractTableLines reads the content of r and parses them as html which is
// then split into a slice of charLines and returned.
func extractTableLines(r io.Reader) ([]charLine, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	lines := make([]charLine, 21800)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "pre" {
			s := strings.Split(cleanNode(n), "\n")
			lines = append(lines, stringsToCharLines(s)...)
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

	sha, err := getShaFromReader(resp.Body)
	if err != nil {
		return charDict{}, err
	}

	cd := charDict{sha: sha}

	l, err := extractTableLines(resp.Body)
	if err != nil {
		return charDict{}, err
	}
	cd.populate(l)
	return cd, nil
	
}
