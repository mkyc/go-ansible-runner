package main

import (
	"fmt"

	gar "github.com/mkyc/go-ansible-runner"
)

type SimpleLogger struct{}

func (s SimpleLogger) Trace(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Debug(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Info(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Warn(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Error(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Fatal(format string, v ...interface{}) {
	justPrint(format, v...)
}

func (s SimpleLogger) Panic(format string, v ...interface{}) {
	justPrint(format, v...)
}

func justPrint(s string, v ...interface{}) {
	if len(v) > 0 {
		fmt.Printf(s, v...)
	} else {
		fmt.Println(s)
	}
}

func main() {
	opts := gar.Options{
		AnsibleRunnerDir: "./tests",
		Playbook:         "test1.yml",
		Ident:            "r1",
		LogsLevel:        gar.L6,
		Logger:           SimpleLogger{},
	}

	println("================")
	println("=== run ========")
	println("================")

	s, err := gar.Run(opts)
	if err != nil {
		panic(err)
	}

	println("================")
	println("=== run = end ==")
	println("================")

	println(s)

	println("================")
	println("=== end ========")
	println("================")
}
