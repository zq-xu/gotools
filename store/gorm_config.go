package store

import (
	"os"
	"strconv"

	"github.com/spf13/pflag"
	"gorm.io/gorm/logger"

	"zq-xu/helper/config"
	"zq-xu/helper/log"
)

const (
	databaseConfigName = "DataBaseConfig"

	databaseAddressEnv     = "DatabaseAddress"
	defaultDatabaseAddress = "localhost"

	databasePortEnv     = "DatabasePort"
	defaultDatabasePort = "3306"

	databaseNameEnv     = "DatabaseName"
	defaultDatabaseName = "warehouse_admin"

	databaseUsernameEnv     = "DatabaseUsername"
	defaultDatabaseUsername = "root"

	databasePasswordEnv     = "DatabasePassword"
	defaultDatabasePassword = "root"

	databaseLogLevelEnv     = "DatabaseLogLevel"
	defaultDatabaseLogLevel = logger.Warn
)

var (
	DatabaseCfg = &DatabaseConfig{}
)

type DatabaseConfig struct {
	DatabaseInfo
	LogLevelStr string

	LogLevel logger.LogLevel
}

func InitDataBaseConfig() {
	config.RegisterCfg(databaseConfigName, DatabaseCfg)
}

func (dbc *DatabaseConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&dbc.Address, "database-address", os.Getenv(databaseAddressEnv), "the database ip")
	fs.StringVar(&dbc.Port, "database-port", os.Getenv(databasePortEnv), "the database port")
	fs.StringVar(&dbc.Username, "database-username", os.Getenv(databaseUsernameEnv), "the database username")
	fs.StringVar(&dbc.Password, "database-password", os.Getenv(databasePasswordEnv), "the database password")
	fs.StringVar(&dbc.DatabaseName, "database-name", os.Getenv(databaseNameEnv), "the database name")
	fs.StringVar(&dbc.LogLevelStr, "database-log-level", os.Getenv(databaseLogLevelEnv), "the database log level")
}

func (dbc *DatabaseConfig) Revise() {
	dbc.reviseAddress()
	dbc.revisePort()
	dbc.reviseUsername()
	dbc.revisePassword()
	dbc.reviseDatabaseName()
	dbc.reviseDatabaseLogLevel()

	log.Logger.Infof("DataBaseConfig is DatabaseName: %v; LogLevel: %v", dbc.DatabaseName, dbc.LogLevel)
}

func (dbc *DatabaseConfig) reviseAddress() {
	if dbc.Address == "" {
		dbc.Address = defaultDatabaseAddress
	}
}

func (dbc *DatabaseConfig) revisePort() {
	if dbc.Port == "" {
		dbc.Port = defaultDatabasePort
	}
}

func (dbc *DatabaseConfig) reviseUsername() {
	if dbc.Username == "" {
		dbc.Username = defaultDatabaseUsername
	}
}

func (dbc *DatabaseConfig) revisePassword() {
	if dbc.Password == "" {
		dbc.Password = defaultDatabasePassword
	}
}

func (dbc *DatabaseConfig) reviseDatabaseName() {
	if dbc.DatabaseName == "" {
		dbc.DatabaseName = defaultDatabaseName
	}
}

func (dbc *DatabaseConfig) reviseDatabaseLogLevel() {
	if dbc.LogLevelStr == "" {
		dbc.LogLevel = defaultDatabaseLogLevel
		return
	}

	lInt, err := strconv.Atoi(dbc.LogLevelStr)
	if err != nil {
		dbc.LogLevel = defaultDatabaseLogLevel
		return
	}

	dbc.LogLevel = logger.LogLevel(lInt)
}
