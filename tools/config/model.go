package config

type Config struct {
	Server  *Server  `yaml:"server" json:"server"`
	Log     *Log     `yaml:"log" json:"log"`
	Sqlite3 *Sqlite3 `yaml:"sqlite3" json:"sqlite3"`
	Mysql   *Mysql   `yaml:"mysql"`
	Rain    *Rain    `yaml:"rain"`
	Model   string   `yaml:"model"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Sqlite3 struct {
	Enable bool `yaml:"enable"`
}

type Mysql struct {
	Port     int    `yaml:"port"`
	Db       string `yaml:"db"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Charset  string `yaml:"charset"`
}

type Rain struct {
	Mode    string `yaml:"mode"`
	WorkDir string `yaml:"workDir"`
	LogDir  string `yaml:"logDir"`
}

type Log struct {
	OutputPaths      []string      `yaml:"outputPaths"`
	EncoderConfig    EncoderConfig `yaml:"encoderConfig"`
	ErrorOutputPaths []string      `yaml:"errorOutputPaths"`

	Compress       bool   `yaml:"compress"`
	Level          string `yaml:"level"`
	Filename       string `yaml:"filename"`
	Maxsize        int    `yaml:"maxsize"`
	MaxAge         int    `yaml:"maxage"`
	Encoding       string `yaml:"encoding"`
	FileMaxBackups int    `yaml:"file_max_backups"`
	Stdout         bool   `yaml:"stdout"`
}
type EncoderConfig struct {
	MessageKey   string `yaml:"messageKey"`
	LevelKey     string `yaml:"levelKey"`
	LevelEncoder string `yaml:"levelEncoder"`
}
