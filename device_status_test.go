package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
