package logger

type LogConfig struct {
	Level            int    `yaml:"level" mapstructure:"level" validate:"omitempty,gte=-1,lte=5" cobra-usage:"the application log level" cobra-default:"1"`
	TimeFormat       string `yaml:"time-format" mapstructure:"time-format" cobra-usage:"the application log time format" cobra-default:"2006-01-02T15:04:05Z07:00"`
	TimestampEnabled bool   `yaml:"timestamp-enabled" mapstructure:"timestamp-enabled" cobra-usage:"specify if the timestamp is enabled"  cobra-default:"false"`
	ServiceName      string `yaml:"service-name" mapstructure:"service-name" cobra-usage:"the application service name" cobra-default:""`
}
