package commons

import (
	"strconv"
	"strings"
)

func FsToString(fs []float32) string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, f := range fs {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmtFloat(f))
	}
	builder.WriteString("]")
	return builder.String()
}

func fmtFloat(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}
