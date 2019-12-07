package order

const alphabetNum = 26
const doubleAlphabetNum = alphabetNum * alphabetNum

// ByAlphabet returns alphabetic order like (aa, ab, ac, ... , zz, zzaa, zzab, ...)
func ByAlphabet(num int) string {
	q, r := divmod(num, doubleAlphabetNum)
	// r < 26 and can be converted to byte type
	return listZZ(q) + toXX(r)
}

func toXX(num int) string {
	if num >= doubleAlphabetNum {
		panic(`invalid argument`)
	}

	q, r := divmod(num, alphabetNum)
	// q, r < 26 and can be converted to byte type
	return toX(byte(q)) + toX(byte(r))
}

func toX(num byte) string {
	return string(byte('a') + num)
}

func listZZ(num int) string {
	zs := ``
	for i := 0; i < num; i++ {
		zs += `zz`
	}
	return zs
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator
	remainder = numerator % denominator
	return
}
