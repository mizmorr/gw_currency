package config

var defaults = []option{
	{
		name:        "logger.level",
		typing:      "string",
		value:       "debug",
		description: "Level of logging",
	},
	{
		name:        "logger.pathFile",
		typing:      "string",
		value:       "./logs/app.log",
		description: "Path to the log file",
	},
	{
		name:        "listen.httpHost",
		typing:      "string",
		value:       "localhost",
		description: "Server host",
	},
	{
		name:        "listen.httpPort",
		typing:      "string",
		value:       "8080",
		description: "Server port",
	},
	{
		name:        "listen.shutdowntimeout",
		typing:      "duration",
		value:       "5s",
		description: "Timeout for graceful shutdown",
	},

	{
		name:        "storage.postgres.URL",
		typing:      "string",
		value:       "postgres://postgres:post@localhost:6432/exchange?sslmode=disable",
		description: "Postgres database URL",
	},
	{
		name:        "storage.postgres.Timeout",
		typing:      "duration",
		value:       "2s",
		description: "Timeout for database connection",
	},
	{
		name:        "storage.postgres.ConnectAttempts",
		typing:      "int",
		value:       10,
		description: "Number of database connection attempts",
	},
	{
		name:        "storage.postgres.MaxIdleTime",
		typing:      "duration",
		value:       "5m",
		description: "Maximum idle time for database connections",
	},
	{
		name:        "storage.postgres.MaxOpenConns",
		typing:      "int",
		value:       100,
		description: "Maximum number of open database connections",
	},
	{
		name:        "storage.postgres.HealthCheckPeriod",
		typing:      "duration",
		value:       "5m",
		description: "Period for database health check",
	},
	{
		name:        "worker.keepAliveTimeout",
		typing:      "duration",
		value:       "5s",
		description: "Timeout for worker keep-alive",
	},
	{
		name:        "storage.redis.host",
		typing:      "string",
		value:       "localhost",
		description: "Redis server host",
	},
	{
		name:        "storage.redis.port",
		typing:      "string",
		value:       "6379",
		description: "Redis server port",
	},
	{
		name:        "storage.redis.password",
		typing:      "string",
		value:       "mor",
		description: "Redis server password",
	},
	{
		name:        "storage.redis.db",
		typing:      "int",
		value:       0,
		description: "Redis database index",
	},
	{
		name:        "storage.redis.ttl",
		typing:      "duration",
		value:       "10m",
		description: "Time-to-live for Redis keys",
	},
	{
		name:        "jwttokens.refreshSecret",
		typing:      "string",
		value:       "1g3VxO3ILNcu2kxIw5166lSIb+y2iVqPuB+hHll1b1tx9QPtaQiF/+PdkYrbG+DMabrSowIk1WzAtaLkQPUcMQ==",
		description: "Secret for refreshing the access token",
	},
	{
		name:        "jwttokens.accessSecret",
		typing:      "string",
		value:       "WNlxDQGlnzsRpK8saejc39A9K7mVVRRyPM3ZDWY7E+z/ANvaDLMhp8dLktvjgLE3l8UZLfjtK1Q047G0gY0zDw==",
		description: "Secret for accessing the API",
	},
	{
		name:        "jwttokens.accessExpiresTime",
		typing:      "duration",
		value:       "1h",
		description: "Expiration time for the access token",
	},
	{
		name:        "jwttokens.refreshExpiresTime",
		typing:      "duration",
		value:       "24h",
		description: "Expiration time for the refresh token",
	},
	{
		name:        "grpc.host",
		typing:      "string",
		value:       "localhost",
		description: "gRPC server host",
	},
	{
		name:        "grpc.port",
		typing:      "string",
		value:       "50051",
		description: "gRPC server port",
	},
	{
		name:        "rates.currencycodes",
		typing:      "slice",
		value:       []string{"USD", "RUB", "EUR"},
		description: "List of supported currencies",
	},
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}
