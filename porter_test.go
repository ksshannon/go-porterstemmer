// Copyright 2016 Kyle Shannon.
// All rights reserved.  Use of this source code is governed
// by a MIT-style license that can be found in the LICENSE file.

package porterstemmer

import (
	"bytes"
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

func TestContainsVowel(t *testing.T) {
	tests := []struct {
		s   []rune
		exp bool
	}{
		{[]rune("apple"), true},
		{[]rune("f"), false},
		{[]rune("a"), true},
		{[]rune("e"), true},
		{[]rune("i"), true},
		{[]rune("o"), true},
		{[]rune("u"), true},
		{[]rune("y"), false},
		{[]rune("cy"), true},
	}
	for _, test := range tests {
		if b := containsVowel(test.s); b != test.exp {
			t.Errorf("Did NOT get what was expected for calling containsVowel() on [%s]. Expect [%t] but got [%t]", string(test.s), test.exp, b)
		}
	}
}

// Test for issue listed here:
// https://github.com/reiver/go-porterstemmer/issues/1
//
// StemString("ion") was causing runtime exception
func TestStemStringIon(t *testing.T) {
	exp := "ion"
	s := "ion"
	stem := StemString(s)
	if stem != exp {
		t.Errorf("Input: [%s] -> Actual: [%s]. Expected: [%s]", s, stem, exp)
	}
}

const maxFuzzLen = 6

// Test inputs of English characters less than maxFuzzLen
// Added to help diagnose https://github.com/reiver/go-porterstemmer/issues/4
func TestStemFuzz(t *testing.T) {
	input := []byte{'a'}
	for len(input) < maxFuzzLen {
		panicked := false
		func() {
			defer func() { panicked = recover() != nil }()
			StemString(string(input))
		}()
		if panicked {
			t.Errorf("StemString panicked for input '%s'", input)
		}
		if allZs(input) {
			input = bytes.Repeat([]byte{'a'}, len(input)+1)
		} else {
			input = incrementBytes(input)
		}
	}
}

func incrementBytes(in []byte) []byte {
	rv := make([]byte, len(in))
	copy(rv, in)
	for i := len(rv) - 1; i >= 0; i-- {
		if rv[i]+1 == '{' {
			rv[i] = 'a'
			continue
		}
		rv[i] = rv[i] + 1
		break
	}
	return rv
}

func allZs(in []byte) bool {
	for _, b := range in {
		if b != 'z' {
			return false
		}
	}
	return true
}

func TestHasDoubleConsonantSuffix(t *testing.T) {
	tests := []struct {
		s   []rune
		exp bool
	}{
		{[]rune("apple"), false},
		{[]rune("hiss"), true},
		{[]rune("fizz"), true},
		{[]rune("fill"), true},
		{[]rune("ahaa"), false},
	}
	for _, test := range tests {
		if b := hasRepeatDoubleConsonantSuffix(test.s); b != test.exp {
			t.Errorf("Did NOT get what was expected for calling hasDoubleConsonantSuffix() on [%s]. Expect [%t] but got [%t]", string(test.s), test.exp, b)
		}
	}
}

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		s, suffix []rune
		exp       bool
	}{
		{[]rune("ran"), []rune("er"), false},
		{[]rune("runner"), []rune("er"), true},
		{[]rune("runnar"), []rune("er"), false},
		{[]rune("runned"), []rune("er"), false},
		{[]rune("runnre"), []rune("er"), false},
		// FIXME marty changed Expected
		// to false here because it seems
		// the contract does not support
		// suffix of same length as input
		// as this test implied
		{[]rune("er"), []rune("er"), false},
		{[]rune("re"), []rune("er"), false},
		{[]rune("ran"), []rune("ER"), false},
		{[]rune("runner"), []rune("ER"), false},
		{[]rune("runnar"), []rune("ER"), false},
		{[]rune("runned"), []rune("ER"), false},
		{[]rune("runnre"), []rune("ER"), false},
		{[]rune("er"), []rune("ER"), false},
		{[]rune("re"), []rune("ER"), false},
		{[]rune(""), []rune("er"), false},
		{[]rune("e"), []rune("er"), false},
		{[]rune("caresses"), []rune("sses"), true},
		{[]rune("ponies"), []rune("ies"), true},
		{[]rune("caress"), []rune("ss"), true},
		{[]rune("cats"), []rune("s"), true},
		{[]rune("feed"), []rune("eed"), true},
		{[]rune("agreed"), []rune("eed"), true},
		{[]rune("plastered"), []rune("ed"), true},
		{[]rune("bled"), []rune("ed"), true},
		{[]rune("motoring"), []rune("ing"), true},
		{[]rune("sing"), []rune("ing"), true},
		{[]rune("conflat"), []rune("at"), true},
		{[]rune("troubl"), []rune("bl"), true},
		{[]rune("siz"), []rune("iz"), true},
		{[]rune("happy"), []rune("y"), true},
		{[]rune("sky"), []rune("y"), true},
		{[]rune("relational"), []rune("ational"), true},
		{[]rune("conditional"), []rune("tional"), true},
		{[]rune("rational"), []rune("tional"), true},
		{[]rune("valenci"), []rune("enci"), true},
		{[]rune("hesitanci"), []rune("anci"), true},
		{[]rune("digitizer"), []rune("izer"), true},
		{[]rune("conformabli"), []rune("abli"), true},
		{[]rune("radicalli"), []rune("alli"), true},
		{[]rune("differentli"), []rune("entli"), true},
		{[]rune("vileli"), []rune("eli"), true},
		{[]rune("analogousli"), []rune("ousli"), true},
		{[]rune("vietnamization"), []rune("ization"), true},
		{[]rune("predication"), []rune("ation"), true},
		{[]rune("operator"), []rune("ator"), true},
		{[]rune("feudalism"), []rune("alism"), true},
		{[]rune("decisiveness"), []rune("iveness"), true},
		{[]rune("hopefulness"), []rune("fulness"), true},
		{[]rune("callousness"), []rune("ousness"), true},
		{[]rune("formaliti"), []rune("aliti"), true},
		{[]rune("sensitiviti"), []rune("iviti"), true},
		{[]rune("sensibiliti"), []rune("biliti"), true},
		{[]rune("triplicate"), []rune("icate"), true},
		{[]rune("formative"), []rune("ative"), true},
		{[]rune("formalize"), []rune("alize"), true},
		{[]rune("electriciti"), []rune("iciti"), true},
		{[]rune("electrical"), []rune("ical"), true},
		{[]rune("hopeful"), []rune("ful"), true},
		{[]rune("goodness"), []rune("ness"), true},
		{[]rune("revival"), []rune("al"), true},
		{[]rune("allowance"), []rune("ance"), true},
		{[]rune("inference"), []rune("ence"), true},
		{[]rune("airliner"), []rune("er"), true},
		{[]rune("gyroscopic"), []rune("ic"), true},
		{[]rune("adjustable"), []rune("able"), true},
		{[]rune("defensible"), []rune("ible"), true},
		{[]rune("irritant"), []rune("ant"), true},
		{[]rune("replacement"), []rune("ement"), true},
		{[]rune("adjustment"), []rune("ment"), true},
		{[]rune("dependent"), []rune("ent"), true},
		{[]rune("adoption"), []rune("ion"), true},
		{[]rune("homologou"), []rune("ou"), true},
		{[]rune("communism"), []rune("ism"), true},
		{[]rune("activate"), []rune("ate"), true},
		{[]rune("angulariti"), []rune("iti"), true},
		{[]rune("homologous"), []rune("ous"), true},
		{[]rune("effective"), []rune("ive"), true},
		{[]rune("bowdlerize"), []rune("ize"), true},
		{[]rune("probate"), []rune("e"), true},
		{[]rune("rate"), []rune("e"), true},
		{[]rune("cease"), []rune("e"), true},
	}
	for _, test := range tests {
		if b := hasSuffix(test.s, test.suffix); b != test.exp {
			t.Errorf("Did NOT get what was expected for calling hasSuffix() on [%s] with suffix [%s]. Expect [%t] but got [%t]", string(test.s), string(test.suffix), test.exp, b)
		}
	}
}

