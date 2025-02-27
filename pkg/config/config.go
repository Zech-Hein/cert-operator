/*
Sample config file

notifiers:
- name: myslacknotifier
  kind: slack
  integration_url: https://hooks.slack.com/services/service_id/auth-token
*/

package config

import (
	"encoding/json"
	"os"

	config "github.com/micro/go-micro"
	"github.com/micro/go-micro/config/source/env"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/config/source/flag"
	"github.com/micro/go-micro/config/source/memory"
	"github.com/Zech-Hein/cert-operator/pkg/certs"
	"github.com/sirupsen/logrus"
)

type Config struct {
	//	Notifiers []notifier.Notifier  `json:"notifiers"`
	Provider certs.ProviderConfig `json:"provider"`
	General  GeneralConfig        `json:"general"`
}

type GeneralConfig struct {
	Annotations AnnotationConfig `json:"annotations"`
}

type AnnotationConfig struct {
	Status        string `json:"status"`
	StatusReason  string `json:"status-reason"`
	Expiry        string `json:"expiry"`
	Format        string `json:"format"`
	NeedCertValue string `json:"need-cert-value"`
	PemFormat     string `json:"pem-format-value"`
	Pkcs12Format  string `json:"pkcs12-format-value"`
}

const (
	defaultConfigFile = "/etc/cert-operator/config.yaml"
	defaultConfig     = `
  {
    "general": {
      "annotations": {
        "status": "openshift.io/cert-ctl-status",
        "status-reason": "openshift.io/cert-ctl-status-reason",
        "expiry": "openshift.io/cert-ctl-expires",
        "format": "openshift.io/cert-ctl-format",
        "need-cert-value": "new",
        "pem-format-value": "PEM",
        "pkcs12-format-value": "PKCS12"
      }
    },
    "provider": {
      "kind": "self-signed",
      "ssl": "false"
    }
  }`
)

func NewConfig() Config {

	tmpConfig := config.NewConfig()

	data := []byte(defaultConfig)

	memorySource := memory.NewSource(
		memory.WithData(data),
	)
	// Load json config file
	tmpConfig.Load(
		memorySource,
		file.NewSource(
			file.WithPath(getConfigFile()),
		),
		env.NewSource(),
		flag.NewSource(),
	)
	var conf Config

	tmpConfig.Scan(&conf)

	return conf
}

func getConfigFile() (configFile string) {
	if value, ok := os.LookupEnv("CERT_OP_CONFIG"); ok {
		logrus.Infof("Loading custom config file from %v", value)
		return value
	}
	logrus.Infof("Loading default config file from %v", defaultConfigFile)
	return defaultConfigFile
}

func (c *Config) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(out)
}
