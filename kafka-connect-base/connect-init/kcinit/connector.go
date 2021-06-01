package kcinit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ConnectorConfig load and parse connector's configuration
type ConnectorConfig struct {
	Name   string            `json:"name"  validate:"required"`
	Config map[string]string `json:"config"  validate:"required"`
}

// ConnectorEnv map for the connector.properties file
var ConnectorEnv map[string]string

// SetConnector Get Connectors configuration from env
func SetConnector(logger logrus.FieldLogger) error {
	var key string
	connectorEnv := make(map[string]string)
	connectorConfig := &ConnectorConfig{}

	connectorPrefix := "CONNECTOR_"
	for _, b := range os.Environ() {
		if strings.HasPrefix(b, connectorPrefix) {
			pair := strings.SplitN(b, "=", 2)
			key = strings.TrimPrefix(pair[0], connectorPrefix)
			key = strings.ToLower(key)
			key = strings.Replace(key, "_", ".", -1)
			connectorEnv[key] = pair[1]
		}
	}
	connectorConfig.Name = connectorEnv["name"]
	connectorConfig.Config = connectorEnv

	if err := validator.New().Struct(connectorConfig); err != nil {
		return err
	}
	SaveConnectorToFile(
		Cfg.HomeDir+"/connector.json",
		connectorConfig,
	)

	opt := ConfigProps{}
	err := GenerateConfigFile(
		Cfg.HomeDir,
		"connector.properties",
		opt,
		connectorEnv,
	)
	if err != nil {
		return err
	}

	if !Cfg.ConnectorInfo.AutoCreate {
		return nil
	}

	err = CheckIfConnectIsUP()
	if err != nil {
		return err
	}
	ifExists, err := CheckIfConnectorExists(connectorConfig)
	if err != nil {
		return err
	}
	if ifExists != http.StatusOK {
		err = CreateConnector(connectorConfig)
	} else {
		err = UpdateConnector(connectorConfig)
	}
	if err != nil {
		return err
	}

	return nil
}

// Marshal is a function that marshals the object into an
// io.Reader.
// By default, it uses the JSON marshaller.
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// SaveConnectorToFile saves connector config to home_dir/connector.json
func SaveConnectorToFile(path string, v interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	r, err := bytes.NewReader(b), nil
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

// CheckIfConnectIsUP check kafka connect rest port for a 200
func CheckIfConnectIsUP() error {
	connectRestPort := os.Getenv("KAFKA_CONNECT_REST")

	if connectRestPort == "" {
		return errors.New("KAFKA_CONNECT_REST env is empty. Can not proceed")
	}

	req, err := http.NewRequest("GET", connectRestPort, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	for i := 0; i <= Cfg.ConnectorInfo.WaitForConnect; i++ {
		time.Sleep(500 * time.Millisecond)
		resp, err := client.Do(req)

		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			time.Sleep(5000 * time.Millisecond)
			break
		}
	}
	return nil
}

// CheckIfConnectorExists check if we have already created the connector.
// In case it exists, we will update instead of creating.
func CheckIfConnectorExists(c *ConnectorConfig) (int, error) {
	connectRestPort := os.Getenv("KAFKA_CONNECT_REST")

	if connectRestPort == "" {
		return 0, errors.New("KAFKA_CONNECT_REST env is empty. Can not proceed")
	}

	req, err := http.NewRequest(
		"GET",
		connectRestPort+"/connectors/"+c.Name+"/config",
		nil,
	)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

// UpdateConnector updated connectors configuration
func UpdateConnector(c *ConnectorConfig) error {
	connectRestPort := os.Getenv("KAFKA_CONNECT_REST")

	if connectRestPort == "" {
		return errors.New("KAFKA_CONNECT_REST env is empty. Can not proceed")
	}

	b, err := json.Marshal(c.Config)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"PUT",
		connectRestPort+"/connectors/"+c.Name+"/config",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	return nil
}

// CreateConnector checks if connect is up and creates the connectors as defined
// under home_dir/connector.json
func CreateConnector(c *ConnectorConfig) error {
	connectRestPort := os.Getenv("KAFKA_CONNECT_REST")

	if connectRestPort == "" {
		return errors.New("KAFKA_CONNECT_REST env is empty. Can not proceed")
	}

	b, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		connectRestPort+"/connectors",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	return nil
}
