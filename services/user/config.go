package user

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
	DBDriver   string `json:"db_driver"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBAddress  string `json:"db_address"`
	DBTable    string `json:"db_table"`
	dbString   string

	EncodingJWT string `json:"encoding_jwt"`
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

	c.dbString = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		c.DBUser,
		c.DBPassword,
		c.DBAddress,
		c.DBTable,
	)

	return err
}