func TestIsConsontant(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []bool
	}{
		{[]rune("apple"), []bool{false, true, true, true, false}},
		{[]rune("cyan"), []bool{true, false, false, true}},
		{[]rune("connects"), []bool{true, false, true, true, false, true, true, true}},
		{[]rune("yellow"), []bool{true, false, true, true, false, true}},
		{[]rune("excellent"), []bool{false, true, true, false, true, true, false, true, true}},
		{[]rune("yuk"), []bool{true, false, true}},
		{[]rune("syzygy"), []bool{true, false, true, false, true, false}},
		{[]rune("school"), []bool{true, true, true, false, false, true}},
		{[]rune("pay"), []bool{true, false, true}},
		{[]rune("golang"), []bool{true, false, true, false, true, true}},
		// NOTE: The Porter Stemmer technical should make a mistake on the second "y".
		//       Really, both the 1st and 2nd "y" are consontants. But
		{[]rune("sayyid"), []bool{true, false, true, false, false, true}},
		{[]rune("ya"), []bool{true, false}},
	}
	for _, test := range tests {
		for i := 0; i < len(test.s); i++ {
			if b := isConsonant(test.s, i); b != test.exp[i] {
				t.Errorf("Did NOT get what was expected for calling isConsonant() on [%s] at [%d] (i.e., [%s]). Expect [%t] but got [%t]", string(test.s), i, string(test.s[i]), test.exp[i], b)
			}
		}
	}
}

