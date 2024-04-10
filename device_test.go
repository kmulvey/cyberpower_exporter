package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModelName(t *testing.T) {
	t.Parallel()

	var name, err = getModelName(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "CP1500PFCLCDa", name)

	name, err = getModelName("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the model name, err: could not find any matches", err.Error())
	assert.Equal(t, "", name)
}

func TestGetFirmwareNumber(t *testing.T) {
	t.Parallel()

	var name, err = getFirmwareNumber(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, "CR01802B7H21", name)

	name, err = getFirmwareNumber("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the firmware number, err: could not find any matches", err.Error())
	assert.Equal(t, "", name)
}

func TestGetRatingVoltage(t *testing.T) {
	t.Parallel()

	var name, err = getRatingVoltage(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, 120, name)

	name, err = getRatingVoltage("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the rating voltage, err: could not find any matches", err.Error())
	assert.Equal(t, 0, name)
}

func TestGetRatingPowerWatts(t *testing.T) {
	t.Parallel()

	var watts, va, err = getRatingPowerWatts(testOutputNormal)
	assert.NoError(t, err)
	assert.Equal(t, 1000, watts)
	assert.Equal(t, 1500, va)

	watts, va, err = getRatingPowerWatts("testOutputNormal") // bad string
	assert.Error(t, err)
	assert.Equal(t, "unable to find the rating power in watts, err: could not find any matches", err.Error())
	assert.Equal(t, 0, watts)
	assert.Equal(t, 0, va)
}
