package utils

// 通过用户ID生成6位随机邀请码

import (
	"errors"
)

const (
	PRIME1    int64 = 11
	PRIME2          = 7
	codeLen         = 6
	SLAT      int64 = 20200810
	sourceLen int64 = 32
)

var (
	source = []string{"F", "L", "G", "W", "5", "X", "C", "3",
		"9", "Z", "M", "6", "7", "Y", "R", "T", "2", "H", "S", "8", "D", "V", "E", "J", "4", "K",
		"Q", "P", "U", "A", "N", "B"}
)

// 用户ID生成邀请码
func GenInviteCode(uid int64) string {
	id := uid*PRIME1 + SLAT
	sourceCode := make([]int64, codeLen)
	sourceCodeTmp := make([]int64, codeLen)
	sourceCode[0] = id
	var code string
	for i := 0; i < codeLen-1; i++ {
		sourceCode[i+1] = sourceCode[i] / sourceLen
		sourceCode[i] = (sourceCode[i] + int64(i)*sourceCode[0]) % sourceLen
	}
	sourceCode[5] = (sourceCode[0] + sourceCode[1] + sourceCode[2] + sourceCode[3] + sourceCode[4]) * PRIME1 % sourceLen
	for i := 0; i < codeLen; i++ {
		sourceCodeTmp[i] = sourceCode[i*PRIME2%codeLen]
	}
	for _, v := range sourceCodeTmp {
		code += string(source[v])
	}
	return code
}

// 邀请码转用户ID

func DecodeInviteCode(code string) (int64, error) {
	if len(code) != codeLen {
		return -1, errors.New("code len error")
	}
	codeTmp := []byte(code)
	a := make([]int64, codeLen)
	b := make([]int64, codeLen)
	for k, v := range codeTmp {
		str := string(v)
		index := findCharIndex(str)
		if index == -1 {
			return -1, errors.New("code char error")
		}
		a[k*PRIME2%codeLen] = int64(index)
	}

	for i := codeLen - 2; i >= 0; i-- {
		b[i] = (a[i] - a[0]*int64(i) + sourceLen*int64(i)) % int64(sourceLen)
	}

	var res int64
	for i := codeLen - 2; i >= 0; i-- {
		res += b[i]
		if i > 0 {
			res *= int64(sourceLen)
		} else {
			res *= int64(1)
		}
	}
	return (res - SLAT) / PRIME1, nil
}

func findCharIndex(c string) int {
	for k, v := range source {
		if v == c {
			return k
		}
	}
	return -1
}
