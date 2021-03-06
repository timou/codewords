// +build ignore

package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const wordnetURL = "http://wordnetcode.princeton.edu/wn3.1.dict.tar.gz"
const minWordLen = 4
const maxWordLen = 15

func main() {
	// Acquire the WordNet database.
	wn, err := NewWNDict()
	if err != nil {
		panic(err)
	}

	// Filter out words that don't make good codewords.
	wn.nouns = filter(wn.nouns)
	wn.adjectives = filter(wn.adjectives)

	// Write out new .go files.
	err = writeList("adjectives.go", "adjectives", wn.adjectives)
	if err != nil {
		panic(err)
	}
	err = writeList("nouns.go", "nouns", wn.nouns)
	if err != nil {
		panic(err)
	}
}

// wndict is an in-memory structure to hold potential codewords.
type wndict struct {
	adjectives []string
	nouns      []string
}

// NewWNDict downloads the WordNet database, parses out the
// adjectives and nouns that are used to form codewords.
func NewWNDict() (*wndict, error) {
	d := &wndict{}
	buf, err := wnbuf()
	if err != nil {
		return nil, err
	}

	gzr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch hdr.Name {
		case "dict/index.adj":
			d.adjectives = parseWords(tr)
		case "dict/index.noun":
			d.nouns = parseWords(tr)
		}
	}
	return d, nil
}

// wnbuf tries to acquire the WordNet database. It first looks
// for the file in the working directory. If this fails, it
// tries to download it from the internet into a memory buffer.
func wnbuf() (*bytes.Buffer, error) {
	var b bytes.Buffer
	slashIdx := strings.LastIndex(wordnetURL, "/")
	if slashIdx == -1 {
		return nil, errors.New("failed to parse WordNet URL")
	}
	if slashIdx+1 >= len(wordnetURL) {
		return nil, errors.New("failed to parse filename from WordNet URL")
	}
	localFile := wordnetURL[slashIdx+1:]

	// First try locally.
	ifp, err := os.Open(localFile)
	if err != nil {
		// Failed to find locally, download instead.
		res, ierr := http.Get(wordnetURL)
		if ierr != nil {
			return nil, fmt.Errorf("failed to download WordNet database, url=%s", wordnetURL)
		}
		defer res.Body.Close()

		_, ierr = io.Copy(&b, res.Body)
		if ierr != nil {
			return nil, fmt.Errorf("failed to copy WordNet database, err=%s", ierr)
		}
		return &b, nil
	}
	defer ifp.Close()
	_, err = io.Copy(&b, ifp)
	if err != nil {
		return nil, fmt.Errorf("failed to copy WordNet database, err=%s", err)
	}
	return &b, nil
}

// parseWords extracts words from the WordNet index files. It
// specifically only looks for adjectives and nouns. Lots of
// fields are ignored.
func parseWords(r io.Reader) []string {
	scn := bufio.NewScanner(r)
	var words []string
	scnRE := regexp.MustCompile("^([^\\s]+)\\s+[an]")
	for scn.Scan() {
		m := scnRE.FindStringSubmatch(scn.Text())
		if len(m) > 1 {
			words = append(words, m[1])
		}
	}
	return words
}

// filter is used to strip out words that make messy codewords.
// It drops any word that contains '_', '.', '/', '-' or any
// digits. It also drops short and long words. Feel free to
// modify this function to generate your own set, drop,
// drop inappropriate or offensive words, etc.
func filter(words []string) []string {
	var res []string
	for _, s := range words {
		if strings.ContainsAny(s, "_-'0123456789/.") {
			continue
		}
		if len(s) < minWordLen || len(s) > maxWordLen {
			continue
		}
		res = append(res, s)
	}
	return res
}

// writeList creates adjectives.go and nouns.go, which makes this
// library easily go-gettable. These generated files should be
// committed.
func writeList(filename string, varname string, words []string) error {
	ofp, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer ofp.Close()

	ofp.WriteString("// Code generated by 'go generate'.\n")
	ofp.WriteString("// Derived from Princeton Wordnet https://wordnet.princeton.edu\n\n")
	ofp.WriteString("package codewords\n\n")
	fmt.Fprintf(ofp, "var %s = []string{\n", varname)
	for _, word := range words {
		fmt.Fprintf(ofp, "    \"%s\",\n", word)
	}
	ofp.WriteString("}\n\n")

	return nil
}
