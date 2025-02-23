package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

var isProd = os.Getenv("APP_ENV") == "production"

func AppLog(v ...interface{}) {
	if isProd {
		return
	}
	// AppLog logs a message.
	// It uses the app.Log function to log the message.
	// It is a wrapper around app.Log.
	// It accepts a variadic number of arguments.
	// It returns nothing.
	token := fmt.Sprintf("[%s] ", time.Now().Format(time.RFC3339))
	args := []any{token}
	args = append(args, v...)
	app.Log(args...)
}

func AppLogf(format string, v ...interface{}) {
	if isProd {
		return
	}
	// AppLogf logs a formatted message.
	// It uses the app.Log function to log the message.
	// It is a wrapper around app.Log.
	// It accepts a format string and a variadic number of arguments.
	// It returns nothing.
	token := fmt.Sprintf("[%s] ", time.Now())
	app.Logf(token+format, v...)
}
