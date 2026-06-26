package logger

import "fmt"

func Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] " + msg + "\n", args...)
}

func Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] " + msg + "\n", args...)
}

func Warn(msg string, args ...interface{}) {
	fmt.Printf("[WARN] " + msg + "\n", args...)
}
