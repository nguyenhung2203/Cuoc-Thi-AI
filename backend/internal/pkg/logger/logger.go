package logger

import "fmt"

func Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] " + msg + "\n", args...)
}
