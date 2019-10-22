package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/egoholic/cfg"
	_ "github.com/lib/pq"
)

var (
	DBConnectionString string
	Domain             string
	Port               int
	dbHost             string
	dbPort             int
	dbUser             string
	dbPwd              string
	DBName             string
	config             *cfg.Cfg
	err                error
	LogFile            *os.File
	PIDFilePath        string
)

func init() {
	defaults := map[string]interface{}{}
	defaults["dbhost"] = "localhost"
	defaults["dbport"] = "5432"
	defaults["dbuser"] = "postgres"
	defaults["dbname"] = "stoa_blogging_development"
	defaults["dbpwd"] = ""
	defaults["port"] = "3000"
	defaults["logpath"] = "stdout"
	defaults["pidpath"] = "tmp/pids/web.pid"
	config := cfg.Config(defaults)
	Port, err = config.IntArg("Web server port", "The port which web server listens to, like: 3000.", "port")
	if err != nil {
		panic(err)
	}
	dbHost, err = config.StringArg("DB Host name", "Database connection host name, like: 'localhost'.", "dbhost")
	if err != nil {
		panic(err)
	}
	dbPort, err = config.IntArg("DB port", "Database connection port, like: 5432.", "dbport")
	if err != nil {
		panic(err)
	}
	dbUser, err = config.StringArg("DB User", "Database connection user name, like: 'postgres'.", "dbuser")
	if err != nil {
		panic(err)
	}
	dbPwd, err = config.StringArg("DB Password", "Database connection password.", "dbpwd")
	if err != nil {
		panic(err)
	}
	DBName, err = config.StringArg("Database name", "Database name, like: 'stoa_blogging_development'.", "dbname")
	if err != nil {
		panic(err)
	}
	logPath, err := config.StringArg("Log path", "A path to log file. Use 'stdout' to print to the terminal.", "logpath")
	if err != nil {
		panic(err)
	}
	if logPath == "stdout" {
		LogFile = os.Stdout
	} else {
		LogFile, err = os.Create(logPath)
		if err != nil {
			panic(err)
		}
	}
	PIDFilePath, err = config.StringArg("PID file path", "A path to the PID file.", "pidpath")
	if err != nil {
		panic(err)
	}
	config.AddHelpCommand()
	var sb strings.Builder
	if len(dbHost) > 0 {
		sb.WriteString(fmt.Sprintf("host=%s ", dbHost))
	}
	sb.WriteString(fmt.Sprintf("port=%d ", dbPort))
	if len(dbPwd) > 0 {
		sb.WriteString(fmt.Sprintf("password=%s ", dbPwd))
	}
	if len(dbUser) > 0 {
		sb.WriteString(fmt.Sprintf("user=%s ", dbUser))
	}
	if len(DBName) > 0 {
		sb.WriteString(fmt.Sprintf("dbname=%s ", DBName))
	}

	sb.WriteString("sslmode=disable")
	DBConnectionString = sb.String()
}
