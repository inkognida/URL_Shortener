package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestGenerateLink(t *testing.T) {
	shortUrl := GenerateLink(rand.Int63())

	if len(shortUrl) != 10 {
		t.Fatalf("failed: len of shortUrl must be 10, got %d", len(shortUrl))
	}
	for c := range shortUrl {
		if !strings.Contains(charset, strconv.Itoa(c)) {
			t.Fatalf("failed: unavaliable char in shortUrl %c", c)
		}
	}
}