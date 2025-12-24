package main

import (
	"fmt"
	"testing"
)

var users []*User

func init() {
	for i := 0; i < 1000; i++ {
		users = append(users, &User{Id: i, Name: fmt.Sprintf("user%d", i)})
	}
}
func BenchmarkGenerateIdsRaw(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateIdsRaw(users)
	}
}

func BenchmarkGenerateIdsBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateIdsBuilder(users)
	}
}

func BenchmarkGenerateIdsStrconv(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateIdsStrconv(users)
	}
}
