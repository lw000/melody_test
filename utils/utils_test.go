package utils

import (
	"log"
	"testing"
)

func TestMaxmin(t *testing.T) {
	max, min := Maxmin(20, 30)
	log.Printf("max: %d, min: %d\n", max, min)
}

func TestMax(t *testing.T) {
	max := Max(20, 30)
	log.Printf("max: %d\n", max)
}

func TestMin(t *testing.T) {
	min := Min(20, 30)
	log.Printf("min: %d\n", min)
}

func TestCompareMapStringString(t *testing.T) {
	m1 := map[string]string{
		"a": "aaaa",
		"b": "bbbb",
	}

	m2 := map[string]string{
		"a": "aaaa1",
		"b": "bbbb",
	}

	if CompareMapStringString(m1, m2) {
		log.Println("m1 == m2")
	}
}

func TestRandomIntger(t *testing.T) {
	log.Println("RandomIntger:", RandomIntger(4))
}

func TestRandomString(t *testing.T) {
	log.Println("RandomString:", RandomString(32))
}

func TestGenerateSID(t *testing.T) {
	log.Println("GenerateSID:", GenerateSID())
}
