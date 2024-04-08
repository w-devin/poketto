package basic

func Btoi(b bool) int {
	if b {
		return 1
	}

	return 0
}

func ItoB(i int) bool {
	if i == 0 {
		return false
	}
	return true
}
