package logger

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/initlove/ocihub/config"
)

func ParseLogLevel(level string) (int, error) {
	switch level {
	case "emergency":
		return logs.LevelEmergency, nil
	case "alert":
		return logs.LevelAlert, nil
	case "critical":
		return logs.LevelCritical, nil
	case "error":
		return logs.LevelError, nil
	case "warn":
		return logs.LevelWarning, nil
	case "notice":
		return logs.LevelNotice, nil
	case "info":
		return logs.LevelInformational, nil
	case "debug":
		return logs.LevelDebug, nil
	}

	return logs.LevelInformational, errors.New("Error in parsing log level.")
}

func ParseAdapter(name string) (string, error) {
	switch name {
	case "":
		return "", errors.New("Log Adapter should not be empty.")
	case "console":
		return logs.AdapterConsole, nil
	case "file":
		return logs.AdapterFile, nil
	case "multifile":
		return logs.AdapterMultiFile, fmt.Errorf("Adapter '%s' is not supported.", name)
	case "smtp":
		return logs.AdapterMail, fmt.Errorf("Adapter '%s' is not supported.", name)
	case "conn":
		return logs.AdapterConn, fmt.Errorf("Adapter '%s' is not supported.", name)
	case "es":
		return logs.AdapterEs, fmt.Errorf("Adapter '%s' is not supported.", name)
	case "jianliao":
		return logs.AdapterJianLiao, fmt.Errorf("Adapter '%s' is not supported.", name)
	case "slack":
		return logs.AdapterSlack, fmt.Errorf("Adapter '%s' is not supported.", name)
	case "alils":
		return logs.AdapterAliLS, fmt.Errorf("Adapter '%s' is not supported.", name)

	}
	return "", fmt.Errorf("Adapter '%s' is invalid.", name)
}

func ParseAdapterArgs(name string, values map[string]interface{}) string {
	// We only turn 'level' to human readable string, so we just need to turn it back.
	v, ok := values["level"]
	if ok {
		if l, err := ParseLogLevel(v.(string)); err == nil {
			values["level"] = l
		} else {
			values["level"] = logs.LevelInformational
		}
	}

	str, _ := json.Marshal(values)
	return string(str)
}

func ParseLog(cfg config.LogConfig) (string, string, error) {
	for n, v := range cfg {
		name, err := ParseAdapter(n)
		if err != nil {
			return "", "", err
		}
		args := ParseAdapterArgs(name, v)
		return name, args, nil
		// If multiple log adapters were detected, use the first one
	}

	return "", "", errors.New("Log is not set.")
}

func InitLogger(cfg config.LogConfig) error {
	n, args, err := ParseLog(cfg)
	if err != nil {
		logs.SetLogger(logs.AdapterConsole, fmt.Sprintf("{\"level\": %d}", logs.LevelInformational))
		return errors.New("Fail to parse logger, fallback to 'console', the debug level 'info'")
	}

	logs.SetLogger(n, args)
	return nil
}
