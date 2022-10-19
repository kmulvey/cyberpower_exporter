package main

import (
	"bytes"
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

/*
Properties:

	Model Name................... CP1500PFCLCDa
	Firmware Number.............. CXXKY2008826
	Rating Voltage............... 120 V
	Rating Power................. 1000 Watt(1500 VA)

Current UPS status:

	State........................ Normal
	Power Supply by.............. Utility Power
	Utility Voltage.............. 122 V
	Output Voltage............... 122 V
	Battery Capacity............. 100 %
	Remaining Runtime............ 55 min.
	Load......................... 140 Watt(14 %)
	Line Interaction............. None
	Test Result.................. Passed at 2022/10/04 13:56:27
	Last Power Event............. Blackout at 2022/09/24 12:12:24 for 3 sec.
*/
var stateRegex = regexp.MustCompile("State\\.+")
var powerSupplyRegex = regexp.MustCompile("Power Supply by\\.+")
var utilityVoltageRegex = regexp.MustCompile("Utility Voltage\\.+\\s\\d+")
var outputVoltageRegex = regexp.MustCompile("Output Voltage\\.+\\s\\d+")
var batteryCapacityRegex = regexp.MustCompile("Battery Capacity\\.+\\s\\d+")

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
		State:         getState(cmdOutput),
		PowerSupplyBy: getState(cmdOutput),
	}

	ds.UtilityVoltage, err = getUtilityVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, err
	}

	ds.OutputVoltage, err = getOutputVoltage(cmdOutput)
	if err != nil {
		return DeviceStatus{}, err
	}

	ds.BatteryCapacity, err = getBatteryCapacity(cmdOutput)
	if err != nil {
		return DeviceStatus{}, err
	}

	return ds, nil
}

func getState(row string) string {
	return strings.TrimSpace(stateRegex.ReplaceAllString(row, ""))
}

func getPowerSupply(row string) string {
	return strings.TrimSpace(powerSupplyRegex.ReplaceAllString(row, ""))
}

func getUtilityVoltage(input string) (int, error) {
	var row = strings.TrimSpace(utilityVoltageRegex.ReplaceAllString(input, ""))
	var re, err = regexp.Compile("\\d+")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getOutputVoltage(input string) (int, error) {
	var row = strings.TrimSpace(outputVoltageRegex.ReplaceAllString(input, ""))
	var re, err = regexp.Compile("\\d+")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}

func getBatteryCapacity(input string) (int, error) {
	var row = strings.TrimSpace(batteryCapacityRegex.ReplaceAllString(input, ""))
	var re, err = regexp.Compile("\\d+")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(re.FindString(row)))
}
