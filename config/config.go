package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type Config struct {
	OauthToken string `json:"oauthToken"`
}

const (
	kConfigFileDir  = ".yfload"
	kConfigFileName = "config.json"
)

func configFilePath() string {
	user, _ := user.Current()
	filePath := path.Join(user.HomeDir, kConfigFileDir, kConfigFileName)
	return filePath
}

func ConfigLoad() (*Config, error) {
	filePath := configFilePath()
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func ConfigSave(config Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	filePath := configFilePath()

	fileDir := path.Dir(filePath)
	os.MkdirAll(fileDir, os.ModeDir)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}

	file.Write(data)
	file.Close()
	return nil
}
