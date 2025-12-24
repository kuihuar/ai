package main

import (
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleRoot(tt.args.w, tt.args.r)
		})
	}
}

var slices []int

func init() {
	slices = getRawSlices()
}
func BenchmarkGetRawSet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getRawSet(slices)
	}
}

func BenchmarkGetEmptyStructSet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getEmptyStructSet(slices)
	}
}

func BenchmarkGetCapacitySet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		getCapacitySet(slices)
	}
}
