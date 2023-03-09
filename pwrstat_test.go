package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testOutput = `
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
`

var testFailOutput = `
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
		Battery Capacity............. 82 %
		Remaining Runtime............ 48 min.
		Load......................... 150 Watt(15 %)
		Line Interaction............. None
		Test Result.................. Passed at 2023/03/06 10:12:45
		Last Power Event............. Blackout at 2023/03/08 14:25:31
`

func TestParsePowerStats(t *testing.T) {
	t.Parallel()

	var status, device, err = parsePowerStats(testOutput)
	assert.NoError(t, err)
	assert.Equal(t, "Normal", status.State)

	assert.Equal(t, "Utility Power", status.PowerSupplyBy)
	assert.Equal(t, 122, status.UtilityVoltage)
	assert.Equal(t, 122, status.OutputVoltage)
	assert.Equal(t, 100, status.BatteryCapacity)
	assert.Equal(t, time.Duration(55)*time.Minute, status.RemainingRuntime)
	assert.Equal(t, 140, status.LoadWatts)
	assert.Equal(t, 14, status.LoadPct)
	assert.Equal(t, 14, status.LoadPct)
	assert.Equal(t, "None", status.LineInteraction)
	assert.Equal(t, "Passed", status.TestResult)
	assert.Equal(t, time.Date(2022, time.October, 4, 13, 56, 27, 0, status.TestResultTime.Location()), status.TestResultTime)
	assert.Equal(t, "Blackout", status.LastPowerEvent)
	assert.Equal(t, time.Date(2022, time.September, 24, 12, 12, 24, 0, status.TestResultTime.Location()), status.LastPowerEventTime)
	assert.Equal(t, time.Duration(3)*time.Second, status.LastPowerEventDuration)
	assert.Equal(t, "CP1500PFCLCDa", device.ModelName)
	assert.Equal(t, "CXXKY2008826", device.FirmwareNumber)
	assert.Equal(t, 120, device.RatingVoltage)
	assert.Equal(t, 1000, device.RatingPowerWatts)
	assert.Equal(t, 1500, device.RatingPowerVA)
}
