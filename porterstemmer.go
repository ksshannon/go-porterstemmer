package porterstemmer

import "unicode"

// isConsonant returns true if the rune represents a constanant.  Y is regarded
// a constanant if it starts the word, or is followed by a vowel.
func isConsonant(s []rune, i int) bool {
	switch s[i] {
	case 'a', 'e', 'i', 'o', 'u':
		return false
	case 'y':
		if i == 0 {
			return true
		} else {
			return !isConsonant(s, i-1)
		}
	default:
		return true
	}
}

func measure(s []rune) uint {
	if len(s) == 0 {
		return 0
	}

	lenS := len(s)
	m := uint(0)

	// Ignore (potential) consonant sequence at the beginning of word.
	i := 0
	for i = 0; i < len(s) && isConsonant(s, i); i++ {
	}
	if i == len(s) {
		return 0
	}

	// For each pair of a vowel sequence followed by a consonant sequence, increment result.
Outer:
	for i < len(s) {
		for !isConsonant(s, i) {
			i++
			if i >= lenS {
				break Outer
			}
		}
		for isConsonant(s, i) {
			i++
			if i >= lenS {
				m++
				break Outer
			}
		}
		m++
	}
	return m
}

// hasSuffix checks if a word has a specific suffix
func hasSuffix(s, suffix []rune) bool {
	lenSMinusOne := len(s) - 1
	lenSuffixMinusOne := len(suffix) - 1

	if lenSMinusOne <= lenSuffixMinusOne {
		return false
	} else if s[lenSMinusOne] != suffix[lenSuffixMinusOne] { // I suspect checking this first should speed this function up in practice.
		return false
	} else {

		for i := 0; i < lenSuffixMinusOne; i++ {

			if suffix[i] != s[lenSMinusOne-lenSuffixMinusOne+i] {
				/////////////// RETURN
				return false
			}

		}

	}

	return true
}

// containsVowel returns true if the string has a vowel
func containsVowel(s []rune) bool {
	for i := 0; i < len(s); i++ {
		if !isConsonant(s, i) {
			return true
		}
	}
	return false
}

func hasRepeatDoubleConsonantSuffix(s []rune) bool {

	// Initialize.
	lenS := len(s)

	result := false

	// Do it!
	if 2 > lenS {
		result = false
	} else if s[lenS-1] == s[lenS-2] && isConsonant(s, lenS-1) { // Will using isConsonant() cause a problem with "YY"?
		result = true
	} else {
		result = false
	}

	// Return,
	return result
}

func hasConsonantVowelConsonantSuffix(s []rune) bool {

	// Initialize.
	lenS := len(s)

	result := false

	// Do it!
	if 3 > lenS {
		result = false
	} else if isConsonant(s, lenS-3) && !isConsonant(s, lenS-2) && isConsonant(s, lenS-1) {
		result = true
	} else {
		result = false
	}

	// Return
	return result
}

func step1a(s []rune) []rune {

	// Initialize.
	var result []rune = s

	lenS := len(s)

	// Do it!
	if suffix := []rune("sses"); hasSuffix(s, suffix) {

		lenTrim := 2

		subSlice := s[:lenS-lenTrim]

		result = subSlice
	} else if suffix := []rune("ies"); hasSuffix(s, suffix) {
		lenTrim := 2

		subSlice := s[:lenS-lenTrim]

		result = subSlice
	} else if suffix := []rune("ss"); hasSuffix(s, suffix) {

		result = s
	} else if suffix := []rune("s"); hasSuffix(s, suffix) {

		lenSuffix := 1

		subSlice := s[:lenS-lenSuffix]

		result = subSlice
	}

	// Return.
	return result
}

