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
	BatteryCapacity        int           // pct
	RemainingRuntime       time.Duration `json:"RemainingRuntimeNano"`
	LoadWatts              int
	LoadPct                int
	LineInteraction        string
	TestResult             string
	TestResultTime         time.Time
	LastPowerEvent         string
	LastPowerEventTime     time.Time
	LastPowerEventDuration time.Duration `json:"LastPowerEventDurationNano"`
	CollectionTime         time.Time     `gorm:"primaryKey"`
}

// DeviceStatus regexs
var stateRegex = regexp.MustCompile(`State\.+.*`)
var powerSupplyRegex = regexp.MustCompile(`Power Supply by\.+.*`)
var utilityVoltageRegex = regexp.MustCompile(`Utility Voltage\.+\s\d+`)
var outputVoltageRegex = regexp.MustCompile(`Output Voltage\.+\s\d+`)
var batteryCapacityRegex = regexp.MustCompile(`Battery Capacity\.+\s\d+`)
var remainingRuntimeRegex = regexp.MustCompile(`Remaining Runtime\.+\s\d{1,3}\smin\.`)
var loadRegex = regexp.MustCompile(`Load\.+.*`)
var lineInteractionRegex = regexp.MustCompile(`Line Interaction\.+.*`)
var testResultRegex = regexp.MustCompile(`Test Result\.+.*`)
var lastPowerEventRegex = regexp.MustCompile(`Last Power Event\.+.*`)

// Device regexs
var ModelNameRegex = regexp.MustCompile(`Model Name\.+\s([a-zA-Z0-9]+)`)
var FirmwareNumberRegex = regexp.MustCompile(`Firmware Number\.+\s([a-zA-Z0-9]+)`)
var RatingVoltageRegex = regexp.MustCompile(`Rating Voltage\.+\s([0-9]+)\sV`)
var RatingPowerWattsRegex = regexp.MustCompile(`Rating Power\.+\s([0-9]+)\sWatt\(([0-9]+)\sVA\)`)

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

func parsePowerStats(cmdOutput string) (DeviceStatus, Device, error) {

	var err error
	var ds = DeviceStatus{
		State:           getState(cmdOutput),
		PowerSupplyBy:   getPowerSupply(cmdOutput),
		LineInteraction: getLineInteraction(cmdOutput),
	}

	ds.UtilityVoltage, err = getUtilityVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getUtilityVoltage err: %w", err)
	}

	ds.OutputVoltage, err = getOutputVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getOutputVoltage err: %w", err)
	}

	ds.BatteryCapacity, err = getBatteryCapacity(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getBatteryCapacity err: %w", err)
	}

	ds.RemainingRuntime, err = getRemainingRuntime(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getRemainingRuntime err: %w", err)
	}

	ds.LoadWatts, ds.LoadPct, err = getLoad(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getLoad err: %w", err)
	}

	ds.TestResult, ds.TestResultTime, err = getTestResult(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getLoad err: %w", err)
	}

	ds.LastPowerEvent, ds.LastPowerEventTime, ds.LastPowerEventDuration, err = getLastPowerEvent(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getLoad err: %w", err)
	}

	ds.CollectionTime = time.Now()

	//////////////// Device
	var device = Device{}

	device.ModelName, err = getModelName(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getModelName err: %w", err)
	}

	device.FirmwareNumber, err = getFirmwareNumber(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getFirmwareNumber err: %w", err)
	}

	device.RatingVoltage, err = getRatingVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getRatingVoltage err: %w", err)
	}

	device.RatingPowerWatts, device.RatingPowerVA, err = getRatingPowerWatts(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getRatingPowerWatts err: %w", err)
	}

	return ds, device, nil
}

//////////////// Device Status

func getState(input string) string {
	var row = strings.TrimSpace(stateRegex.FindString(input))
	var re = regexp.MustCompile(`State\.+`)
	return strings.TrimSpace(re.ReplaceAllString(row, ""))
}

func getPowerSupply(input string) string {
	var row = strings.TrimSpace(powerSupplyRegex.FindString(input))
	var re = regexp.MustCompile(`Power Supply by\.+`)
	return strings.TrimSpace(re.ReplaceAllString(row, ""))
}

