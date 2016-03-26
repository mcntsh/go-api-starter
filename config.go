package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
)

// Config represents individual configuration definitions
// derrived from the `config.json` file.
type Config struct {
	HTTPAddress string `json:"http_address"`
	HTTPTimeout int    `json:"http_timeout"`
}

var config = &Config{}

func (c *Config) LoadJSON() error {
	_, curPath, _, _ := runtime.Caller(0)
	curDir := path.Dir(curPath)

	cFile, err := os.Open(path.Join(curDir, "config.json"))
	if err != nil {
		return fmt.Errorf("Could not load config file")
	}
	defer cFile.Close()

	cParser := json.NewDecoder(cFile)
	err = cParser.Decode(c)
	if err != nil {
		return fmt.Errorf("Error unmarshalling config into struct")
	}

	return err
}
