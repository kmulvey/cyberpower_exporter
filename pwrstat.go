package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Device struct {
	ModelName        string
	FirmwareNumber   string
	RatingVoltage    int
	RatingPowerWatts int
	RatingPowerVA    int
}

type DeviceStatus struct {
	State                  string
	PowerSupplyBy          string // really enum
	UtilityVoltage         int
	OutputVoltage          int
	BatteryCapacity        int // pct
	RemainingRuntime       time.Duration
	LoadWatts              int
	LoadPct                int
	LineInteraction        string
	TestResult             string
	TestResultTime         time.Time
	LastPowerEvent         string
	LastPowerEventTime     time.Time
	LastPowerEventDuration time.Duration
}

var dateFormat = "2006/01/02 15:04:05"
var stateRegex = regexp.MustCompile("State\\.+.*")
var powerSupplyRegex = regexp.MustCompile("Power Supply by\\.+.*")
var utilityVoltageRegex = regexp.MustCompile("Utility Voltage\\.+\\s\\d+")
var outputVoltageRegex = regexp.MustCompile("Output Voltage\\.+\\s\\d+")
var batteryCapacityRegex = regexp.MustCompile("Battery Capacity\\.+\\s\\d+")
var remainingRuntimeRegex = regexp.MustCompile("Remaining Runtime\\.+\\s\\d{1,3}\\smin\\.")
var loadRegex = regexp.MustCompile("Load\\.+.*")
var lineInteractionRegex = regexp.MustCompile("Line Interaction\\.+.*")
var testResultRegex = regexp.MustCompile("Test Result\\.+.*")
var lastPowerEventRegex = regexp.MustCompile("Last Power Event\\.+.*")

func getPowerStats(cmdPath string) (string, error) {

	var cmd = exec.Command(cmdPath, "-status")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	var err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running command, stderr: %s, go err: %w", stderr.String(), err)
	}

	return out.String(), nil
}

func parsePowerStats(cmdOutput string) (DeviceStatus, error) {

	var err error
	var ds = DeviceStatus{
		State:           getState(cmdOutput),
		PowerSupplyBy:   getPowerSupply(cmdOutput),
		LineInteraction: getLineInteraction(cmdOutput),
	}

	ds.UtilityVoltage, err = getUtilityVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getUtilityVoltage err: %w", err)
	}

	ds.OutputVoltage, err = getOutputVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getOutputVoltage err: %w", err)
	}

	ds.BatteryCapacity, err = getBatteryCapacity(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getBatteryCapacity err: %w", err)
	}

	ds.RemainingRuntime, err = getRemainingRuntime(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getRemainingRuntime err: %w", err)
	}

	ds.LoadWatts, ds.LoadPct, err = getLoad(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getLoad err: %w", err)
	}

	ds.TestResult, ds.TestResultTime, err = getTestResult(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getLoad err: %w", err)
	}

	ds.LastPowerEvent, ds.LastPowerEventTime, ds.LastPowerEventDuration, err = getLastPowerEvent(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getLoad err: %w", err)
	}

	return ds, nil
}

func getState(input string) string {
	var row = strings.TrimSpace(stateRegex.FindString(input))
	var re = regexp.MustCompile("State\\.+")
	return strings.TrimSpace(re.ReplaceAllString(row, ""))
}

func getPowerSupply(input string) string {
	var row = strings.TrimSpace(powerSupplyRegex.FindString(input))
	var re = regexp.MustCompile("Power Supply by\\.+")
	return strings.TrimSpace(re.ReplaceAllString(row, ""))
}

func getUtilityVoltage(input string) (int, error) {
	var row = strings.TrimSpace(utilityVoltageRegex.FindString(input))
	var re, err = regexp.Compile("\\d+")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getOutputVoltage(input string) (int, error) {
	var row = strings.TrimSpace(outputVoltageRegex.FindString(input))
	var re, err = regexp.Compile("\\d+")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getBatteryCapacity(input string) (int, error) {
	var row = strings.TrimSpace(batteryCapacityRegex.FindString(input))
	var re, err = regexp.Compile("\\d+")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getRemainingRuntime(input string) (time.Duration, error) {
	var row = strings.TrimSpace(remainingRuntimeRegex.FindString(input))
	var re, err = regexp.Compile("\\d{1,3}")
	if err != nil {
		return 0, err
	}
	mins, err := strconv.Atoi(strings.TrimSpace(re.FindString(row)))
	if err != nil {
		return 0, err
	}
	return time.Duration(mins) * time.Minute, nil
}

func getLoad(input string) (int, int, error) {
	var row = strings.TrimSpace(loadRegex.FindString(input))
	var watt = regexp.MustCompile("\\d{1,4}\\sWatt")
	var pctRe = regexp.MustCompile("\\(\\d{1,3}\\s\\%\\)")

	var wattStr = watt.FindString(row)
	watts, err := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(wattStr, " Watt", "")))
	if err != nil {
		return 0, 0, err
	}

	var pctStr = pctRe.FindString(row)
	var pctNum = regexp.MustCompile("\\d{1,3}")
	pct, err := strconv.Atoi(strings.TrimSpace(pctNum.FindString(pctStr)))
	if err != nil {
		return 0, 0, err
	}

	return watts, pct, nil
}

func getLineInteraction(input string) string {
	var row = strings.TrimSpace(lineInteractionRegex.FindString(input))
	var re = regexp.MustCompile("Line Interaction\\.+")
	return strings.TrimSpace(re.ReplaceAllString(row, ""))
}

func getTestResult(input string) (string, time.Time, error) {
	var testResultRegex = regexp.MustCompile("Test Result\\.+.*")
	var row = strings.TrimSpace(testResultRegex.FindString(input))
	var leftSide = regexp.MustCompile("^Test Result\\.*\\s")
	row = leftSide.ReplaceAllString(row, "")
	var result = regexp.MustCompile("^.*\\sat")

	var dateStr = strings.TrimSpace(result.ReplaceAllString(row, ""))
	var date, err = time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", time.Time{}, err
	}

	return strings.ReplaceAll(result.FindString(row), " at", ""), date, nil
}

func getLastPowerEvent(input string) (string, time.Time, time.Duration, error) {
	var testResultRegex = regexp.MustCompile("Last Power Event\\.+.*")
	var row = strings.TrimSpace(testResultRegex.FindString(input))
	var leftSide = regexp.MustCompile("^Last Power Event\\.*\\s")
	row = leftSide.ReplaceAllString(row, "")

	var result, right, found = strings.Cut(row, " at ")
	if !found {
		return "", time.Time{}, 0, errors.New("did not find substring: ' at '")
	}

	dateStr, durationStr, found := strings.Cut(right, " for ")
	if !found {
		return "", time.Time{}, 0, errors.New("did not find substring: ' for '")
	}

	var date, err = time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", time.Time{}, 0, err
	}

	durationNumStr, durationUnit, found := strings.Cut(durationStr, " ")
	if !found {
		return "", time.Time{}, 0, errors.New("did not find substring: ' '")
	}
	durationNum, err := strconv.Atoi(durationNumStr)
	if err != nil {
		return "", time.Time{}, 0, err
	}

	var duration time.Duration
	switch durationUnit {
	case "sec.":
		duration = time.Duration(durationNum) * time.Second
	case "min.":
		duration = time.Duration(durationNum) * time.Minute
	default:
		return "", time.Time{}, 0, fmt.Errorf("event duration has an unknown unit, input: %s", input)
	}

	return result, date, duration, nil
}
