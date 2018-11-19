package utils

func UcFirst(str string) (string) {
	strRune := []rune(str)
	if strRune[0] >= 97 && strRune[0] <= 122 {
		strRune[0] -= 32
	}
	return string(strRune)
}

func LcFirst(str string) (string) {
	strRune := []rune(str)
	if strRune[0] >= 65 && strRune[0] <= 90 {
		strRune[0] += 32
	}
	return string(strRune)
}