func step1b(s []rune) []rune {

	// Initialize.
	var result []rune = s

	lenS := len(s)

	// Do it!
	if suffix := []rune("eed"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 0 < m {
			lenTrim := 1

			result = s[:lenS-lenTrim]
		}
	} else if suffix := []rune("ed"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		if containsVowel(subSlice) {

			if suffix2 := []rune("at"); hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]
			} else if suffix2 := []rune("bl"); hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]
			} else if suffix2 := []rune("iz"); hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]
			} else if c := subSlice[len(subSlice)-1]; 'l' != c && 's' != c && 'z' != c && hasRepeatDoubleConsonantSuffix(subSlice) {
				lenTrim := 1

				lenSubSlice := len(subSlice)

				result = subSlice[:lenSubSlice-lenTrim]
			} else if c := subSlice[len(subSlice)-1]; 1 == measure(subSlice) && hasConsonantVowelConsonantSuffix(subSlice) && 'w' != c && 'x' != c && 'y' != c {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else {
				result = subSlice
			}

		}
	} else if suffix := []rune("ing"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		if containsVowel(subSlice) {

			if suffix2 := []rune("at"); hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else if suffix2 := []rune("bl"); hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else if suffix2 := []rune("iz"); hasSuffix(subSlice, suffix2) {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else if c := subSlice[len(subSlice)-1]; 'l' != c && 's' != c && 'z' != c && hasRepeatDoubleConsonantSuffix(subSlice) {
				lenTrim := 1

				lenSubSlice := len(subSlice)

				result = subSlice[:lenSubSlice-lenTrim]
			} else if c := subSlice[len(subSlice)-1]; 1 == measure(subSlice) && hasConsonantVowelConsonantSuffix(subSlice) && 'w' != c && 'x' != c && 'y' != c {
				lenTrim := -1

				result = s[:lenS-lenSuffix-lenTrim]

				result[len(result)-1] = 'e'
			} else {
				result = subSlice
			}

		}
	}

	// Return.
	return result
}

func step1c(s []rune) []rune {

	// Initialize.
	lenS := len(s)

	result := s

	// Do it!
	if 2 > lenS {
		/////////// RETURN
		return result
	}

	if 'y' == s[lenS-1] && containsVowel(s[:lenS-1]) {

		result[lenS-1] = 'i'

	} else if 'Y' == s[lenS-1] && containsVowel(s[:lenS-1]) {

		result[lenS-1] = 'I'

	}

	// Return.
	return result
}

func step2(s []rune) []rune {

	// Initialize.
	lenS := len(s)

	result := s

	// Do it!
	if suffix := []rune("ational"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-5] = 'e'
			result = result[:lenS-4]
		}
	} else if suffix := []rune("tional"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = result[:lenS-2]
		}
	} else if suffix := []rune("enci"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-1] = 'e'
		}
	} else if suffix := []rune("anci"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-1] = 'e'
		}
	} else if suffix := []rune("izer"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-1]
		}
	} else if suffix := []rune("bli"); hasSuffix(s, suffix) { // --DEPARTURE--
		//		} else if suffix := []rune("abli") ; hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-1] = 'e'
		}
	} else if suffix := []rune("alli"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("entli"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("eli"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("ousli"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-2]
		}
	} else if suffix := []rune("ization"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-5] = 'e'

			result = s[:lenS-4]
		}
	} else if suffix := []rune("ation"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-3] = 'e'

			result = s[:lenS-2]
		}
	} else if suffix := []rune("ator"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-2] = 'e'

			result = s[:lenS-1]
		}
	} else if suffix := []rune("alism"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-3]
		}
	} else if suffix := []rune("iveness"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-4]
		}
	} else if suffix := []rune("fulness"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-4]
		}
	} else if suffix := []rune("ousness"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-4]
		}
	} else if suffix := []rune("aliti"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result = s[:lenS-3]
		}
	} else if suffix := []rune("iviti"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-3] = 'e'

			result = result[:lenS-2]
		}
	} else if suffix := []rune("biliti"); hasSuffix(s, suffix) {
		if 0 < measure(s[:lenS-len(suffix)]) {
			result[lenS-5] = 'l'
			result[lenS-4] = 'e'

			result = result[:lenS-3]
		}
	} else if suffix := []rune("logi"); hasSuffix(s, suffix) { // --DEPARTURE--
		if 0 < measure(s[:lenS-len(suffix)]) {
			lenTrim := 1

			result = s[:lenS-lenTrim]
		}
	}

	// Return.
	return result
}

