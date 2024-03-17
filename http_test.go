package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kmulvey/goutils"
	"github.com/stretchr/testify/assert"
)

func TestWebServer(t *testing.T) {
	t.Parallel()

	var status, _, err = parsePowerStats(testOutputBlackout)
	assert.NoError(t, err)

	var fileName = goutils.RandomString(10) + ".db"

	db, err := dbConnect(fileName)
	assert.NoError(t, err)

	err = insert(db, status)
	assert.NoError(t, err)

	go webServer(":8080", db)

	res, err := http.Get("http://localhost:8080/latest")
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	var getStatus DeviceStatus
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(body, &getStatus))
	assert.NoError(t, res.Body.Close())

	// tz are a little off, dont really care so null them out
	status.CollectionTime = time.Time{}
	getStatus.CollectionTime = time.Time{}

	assert.EqualValues(t, status, getStatus)

	ddb, err := db.DB()
	assert.NoError(t, err)
	assert.NoError(t, ddb.Close())

	assert.NoError(t, os.RemoveAll(fileName))
}
