/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"

	"github.com/hashicorp/go-version"
)

type OpsManager struct {
	Username   string
	Password   string
	URL        url.URL
	IP         net.IP
	PrivateKey string
}

type PKSApi struct {
	Username string
	Password string
	URL      url.URL
}

type Config struct {
	Name          string
	Version       version.Version
	CFDomain      string
	AppsDomain    string
	OpsManager    OpsManager
	PKSApi        PKSApi
	PasSubnet     string
	ServiceSubnet string
	AZs           []string
}

type environmentReader struct {
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	SysDomain     string   `json:"sys_domain"`
	AppsDomain    string   `json:"apps_domain"`
	PrivateKey    string   `json:"ops_manager_private_key"`
	IP            string   `json:"ops_manager_public_ip"`
	PasSubnet     string   `json:"ert_subnet"`
	ServiceSubnet string   `json:"service_subnet_name"`
	AZs           []string `json:"azs"`
	OpsManager    struct {
		URL      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"ops_manager"`
	PKSApi struct {
		Username string `json:"uaa_admin_user"`
		Password string `json:"uaa_admin_password"`
		URL      string `json:"url"`
	} `json:"pks_api"`
}

func FromFile(path string) (Config, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var data environmentReader

	if err := json.Unmarshal(contents, &data); err != nil {
		return Config{}, err
	}

	return newLockfile(data)
}

func newLockfile(data environmentReader) (Config, error) {
	var err error

	parsedVersion := &version.Version{}
	if data.Version != "" {
		parsedVersion, err = version.NewVersion(data.Version)
		if err != nil {
			return Config{}, err
		}
	}

	parsedOpsManagerURL, err := url.Parse(data.OpsManager.URL)
	if err != nil {
		return Config{}, err
	}

	opsManagerIp := net.ParseIP(data.IP)
	if opsManagerIp == nil {
		return Config{}, fmt.Errorf("Could not parse IP address: %s", data.IP)
	}

	parsedPKSApiURL, err := url.Parse(data.PKSApi.URL)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Name:          data.Name,
		Version:       *parsedVersion,
		CFDomain:      data.SysDomain,
		AppsDomain:    data.AppsDomain,
		PasSubnet:     data.PasSubnet,
		ServiceSubnet: data.ServiceSubnet,
		AZs:           data.AZs,
		OpsManager: OpsManager{
			Username:   data.OpsManager.Username,
			Password:   data.OpsManager.Password,
			URL:        *parsedOpsManagerURL,
			IP:         opsManagerIp,
			PrivateKey: data.PrivateKey,
		},
		PKSApi: PKSApi{
			Username: data.PKSApi.Username,
			Password: data.PKSApi.Password,
			URL:      *parsedPKSApiURL,
		},
	}, nil
}
