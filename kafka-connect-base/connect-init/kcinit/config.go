package kcinit

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Livenes for enabling livness thread
type Livenes struct {
	Enabled bool
	Port    string
}

// ConnectorInfo for enabling connector autocreate
type ConnectorInfo struct {
	AutoCreate     bool `mapstructure:"auto_create"`
	WaitForConnect int  `mapstructure:"wait_for_connect"  validate:"required"`
}

// Config describes the available configuration
// of the running service
type Config struct {
	Livenes            Livenes
	Environment        string
	Debug              bool
	KafkaConnectorsDir string        `mapstructure:"kafka_connectors_dir"  validate:"required"`
	KafkaConnectDir    string        `mapstructure:"kafka_connect_dir"  validate:"required"`
	HomeDir            string        `mapstructure:"home_dir"  validate:"required"`
	PluginsDir         string        `mapstructure:"plugins_dir"  validate:"required"`
	LibDir             string        `mapstructure:"lib_dir"  validate:"required"`
	ConnectorInfo      ConnectorInfo `mapstructure:"connector_info"  validate:"required"`
	Logger             logrus.FieldLogger
}

// Cfg global configuration across the whole
// services
var Cfg Config

// Set the file name of the configurations file
func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("streamreactor")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	defaults := map[string]interface{}{
		"kafka_connectors_dir":            "/connectors",
		"kafka_connect_dir":               "/etc/kafka-connect",
		"home_dir":                        "/opt/lenses",
		"plugins_dir":                     "/opt/lenses/kafka-connect/plugins",
		"lib_dir":                         "/opt/lenses/lib",
		"connector_info.auto_create":      true,
		"connector_info.wait_for_connect": 60,
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

// Validate makes sure that the config makes sense
func (c *Config) Validate() error {
	// home
	if err := validator.New().Struct(c.ConnectorInfo); err != nil {
		return err
	}
	if err := validator.New().Struct(c); err != nil {
		return err
	}
	return nil
}

// LoadConfig checks file and environment variables
func LoadConfig(logger logrus.FieldLogger) error {
	setLog := LogInfo{
		Key:   "Component",
		Value: "Init",
	}
	logger.WithFields(logrus.Fields{
		"Component": "Init",
		"Stage":     "Loead Configs",
	}).Info("Load process environment")
	err := viper.ReadInConfig()
	// if err != nil {
	// 	return setLog.LogError(
	// 		fmt.Sprintf(
	// 			"%s -- %s :: %s",
	// 			GetFunctionName(LoadConfig),
	// 			GetCallerInfo(),
	// 			GetFunctionName(viper.ReadInConfig),
	// 		),
	// 		err,
	// 	)
	// }
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return setLog.LogError(
			fmt.Sprintf(
				"%s -- %s :: %s",
				GetFunctionName(LoadConfig),
				GetCallerInfo(),
				GetFunctionName(viper.Unmarshal),
			),
			err,
		)
	}
	err = Cfg.Validate()
	if err != nil {
		return setLog.LogError(
			fmt.Sprintf(
				"%s -- %s :: %s",
				GetFunctionName(LoadConfig),
				GetCallerInfo(),
				GetFunctionName(Cfg.Validate),
			),
			err,
		)
	}
	return nil
}
