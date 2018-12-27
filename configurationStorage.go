package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

const configStorageLocation = "/.config/podman/"
const configName = "config.json"

func getConfigStorage() string {
	usr, err := user.Current()
	var defaultStorage string
	if err == nil {
		defaultStorage = usr.HomeDir + configStorageLocation
	} else {
		defaultStorage = "."
	}
	return defaultStorage
}

//read config in
func ReadConfig(c Configuration) Configuration {
	// make the config location if needed
	if _, err := os.Stat(c.StorageLocation); os.IsNotExist(err) {
		// try to make path first,
		fmt.Printf("making folder at:%s\n", c.StorageLocation)
		err = os.MkdirAll(c.StorageLocation, 0700)
		if err != nil {
			fmt.Println("cannot make folder, defaulting storage to local directory")
			//failed to create folder to store, store files in same directory as program
		}
		fmt.Printf("making folder at:%s\n", getConfigStorage())
		err = os.MkdirAll(getConfigStorage(), 0700)
		WriteConfig(&c)
		return c
	}
	config, err := os.Open(getConfigStorage() + configName)
	if err != nil {
		panic("could not read config file")
	}
	defer config.Close()
	buffer, err := ioutil.ReadAll(config)
	if err != nil {
		panic("could not read config file")
	}
	json.Unmarshal(buffer, &c)
	//now read in the settings and write it to the configuration object
	return c
}

//save current config to file
func WriteConfig(c *Configuration) {

	config, err := os.Create(getConfigStorage() + configName)
	if err != nil {
		//using default settings because cannot write settings
		panic("could not save settings, cannot continue")
	}
	defer func() {
		err = config.Close()
		if err != nil {
			fmt.Printf("Unable to save config!")
		}
	}()
	jsonSettings, err := json.Marshal(&c)
	if err != nil {
		panic("could not save settings, cannot continue")
	}
	var jsonSettingsPretty bytes.Buffer
	err = json.Indent(&jsonSettingsPretty, jsonSettings, "", "    ")
	if err != nil {
		config.Write(jsonSettings)
		return
	}
	config.Write(jsonSettingsPretty.Bytes())
}