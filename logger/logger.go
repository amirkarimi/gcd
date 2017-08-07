package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type (
	Logger interface {
		Info(mesage ...interface{})
		Error(message ...interface{})
		OK(message ...interface{})
		Exit(code int, message ...interface{})
		Skip(message ...interface{})
	}

	LoggerAsService struct {
		info func(a ...interface{}) string
		err  func(a ...interface{}) string
		ok   func(a ...interface{}) string
		exit func(a ...interface{}) string
		skip func(a ...interface{}) string
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

func (las LoggerAsService) OK(message ...interface{}) {
	text := []interface{}{las.ok("OK")}
	text = append(text, message...)
	fmt.Println(text...)
}

func (las LoggerAsService) Exit(code int, message ...interface{}) {
	fmt.Println("")
	text := []interface{}{las.exit("EXIT")}
	text = append(text, message...)
	fmt.Println(text...)
	os.Exit(code)
}

func (las LoggerAsService) Skip(message ...interface{}) {
	text := []interface{}{las.skip("SKIP")}
	text = append(text, message...)
	fmt.Println(text...)
}

func New() Logger {
	return LoggerAsService{
		info: color.New(color.FgCyan).Add(color.Bold).SprintFunc(),
		err:  color.New(color.FgRed).Add(color.Bold).SprintFunc(),
		ok:   color.New(color.FgGreen).Add(color.Bold).SprintFunc(),
		exit: color.New(color.FgWhite).Add(color.Bold).SprintFunc(),
		skip: color.New(color.FgYellow).Add(color.Bold).SprintFunc(),
	}
}
