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

var (
	errUnknownDurationUnit  = errors.New("event duration has an unknown unit")
	errGroupIDExceedsLength = errors.New("groupID exceeds array length")
	errNoMatchesFound       = errors.New("could not find any matches")
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
	BatteryCapacity        int // pct out of 100
	RemainingRuntime       time.Duration
	LoadWatts              int
	LoadPct                int
	LineInteraction        string
	TestResult             string
	TestResultTime         time.Time
	LastPowerEvent         string
	LastPowerEventTime     time.Time
	LastPowerEventDuration time.Duration
	CollectionTime         time.Time
}

// $ does not work to find EOL, idk why, so we use \n.
var stateRegex = regexp.MustCompile(`State\.+\s([a-zA-Z\s]+)\n`)
var powerSupplyRegex = regexp.MustCompile(`Power Supply by\.+\s([a-zA-Z\s]+)\n`)
var utilityVoltageRegex = regexp.MustCompile(`Utility Voltage\.+\s(\d+)\sV\n`)
var outputVoltageRegex = regexp.MustCompile(`Output Voltage\.+\s(\d+)\sV\n`)
var batteryCapacityRegex = regexp.MustCompile(`Battery Capacity\.+\s(\d+)\s\%\n`)
var remainingRuntimeRegex = regexp.MustCompile(`Remaining Runtime\.+\s(\d{1,3})\smin\.\n`)
var loadRegex = regexp.MustCompile(`Load\.+\s(\d+)\sWatt\((\d+)\s\%\)\n`)
var lineInteractionRegex = regexp.MustCompile(`Line Interaction\.+\s([a-zA-Z]+)\n`)
var testResultRegex = regexp.MustCompile(`Test Result\.+\s([a-zA-Z]+)\sat\s(.*)\n`)
var lastPowerEventRegex = regexp.MustCompile(`Last Power Event\.+\s([a-zA-Z]+)\sat\s(\d+\/\d+\/\d+\s+\d+\:\d+\d\:\d+)(?:\sfor\s(\d+)\s([a-zA-Z]+)\.|\n)`)

// regexs for Device.
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

func parsePowerStatus(cmdOutput string) (DeviceStatus, error) {
	var err error
	var status = DeviceStatus{}

	status.State, err = getState(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getState err: %w", err)
	}

	status.TestResult, status.TestResultTime, err = getTestResult(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getTestResult err: %w", err)
	}

	status.LastPowerEvent, status.LastPowerEventTime, status.LastPowerEventDuration, err = getLastPowerEvent(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getLastPowerEvent err: %w", err)
	}

	// if we lost communication, we dont need to parse the rest.
	// see test data for an example.
	if strings.Contains(cmdOutput, "Lost Communication") {
		return status, nil
	}

	status.PowerSupplyBy, err = getPowerSupply(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getPowerSupply err: %w", err)
	}

	status.LineInteraction, err = getLineInteraction(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getLineInteraction err: %w", err)
	}

	status.UtilityVoltage, err = getUtilityVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getUtilityVoltage err: %w", err)
	}

	status.OutputVoltage, err = getOutputVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getOutputVoltage err: %w", err)
	}

	status.BatteryCapacity, err = getBatteryCapacity(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getBatteryCapacity err: %w", err)
	}

	status.RemainingRuntime, err = getRemainingRuntime(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getRemainingRuntime err: %w", err)
	}

	status.LoadWatts, status.LoadPct, err = getLoad(cmdOutput)
	if err != nil {
		return DeviceStatus{}, fmt.Errorf("getLoad err: %w", err)
	}

	status.CollectionTime = time.Now()
	return status, nil
}

// parseDeviceProperties parses USP properties.
func parseDeviceProperties(cmdOutput string) (Device, error) {
	var device = Device{}
	var err error

	device.ModelName, err = getModelName(cmdOutput)
	if err != nil {
		return Device{}, fmt.Errorf("getModelName err: %w", err)
	}

	device.FirmwareNumber, err = getFirmwareNumber(cmdOutput)
	if err != nil {
		return Device{}, fmt.Errorf("getFirmwareNumber err: %w", err)
	}

	device.RatingVoltage, err = getRatingVoltage(cmdOutput)
	if err != nil {
		return Device{}, fmt.Errorf("getRatingVoltage err: %w", err)
	}

	device.RatingPowerWatts, device.RatingPowerVA, err = getRatingPowerWatts(cmdOutput)
	if err != nil {
		return Device{}, fmt.Errorf("getRatingPowerWatts err: %w", err)
	}

	return device, nil
}

func getState(input string) (string, error) {
	var val, err = getDeviceInfoAsString(stateRegex, input, 1)
	if err != nil {
		return "", fmt.Errorf("unable to get the state, err: %w", err)
	}
	return val, nil
}

func getPowerSupply(input string) (string, error) {
	var val, err = getDeviceInfoAsString(powerSupplyRegex, input, 1)
	if err != nil {
		return "", fmt.Errorf("unable to get the power supply, err: %w", err)
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
	var capacity, err = getDeviceInfoAsInt(batteryCapacityRegex, input, 1)
	if err != nil {
		return 0, fmt.Errorf("unable to find the battery capacity, err: %w", err)
	}
	return capacity, nil
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
		return "", fmt.Errorf("unable to find the line interaction, err: %w", err)
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
		return "", time.Time{}, fmt.Errorf("unable to find the last test result date, err: %w", err)
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("unable to pasre date: %s, err: %w", dateStr, err)
	}

	return result, date, nil
}

func getLastPowerEvent(input string) (string, time.Time, time.Duration, error) {

	// lack of a power event is not in the main regex as it is already too complicated
	if regexp.MustCompile(`Last Power Event............. None`).MatchString(input) {
		return "None", time.Time{}, 0, nil
	}

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
		return "", time.Time{}, 0, fmt.Errorf("unable to parse date: %s, err: %w", dateStr, err)
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
			return "", time.Time{}, 0, fmt.Errorf("%w: %s", errUnknownDurationUnit, durationUnit)
		}
	}

	return result, date, duration, nil
}

//////////////// Device

func getModelName(input string) (string, error) {
	var val, err = getDeviceInfoAsString(modelNameRegex, input, 1)
	if err != nil {
		return "", fmt.Errorf("unable to find the model name, err: %w", err)
	}
	return val, nil
}

func getFirmwareNumber(input string) (string, error) {
	var val, err = getDeviceInfoAsString(firmwareNumberRegex, input, 1)
	if err != nil {
		return "", fmt.Errorf("unable to find the firmware number, err: %w", err)
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
				return "", fmt.Errorf("%w: groupID %d exceeds array length", errGroupIDExceedsLength, groupID)
			}
			return strings.TrimSpace(match[0][groupID]), nil
		}
	}

	return "", errNoMatchesFound
}

func getDeviceInfoAsInt(re *regexp.Regexp, input string, groupID int) (int, error) {
	var match = re.FindAllStringSubmatch(input, -1)
	if len(match) == 1 && len(match[0]) >= 2 {
		if len(match[0]) <= groupID {
			return 0, fmt.Errorf("%w: groupID %d exceeds array length", errGroupIDExceedsLength, groupID)
		}

		var val, err = strconv.Atoi(strings.TrimSpace(match[0][groupID]))
		if err != nil {
			return 0, fmt.Errorf("unable to convert string: %s to int, err: %w", strings.TrimSpace(match[0][groupID]), err)
		}

		return val, nil
	}
	return 0, errNoMatchesFound
}
