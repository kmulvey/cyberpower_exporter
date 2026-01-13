package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// nolint: gochecknoglobals
var (
	testOutputNormal = `
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

	testOutputBlackout = `
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

	testLostConnection = `
The UPS information shows as following:

	Properties:
		Model Name................... CP1500PFCLCDa
		Firmware Number.............. CR01802B7H21
		Rating Voltage............... 120 V
		Rating Power................. 1000 Watt(1500 VA)

	Current UPS status:
		State........................ Lost Communication
		Test Result.................. Passed at 2025/01/21 13:13:05
		Last Power Event............. Blackout at 2025/01/23 12:33:09

`
)

func TestParsePowerStatusNormal(t *testing.T) {
	t.Parallel()

	var status, err = parsePowerStatus(testOutputNormal)
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
}

func TestParsePowerStatusBlackout(t *testing.T) {
	t.Parallel()

	var status, err = parsePowerStatus(testOutputBlackout)
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
}

func TestParsePowerStatusLostConn(t *testing.T) {
	t.Parallel()

	var status, err = parsePowerStatus(testLostConnection)
	assert.NoError(t, err)

	assert.Equal(t, "Lost Communication", status.State)
	assert.Empty(t, status.PowerSupplyBy)
	assert.Equal(t, 0, status.UtilityVoltage)
	assert.Equal(t, 0, status.OutputVoltage)
	assert.Equal(t, 0, status.BatteryCapacity)
	assert.Equal(t, time.Duration(0), status.RemainingRuntime)
	assert.Equal(t, 0, status.LoadWatts)
	assert.Equal(t, 0, status.LoadPct)
	assert.Empty(t, status.LineInteraction)
	assert.Equal(t, "Passed", status.TestResult)
	assert.Equal(t, time.Date(2025, time.January, 21, 13, 13, 5, 0, status.TestResultTime.Location()), status.TestResultTime)
	assert.Equal(t, "Blackout", status.LastPowerEvent)
	assert.Equal(t, time.Date(2025, time.January, 23, 12, 33, 9, 0, status.LastPowerEventTime.Location()), status.LastPowerEventTime)
	assert.Equal(t, time.Duration(0), status.LastPowerEventDuration)
}

func TestParseDevicePropertiesNormal(t *testing.T) {
	t.Parallel()

	var device, err = parseDeviceProperties(testOutputNormal)
	assert.NoError(t, err)

	assert.Equal(t, "CP1500PFCLCDa", device.ModelName)
	assert.Equal(t, "CR01802B7H21", device.FirmwareNumber)
	assert.Equal(t, 120, device.RatingVoltage)
	assert.Equal(t, 1000, device.RatingPowerWatts)
	assert.Equal(t, 1500, device.RatingPowerVA)
}

func TestParseDevicePropertiesBlackout(t *testing.T) {
	t.Parallel()

	var device, err = parseDeviceProperties(testOutputBlackout)
	assert.NoError(t, err)

	assert.Equal(t, "CP1500PFCLCDa", device.ModelName)
	assert.Equal(t, "CR01802B7H21", device.FirmwareNumber)
	assert.Equal(t, 120, device.RatingVoltage)
	assert.Equal(t, 1000, device.RatingPowerWatts)
	assert.Equal(t, 1500, device.RatingPowerVA)
}

func TestParseDevicePropertiesLostConn(t *testing.T) {
	t.Parallel()

	var device, err = parseDeviceProperties(testLostConnection)
	assert.NoError(t, err)

	assert.Equal(t, "CP1500PFCLCDa", device.ModelName)
	assert.Equal(t, "CR01802B7H21", device.FirmwareNumber)
	assert.Equal(t, 120, device.RatingVoltage)
	assert.Equal(t, 1000, device.RatingPowerWatts)
	assert.Equal(t, 1500, device.RatingPowerVA)
}
