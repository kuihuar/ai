package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	slices := getCapacitySlices()
	getEmptyStructSet(slices)
	fmt.Fprintln(w, "Hello, Root!")

}
func main() {
	http.HandleFunc("/", HandleRoot)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}

func getRawSlices() []int {
	// fmt.Println("exec getRawSlices")
	n := 100000
	slices := make([]int, 0)
	for i := 0; i < n; i++ {
		slices = append(slices, i)

	}
	return slices
}
func getRawSet(slices []int) map[int]bool {
	// fmt.Println("exec getRawSet")
	set := make(map[int]bool, 0)
	for _, item := range slices {
		set[item] = true
	}
	return set
}
func getCapacitySlices() []int {
	n := 100000
	slices := make([]int, 0, n)
	for i := 0; i < n; i++ {
		slices = append(slices, i)
	}
	return slices
}

func getEmptyStructSet(slices []int) map[int]struct{} {
	set := make(map[int]struct{}, 0)
	for _, item := range slices {
		set[item] = struct{}{}
	}
	return set
}
func getCapacitySet(slices []int) map[int]struct{} {
	set := make(map[int]struct{}, len(slices))
	for _, item := range slices {
		set[item] = struct{}{}
	}
	return set
}
