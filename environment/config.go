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
	CIDR       net.IPNet
	PrivateKey string
}

type Config struct {
	Name          string
	Version       version.Version
	CFDomain      string
	AppsDomain    string
	OpsManager    OpsManager
	PasCIDR       net.IPNet
	PasSubnet     string
	ServicesCIDR  net.IPNet
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
	OpsManCIDR    string   `json:"ops_manager_cidr"`
	PasCIDR       string   `json:"ert_cidr"`
	PasSubnet     string   `json:"ert_subnet"`
	ServicesCIDR  string   `json:"services_cidr"`
	ServiceSubnet string   `json:"service_subnet_name"`
	AZs           []string `json:"azs"`
	OpsManager    struct {
		URL      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"ops_manager"`
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
	version, err := version.NewVersion(data.Version)
	if err != nil {
		return Config{}, err
	}

	url, err := url.Parse(data.OpsManager.URL)
	if err != nil {
		return Config{}, err
	}

	ip := net.ParseIP(data.IP)
	if ip == nil {
		return Config{}, fmt.Errorf("Could not parse IP address: %s", data.IP)
	}

	_, opsManCIDR, err := net.ParseCIDR(data.OpsManCIDR)
	if err != nil {
		return Config{}, err
	}

	_, pasCIDR, err := net.ParseCIDR(data.PasCIDR)
	if err != nil {
		return Config{}, err
	}

	_, servicesCIDR, err := net.ParseCIDR(data.ServicesCIDR)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Name:          data.Name,
		Version:       *version,
		CFDomain:      data.SysDomain,
		AppsDomain:    data.AppsDomain,
		PasCIDR:       *pasCIDR,
		PasSubnet:     data.PasSubnet,
		ServicesCIDR:  *servicesCIDR,
		ServiceSubnet: data.ServiceSubnet,
		AZs:           data.AZs,
		OpsManager: OpsManager{
			Username:   data.OpsManager.Username,
			Password:   data.OpsManager.Password,
			URL:        *url,
			IP:         ip,
			CIDR:       *opsManCIDR,
			PrivateKey: data.PrivateKey,
		},
	}, nil
}
