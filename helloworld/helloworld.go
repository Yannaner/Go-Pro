package helloworld

import "fmt"

// BubbleSort sorts an array of integers using the bubble sort algorithm
func BubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				// Swap arr[j] and arr[j+1]
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// PrintArray prints the elements of an array
func PrintArray(arr []int) {
	for _, value := range arr {
		fmt.Printf("%d ", value)
	}
	fmt.Println()
}
