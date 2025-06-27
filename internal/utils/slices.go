package utils

func IsSliceWithNilValues(value []any) bool {
	for i := 0; i < len(value)-1; i++ {
		if value[i] == nil {
			return true
		}
	}
	return false
}
