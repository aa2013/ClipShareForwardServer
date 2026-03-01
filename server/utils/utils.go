package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
)

func IntToSizeStr(size int64) string {
	if size < 0 {
		return "-"
	}

	const kb = 1024
	const mb = kb * 1024
	const gb = mb * 1024

	switch {
	case size >= gb:
		return strconv.FormatFloat(float64(size)/float64(gb), 'f', 2, 64) + " GB"
	case size >= mb:
		return strconv.FormatFloat(float64(size)/float64(mb), 'f', 2, 64) + " MB"
	case size >= kb:
		return strconv.FormatFloat(float64(size)/float64(kb), 'f', 2, 64) + " KB"
	default:
		return strconv.FormatInt(size, 10) + " B"
	}
}

func GenRandomString(length int) (string, error) {
	bytes := make([]byte, length/2) // 每个字节对应两个十六进制字符
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SimpleIf 泛型三元函数，支持不同类型
func SimpleIf[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
