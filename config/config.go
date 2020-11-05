package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AwesomeHomeConfig struct {
	ConfYaml
	defaultConf []byte
}

type ConfYaml struct {
	BroadcastPort          int
	BroadcastPacketMaxSize int
	BroadcastSendTime      time.Duration

	HttpServerPort int
}

func NewAwesomeHomeConfig(serviceName string, path string) *AwesomeHomeConfig {
	conf := &AwesomeHomeConfig{defaultConf: []byte(`
broadcast:
  packet_max_size: 10240
  port: 60504
  send_time: 5 seconds
http:
  server:
    port: 60503
`)}
	err := conf.loadConf(serviceName, path)
	if err != nil {
		log.Fatalf("Load yaml config file error: '%v'", err)
		return nil
	}
	return conf
}

// LoadConf load config from file and read in environment variables that match
func (conf *AwesomeHomeConfig) loadConf(serviceName string, path string) error {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix(serviceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if path != "" {
		content, err := ioutil.ReadFile(path)

		if err != nil {
			log.Errorf("File does not exist : %s", path)
			return err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return err
		}
	} else {
		// Search config in home directory with name ".pkg" (without extension).
		viper.AddConfigPath("..")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./volumes/conf/")
		viper.SetConfigName("config")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(conf.defaultConf)); err != nil {
				return err
			}
		}
	}

	conf.BroadcastPort = viper.GetInt("broadcast.port")
	conf.BroadcastPacketMaxSize = viper.GetInt("broadcast.packet_max_size")
	conf.BroadcastSendTime = viper.GetDuration("broadcast.send_time")
	conf.HttpServerPort = viper.GetInt("http.server.port")
	return nil
}
