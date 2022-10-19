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

func TestParsePowerStats(t *testing.T) {
	t.Parallel()

	var ds, err = parsePowerStats(testOutput)
	assert.NoError(t, err)
	assert.Equal(t, "Normal", ds.State)
	assert.Equal(t, "Utility Power", ds.PowerSupplyBy)
	assert.Equal(t, 122, ds.UtilityVoltage)
	assert.Equal(t, 122, ds.OutputVoltage)
	assert.Equal(t, 100, ds.BatteryCapacity)
	assert.Equal(t, time.Duration(55)*time.Minute, ds.RemainingRuntime)
	assert.Equal(t, 140, ds.LoadWatts)
	assert.Equal(t, 14, ds.LoadPct)
	assert.Equal(t, 14, ds.LoadPct)
	assert.Equal(t, "None", ds.LineInteraction)
	assert.Equal(t, "Passed", ds.TestResult)
	assert.Equal(t, time.Date(2022, time.October, 4, 13, 56, 27, 0, ds.TestResultTime.Location()), ds.TestResultTime)
	assert.Equal(t, "Blackout", ds.LastPowerEvent)
	assert.Equal(t, time.Date(2022, time.September, 24, 12, 12, 24, 0, ds.TestResultTime.Location()), ds.LastPowerEventTime)
	assert.Equal(t, time.Duration(3)*time.Second, ds.LastPowerEventDuration)
}
