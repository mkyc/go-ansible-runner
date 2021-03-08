package gar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Run(originalOptions Options) (string, error) {
	options, args := getCommonOptions(originalOptions, "run", ".")
	_, err := runAnsibleRunnerCommand(options, args...)
	if err != nil {
		return options.Ident, err
	}
	return options.Ident, nil
}

func GetOutput(originalOptions Options) ([]byte, error) {
	options, _ := getCommonOptions(originalOptions)
	return ioutil.ReadFile(filepath.Join(options.AnsibleRunnerDir, "artifacts", options.Ident, "stdout"))
}

func GetStatus(originalOptions Options) (int, string, error) {
	options, _ := getCommonOptions(originalOptions)
	rcBytes, err := ioutil.ReadFile(filepath.Join(options.AnsibleRunnerDir, "artifacts", options.Ident, "rc"))
	if err != nil {
		return -1, "", err
	}
	statusBytes, err := ioutil.ReadFile(filepath.Join(options.AnsibleRunnerDir, "artifacts", options.Ident, "status"))
	if err != nil {
		return -2, "", err
	}
	rc, err := strconv.Atoi(string(rcBytes))
	if err != nil {
		return -3, "", err
	}
	return rc, string(statusBytes), nil
}

type PlayRecap struct {
	Changed   map[string]int `json:"changed"`
	Failures  map[string]int `json:"failures"`
	Ignored   map[string]int `json:"ignored"`
	Ok        map[string]int `json:"ok"`
	Processed map[string]int `json:"processed"`
	Rescued   map[string]int `json:"rescued"`
	Skipped   map[string]int `json:"skipped"`
}

func GetPlayRecap(originalOptions Options) (*PlayRecap, error) {
	options, _ := getCommonOptions(originalOptions)
	var lastEventNumber int
	var lastEventFilePath string
	err := filepath.Walk(filepath.Join(options.AnsibleRunnerDir, "artifacts", options.Ident, "job_events"),
		func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".json") {
				eventNumber, err := strconv.Atoi(strings.Split(filepath.Base(path), "-")[0])
				if err != nil {
					return err
				}
				if eventNumber > lastEventNumber {
					lastEventNumber = eventNumber
					lastEventFilePath = path
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	eventFileBytes, err := ioutil.ReadFile(lastEventFilePath)
	if err != nil {
		return nil, err
	}

	type event struct {
		Event     string    `json:"event"`
		EventData PlayRecap `json:"event_data"`
	}

	var e event
	if err := json.Unmarshal(eventFileBytes, &e); err != nil {
		return nil, err
	}
	if e.Event != "playbook_on_stats" {
		return nil, fmt.Errorf("incorrect event type found: %s", e.Event)
	}
	return &e.EventData, nil
}

type TasksCount struct {
	Changed   int
	Failures  int
	Ignored   int
	Ok        int
	Processed int
	Rescued   int
	Skipped   int
}

func Count(recap PlayRecap) TasksCount {
	cnt := TasksCount{}
	for _, v := range recap.Changed {
		cnt.Changed = cnt.Changed + v
	}
	for _, v := range recap.Failures {
		cnt.Failures = cnt.Failures + v
	}
	for _, v := range recap.Ignored {
		cnt.Ignored = cnt.Ignored + v
	}
	for _, v := range recap.Ok {
		cnt.Ok = cnt.Ok + v
	}
	for _, v := range recap.Processed {
		cnt.Processed = cnt.Processed + v
	}
	for _, v := range recap.Rescued {
		cnt.Rescued = cnt.Rescued + v
	}
	for _, v := range recap.Skipped {
		cnt.Skipped = cnt.Skipped + v
	}
	return cnt
}
