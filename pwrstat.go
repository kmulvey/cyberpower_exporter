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
	BatteryCapacity        int           // pct out of 100
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
var stateRegex = regexp.MustCompile(`State\.+\s([a-zA-Z]+)`)
var powerSupplyRegex = regexp.MustCompile(`Power Supply by\.+\s([a-zA-Z ]+)`)
var utilityVoltageRegex = regexp.MustCompile(`Utility Voltage\.+\s(\d+)\sV`)
var outputVoltageRegex = regexp.MustCompile(`Output Voltage\.+\s(\d+)\sV`)
var batteryCapacityRegex = regexp.MustCompile(`Battery Capacity\.+\s(\d+)\s\%`)
var remainingRuntimeRegex = regexp.MustCompile(`Remaining Runtime\.+\s(\d{1,3})\smin\.`)
var loadRegex = regexp.MustCompile(`Load\.+\s(\d+)\sWatt\((\d+)\s\%\)`)
var lineInteractionRegex = regexp.MustCompile(`Line Interaction\.+\s([a-zA-Z]+)`)
var testResultRegex = regexp.MustCompile(`Test Result\.+\s([a-zA-Z]+)\sat\s(.*)`)
var lastPowerEventRegex = regexp.MustCompile(`Last Power Event\.+\s([a-zA-Z]+)\sat\s(.*)(\sfor\s(\d+)\s([a-zA-Z]+)\.)?`) //Last Power Event\.+\s([a-zA-Z]+)\sat\s(.*)\sfor\s(\d+)\s([a-zA-Z]+)\.`)

// Device regexs
var modelNameRegex = regexp.MustCompile(`Model Name\.+\s([a-zA-Z0-9]+)`)
var firmwareNumberRegex = regexp.MustCompile(`Firmware Number\.+\s([a-zA-Z0-9]+)`)
var ratingVoltageRegex = regexp.MustCompile(`Rating Voltage\.+\s(\d+)\sV`)
var ratingPowerWattsRegex = regexp.MustCompile(`Rating Power\.+\s(\d+)\sWatt\((\d+)\sVA\)`)

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
	var ds = DeviceStatus{}

	ds.State, err = getState(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getState err: %w", err)
	}

	ds.PowerSupplyBy, err = getPowerSupply(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getPowerSupply err: %w", err)
	}

	ds.LineInteraction, err = getLineInteraction(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getLineInteraction err: %w", err)
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
		return DeviceStatus{}, Device{}, fmt.Errorf("getTestResult err: %w", err)
	}

	ds.LastPowerEvent, ds.LastPowerEventTime, ds.LastPowerEventDuration, err = getLastPowerEvent(cmdOutput)
	if err != nil {
		return DeviceStatus{}, Device{}, fmt.Errorf("getLastPowerEvent err: %w", err)
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

func getState(input string) (string, error) {
	var val, err = getDeviceInfoAsString(stateRegex, input, 1)
	if err != nil {
		return "", fmt.Errorf("unable to get the state, err %w", err)
	}
	return val, nil
}

func getPowerSupply(input string) (string, error) {
	var val, err = getDeviceInfoAsString(powerSupplyRegex, input, 1)
	if err != nil {
		return "", fmt.Errorf("unable to get the power supply, err %w", err)
	}
	return val, nil
}

func getUtilityVoltage(input string) (int, error) {
	var output, err = getDeviceInfoAsInt(utilityVoltageRegex, input, 1)
	if err != nil {
		return 0, fmt.Errorf("unable to find the utility voltage, err: %w", err)
	}
	return output, nil
}

func getOutputVoltage(input string) (int, error) {
	var output, err = getDeviceInfoAsInt(outputVoltageRegex, input, 1)
	if err != nil {
		return 0, fmt.Errorf("unable to find the output voltage, err: %w", err)
	}
	return output, nil
}

func getBatteryCapacity(input string) (int, error) {
	var cap, err = getDeviceInfoAsInt(batteryCapacityRegex, input, 1)
	if err != nil {
		return 0, fmt.Errorf("unable to find the battery capacity, err: %w", err)
	}
	return cap, nil
}

func getRemainingRuntime(input string) (time.Duration, error) {
	var mins, err = getDeviceInfoAsInt(remainingRuntimeRegex, input, 1)
	if err != nil {
		return 0, fmt.Errorf("unable to find the remaining runtime, err: %w", err)
	}
	return time.Duration(mins) * time.Minute, nil
}

func getLoad(input string) (int, int, error) {
	var watts, err = getDeviceInfoAsInt(loadRegex, input, 1)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to find the load, err: %w", err)
	}

	pct, err := getDeviceInfoAsInt(loadRegex, input, 2)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to find the load, err: %w", err)
	}

	return watts, pct, nil
}

