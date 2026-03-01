package types

import (
	"clipshare/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	_ "gopkg.in/yaml.v3"
	"os"
	"time"
)

var AppConfig Config
var DefaultConfig = Config{
	PublicMode: true,
	Web: WebConfig{
		Port:                func(s int) *int { return &s }(80),
		LoginExpiredSeconds: 1800,
		Admin: AdminUserConfig{
			Username: func(s string) *string { return &s }("admin"),
			Password: func(s string) *string { return &s }("1234567"),
		},
	},
	Forward: ForwardConfig{
		Port:             func(s int) *int { return &s }(9283),
		UnlimitedDevices: make([]DeviceBaseInfo, 0),
		FileTransferLimit: FileTransferLimitConfig{
			Enabled: true,
			Rate:    func(s int) *int { return &s }(100),
		},
	},
	Log: &LogConfig{
		MemoryBufferSize: 1000,
	},
}

type OnConfigChanged func(e fsnotify.Event)

func ReadConfig() Config {
	if !utils.FileExists("./data/config.yaml") {
		_ = DefaultConfig.Save("./data/config.yaml")
	}
	viper.AddConfigPath("./data")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	// 读取配置文件
	if err != nil {
		panic(err)
	}
	var config Config
	err = viper.Unmarshal(&config)
	config.CheckValues()
	if err != nil {
		panic(err)
	}
	return config
}
func WatchConfig(onConfigChanged OnConfigChanged) {
	var lastCall time.Time
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		now := time.Now()
		// 设置 100ms 间隔来防止短时间内的重复触发
		if now.Sub(lastCall) <= 100*time.Millisecond {
			return
		}
		lastCall = now
		onConfigChanged(e)
	})
}
func (config *Config) CheckValues() {
	if config.Web.Port == nil {
		panic("Missing `web.port` in config.yaml")
	}
	if config.Web.Admin.Username == nil {
		panic("Missing `username` in config.yaml")
	}
	if config.Web.Admin.Password == nil {
		panic("Missing `password` in config.yaml")
	}
	if config.Forward.Port == nil {
		panic("Missing `forward.port` in config.yaml")
	}
	if config.Forward.FileTransferLimit.Enabled {
		if config.Forward.FileTransferLimit.Rate == nil {
			panic("Missing `file-transfer-limit` in config.yaml")
		}
		if *config.Forward.FileTransferLimit.Rate <= 0 {
			panic("Limit the enabled speed rate to not be less than 0")
		}
	}
	if config.Web.LoginExpiredSeconds <= 0 {
		config.Web.LoginExpiredSeconds = 30 * 60
	}
	if config.Forward.UnlimitedDevices == nil {
		config.Forward.UnlimitedDevices = make([]DeviceBaseInfo, 0)
	}
	if config.Log != nil {
		if config.Log.MemoryBufferSize < 10 {
			panic("MemoryBufferSize not be less than 10")
		}
	} else {
		config.Log = &LogConfig{
			MemoryBufferSize: 1000,
		}
	}
}
func (config *Config) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		_ = file.Close()
		return err
	}
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	err = encoder.Encode(config)
	if err != nil {
		_ = file.Close()
		return err
	}
	return nil
}

func (config *Config) ToDto() ConfigDto {
	return ConfigDto{
		LoginExpiredSeconds:   &config.Web.LoginExpiredSeconds,
		UnlimitedDevices:      &config.Forward.UnlimitedDevices,
		FileTransferEnabled:   &config.Forward.FileTransferLimit.Enabled,
		FileTransferRateLimit: config.Forward.FileTransferLimit.Rate,
		Log:                   config.Log,
		PublicMode:            &config.PublicMode,
	}
}
func (config *Config) GetUnlimitedDeviceIds() []string {
	arr := make([]string, len(config.Forward.UnlimitedDevices))
	for _, dev := range config.Forward.UnlimitedDevices {
		arr = append(arr, dev.Id)
	}
	return arr
}

type Config struct {
	PublicMode bool          `mapstructure:"public-mode" yaml:"public-mode"`
	Web        WebConfig     `mapstructure:"web" json:"web" yaml:"web"`
	Forward    ForwardConfig `mapstructure:"forward" json:"forward" yaml:"forward"`
	Log        *LogConfig    `mapstructure:"log" json:"log" yaml:"log"`
}
type WebConfig struct {
	Port                *int            `mapstructure:"port" yaml:"port"`
	LoginExpiredSeconds int             `mapstructure:"login-expired-seconds" yaml:"login-expired-seconds"`
	Admin               AdminUserConfig `mapstructure:"admin" yaml:"admin"`
}
type AdminUserConfig struct {
	Username *string `mapstructure:"username" yaml:"username"`
	Password *string `mapstructure:"password" yaml:"password"`
}
type ForwardConfig struct {
	Port              *int                    `mapstructure:"port" yaml:"port"`
	UnlimitedDevices  []DeviceBaseInfo        `mapstructure:"unlimited-devices" yaml:"unlimited-devices"`
	FileTransferLimit FileTransferLimitConfig `mapstructure:"file-transfer-limit" yaml:"file-transfer-limit"`
}
type FileTransferLimitConfig struct {
	Enabled bool `mapstructure:"enabled" yaml:"enabled"`
	Rate    *int `mapstructure:"rate" yaml:"rate"`
}
type DeviceBaseInfo struct {
	Id   string `mapstructure:"id" yaml:"id" json:"id" binding:"required"`
	Name string `mapstructure:"name" yaml:"name" json:"name" binding:"required"`
	Desc string `mapstructure:"desc" yaml:"desc" json:"desc"`
}
type LogConfig struct {
	MemoryBufferSize int `mapstructure:"memory-buffer-size" yaml:"memory-buffer-size" json:"memoryBufferSize"`
}
