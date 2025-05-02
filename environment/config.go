/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms of the under the Apache License, Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/

package environment

import (
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/hashicorp/go-version"
	"gopkg.in/yaml.v2"
)

const defaultSSHUser = "ubuntu"

type OpsManager struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	URL          url.URL
	IP           net.IP
	PrivateKey   string
	SshUser      string
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
	Name          string   `yaml:"name"`
	Version       string   `yaml:"version"`
	SysDomain     string   `yaml:"sys_domain"`
	AppsDomain    string   `yaml:"apps_domain"`
	PrivateKey    string   `yaml:"ops_manager_private_key"`
	IP            string   `yaml:"ops_manager_public_ip"`
	SshUser       string   `yaml:"ops_manager_ssh_user"`
	PasSubnet     string   `yaml:"ert_subnet"`
	ServiceSubnet string   `yaml:"service_subnet_name"`
	AZs           []string `yaml:"azs"`
	OpsManager    struct {
		URL          string `yaml:"url"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"ops_manager"`
	PKSApi struct {
		Username string `yaml:"uaa_admin_user"`
		Password string `yaml:"uaa_admin_password"`
		URL      string `yaml:"url"`
	} `yaml:"pks_api"`
}

func FromFile(path, environmentName string) (Config, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var selectedEnvironment environmentReader
	err = yaml.Unmarshal(contents, &selectedEnvironment)
	if err == nil {
		if environmentName != "" && environmentName != selectedEnvironment.Name {
			return Config{}, fmt.Errorf("Environment name '%s' specified but does not match environment in config", environmentName)
		}
		return newLockfile(selectedEnvironment)
	}

	var environments []environmentReader
	if arrayErr := yaml.Unmarshal(contents, &environments); arrayErr != nil {
		return Config{}, fmt.Errorf("Unable to unmarshal specified config as either single environment, '%s' or multiple environments, '%s'", err.Error(), arrayErr.Error())
	}

	if len(environments) == 0 {
		return Config{}, fmt.Errorf("Target config is an empty array")
	}

	if environmentName == "" {
		return newLockfile(environments[0])
	}

	for _, environment := range environments {
		if environmentName == environment.Name {
			selectedEnvironment = environment
		}
	}

	if selectedEnvironment.Name == "" {
		return Config{}, fmt.Errorf("Environment name '%s' specified but does not match environment in config", environmentName)
	}
	return newLockfile(selectedEnvironment)

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

	sshUser := data.SshUser
	if sshUser == "" {
		sshUser = defaultSSHUser
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
			Username:     data.OpsManager.Username,
			Password:     data.OpsManager.Password,
			ClientID:     data.OpsManager.ClientID,
			ClientSecret: data.OpsManager.ClientSecret,
			URL:          *parsedOpsManagerURL,
			IP:           opsManagerIp,
			PrivateKey:   data.PrivateKey,
			SshUser:      sshUser,
		},
		PKSApi: PKSApi{
			Username: data.PKSApi.Username,
			Password: data.PKSApi.Password,
			URL:      *parsedPKSApiURL,
		},
	}, nil
}
