package utils

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			luhn += cur
			number = number / 10
			continue
		}

		cur = cur * 2

		if cur > 9 {
			cur = cur - 9
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func ValidateLuhn(number int) bool {
	return checksum(number) == 0
}