func step3(s []rune) []rune {

	// Initialize.
	lenS := len(s)
	result := s

	// Do it!
	if suffix := []rune("icate"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if 0 < measure(s[:lenS-lenSuffix]) {
			result = result[:lenS-3]
		}
	} else if suffix := []rune("ative"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 0 < m {
			result = subSlice
		}
	} else if suffix := []rune("alize"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if 0 < measure(s[:lenS-lenSuffix]) {
			result = result[:lenS-3]
		}
	} else if suffix := []rune("iciti"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if 0 < measure(s[:lenS-lenSuffix]) {
			result = result[:lenS-3]
		}
	} else if suffix := []rune("ical"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		if 0 < measure(s[:lenS-lenSuffix]) {
			result = result[:lenS-2]
		}
	} else if suffix := []rune("ful"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 0 < m {
			result = subSlice
		}
	} else if suffix := []rune("ness"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 0 < m {
			result = subSlice
		}
	}

	// Return.
	return result
}

func step4(s []rune) []rune {

	// Initialize.
	lenS := len(s)
	result := s

	// Do it!
	if suffix := []rune("al"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = result[:lenS-lenSuffix]
		}
	} else if suffix := []rune("ance"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = result[:lenS-lenSuffix]
		}
	} else if suffix := []rune("ence"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = result[:lenS-lenSuffix]
		}
	} else if suffix := []rune("er"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ic"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("able"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ible"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ant"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ement"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ment"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ent"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ion"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		c := subSlice[len(subSlice)-1]

		if 1 < m && ('s' == c || 't' == c) {
			result = subSlice
		}
	} else if suffix := []rune("ou"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ism"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ate"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("iti"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ous"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ive"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	} else if suffix := []rune("ize"); hasSuffix(s, suffix) {
		lenSuffix := len(suffix)

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	}

	// Return.
	return result
}

func step5a(s []rune) []rune {

	// Initialize.
	lenS := len(s)
	result := s

	// Do it!
	if 'e' == s[lenS-1] {
		lenSuffix := 1

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		} else if 1 == m {
			if c := subSlice[len(subSlice)-1]; !(hasConsonantVowelConsonantSuffix(subSlice) && 'w' != c && 'x' != c && 'y' != c) {
				result = subSlice
			}
		}
	}

	// Return.
	return result
}

func step5b(s []rune) []rune {

	lenS := len(s)
	result := s

	// Do it!
	if 2 < lenS && 'l' == s[lenS-2] && 'l' == s[lenS-1] {

		lenSuffix := 1

		subSlice := s[:lenS-lenSuffix]

		m := measure(subSlice)

		if 1 < m {
			result = subSlice
		}
	}

	// Return.
	return result
}

// StemString converts a string to a rune array, then stems the result.
func StemString(s string) string {
	ra := []rune(s)
	ra = Stem(ra)
	return string(ra)
}

// Stem converts the runes to lower case, then stems the lowercase runes
func Stem(s []rune) []rune {
	if len(s) == 0 {
		return s
	}
	for i := 0; i < len(s); i++ {
		s[i] = unicode.ToLower(s[i])
	}
	result := StemWithoutLowerCasing(s)
	return result
}

// StemWithoutLowerCasing applies the stemming assuming that the runes are
// lowercase.
func StemWithoutLowerCasing(s []rune) []rune {
	if len(s) <= 2 {
		return s
	}
	s = step1a(s)
	s = step1b(s)
	s = step1c(s)
	s = step2(s)
	s = step3(s)
	s = step4(s)
	s = step5a(s)
	s = step5b(s)
	return s
}
