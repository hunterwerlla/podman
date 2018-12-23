package main

import (
	"bytes"
	"encoding/json" //for settings
	"fmt"
	"io/ioutil"
	"os"
)

//read config in
func readConfig(c Configuration) Configuration {
	//check if there is a config file
	config, err := os.Open("./config.json")
	defer config.Close()
	if err != nil {
		//config does not exist so build one out of the defult settings
		//first check if the storage location is ok
		if _, err := os.Stat(c.StorageLocation); os.IsNotExist(err) {
			//path does not exist try to make
			fmt.Printf("making folder at:%s\n", c.StorageLocation)
			err := os.MkdirAll(c.StorageLocation, 0700)
			if err != nil {
				fmt.Println("cannot make folder everyting kill")
				//failed to create folder to store, store files in same directory as program
				c.StorageLocation = "."
			}
		}
		return c
	}
	buffer, err := ioutil.ReadAll(config)
	if err != nil {
		panic("could not read config file")
	}
	json.Unmarshal(buffer, &c)
	//now read in the settings and write it to the configuration object
	return c
}

//save current config to file
func writeConfig(c Configuration) {
	config, err := os.Create("./config.json")
	if err != nil {
		//using default settings because cannot write settings
		panic("could not save settings, cannot continue")
	}
	defer func() {
		err := config.Close()
		if err != nil {

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