func TestMeasure(t *testing.T) {
	tests := []struct {
		s   []rune
		exp uint
	}{
		{[]rune("ya"), 0},
		{[]rune("cyan"), 1},
		{[]rune("connects"), 2},
		{[]rune("yellow"), 2},
		{[]rune("excellent"), 3},
		{[]rune("yuk"), 1},
		{[]rune("syzygy"), 2},
		{[]rune("school"), 1},
		{[]rune("pay"), 1},
		{[]rune("golang"), 2},
		// NOTE: The Porter Stemmer technical should make a mistake on the second "y".
		//       Really, both the 1st and 2nd "y" are consontants. But
		{[]rune("sayyid"), 2},
		{[]rune("ya"), 0},
		{[]rune(""), 0},
		{[]rune("tr"), 0},
		{[]rune("ee"), 0},
		{[]rune("tree"), 0},
		{[]rune("t"), 0},
		{[]rune("by"), 0},
		{[]rune("trouble"), 1},
		{[]rune("oats"), 1},
		{[]rune("trees"), 1},
		{[]rune("ivy"), 1},
		{[]rune("troubles"), 2},
		{[]rune("private"), 2},
		{[]rune("oaten"), 2},
		{[]rune("orrery"), 2},
	}
	for _, test := range tests {
		if n := measure(test.s); n != test.exp {
			t.Errorf("Did NOT get what was expected for calling measure() on [%s]. Expect [%d] but got [%d]", string(test.s), test.exp, n)
		}
	}
}

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

func TestStemWithoutLowerCasing(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("controll"), []rune("control")},
		{[]rune("roll"), []rune("roll")},
	}

	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = StemWithoutLowerCasing(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling StemWithoutLowerCasing() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}

func TestStep1a(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("caresses"), []rune("caress")},
		{[]rune("ponies"), []rune("poni")},
		{[]rune("ties"), []rune("ti")},
		{[]rune("caress"), []rune("caress")},
		{[]rune("cats"), []rune("cat")},
	}
	for _, test := range tests {
		for i := 0; i < len(test.s); i++ {
			stem := make([]rune, len(test.s))
			copy(stem, test.s)
			stem = step1a(stem)
			equal := true
			if 0 == len(stem) && 0 == len(test.exp) {
				equal = true
			} else if len(stem) != len(test.exp) {
				equal = false
			} else if stem[0] != test.exp[0] {
				equal = false
			} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
				equal = false
			} else {
				for j := 0; j < len(stem); j++ {
					if stem[j] != test.exp[j] {
						equal = false
					}
				}
			}
			if !equal {
				t.Errorf("Did NOT get what was expected for calling step1a() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
			}
		}
	}
}

func TestStep1b(t *testing.T) {

	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("feed"), []rune("feed")},
		{[]rune("agreed"), []rune("agree")},
		{[]rune("plastered"), []rune("plaster")},
		{[]rune("bled"), []rune("bled")},
		{[]rune("motoring"), []rune("motor")},
		{[]rune("sing"), []rune("sing")},
		{[]rune("conflated"), []rune("conflate")},
		{[]rune("troubled"), []rune("trouble")},
		{[]rune("sized"), []rune("size")},
		{[]rune("hopping"), []rune("hop")},
		{[]rune("tanned"), []rune("tan")},
		{[]rune("falling"), []rune("fall")},
		{[]rune("hissing"), []rune("hiss")},
		{[]rune("fizzed"), []rune("fizz")},
		{[]rune("failing"), []rune("fail")},
		{[]rune("filing"), []rune("file")},
	}
	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = step1b(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling step1b() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}

