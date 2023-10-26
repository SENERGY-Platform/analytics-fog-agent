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
	"path/filepath"

	agentEntities "github.com/SENERGY-Platform/analytics-fog-lib/lib/agent"

	"github.com/SENERGY-Platform/analytics-fog-agent/lib/constants"
	"github.com/SENERGY-Platform/analytics-fog-agent/lib/logging"

	"github.com/docker/distribution/uuid"
)

var CONF agentEntities.Configuration

func InitConf(dataDir string) {
	configPath := filepath.Join(dataDir, constants.ConfFileName)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		_ = os.Mkdir(dataDir, 0700)
	}

	if _, err := os.Stat(configPath); err == nil {
		logging.Logger.Debug("Read agent config from ", configPath)

		conf := readConf(dataDir)
		if conf.Id == "" {
			WriteConf(dataDir, agentEntities.Configuration{Id: uuid.Generate().String()})
		}

	} else if os.IsNotExist(err) {
		logging.Logger.Debug("Creating agent config at ", configPath)
		f, err := os.Create(configPath)
		if err != nil {
			fmt.Println("error:", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		WriteConf(dataDir, agentEntities.Configuration{Id: uuid.Generate().String()})

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
	}
	CONF = readConf(dataDir)
}

func WriteConf(dataDir string, confNew agentEntities.Configuration) {
	configPath := filepath.Join(dataDir, constants.ConfFileName)

	confJson, _ := json.Marshal(confNew)
	err := ioutil.WriteFile(configPath, confJson, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetConf() (conf agentEntities.Configuration) {
	return CONF
}

func readConf(dataDir string) (configuration agentEntities.Configuration) {
	configPath := filepath.Join(dataDir, constants.ConfFileName)

	f, _ := os.Open(configPath)
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
