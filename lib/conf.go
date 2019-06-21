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

package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/distribution/uuid"
)

const CONFPATH = "./conf.json"

var CONF Configuration

func InitConf() {
	if _, err := os.Stat(CONFPATH); err == nil {
		conf := readConf()
		if conf.Id == "" {
			WriteConf(Configuration{Id: uuid.Generate().String()})
		}

	} else if os.IsNotExist(err) {
		fmt.Println("Creating config")
		f, err := os.Create(CONFPATH)
		if err != nil {
			fmt.Println("error:", err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		WriteConf(Configuration{Id: uuid.Generate().String()})

	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
	}
	CONF = readConf()
}

func WriteConf(confNew Configuration) {
	confJson, _ := json.Marshal(confNew)
	err := ioutil.WriteFile(CONFPATH, confJson, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetConf() (conf Configuration) {
	return CONF
}

func readConf() (configuration Configuration) {
	f, _ := os.Open(CONFPATH)
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
