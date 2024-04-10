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

func TestGetState(t *testing.T) {
	t.Parallel()

	var state, err = getState(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "Normal", state)

	state, err = getState("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to get the state, err: could not find any matches", err.Error())
	assert.Equal(t, "", state)
}

func TestGetPowerSupply(t *testing.T) {
	t.Parallel()

	var powerSupply, err = getPowerSupply(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "Utility Power", powerSupply)

	powerSupply, err = getPowerSupply("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to get the power supply, err: could not find any matches", err.Error())
	assert.Equal(t, "", powerSupply)
}

func TestGetUtilityVoltage(t *testing.T) {
	t.Parallel()

	var utilityVoltage, err = getUtilityVoltage(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, 122, utilityVoltage)

	utilityVoltage, err = getUtilityVoltage("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the utility voltage, err: could not find any matches", err.Error())
	assert.Equal(t, 0, utilityVoltage)
}

func TestGetOutputVoltage(t *testing.T) {
	t.Parallel()

	var outputVoltage, err = getOutputVoltage(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, 122, outputVoltage)

	outputVoltage, err = getOutputVoltage("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the output voltage, err: could not find any matches", err.Error())
	assert.Equal(t, 0, outputVoltage)
}

func TestGetBatteryCapacity(t *testing.T) {
	t.Parallel()

	var batteryCapacity, err = getBatteryCapacity(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, 46, batteryCapacity)

	batteryCapacity, err = getBatteryCapacity("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the battery capacity, err: could not find any matches", err.Error())
	assert.Equal(t, 0, batteryCapacity)
}

func TestGetRemainingRuntime(t *testing.T) {
	t.Parallel()

	var remainingRuntime, err = getRemainingRuntime(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(1680000000000), remainingRuntime)

	remainingRuntime, err = getRemainingRuntime("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the remaining runtime, err: could not find any matches", err.Error())
	assert.Equal(t, time.Duration(0), remainingRuntime)
}

func TestGetLoad(t *testing.T) {
	t.Parallel()

	var watts, pct, err = getLoad(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, 120, watts)
	assert.Equal(t, 12, pct)

	watts, pct, err = getLoad("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the load, err: could not find any matches", err.Error())
	assert.Equal(t, 0, watts)
	assert.Equal(t, 0, pct)
}

func TestGetLineInteraction(t *testing.T) {
	t.Parallel()

	var lineInteraction, err = getLineInteraction(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "None", lineInteraction)

	lineInteraction, err = getLineInteraction("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the line interaction, err: could not find any matches", err.Error())
	assert.Equal(t, "", lineInteraction)
}

func TestGetTestResult(t *testing.T) {
	t.Parallel()

	var result, date, err = getTestResult(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "Passed", result)
	assert.Equal(t, time.Time(time.Date(2023, time.March, 9, 13, 25, 33, 0, time.UTC)), date)

	/* getDeviceInfoAsString needs to return a slice in order to test this
	result, date, err = getTestResult("Test Result.................. Passed at \n") // missing date string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the last test result, err: could not find any matches", err.Error())
	assert.Equal(t, "", result)
	assert.Equal(t, time.Time{}, date)
	*/

	result, date, err = getTestResult("Test Result.................. Passed at 2023/03/\n") // bad date string
	assert.Error(t, err)
	assert.Equal(t, `unable to pasre date: 2023/03/, err: parsing time "2023/03/" as "2006/01/02 15:04:05": cannot parse "" as "02"`, err.Error())
	assert.Equal(t, "", result)
	assert.Equal(t, time.Time{}, date)

	result, date, err = getTestResult("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the last test result, err: could not find any matches", err.Error())
	assert.Equal(t, "", result)
	assert.Equal(t, time.Time{}, date)
}

func TestGetLastPowerEvent(t *testing.T) {
	t.Parallel()

	var result, date, duration, err = getLastPowerEvent(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "Blackout", result)
	assert.Equal(t, time.Date(2023, time.March, 9, 12, 55, 9, 0, time.UTC), date)
	assert.Equal(t, time.Duration(3000000000), duration)

	result, date, duration, err = getLastPowerEvent("Last Power Event............. None\n") // no power event
	assert.NoError(t, err)
	assert.Equal(t, "None", result)
	assert.Equal(t, time.Time{}, date)
	assert.Equal(t, time.Duration(0), duration)

	result, date, duration, err = getLastPowerEvent("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, time.Time{}, date)
	assert.Equal(t, time.Duration(0), duration)
}