func getLineInteraction(input string) (string, error) {
	var val, err = getDeviceInfoAsString(lineInteractionRegex, input, 1)
	if err != nil {
		return "", fmt.Errorf("unable to get the line interaction, err %w", err)
	}
	return val, nil
}

func getTestResult(input string) (string, time.Time, error) {

	var result, err = getDeviceInfoAsString(testResultRegex, input, 1)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("unable to find the last test result, err: %w", err)
	}

	dateStr, err := getDeviceInfoAsString(testResultRegex, input, 2)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("unable to find the last test result, err: %w", err)
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", time.Time{}, err
	}

	return result, date, nil
}

func getLastPowerEvent(input string) (string, time.Time, time.Duration, error) {

	var result, err = getDeviceInfoAsString(lastPowerEventRegex, input, 1)
	if err != nil {
		return "", time.Time{}, 0, fmt.Errorf("unable to find the last power event, err: %w", err)
	}

	dateStr, err := getDeviceInfoAsString(lastPowerEventRegex, input, 2)
	if err != nil {
		return "", time.Time{}, 0, fmt.Errorf("unable to find the last power event, err: %w", err)
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", time.Time{}, 0, err
	}

	var duration time.Duration
	if regexp.MustCompile(`\sfor\s(\d+)\s([a-zA-Z]+)\.`).MatchString(input) { // we dont always get the duration part " for 3 sec."

		durationInt, err := getDeviceInfoAsInt(lastPowerEventRegex, input, 3)
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("unable to find the last power event, err: %w", err)
		}

		durationUnit, err := getDeviceInfoAsString(lastPowerEventRegex, input, 4)
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("unable to find the last power event, err: %w", err)
		}

		switch strings.TrimSpace(durationUnit) {
		case "sec":
			duration = time.Duration(durationInt) * time.Second
		case "min":
			duration = time.Duration(durationInt) * time.Minute
		default:
			return "", time.Time{}, 0, fmt.Errorf("event duration has an unknown unit, input: %s", durationUnit)
		}
	}

	return result, date, duration, nil
}

//////////////// Device

func getModelName(input string) (string, error) {
	var val, err = getDeviceInfoAsString(modelNameRegex, input, 1)
	if err != nil {
		return "", errors.New("unable to find the model name")
	}
	return val, nil
}

func getFirmwareNumber(input string) (string, error) {
	var val, err = getDeviceInfoAsString(firmwareNumberRegex, input, 1)
	if err != nil {
		return "", errors.New("unable to find the firmware number")
	}
	return val, nil
}

func getRatingVoltage(input string) (int, error) {

	var volts, err = getDeviceInfoAsInt(ratingVoltageRegex, input, 1)
	if err != nil {
		return 0, fmt.Errorf("unable to find the rating voltage, err: %w", err)
	}
	return volts, nil
}

func getRatingPowerWatts(input string) (int, int, error) {

	var watts, err = getDeviceInfoAsInt(ratingPowerWattsRegex, input, 1)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to find the rating power in watts, err: %w", err)
	}

	va, err := getDeviceInfoAsInt(ratingPowerWattsRegex, input, 2)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to find the rating power in va, err: %w", err)
	}

	return watts, va, nil
}

//////////////// Common

func getDeviceInfoAsString(re *regexp.Regexp, input string, groupID int) (string, error) {
	var match = re.FindAllStringSubmatch(input, -1)

	if len(match) == 1 {
		if len(match[0]) >= 2 {
			if len(match[0]) <= groupID {
				return "", errors.New("groupID exceeds arr length")
			}
			return strings.TrimSpace(match[0][groupID]), nil
		}
	}

	return "", errors.New("could not find any matches")
}

func getDeviceInfoAsInt(re *regexp.Regexp, input string, groupID int) (int, error) {
	var match = re.FindAllStringSubmatch(input, -1)
	if len(match) == 1 {
		if len(match[0]) >= 2 {
			if len(match[0]) <= groupID {
				return 0, errors.New("groupID exceeds arr length")
			}

			var val, err = strconv.Atoi(strings.TrimSpace(match[0][groupID]))
			if err != nil {
				return 0, err
			}
			return val, nil
		}
	}
	return 0, errors.New("could not find any matches")
}
