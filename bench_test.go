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

func BenchmarkStep1A(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step1a(s)
			_ = stem
		}
	}
}

func BenchmarkStep1B(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step1b(s)
			_ = stem
		}
	}
}

func BenchmarkStep1C(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step1c(s)
			_ = stem
		}
	}
}

func BenchmarkStep2(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step2(s)
			_ = stem
		}
	}
}

func BenchmarkStep3(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step3(s)
			_ = stem
		}
	}
}

func BenchmarkStep4(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step1a(s)
			_ = stem
		}
	}
}

func BenchmarkStep5A(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step5a(s)
			_ = stem
		}
	}
}

func BenchmarkStep5B(b *testing.B) {
	ss := getVoc()
	rs := make([][]rune, len(ss))
	for i, s := range ss {
		rs[i] = []rune(s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range rs {
			stem := step5b(s)
			_ = stem
		}
	}
}
