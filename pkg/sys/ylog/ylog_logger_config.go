package ylog

import "time"

const (
	LogTypeConsole = "console"
	LogTypeJson    = "json"
)

// Config is the configuration object for logger.
type Config struct {
	TimeFormat string       `json:"timeFormat" yaml:"time_format"` // Logging time format.
	Encoding   string       `json:"encoding" yaml:"encoding"`      // Logging type: console, json.
	Path       string       `json:"path" yaml:"path"`              // Logging directory path.
	Filename   string       `json:"filename" yaml:"filename"`      // Filename is the file to write logs to.
	Level      string       `json:"level" yaml:"level"`            // Output level.
	Prefix     string       `json:"prefix" yaml:"prefix"`          // Prefix string for every logging content.
	CtxKeys    []string     `json:"ctxKeys" yaml:"ctx_keys"`       // Context keys for logging, which is used for value retrieving from context.
	Stdout     bool         `json:"stdout" yaml:"stdout"`          // Output to stdout or not(true in default).
	Rotate     ConfigRotate `json:"rotate" yaml:"rotate"`
}

type ConfigRotate struct {
	MaxSize    int  `json:"maxSize" yaml:"max_size"`       // MaxSize is the maximum size in megabytes of the log file before it gets rotated. It defaults to 100 megabytes.
	MaxAge     int  `json:"maxAge" yaml:"max_age"`         // MaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
	MaxBackups int  `json:"maxBackups" yaml:"max_backups"` // MaxBackups is the maximum number of old log files to retain.
	LocalTime  bool `json:"localTime" yaml:"local_time"`   //
	Compress   bool `json:"compress" yaml:"compress"`      // Compress determines if the rotated log files should be compressed using gzip. The default is not to perform compression.
}

// DefaultConfig returns the default configuration for logger.
func DefaultConfig() Config {
	c := Config{
		TimeFormat: time.DateTime + ".000000",
		Encoding:   LogTypeConsole,
		Path:       "./storage/logs",
		Filename:   "log.log",
		Level:      "debug",
		Prefix:     "YAFGO",
		CtxKeys:    []string{},
		Stdout:     true,
		Rotate: ConfigRotate{
			MaxSize:    4,
			MaxAge:     7,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   false,
		},
	}

	return c
}
