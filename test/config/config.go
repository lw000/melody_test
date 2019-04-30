package config

import (
	"encoding/json"
	"io/ioutil"
)

type JsonConfig struct {
	Host        string `json:"host"`
	Scheme      string `json:"scheme"`
	Path        string `json:"path"`
	Count       int    `json:"count"`
	Millisecond int    `json:"millisecond"`
	Send        bool   `json:"send"`
}

func NewJsonConfig() *JsonConfig {
	return &JsonConfig{}
}

func LoadJsonConfig(file string) (*JsonConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfgStruct JsonConfig
	if err = json.Unmarshal(data, &cfgStruct); err != nil {
		return nil, err
	}

	return &cfgStruct, err
}
