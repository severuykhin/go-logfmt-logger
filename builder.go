package logfmt

import (
	"strconv"
	"strings"
)

type stringBuilder struct {
	builder strings.Builder
}

func newStringBuilder() stringBuilder {
	return stringBuilder{
		builder: strings.Builder{},
	}
}

func (sb *stringBuilder) Grow(length int) {
	sb.builder.Grow(length)
}

func (sb *stringBuilder) StringFrom(keyVals []interface{}) string {

	keyValsLength := len(keyVals)

	if keyValsLength == 0 {
		return ""
	}

	for i := 0; i < keyValsLength; i++ {
		if i > 0 && i%2 != 0 {
			continue
		}

		key := sb.valueToString(keyVals[i])
		var value string

		if i == keyValsLength-1 {
			value = ""
		} else {
			value = sb.valueToString(keyVals[i+1])
		}

		if value == "" {
			continue
		}

		sb.builder.WriteString(key)
		sb.builder.WriteString("=")
		sb.builder.WriteString(value)
		sb.builder.WriteString("\n")
	}

	return sb.builder.String()
}

func (sb *stringBuilder) valueToString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case error:
		return v.Error()
	default:
		return "unknowntype"
	}
}
