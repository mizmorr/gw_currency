package config

var Defaults = []option{
	{
		Name:        "loglevel",
		Typing:      "string",
		Value:       "info",
		Description: "Level of logging",
	},
	{
		Name:        "listen.host",
		Typing:      "string",
		Value:       "localhost",
		Description: "Server host",
	},
	{
		Name:        "listen.port",
		Typing:      "int",
		Value:       8080,
		Description: "Server port",
	},

	{
		Name:        "storage.postgresURL",
		Typing:      "string",
		Value:       "postgres://postgres:post@localhost:6432/exchange?sslmode=disable",
		Description: "Postgres database URL",
	},
	{
		Name:        "storage.postgresTimeout",
		Typing:      "duration",
		Value:       "2s",
		Description: "Timeout for database connection",
	},
	{
		Name:        "storage.postgresConnectAttempts",
		Typing:      "int",
		Value:       10,
		Description: "Number of database connection attempts",
	},
	{
		Name:        "storage.postgresMaxIdleTime",
		Typing:      "duration",
		Value:       "5m",
		Description: "Maximum idle time for database connections",
	},
	{
		Name:        "storage.postgresMaxOpenConns",
		Typing:      "int",
		Value:       100,
		Description: "Maximum number of open database connections",
	},
	{
		Name:        "storage.postgresHealthCheckPeriod",
		Typing:      "duration",
		Value:       "5m",
		Description: "Period for database health check",
	},
}

type option struct {
	Name        string
	Typing      string
	Value       interface{}
	Description string
}
