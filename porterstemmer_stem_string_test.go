package porterstemmer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const testDir = "testdata"

const (
	vocURL    = "http://tartarus.org/martin/PorterStemmer/voc.txt"
	outputURL = "http://tartarus.org/martin/PorterStemmer/output.txt"
)

// init downloads the appropriate files, if necessary.
func init() {
	_, err := os.Stat(testDir)
	if nil != err {
		_ = os.Mkdir(testDir, 0755)
	}
	_, err = os.Stat(testDir)
	if nil != err {
		panic(fmt.Sprintf("The test data folder ([%s]) does not exists (and could not create it).", testDir))
	}
	for _, u := range []string{vocURL, outputURL} {
		fname := filepath.Join(testDir, path.Base(u))
		_, err = os.Stat(fname)
		if err == nil {
			continue
		}
		resp, err := http.Get(u)
		if err != nil {
			panic(fmt.Sprintf("could not download %s", u))
		}
		defer resp.Body.Close()
		fout, err := os.Create(fname)
		if err != nil {
			panic(fmt.Sprintf("could not download %s", u))
		}
		defer fout.Close()
		_, err = io.Copy(fout, resp.Body)
		if err != nil {
			panic(fmt.Sprintf("could not download %s", u))
		}
	}
}

func TestStemString(t *testing.T) {

	v, err := os.Open(filepath.Join(testDir, path.Base(vocURL)))
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer v.Close()
	o, err := os.Open(filepath.Join(testDir, path.Base(outputURL)))
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer o.Close()

	data, err := ioutil.ReadAll(v)
	if err != nil {
		t.Fatalf("%s", err)
	}
	vs := strings.Fields(string(data))

	data, err = ioutil.ReadAll(o)
	if err != nil {
		t.Fatalf("%s", err)
	}
	os := strings.Fields(string(data))

	for i, word := range vs {
		stem := StemString(word)
		if stem != os[i] {
			t.Errorf("Input: [%s] -> Actual: [%s]. Expected: [%s]", word, stem, os[i])
		}
	}
}
