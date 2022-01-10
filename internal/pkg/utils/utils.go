package utils

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

func ExtractBoolValue(ctx *fasthttp.RequestCtx, valueName string) (bool, error) {
	valueStr := string(ctx.QueryArgs().Peek(valueName))
	if valueStr == "" {
		return false, nil
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, err
	}
	return value, nil
}

func ExtractIntValue(ctx *fasthttp.RequestCtx, valueName string) (int, error) {
	ValueStr := string(ctx.QueryArgs().Peek(valueName))
	if ValueStr == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(ValueStr)
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
