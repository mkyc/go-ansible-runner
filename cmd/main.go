package main

import (
	"encoding/json"
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
		LogsLevel:        gar.L0,
		Logger:           SimpleLogger{},
	}

	println("================")
	println("=== run ========")
	println("================")

	ident, err := gar.Run(opts)
	if err != nil {
		panic(err)
	}

	println("================")
	println("=== run = end ==")
	println("================")
	println("==== ident =====")
	println("================")

	println(ident)

	println("================")
	println("==== rc ========")
	println("================")

	rc, status, err := gar.GetStatus(opts)
	if err != nil {
		panic(err)
	}

	println(rc)

	println("================")
	println("==== status ====")
	println("================")

	println(status)

	println("================")
	println("==== recap  ====")
	println("================")

	pr, err := gar.GetPlayRecap(opts)
	if err != nil {
		panic(err)
	}

	s, err := json.MarshalIndent(pr, "", "\t")
	if err != nil {
		panic(err)
	}
	println(string(s))

	println("================")
	println("==== count =====")
	println("================")

	c := gar.Count(*pr)
	cj, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		panic(err)
	}
	println(string(cj))

	println("================")
	println("==== output ====")
	println("================")

	output, err := gar.GetOutput(opts)
	if err != nil {
		panic(err)
	}

	println(string(output))

	println("================")
	println("=== end ========")
	println("================")
}
