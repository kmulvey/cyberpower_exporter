package main

import (
	"os"
	"testing"
	"time"

	"github.com/kmulvey/goutils"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	t.Parallel()

	var status, _, err = parsePowerStats(testOutputBlackout)
	assert.NoError(t, err)

	var fileName = goutils.RandomString(10) + ".db"

	db, err := dbConnect(fileName)
	assert.NoError(t, err)

	err = insert(db, status)
	assert.NoError(t, err)

	dbStatus, err := getLatest(db)
	assert.NoError(t, err)

	// tz are a little off, dont really care so null them out
	status.CollectionTime = time.Time{}
	dbStatus.CollectionTime = time.Time{}

	assert.EqualValues(t, status, dbStatus)

	assert.NoError(t, os.RemoveAll(fileName))
}
