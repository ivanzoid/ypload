package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"
)

type Config struct {
	OauthToken         string `json:"oauthToken"`
	ExpirationDateTime int64  `json:"expirationDateTime"`
	MainAlbumUrl       string `json:"mainAlbumUrl"`
}

const (
	kConfigFileDir  = ".ypload"
	kConfigFileName = "config.json"
)

func configFilePath() string {
	user, _ := user.Current()
	filePath := path.Join(user.HomeDir, kConfigFileDir, kConfigFileName)
	return filePath
}

func Load() (*Config, error) {
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

func (config *Config) Save() error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	filePath := configFilePath()

	fileDir := path.Dir(filePath)
	os.MkdirAll(fileDir, os.ModeDir|0700)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	file.Write(data)
	file.Close()
	return nil
}

func (config *Config) TokenExpired() bool {
	if config.OauthToken == "" {
		return true
	}

	expired := time.Now().After(time.Unix(config.ExpirationDateTime, 0))
	return expired
}

func (config *Config) UpdateExpirationDateTime(expiresIn int64) {
	config.ExpirationDateTime = time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()
}
