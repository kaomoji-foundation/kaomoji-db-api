package filtering

func Similar(str1 string, str2 string, maxDistance ...int) (similar bool) {
	if len(maxDistance) == 0 {
		maxDistance = append(maxDistance, 2)
	}

	distance := Levenshtein([]rune(str1), []rune(str2))
	similar = distance <= maxDistance[0]
	return
}

/*
Returns true when str1 and str2 are within max distance (default: 2),
will also return the distance between theese two independently
*/
func RankedSimilar(str1 string, str2 string, maxDistance ...int) (similar bool, distance int) {
	if len(maxDistance) == 0 {
		maxDistance = append(maxDistance, 2)
	}

	distance = Levenshtein([]rune(str1), []rune(str2))
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

/*
returns the strings matched by at least one of the the fuzy patterns,
distances arrays go from 0 to maxDistance
*/
func SearchStrings(strings, patterns []string, maxDistance ...int) (results []string, distances []int) {
	results = make([]string, 0, len(strings))
	distances = make([]int, 0, len(strings))

	for i := 0; i < len(strings); i++ {
		for j := 0; j < len(patterns); j++ {
			similar, distance := RankedSimilar(strings[i], patterns[j], maxDistance...)
			if similar {
				results = append(results, strings[i])
				distances = append(distances, distance)
				if distance == 0 {
					break
				}
			}
		}
	}

	return
}
