package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Payload struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

type Response struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
	Content string `json:"content"`
}

func ensureTableExists(ctx context.Context, logger runtime.Logger, db *sql.DB) error {
	dropTableQuery := `DROP TABLE IF EXISTS files;`

	createTableQuery := `
    CREATE TABLE files (
        id SERIAL PRIMARY KEY,
        type VARCHAR(255) NOT NULL,
        version VARCHAR(255) NOT NULL,
        hash VARCHAR(255) NOT NULL,
        content BYTEA NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	logger.Debug("Dropping table if exists...")
	_, err := db.ExecContext(ctx, dropTableQuery)
	if err != nil {
		logger.Error("Failed to drop table: %v", err)
		return err
	}

	logger.Debug("Creating table...")
	_, err2 := db.ExecContext(ctx, createTableQuery)
	if err2 != nil {
		logger.Error("Failed to create table: %v", err2)
		return err2
	}
	logger.Info("Table 'files' ensured to exist")
	return nil
}
func RpcFunction(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	defaultType := "core"
	defaultVersion := "1.0.0"
	var p Payload

	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return "", err
	}

	if p.Type == "" {
		p.Type = defaultType
	}
	if p.Version == "" {
		p.Version = defaultVersion
	}

	filePath := filepath.Join("/nakama/data", p.Type, p.Version+".json")
	//logger.Debug("file path {}", filePath) commeted due to test

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("file does not exist")
			return "", errors.New("file does not exist")
		}
		return "", err
	}

	hash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(hash[:])
	//logger.Info("hash {}", hashStr)

	var response Response
	response.Type = p.Type
	response.Version = p.Version
	response.Hash = hashStr
	//logger.Info("response {}", response)

	if p.Hash != "" && p.Hash != hashStr {
		response.Content = ""
	} else {
		response.Content = string(content)
	}

	if _, err := db.ExecContext(ctx, "INSERT INTO files (type, version, hash, content) VALUES ($1, $2, $3, $4)", p.Type, p.Version, hashStr, content); err != nil {
		return "", err
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {

	logger.Debug("Initializing module...")

	if err := ensureTableExists(ctx, logger, db); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("my_rpc_function", RpcFunction); err != nil {
		return err
	}
	return nil
}

// func main(){

// }
