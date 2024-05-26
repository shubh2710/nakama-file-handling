package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// test to cover RPC funtion
func TestRpcFunction(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()

	payload := `{"type": "core", "version": "1.0.0"}`
	content := `{"key": "value"}`

	filePath := "/nakama/data/core/1.0.0.json"
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	require.NoError(t, err)
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)
	defer os.RemoveAll("/nakama/data/core")

	hash := sha256.Sum256([]byte(content))
	expectedHash := hex.EncodeToString(hash[:])

	mock.ExpectExec("INSERT INTO files").
		WithArgs("core", "1.0.0", expectedHash, []byte(content)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	response, err := RpcFunction(ctx, nil, db, nil, payload)
	require.NoError(t, err)

	expectedResponse := Response{
		Type:    "core",
		Version: "1.0.0",
		Hash:    expectedHash,
		Content: content,
	}
	expectedResponseBytes, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	assert.JSONEq(t, string(expectedResponseBytes), response)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
