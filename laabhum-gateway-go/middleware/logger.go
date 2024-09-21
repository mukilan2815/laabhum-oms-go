package logger

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin" // Add this import
	"os"
)

// Logger wraps the logrus.Logger
type Logger struct {
	*logrus.Logger
}

// New creates a new Logger instance
func New(level logrus.Level) *Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(level)

	return &Logger{log}
}

// Infof logs info level messages
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Errorf logs error level messages
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Logger middleware for Gin
func Middleware(log *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Infof("%s | %3d | %13v | %15s | %s %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