func getUtilityVoltage(input string) (int, error) {
	var row = strings.TrimSpace(utilityVoltageRegex.FindString(input))
	var re, err = regexp.Compile(`\d+`)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getOutputVoltage(input string) (int, error) {
	var row = strings.TrimSpace(outputVoltageRegex.FindString(input))
	var re, err = regexp.Compile(`\d+`)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getBatteryCapacity(input string) (int, error) {
	var row = strings.TrimSpace(batteryCapacityRegex.FindString(input))
	var re, err = regexp.Compile(`\d+`)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getRemainingRuntime(input string) (time.Duration, error) {
	var row = strings.TrimSpace(remainingRuntimeRegex.FindString(input))
	var re, err = regexp.Compile(`\d{1,3}`)
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
	var watt = regexp.MustCompile(`\d{1,4}\sWatt`)
	var pctRe = regexp.MustCompile(`\(\d{1,3}\s\%\)`)

	var wattStr = watt.FindString(row)
	watts, err := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(wattStr, " Watt", "")))
	if err != nil {
		return 0, 0, err
	}

	var pctStr = pctRe.FindString(row)
	var pctNum = regexp.MustCompile(`\d{1,3}`)
	pct, err := strconv.Atoi(strings.TrimSpace(pctNum.FindString(pctStr)))
	if err != nil {
		return 0, 0, err
	}

	return watts, pct, nil
}

func getLineInteraction(input string) string {
	var row = strings.TrimSpace(lineInteractionRegex.FindString(input))
	var re = regexp.MustCompile(`Line Interaction\.+`)
	return strings.TrimSpace(re.ReplaceAllString(row, ""))
}

func getTestResult(input string) (string, time.Time, error) {
	var row = strings.TrimSpace(testResultRegex.FindString(input))
	var leftSide = regexp.MustCompile(`^Test Result\.*\s`)
	row = leftSide.ReplaceAllString(row, "")
	var result = regexp.MustCompile(`^.*\sat`)

	var dateStr = strings.TrimSpace(result.ReplaceAllString(row, ""))
	var date, err = time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", time.Time{}, err
	}

	return strings.ReplaceAll(result.FindString(row), " at", ""), date, nil
}

func getLastPowerEvent(input string) (string, time.Time, time.Duration, error) {
	var row = strings.TrimSpace(lastPowerEventRegex.FindString(input))
	var leftSide = regexp.MustCompile(`^Last Power Event\.*\s`)
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

//////////////// Device

func getDeviceInfoAsString(re *regexp.Regexp, input string) (string, error) {
	var match = re.FindAllStringSubmatch(input, 1)
	if len(match) == 1 {
		if len(match[0]) == 2 {
			return strings.TrimSpace(match[0][1]), nil
		}
	}
	return "", errors.New("") // the real error is generated by the caller, this is just a signal
}

func getModelName(input string) (string, error) {
	var val, err = getDeviceInfoAsString(ModelNameRegex, input)
	if err != nil {
		return "", errors.New("unable to find the model name")
	}
	return val, nil
}

func getFirmwareNumber(input string) (string, error) {
	var val, err = getDeviceInfoAsString(FirmwareNumberRegex, input)
	if err != nil {
		return "", errors.New("unable to find the firmware number")
	}
	return val, nil
}

func getRatingVoltage(input string) (int, error) {
	var match = RatingVoltageRegex.FindAllStringSubmatch(input, 1)
	if len(match) == 1 {
		if len(match[0]) == 2 {
			var volts, err = strconv.Atoi(strings.TrimSpace(match[0][1]))
			if err != nil {
				return 0, fmt.Errorf("unable to find the rating power in watts, err: %w", err)
			}

			return volts, nil
		}
	}
	return 0, errors.New("unable to find the rating power in watts")
}

func getRatingPowerWatts(input string) (int, int, error) {
	var match = RatingPowerWattsRegex.FindAllStringSubmatch(input, 2)
	if len(match) == 1 {
		if len(match[0]) == 3 {
			var watts, err = strconv.Atoi(strings.TrimSpace(match[0][1]))
			if err != nil {
				return 0, 0, fmt.Errorf("unable to find the rating power in watts, err: %w", err)
			}

			va, err := strconv.Atoi(strings.TrimSpace(match[0][2]))
			if err != nil {
				return 0, 0, fmt.Errorf("unable to find the rating power in va, err: %w", err)
			}

			return watts, va, nil
		}
	}
	return 0, 0, errors.New("unable to find the rating power")
}
