package util

import (
	"math/rand"
	"time"
)

var (
	strLetters       = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numLetters       = []rune("0123456789")
	numAndStrLetters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

/**
生成指定长度随机数
*/
func GetRandomNumber(size int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, size)
	for i := range b {
		b[i] = numLetters[rand.Intn(len(numLetters))]
	}
	return string(b)
}
