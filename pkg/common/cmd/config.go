package cmd

type DatabaseConfig struct {
	DatabaseDriver    string `envconfig:"database_driver"`
	DatabaseName      string `envconfig:"database_name"`
	DatabaseAddress   string `envconfig:"database_address"`
	DatabaseUser      string `envconfig:"database_user"`
	DatabasePassword  string `envconfig:"database_password"`
	DatabaseArguments string `envconfig:"database_arguments"`
}

type WebConfig struct {
	ServerPort string `envconfig:"port"`
}
