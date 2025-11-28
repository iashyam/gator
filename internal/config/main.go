package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding the home direcotre %v", err)
	}
	fullPath := fmt.Sprintf("%s/%s", homeDir, configFileName)
	return fullPath, nil
}

func ReadConfig() (Config, error) {

	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	jsonContent, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading the json %v", err)
	}

	var ret Config
	err = json.Unmarshal(jsonContent, &ret)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshalling jsong %v", err)
	}

	return ret, nil
}

func (conf *Config) SetUser(user string) error {
	conf.Current_user_name = user
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	jsonFile, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	if err := encoder.Encode(conf); err != nil {
		return fmt.Errorf("error encoding json %v", err)
	}

	return nil

}

func (conf *Config) PrintCofig() {
	fmt.Printf("DB URL : %s\n", conf.Db_url)
	fmt.Printf("Current User : %s\n", conf.Current_user_name)
}
