package kcinit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// ServiceInfo for managing kafka connect
type ServiceInfo struct {
	PID       chan int
	Terminate chan int
	StdOut    chan string
	StdErr    chan error
	Command   chan string
	Runtime   map[string]string
	Logger    logrus.FieldLogger
}

// RequestResponse send service response
type RequestResponse struct {
	Status  int    `json:"status,omitempty"`
	Content string `json:"content"`
	PID     int    `json:"pid,omitempty"`
}

// Service Global ServiceInfo variable
var Service ServiceInfo

// SetupService setting up service
func (s *ServiceInfo) SetupService(logger logrus.FieldLogger) (*ServiceInfo, error) {
	setLog := LogInfo{
		Key:   "Component",
		Value: "Kafka Connect Service",
	}
	// --------- SETTUP HTTP SERVICE ---------
	logger.WithFields(logrus.Fields{
		"Component": "Kafka Connect Service",
		"Stage":     "ServiceManager",
	}).Info("Configuring ServiceManager")

	cmdPid, cmdTerm, cmdErr, cmdOut, cmdCMD := make(chan int),
		make(chan int),
		make(chan error),
		make(chan string),
		make(chan string)

	Service.PID = cmdPid
	Service.Terminate = cmdTerm
	Service.StdOut = cmdOut
	Service.StdErr = cmdErr
	Service.Command = cmdCMD
	Service.Logger = logger

	osEnviron := os.Environ()
	kafkaEnv := make(map[string]string)
	kafkaPrefix := "KAFKA_"

	for _, b := range osEnviron {
		if strings.HasPrefix(b, kafkaPrefix) {
			pair := strings.SplitN(b, "=", 2)
			kafkaEnv[pair[0]] = pair[1]
		}
	}

	if kafkaEnv["KAFKA_JMX_OPTS"] != "" {
		kakfaJMXOpts := "-Dcom.sun.management.jmxremote=true"
		kakfaJMXOpts = kakfaJMXOpts + " -Dcom.sun.management.jmxremote.authenticate=false"
		kakfaJMXOpts = kakfaJMXOpts + " -Dcom.sun.management.jmxremote.ssl=false "
		os.Setenv(
			"KAFKA_JMX_OPTS",
			kakfaJMXOpts,
		)
	}

	if kafkaEnv["KAFKA_JMX_PORT"] != "" {
		os.Setenv(
			"JMX_PORT",
			kafkaEnv["KAFKA_JMX_PORT"],
		)
		jmxOpts := os.ExpandEnv("$KAFKA_JMX_OPTS")
		jmxOpts = jmxOpts + " -Djava.rmi.server.hostname=" + os.ExpandEnv("$KAFKA_JMX_HOSTNAME")
		jmxOpts = jmxOpts + " -Dcom.sun.management.jmxremote.local.only=false"
		jmxOpts = jmxOpts + " -Dcom.sun.management.jmxremote.rmi.port=" + os.ExpandEnv("$JMX_PORT")
		jmxOpts = jmxOpts + " -Dcom.sun.management.jmxremote.por" + os.ExpandEnv("$JMX_PORT")
		os.Setenv(
			"KAFKA_JMX_OPTS",
			jmxOpts,
		)
	}

	go s.ManageService(cmdPid, cmdTerm, cmdOut, cmdCMD, cmdErr)
	Service.Command <- "Start"
	err := <-Service.StdErr
	if err != nil {
		return nil, setLog.LogError(
			fmt.Sprintf(
				"%s -- %s :: %s",
				GetFunctionName(s.SetupService),
				GetCallerInfo(),
				GetFunctionName(s.ManageService),
			),
			err,
		)
	}
	return &Service, nil
}

