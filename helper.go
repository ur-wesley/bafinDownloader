package main

func RemoveEmptyLines(array []string) []string {
	for i := 0; i < len(array); i++ {
		if array[i] == "" {
			array = append(array[:i], array[i+1:]...)
			i--
		}
	}
	return array
}