package hasher

import (
	"math"
	"strings"
)

/*
сгенерировать случайное число из uint64 и из alphanumeric собрать строку с перестановками
Ресурсы wiki https://ru.wikipedia.org/wiki/Перестановка
*/

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenHash(num uint64) []byte {
	var hash []byte
	for ; num > 0; num = num / uint64(len(alphanumeric)) {
		hash = append(hash, alphanumeric[(num%uint64(len(alphanumeric)))])
	}
	return hash
}

func GenClear(s string) (uint64, error) {
	var num uint64
	for i, ch := range s {
		num += uint64(strings.IndexRune(alphanumeric, ch)) * uint64(math.Pow(float64(len(alphanumeric)), float64(i)))
	}
	return num, nil
}
