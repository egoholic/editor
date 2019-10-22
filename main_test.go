package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	. "github.com/egoholic/editor/config"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/sclevine/agouti"
	. "github.com/stoa-bd/blogging-schema/seed"
)

var (
	logger = log.New(os.Stdout, "editor-test", 0)
	page   *agouti.Page
	err    error
	driver = agouti.ChromeDriver()
	opt    = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "cucumber",
	}
	cmd = exec.Command("go", "run", "targets/web/main.go", "-logpath", "tmp/log/test.log", "-port", strconv.Itoa(Port), "-dbname", "stoa_blogging_test_acceptance", "-pidpath", "tmp/pids/web.pid", "-domain", "localhost")
)

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
	DB, err = sql.Open("postgres", DBConnectionString)
	if err != nil {
		panic(err)
	}
}

func FeatureContext(s *godog.Suite) {
	s.BeforeSuite(func() {
		err = driver.Start()
		if err != nil {
			logger.Printf("Error: Can't run driver: %s\n", err.Error())
		}
		logger.Println("Driver runned!")
		page, err = driver.NewPage(agouti.Browser("firefox"))
		if err != nil {
			logger.Printf("Error: Can't run client: %s\n", err.Error())
		}
		logger.Println("Client runned!")
		err = cmd.Start()
		if err != nil {
			logger.Printf("Error: Can't run blog server: %s\n", err.Error())
		}
	})
	s.AfterSuite(func() {
		err = driver.Stop()
		if err != nil {
			logger.Println("Error: ", err.Error())
		}
		err = stopBlogApp()
		if err != nil {
			logger.Println("Error: ", err.Error())
		}
	})
	s.BeforeScenario(func(interface{}) {
		err = Truncate("blogs", "accounts", "rubrics", "publications", "publication_authors")
		if err != nil {
			logger.Printf("Error: Can't clean up DB: %s\n", err.Error())
		}
		page, err = driver.NewPage(agouti.Browser("firefox"))
		if err != nil {
			logger.Printf("Error: Can't run client: %s\n", err.Error())
		}
	})
	s.AfterScenario(func(_ interface{}, err error) {
		if err != nil {
			url, err2 := page.URL()
			if err2 != nil {
				url = ""
			}
			url = strings.ReplaceAll(url, "/", "--")
			page.Screenshot(fmt.Sprintf("tmp/screenshots/screenshot-%d[%s].png", time.Now().Unix(), url))
			logger.Println("Error: ", err.Error())
		}
		err = page.Destroy()
		if err != nil {
			logger.Println("Error: ", err.Error())
		}
	})

	s.Step(``, func() {})
}

func stopBlogApp() (err error) {
	fmt.Println("stopping blog app")
	file, err := os.Open(PIDFilePath)
	if err != nil {
		return
	}
	pidStr := make([]byte, 16)
	_, err = file.Read(pidStr)
	if err != nil {
		return
	}
	pid, err := strconv.Atoi(strings.TrimRight(string(pidStr), "\x00"))
	if err != nil {
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return
	}
	err = process.Kill()
	if err != nil {
		return
	}
	fmt.Println("pid file dropped")
	return
}