func TestStep1c(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("happy"), []rune("happi")},
		{[]rune("sky"), []rune("sky")},
		{[]rune("apology"), []rune("apologi")},
	}
	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = step1c(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling step1c() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}

func TestStep2(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("relational"), []rune("relate")},
		{[]rune("conditional"), []rune("condition")},
		{[]rune("rational"), []rune("rational")},
		{[]rune("valenci"), []rune("valence")},
		{[]rune("hesitanci"), []rune("hesitance")},
		{[]rune("digitizer"), []rune("digitize")},
		{[]rune("conformabli"), []rune("conformable")},
		{[]rune("radicalli"), []rune("radical")},
		{[]rune("differentli"), []rune("different")},
		{[]rune("vileli"), []rune("vile")},
		{[]rune("analogousli"), []rune("analogous")},
		{[]rune("vietnamization"), []rune("vietnamize")},
		{[]rune("predication"), []rune("predicate")},
		{[]rune("operator"), []rune("operate")},
		{[]rune("feudalism"), []rune("feudal")},
		{[]rune("decisiveness"), []rune("decisive")},
		{[]rune("hopefulness"), []rune("hopeful")},
		{[]rune("callousness"), []rune("callous")},
		{[]rune("formaliti"), []rune("formal")},
		{[]rune("sensitiviti"), []rune("sensitive")},
		{[]rune("sensibiliti"), []rune("sensible")},
	}
	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = step2(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling step2() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}

func TestStep3(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("triplicate"), []rune("triplic")},
		{[]rune("formative"), []rune("form")},
		{[]rune("formalize"), []rune("formal")},
		{[]rune("electriciti"), []rune("electric")},
		{[]rune("electrical"), []rune("electric")},
		{[]rune("hopeful"), []rune("hope")},
		{[]rune("goodness"), []rune("good")},
	}
	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = step3(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling step3() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}

func TestStep4(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("revival"), []rune("reviv")},
		{[]rune("allowance"), []rune("allow")},
		{[]rune("inference"), []rune("infer")},
		{[]rune("airliner"), []rune("airlin")},
		{[]rune("gyroscopic"), []rune("gyroscop")},
		{[]rune("adjustable"), []rune("adjust")},
		{[]rune("defensible"), []rune("defens")},
		{[]rune("irritant"), []rune("irrit")},
		{[]rune("replacement"), []rune("replac")},
		{[]rune("adjustment"), []rune("adjust")},
		{[]rune("dependent"), []rune("depend")},
		{[]rune("adoption"), []rune("adopt")},
		{[]rune("homologou"), []rune("homolog")},
		{[]rune("communism"), []rune("commun")},
		{[]rune("activate"), []rune("activ")},
		{[]rune("angulariti"), []rune("angular")},
		{[]rune("homologous"), []rune("homolog")},
		{[]rune("effective"), []rune("effect")},
		{[]rune("bowdlerize"), []rune("bowdler")},
	}
	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = step4(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling step4() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}

func TestStep5a(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("probate"), []rune("probat")},
		{[]rune("rate"), []rune("rate")},
		{[]rune("cease"), []rune("ceas")},
	}
	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = step5a(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling step5a() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}

func TestStep5b(t *testing.T) {
	tests := []struct {
		s   []rune
		exp []rune
	}{
		{[]rune("controll"), []rune("control")},
		{[]rune("roll"), []rune("roll")},
	}
	for _, test := range tests {
		stem := make([]rune, len(test.s))
		copy(stem, test.s)
		stem = step5b(stem)
		equal := true
		if 0 == len(stem) && 0 == len(test.exp) {
			equal = true
		} else if len(stem) != len(test.exp) {
			equal = false
		} else if stem[0] != test.exp[0] {
			equal = false
		} else if stem[len(stem)-1] != test.exp[len(test.exp)-1] {
			equal = false
		} else {
			for j := 0; j < len(stem); j++ {
				if stem[j] != test.exp[j] {
					equal = false
				}
			}
		}
		if !equal {
			t.Errorf("Did NOT get what was expected for calling step5b() on [%s]. Expect [%s] but got [%s]", string(test.s), string(test.exp), string(stem))
		}
	}
}
