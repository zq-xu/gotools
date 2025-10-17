package utilsx

// MergeStringMaps merges all the string maps received as argument into a single new label map.
func MergeStringMaps(ms ...map[string]string) map[string]string {
	res := map[string]string{}

	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}

	return res
}

func CopyStringMap(m map[string]string) map[string]string {
	ls := make(map[string]string, len(m))

	for k, v := range m {
		ls[k] = v
	}

	return ls
}

func RemoveEmptyStringValueFromMap(m map[string]string) {
	for k, v := range m {
		if v == "" {
			delete(m, k)
		}
	}
}

func GetStringMapItemWithDefault(m map[string]string, key, defaultValue string) string {
	if m == nil {
		return defaultValue
	}

	v, ok := m[key]
	if !ok || v == "" {
		return defaultValue
	}

	return v
}
