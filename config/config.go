package config

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	DBName   string `json:"db_name"`
	DBSource string `json:"db_source"`
	EnableOrmLog bool `json:"enable_orm_log"`
}

var c config

func LoadConfig(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		println("not find file")
		return err
	}
	println(path)

	var ctmp config
	err = json.Unmarshal(b, &ctmp)
	if err != nil {
		print("1")
		return err
	}

	c = ctmp
	return nil
}

func GetDBName() string {
	return c.DBName
}

func GetDBSource() string {
	return c.DBSource
}

func IsOrmLogEnabled() bool {
	return c.EnableOrmLog
}
