// Copyright 2016 Kyle Shannon.
// All rights reserved.  Use of this source code is governed
// by a MIT-style license that can be found in the LICENSE file.

package porterstemmer

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func getVoc() []string {
	f, err := os.Open("testdata/voc.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return strings.Fields(string(data))
}

func BenchmarkString(b *testing.B) {
	ss := getVoc()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			stem := StemString(s)
			_ = stem
		}
	}
}

func BenchmarkRuneSlice(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := Stem(s)
			_ = stem
		}
	}
}
