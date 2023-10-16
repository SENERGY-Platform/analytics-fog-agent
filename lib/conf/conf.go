/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/entities"
	"github.com/docker/distribution/uuid"
)

const ConfigPath = "./data/conf.json"

var CONF entities.Configuration

func InitConf() {
	if _, err := os.Stat("./data/"); os.IsNotExist(err) {
		_ = os.Mkdir("./data/", 0700)
	}
	if _, err := os.Stat(ConfigPath); err == nil {
		conf := readConf()
		if conf.Id == "" {
			WriteConf(entities.Configuration{Id: uuid.Generate().String()})
		}

	} else if os.IsNotExist(err) {
		fmt.Println("Creating config")
		f, err := os.Create(ConfigPath)
		if err != nil {
			fmt.Println("error:", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		WriteConf(entities.Configuration{Id: uuid.Generate().String()})

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
	}
	CONF = readConf()
}

func WriteConf(confNew entities.Configuration) {
	confJson, _ := json.Marshal(confNew)
	err := ioutil.WriteFile(ConfigPath, confJson, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetConf() (conf entities.Configuration) {
	return CONF
}

func readConf() (configuration entities.Configuration) {
	f, _ := os.Open(ConfigPath)
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	decoder := json.NewDecoder(f)
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}
