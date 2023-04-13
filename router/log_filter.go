package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	LogInterrupt = "interrupt"
)

var (
	methodFilterSet = map[string]string{
		http.MethodGet:     LogInterrupt,
		http.MethodOptions: LogInterrupt,
	}
)

// GetMethodFilter filter the Get request
func GetMethodFilter(param gin.LogFormatterParams) string {
	return methodFilterSet[param.Method]
}

func LoggerFilter(skipPaths []string, formatters ...gin.LogFormatter) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: skipPaths,
		Formatter: FormatterMiddleWare(formatters...),
		Output:    gin.DefaultWriter,
	})
}

// FormatterMiddleWare is for log format
func FormatterMiddleWare(formatters ...gin.LogFormatter) gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		for _, fn := range formatters {
			if fn(param) == LogInterrupt {
				return ""
			}
		}

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency = param.Latency - param.Latency%time.Second
		}
		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}
}
