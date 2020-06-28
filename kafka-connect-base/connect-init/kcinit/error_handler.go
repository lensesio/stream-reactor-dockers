package kcinit

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LogInfo struct {
	Key   string
	Value string
}

// LogError handle log error and return back the original error
func (c *LogInfo) LogError(message string, err error) error {
	if err != nil {
		err = errors.Wrap(err, message)
		logrus.WithField(c.Key, c.Value).Error(err.Error())
		return err
	}

	logrus.WithField(c.Key, c.Value).Error(message)
	return errors.New(message)
}

// GetFunctionName for getting the name of a function
func GetFunctionName(i interface{}) string {
	funcNameFull := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	funcNameArray := strings.Split(funcNameFull, "/")
	funcName := funcNameArray[len(funcNameArray)-1]
	return funcName
}

// GetCallerInfo for getting info about the caller function, name, line
func GetCallerInfo() string {
	_, fn, line, _ := runtime.Caller(1)
	callerNameArray := strings.Split(fn, "/")
	callerName := callerNameArray[len(callerNameArray)-1]
	callerLine := fmt.Sprintf("%s:%d", callerName, line)
	return callerLine
}
