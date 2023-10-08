package helper

func PadLeft(input string, length int) string {
	for len(input) < length {
		input = "0" + input
	}
	return input
}
