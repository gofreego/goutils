package utils

func NotIn[T comparable](value T, values ...T) bool {
	for _, v := range values {
		if v == value {
			return false
		}
	}
	return true
}
