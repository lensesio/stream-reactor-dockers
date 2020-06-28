package kcinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SchemaRegistry for schema registry url
type SchemaRegistry struct {
	KeyConverterSchemaRegistryURL   string
	ValueConverterSchemaRegistryURL string
}

// ConnectInfo struct of the mandatory connect otions
type ConnectInfo struct {
	BootStrapServers       string `validate:"required"`
	GroupID                string `validate:"required"`
	ConfigStorageTopic     string `validate:"required"`
	OffsetStorageTopic     string `validate:"required"`
	StatusStorageTopic     string `validate:"required"`
	KeyConverter           string `validate:"required"`
	ValueConverter         string `validate:"required"`
	RestAdvertisedHostName string `validate:"required"`
	InternalKeyConverter   string
	InternalValueConverter string
	RestPort               string //`validate:"required"`
	SchemaRegistry         SchemaRegistry
	ClassPath              string `validate:"required"`
	PluginPath             string `validate:"required"`
}

// ConnectEnv map for the connector.properties file
var ConnectEnv map[string]string

// Connect global var for connectInfo
var Connect ConnectInfo

// ConfigProps for populating x.properties
type ConfigProps struct {
	Name  string
	Value string
}

// GenerateConfigFile for parsing env and populating connect.properties
func GenerateConfigFile(path, fileName string, opt ConfigProps, connectEnv map[string]string) error {
	const connectPropsTempl = `{{.Name}}={{.Value}}
`
	t := template.Must(template.New("ConfigProps").Parse(connectPropsTempl))
	outLogFile, err := os.Create(path + "/" + fileName)
	if err != nil {
		return err
	}
	defer outLogFile.Close()

	for n, v := range connectEnv {
		opt.Name = n
		opt.Value = v

		err := t.Execute(outLogFile, opt)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetupKafkaConnect Get Connectors configuration from env
func SetupKafkaConnect(logger logrus.FieldLogger) (*map[string]string, error) {
	setLog := LogInfo{
		Key:   "Component",
		Value: "Kafka Connect",
	}
	// --------- LOADING ENV ---------
	logger.WithFields(logrus.Fields{
		"Component": "Kafka Connect",
		"Stage":     "Environment",
	}).Info("Loading and parsing environment")

	var key string
	connectInfo := &ConnectInfo{}
	connectEnv := make(map[string]string)

	osEnviron := os.Environ()
	connectPrefix := "CONNECT_"

	for _, b := range osEnviron {
		if strings.HasPrefix(b, connectPrefix) {
			pair := strings.SplitN(b, "=", 2)
			key = strings.TrimPrefix(pair[0], connectPrefix)
			key = strings.ToLower(key)
			key = strings.Replace(key, "_", ".", -1)
			connectEnv[key] = pair[1]
		}
	}

	connectInfo.BootStrapServers = connectEnv["bootstrap.servers"]
	connectInfo.GroupID = connectEnv["group.id"]
	connectInfo.ConfigStorageTopic = connectEnv["config.storage.topic"]
	connectInfo.OffsetStorageTopic = connectEnv["offset.storage.topic"]
	connectInfo.StatusStorageTopic = connectEnv["status.storage.topic"]
	connectInfo.KeyConverter = connectEnv["key.converter"]
	connectInfo.ValueConverter = connectEnv["value.converter"]
	connectInfo.RestPort = connectEnv["rest.port"]
	connectInfo.RestAdvertisedHostName = connectEnv["rest.advertised.host.name"]
	connectInfo.InternalKeyConverter = connectEnv["internal.key.converter"]
	connectInfo.InternalValueConverter = connectEnv["internal.value.converter"]

	if connectInfo.RestPort == "" {
		connectInfo.RestPort = "8083"
	}

	if connectInfo.KeyConverter == "io.confluent.connect.avro.AvroConverter" {
		key, ok := connectEnv["key.converter.schema.registry.url"]
		if !ok {
			return nil, setLog.LogError(
				fmt.Sprintf(
					"%s -- %s :: %s",
					GetFunctionName(SetupKafkaConnect),
					GetCallerInfo(),
					GetFunctionName(connectEnv),
				),
				errors.New("Key converts was set to avro but no key schema registry url was given"),
			)
		}
		value, ok := connectEnv["key.converter.schema.registry.url"]
		if !ok {
			return nil, setLog.LogError(
				fmt.Sprintf(
					"%s -- %s :: %s",
					GetFunctionName(SetupKafkaConnect),
					GetCallerInfo(),
					GetFunctionName(connectEnv),
				),
				errors.New("Value converts was set to avro but no value schema registry url was given"),
			)
		}
		connectInfo.SchemaRegistry.KeyConverterSchemaRegistryURL = key
		connectInfo.SchemaRegistry.ValueConverterSchemaRegistryURL = value
	}
	interJSON := "org.apache.kafka.connect.json.JsonConverter"
	if connectInfo.InternalKeyConverter == interJSON || connectInfo.InternalValueConverter == interJSON {
		os.Setenv("CONNECT_INTERNAL_KEY_CONVERTER_SCHEMAS_ENABLE", "false")
		os.Setenv("CONNECT_INTERNAL_VALUE_CONVERTER_SCHEMAS_ENABLE", "false")
	}

	pluginPath := connectEnv["plugin.path"]

	if pluginPath == "" {
		os.Setenv("CONNECT_PLUGIN_PATH", Cfg.KafkaConnectorsDir)
	}
	connectInfo.PluginPath = pluginPath

	newClassPathh := Cfg.PluginsDir + "/*" +
		":/etc/kafka-connect/jars/*" + ":/opt/calcite/*"

	classpath := os.Getenv("CLASSPATH")

	if classpath == "" {
		os.Setenv(
			"CLASSPATH",
			newClassPathh,
		)
	} else {
		os.Setenv(
			"CLASSPATH",
			os.ExpandEnv("${CLASSPATH}")+":"+newClassPathh,
		)
	}

	classpath = os.Getenv("CLASSPATH")
	connectInfo.ClassPath = classpath

	if err := validator.New().Struct(connectInfo); err != nil {
		return nil, setLog.LogError(
			fmt.Sprintf(
				"%s -- %s :: %s",
				GetFunctionName(SetupKafkaConnect),
				GetCallerInfo(),
				GetFunctionName(validator.New),
			),
			err,
		)
	}

	// --------- END OF LOADING ENV ---------
	// --------- Configuring System Items ---------
	logger.WithFields(logrus.Fields{
		"Component": "Kafka Connect",
		"Stage":     "Preparing the system",
	}).Info("Setting up plugin.path, creating symlinks and creating conigs")

	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		os.MkdirAll(pluginPath, os.ModePerm)
	}
	archiveFile := os.ExpandEnv("${ARCHIVE}")
	connectJar := strings.TrimSuffix(archiveFile, ".tar.gz") + ".jar"

	if _, err := os.Stat(pluginPath + "/" + connectJar); err != nil {
		source := Cfg.LibDir + "/" + connectJar
		target := pluginPath + "/" + connectJar

		err = os.Symlink(source, target)
		if err != nil {
			return nil, setLog.LogError(
				fmt.Sprintf(
					"%s -- %s :: %s",
					GetFunctionName(SetupKafkaConnect),
					GetCallerInfo(),
					GetFunctionName(os.Symlink),
				),
				err,
			)
		}
	}

	files, err := ioutil.ReadDir(Cfg.PluginsDir)
	if err != nil {
		return nil, setLog.LogError(
			fmt.Sprintf(
				"%s -- %s :: %s",
				GetFunctionName(SetupKafkaConnect),
				GetCallerInfo(),
				GetFunctionName(ioutil.ReadDir),
			),
			err,
		)
	}
	for _, f := range files {
		if _, err := os.Stat(pluginPath + "/" + f.Name()); os.IsNotExist(err) {
			continue
		}

		source := Cfg.PluginsDir + "/" + f.Name()
		target := pluginPath + "/" + f.Name()

		err = os.Symlink(source, target)
		if err != nil {
			return nil, setLog.LogError(
				fmt.Sprintf(
					"%s -- %s :: %s",
					GetFunctionName(SetupKafkaConnect),
					GetCallerInfo(),
					GetFunctionName(os.Symlink),
				),
				err,
			)
		}

	}

	kafkaConnectRest := os.Getenv("KAFKA_CONNECT_REST")
	if kafkaConnectRest == "" {
		kafkaConnectRest = "http://" +
			os.ExpandEnv("${CONNECT_REST_ADVERTISED_HOST_NAME}") +
			":" + os.ExpandEnv("${CONNECT_REST_PORT}")
		os.Setenv("KAFKA_CONNECT_REST", kafkaConnectRest)
	}

	opt := ConfigProps{}
	err = GenerateConfigFile(
		Cfg.KafkaConnectDir,
		"kafka-connect.properties",
		opt,
		connectEnv,
	)
	if err != nil {
		return nil, setLog.LogError(
			fmt.Sprintf(
				"%s -- %s :: %s",
				GetFunctionName(SetupKafkaConnect),
				GetCallerInfo(),
				GetFunctionName(GenerateConfigFile),
			),
			err,
		)
	}

	// --------- End Configuring System Items ---------
	logger.WithFields(logrus.Fields{
		"Component": "Kafka Connect",
		"Stage":     "-",
	}).Info("Successfully configured kafka connect")

	return &connectEnv, nil
}
