package utils

import (
	"net/http"
	"strconv"
)

func ExtractBoolValue(r *http.Request, valueName string) (bool, error) {
	valueStr := r.URL.Query().Get(valueName)
	if valueStr == "" {
		return false, nil
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, err
	}
	return value, nil
}

func ExtractIntValue(r *http.Request, valueName string) (int, error) {
	valueStr := r.URL.Query().Get(valueName)
	if valueStr == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return -1, err
	}

	return value, nil
}

func Find(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}
