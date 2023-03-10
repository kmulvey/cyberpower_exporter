package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testOutputNormal = `
The UPS information shows as following:

	Properties:
		Model Name................... CP1500PFCLCDa
		Firmware Number.............. CR01802B7H21
		Rating Voltage............... 120 V
		Rating Power................. 1000 Watt(1500 VA)

	Current UPS status:
		State........................ Normal
		Power Supply by.............. Utility Power
		Utility Voltage.............. 122 V
		Output Voltage............... 122 V
		Battery Capacity............. 46 %
		Remaining Runtime............ 28 min.
		Load......................... 120 Watt(12 %)
		Line Interaction............. None
		Test Result.................. Passed at 2023/03/09 13:25:33
		Last Power Event............. Blackout at 2023/03/09 12:55:09 for 3 sec.

`

var testOutputBlackout = `
The UPS information shows as following:

	Properties:
		Model Name................... CP1500PFCLCDa
		Firmware Number.............. CR01802B7H21
		Rating Voltage............... 120 V
		Rating Power................. 1000 Watt(1500 VA)

	Current UPS status:
		State........................ Power Failure
		Power Supply by.............. Battery Power
		Utility Voltage.............. 0 V
		Output Voltage............... 120 V
		Battery Capacity............. 39 %
		Remaining Runtime............ 24 min.
		Load......................... 120 Watt(12 %)
		Line Interaction............. None
		Test Result.................. Passed at 2023/03/09 13:25:33
		Last Power Event............. Blackout at 2023/03/09 13:38:21

`

func TestParsePowerStats(t *testing.T) {
	t.Parallel()

	var status, device, err = parsePowerStats(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "Normal", status.State)
	assert.Equal(t, "Utility Power", status.PowerSupplyBy)
	assert.Equal(t, 122, status.UtilityVoltage)
	assert.Equal(t, 122, status.OutputVoltage)
	assert.Equal(t, 46, status.BatteryCapacity)
	assert.Equal(t, time.Duration(28)*time.Minute, status.RemainingRuntime)
	assert.Equal(t, 120, status.LoadWatts)
	assert.Equal(t, 12, status.LoadPct)
	assert.Equal(t, "None", status.LineInteraction)
	assert.Equal(t, "Passed", status.TestResult)
	assert.Equal(t, time.Date(2023, time.March, 9, 13, 25, 33, 0, status.TestResultTime.Location()), status.TestResultTime)
	assert.Equal(t, "Blackout", status.LastPowerEvent)
	assert.Equal(t, time.Date(2023, time.March, 9, 12, 55, 9, 0, status.LastPowerEventTime.Location()), status.LastPowerEventTime)
	assert.Equal(t, time.Duration(3)*time.Second, status.LastPowerEventDuration)
	assert.Equal(t, "CP1500PFCLCDa", device.ModelName)
	assert.Equal(t, "CR01802B7H21", device.FirmwareNumber)
	assert.Equal(t, 120, device.RatingVoltage)
	assert.Equal(t, 1000, device.RatingPowerWatts)
	assert.Equal(t, 1500, device.RatingPowerVA)
}

func TestParsePowerStatsBlackout(t *testing.T) {
	t.Parallel()

	var status, device, err = parsePowerStats(testOutputBlackout)
	assert.NoError(t, err)
	assert.Equal(t, "Power Failure", status.State)
	assert.Equal(t, "Battery Power", status.PowerSupplyBy)
	assert.Equal(t, 0, status.UtilityVoltage)
	assert.Equal(t, 120, status.OutputVoltage)
	assert.Equal(t, 39, status.BatteryCapacity)
	assert.Equal(t, time.Duration(24)*time.Minute, status.RemainingRuntime)
	assert.Equal(t, 120, status.LoadWatts)
	assert.Equal(t, 12, status.LoadPct)
	assert.Equal(t, "None", status.LineInteraction)
	assert.Equal(t, "Passed", status.TestResult)
	assert.Equal(t, time.Date(2023, time.March, 9, 13, 25, 33, 0, status.TestResultTime.Location()), status.TestResultTime)
	assert.Equal(t, "Blackout", status.LastPowerEvent)
	assert.Equal(t, time.Date(2023, time.March, 9, 13, 38, 21, 0, status.LastPowerEventTime.Location()), status.LastPowerEventTime)
	assert.Equal(t, time.Duration(0), status.LastPowerEventDuration)
	assert.Equal(t, "CP1500PFCLCDa", device.ModelName)
	assert.Equal(t, "CR01802B7H21", device.FirmwareNumber)
	assert.Equal(t, 120, device.RatingVoltage)
	assert.Equal(t, 1000, device.RatingPowerWatts)
	assert.Equal(t, 1500, device.RatingPowerVA)
}

// func TestGetLastPowerEventBlackout(t *testing.T) {
// 	t.Parallel()

// 	var lastPowerEvent, lastPowerEventTime, lastPowerEventDuration, err = getLastPowerEvent("Last Power Event............. Blackout at 2023/03/09 11:43:36")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Blackout", lastPowerEvent)
// 	assert.Equal(t, time.Date(2023, time.March, 9, 11, 43, 36, 0, lastPowerEventTime.Location()), lastPowerEventTime)
// 	assert.Equal(t, time.Duration(0), lastPowerEventDuration)

// 	lastPowerEvent, lastPowerEventTime, lastPowerEventDuration, err = getLastPowerEvent("	Last Power Event............. Blackout at 2023/03/09 11:43:36")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Blackout", lastPowerEvent)
// 	assert.Equal(t, time.Date(2023, time.March, 9, 11, 43, 36, 0, lastPowerEventTime.Location()), lastPowerEventTime)
// 	assert.Equal(t, time.Duration(0), lastPowerEventDuration)

// 	lastPowerEvent, lastPowerEventTime, lastPowerEventDuration, err = getLastPowerEvent("		Last Power Event............. Blackout at 2023/03/09 11:43:36")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Blackout", lastPowerEvent)
// 	assert.Equal(t, time.Date(2023, time.March, 9, 11, 43, 36, 0, lastPowerEventTime.Location()), lastPowerEventTime)
// 	assert.Equal(t, time.Duration(0), lastPowerEventDuration)
// }
