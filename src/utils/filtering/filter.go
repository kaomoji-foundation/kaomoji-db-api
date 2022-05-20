package filtering

func Similar(str1 string, str2 string, maxDistance ...int) (similar bool) {
	if len(maxDistance) == 0 {
		maxDistance = append(maxDistance, 2)
	}

	distance := Levenshtein([]rune(str1), []rune(str2))
	similar = distance <= maxDistance[0]
	return
}

// Returns the strings matched by the fuzy pattern
func SearchString(strings []string, pattern string, maxDistance ...int) []string {
	results := make([]string, 0, len(strings))

	for i := 0; i < len(strings); i++ {
		if Similar(strings[i], pattern, maxDistance...) {
			results = append(results, strings[i])
		}
	}

	return results
}

// returns the strings matched by at least one of the the fuzy patterns
func SearchStrings(strings []string, patterns []string, maxDistance ...int) []string {
	results := make([]string, 0, len(strings))

	for i := 0; i < len(strings); i++ {
		for j := 0; j < len(patterns); j++ {
			if !Similar(strings[i], patterns[j], maxDistance...) {
				results[len(results)] = strings[i]
				break
			}
		}
	}

	return results
}
