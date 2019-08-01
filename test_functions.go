package gol

import "fmt"

func matchSlice(a [][]uint8, b [][]uint8) bool {
	for y := range a {
		for x := range a[y] {
			if a[y][x] != b[y][x] {
				fmt.Println(x, y)
				return false
			}
		}
	}
	return true
}

func printArray(array [][]uint8) {
	for idx := range array {
		fmt.Println(array[idx])
	}
}

func mismatchCheck(expected [][]uint8, got [][]uint8) bool {
	if !matchSlice(expected, got) {
		fmt.Println("Expected:")
		printArray(expected)
		fmt.Println("Got:")
		printArray(got)
		return true
	}
	return false
}
