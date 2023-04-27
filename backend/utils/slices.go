package utils

func Includes[Type comparable](arr []Type, item Type) bool {
	for _, el := range arr {
		if el == item {
			return true
		}
	}

	return false
}
