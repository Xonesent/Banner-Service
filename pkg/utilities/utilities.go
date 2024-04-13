package utilities

func InStringSlice(target string, src []string) bool {
	for _, el := range src {
		if el == target {
			return true
		}
	}
	return false
}

func RemoveDuplicates[T comparable](slice []T) []T {
	keys := make(map[T]bool)
	var list []T
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func AreSlicesEqual[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	counts := make(map[T]int)
	for _, val := range slice1 {
		counts[val]++
	}

	for _, val := range slice2 {
		counts[val]--
		if counts[val] < 0 {
			return false
		}
	}

	return true
}

func FindUniqueElements[T comparable](firstSlice, secondSlice []T) []T {
	secondMap := make(map[T]bool)
	for _, value := range secondSlice {
		secondMap[value] = true
	}

	var uniqueElements []T
	for _, value := range firstSlice {
		if _, found := secondMap[value]; !found {
			uniqueElements = append(uniqueElements, value)
		}
	}

	return uniqueElements
}
