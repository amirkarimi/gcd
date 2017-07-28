package logger

import (
	"fmt"

	"github.com/fatih/color"
)

type (
	Logger interface {
		Info(mesage ...interface{})
		Error(message ...interface{})
	}

	LoggerAsService struct {
		info func(a ...interface{}) string
		err  func(a ...interface{}) string
	}
)

func (las LoggerAsService) Info(message ...interface{}) {
	text := []interface{}{las.info("INFO")}
	text = append(text, message...)
	fmt.Println(text...)
}

func (las LoggerAsService) Error(message ...interface{}) {
	text := []interface{}{las.err("ERROR")}
	text = append(text, message...)
	fmt.Println(text...)
}

func New() Logger {
	return LoggerAsService{
		color.New(color.FgCyan).Add(color.Bold).SprintFunc(),
		color.New(color.FgRed).Add(color.Bold).SprintFunc(),
	}
}