// ManageService function for executing and managing kafka connect process
func (s *ServiceInfo) ManageService(cmdPid, cmdTerm chan int, cmdOut, cmdCMD chan string, cmdErr chan error) {
	classpath := os.Getenv("CLASSPATH")

	cmd := exec.Command(
		"export CLASSPATH="+classpath+";connect-distributed",
		"/etc/kafka-connect/kafka-connect.properties",
	)

	outLogFile, err := os.OpenFile(
		"/kafka-connect.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	defer outLogFile.Close()

	cmd.Stdout = outLogFile
	cmd.Stderr = outLogFile

	for {
		select {
		case status := <-cmdCMD:
			if status == "Start" {
				cmd = exec.Command(
					"connect-distributed",
					"/etc/kafka-connect/kafka-connect.properties",
				)
				cmd.Stdout = outLogFile
				cmd.Stderr = outLogFile
				err = cmd.Start()
				cmdErr <- err
				continue
			}

			if status == "Status" {
				cmdPid <- cmd.Process.Pid
				continue
			}

		case sig := <-cmdTerm:
			if sig == 1 {
				err := cmd.Process.Signal(syscall.SIGTERM)
				cmdErr <- err
				err = cmd.Wait()
				cmdOut <- err.Error()
				continue
			}
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

//requestErrorHandler for HTTP Error Requests on Lenses Portal Endpoint
func requestErrorHandler(funcName string, msg string, err error, httpErr int, w http.ResponseWriter) {
	w.WriteHeader(httpErr)

	payload := &RequestResponse{
		Content: fmt.Sprintf("%s - %s", strconv.Itoa(httpErr), msg),
	}

	// Serialize to bytes the payload
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(
			w,
			fmt.Sprintf(
				"%s - %s",
				strconv.Itoa(httpErr),
				msg,
			),
		)
		return
	}

	// Send payload response back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
	logrus.WithError(err).Error(msg)
}

// StatusListener check connect distrubuted status
func StatusListener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		requestErrorHandler(
			"StatusListener",
			"Method not allowed",
			fmt.Errorf("Only get method is allowed for /api/status"),
			http.StatusMethodNotAllowed,
			w,
		)
		return
	}

	payload := RequestResponse{}

	Service.Command <- "Status"
	p := <-Service.PID

	_, err := os.FindProcess(int(p))
	if err != nil {
		requestErrorHandler(
			"StatusListener",
			"Could not find kafka connect process",
			err,
			http.StatusInternalServerError,
			w,
		)
		return
	}

	if _, err := os.Stat("/proc/" + strconv.Itoa(p) + "/exe"); err == nil {
		payload.Content = "Kafka connect is running"
		payload.PID = p
	} else {
		payload.Content = "Kafka connect is not running"
	}

	// Serialize to bytes the payload
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		requestErrorHandler(
			"StatusListener",
			"Failed to create json response",
			err,
			http.StatusInternalServerError,
			w,
		)
		return
	}

	// Send payload response back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// StartServiceListener check connect distrubuted status
func StartServiceListener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		requestErrorHandler(
			"StartServiceListener",
			"Method not allowed",
			fmt.Errorf("Only get method is allowed for /api/start"),
			http.StatusMethodNotAllowed,
			w,
		)
		return
	}

	payload := RequestResponse{}

	Service.Command <- "Status"
	p := <-Service.PID
	_, err := os.FindProcess(int(p))
	if err == nil {
		requestErrorHandler(
			"StopServiceListener",
			"Kafka Connect is already running",
			err,
			http.StatusOK,
			w,
		)
		return
	}

	Service.Command <- "Start"
	err = <-Service.StdErr
	if err != nil {
		requestErrorHandler(
			"StatusListener",
			"Could not start kafka connect",
			err,
			http.StatusInternalServerError,
			w,
		)
		return
	}
	payload.Content = "Kafka connect started!"

	Service.Command <- "Status"
	p = <-Service.PID
	_, err = os.FindProcess(int(p))
	if err != nil {
		requestErrorHandler(
			"StopServiceListener",
			"Could not locate kafka connect process after startup",
			err,
			http.StatusOK,
			w,
		)
		return
	}
	payload.PID = p

	// Serialize to bytes the payload
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		requestErrorHandler(
			"StatusListener",
			"Failed to create json response",
			err,
			http.StatusInternalServerError,
			w,
		)
		return
	}

	// Send payload response back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// StopServiceListener terminate connect distrubuted process
func StopServiceListener(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		requestErrorHandler(
			"StopServiceListener",
			"Method not allowed",
			fmt.Errorf("Only post method is allowed for /api/stop"),
			http.StatusMethodNotAllowed,
			w,
		)
		return
	}

	payload := RequestResponse{}

	Service.Command <- "Status"
	p := <-Service.PID
	_, err := os.FindProcess(int(p))
	if err != nil {
		requestErrorHandler(
			"StopServiceListener",
			"Could not find kafka connect process",
			err,
			http.StatusOK,
			w,
		)
		return
	}

	Service.Terminate <- 1
	err = <-Service.StdErr
	if err != nil {
		requestErrorHandler(
			"StopServiceListener",
			"Could not terminate Kafka Connect service",
			err,
			http.StatusInternalServerError,
			w,
		)
		return
	}

	waitOutput := <-Service.StdOut
	payload.Content = "Send SIGTERM to Kafka Connect process: " + waitOutput

	// Serialize to bytes the payload
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		requestErrorHandler(
			"StopServiceListener",
			"Failed to create json response",
			err,
			http.StatusInternalServerError,
			w,
		)
		return
	}

	// Send payload response back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
