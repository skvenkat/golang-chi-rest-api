package app

type Config struct {
	Server      ServerConfig
	Cache       CacheConfig
	Deployment  string
	Credentials CredentialsConfig
	Database    DatabaseConfig
}

type CredentialsConfig struct {
	Key    string
	Secret string
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Filename string
}

type CacheConfig struct {
	Type string
}
