package utils_test

import (
	"testing"

	"url-shortener/utils"
)

func TestRandomStringLength(t *testing.T) {
	length := 5
	str := utils.RandomString(length)
	if len(str) != length {
		t.Fatalf("Random string generated of length %d instead of %d", len(str), length)
	}
}

func TestRandomStringsUnique(t *testing.T) {
	length := 5
	str1 := utils.RandomString(length)
	str2 := utils.RandomString(length)

	if str1 == str2 {
		t.Fatalf("Random generated strings are not unique, %s == %s", str1, str2)
	}
}
