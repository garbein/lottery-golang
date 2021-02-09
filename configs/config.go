package configs

type Config struct {
	DB           DB
	Redis        Redis
	ServerConfig ServerConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type Redis struct {
	Host string
	Port string
}

type DB struct {
	Host     string
	Port     string
	Username string
	Password string
	Dbname   string
	Charset  string
}
