package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	. "github.com/egoholic/editor/config"
	"github.com/egoholic/editor/lib/pwd"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/DATA-DOG/godog/gherkin"
	. "github.com/stoa-bd/blogging-schema/factory"
)

type ResponseSaver struct {
	Response *http.Response
}

var (
	logger = log.New(os.Stdout, "editor-test", 0)
	domain = "localhost"
	err    error
	opt    = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "cucumber",
	}
	cmd    = exec.Command("go", "run", "targets/web/main.go", "-logpath", "tmp/log/test.log", "-port", strconv.Itoa(Port), "-dbname", "stoa_blogging_test_acceptance", "-pidpath", "tmp/pids/web.pid", "-domain", domain)
	client = &http.Client{}
)

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
	DB, err = sql.Open("postgres", DBConnectionString)
	if err != nil {
		panic(err)
	}
}
func thereIsAClientCompany(company *gherkin.DataTable) error {
	cnames := company.Rows[0].Cells
	cvals := company.Rows[1].Cells
	mapping := map[string]*Value{}
	for i, cname := range cnames {
		mapping[cname.Value] = V(cvals[i].Value)
	}
	return Insert(Must(NewClientCompany(mapping)))
}
func theCompanyHasTheFollowingEditors(cid string, editors *gherkin.DataTable) error {
	var (
		cnames   = editors.Rows[0].Cells
		accounts = []*Tuple{}
		acc      *Tuple
		password string
		login    string
	)
	for _, row := range editors.Rows[1:] {
		cvals := row.Cells
		mapping := map[string]*Value{}
		for i, cname := range cnames {
			if cname.Value == "password" {
				password = cvals[i].Value
				continue
			}
			if cname.Value == "login" {
				login = cvals[i].Value
			}
			mapping[cname.Value] = V(cvals[i].Value)
		}
		encryptedPassword, err := pwd.Encrypt([]byte(password), []byte(login))
		if err != nil {
			return err
		}
		mapping["encrypted_password"] = &Value{
			Literal: string(encryptedPassword),
			Type:    "bytea",
		}
		mapping["cid"] = V(cid)
		acc, err = NewAccount(mapping)
		if err != nil {
			return err
		}
		accounts = append(accounts, acc)
	}
	return InsertList(accounts)
}
func (rs *ResponseSaver) iSubmitedTheFollowingCredentials(creds *gherkin.DataTable) error {
	rs.Response = nil
	payload := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    "",
		Password: "",
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	rs.Response, err = client.Post(
		fmt.Sprintf("http://%s:%d", domain, Port),
		"application/json",
		bytes.NewReader(jsonBytes),
	)
	if err != nil {
		return err
	}
	expectedCode := 200
	if rs.Response.StatusCode != expectedCode {
		return fmt.Errorf("Wrong response code. Expected: %d, got: %d", expectedCode, rs.Response.StatusCode)
	}
	return nil
}
func (rs *ResponseSaver) iGotAnAccessToken() error {
	if rs.Response.Header.Get("Access-Token") == "" {
		return fmt.Errorf("Expected Access-Token header to not be empty.\n%#v\n\n", map[string][]string(rs.Response.Header))
	}
	return nil
}
func iGotResponse(arg1 string) error {
	return godog.ErrPending
}
func FeatureContext(s *godog.Suite) {
	responseSaver := &ResponseSaver{}

	s.BeforeSuite(func() {
		err = cmd.Start()
		if err != nil {
			logger.Printf("Error: Can't run blog server: %s\n", err.Error())
		}
	})
	s.AfterSuite(func() {
		err = stopBlogApp()
		if err != nil {
			logger.Println("Error: ", err.Error())
		}
	})
	s.BeforeScenario(func(interface{}) {
		err = Truncate("client_companies", "blogs", "accounts", "signins", "rubrics", "publications", "publication_authors")
		if err != nil {
			logger.Printf("Error: Can't clean up DB: %s\n", err.Error())
		}
	})

	s.Step(`^there is a client company$`, thereIsAClientCompany)
	s.Step(`^the company "([^"]*)" has the following editors:$`, theCompanyHasTheFollowingEditors)
	s.Step(`^I submited the following credentials:$`, responseSaver.iSubmitedTheFollowingCredentials)
	s.Step(`^I got an access token$`, responseSaver.iGotAnAccessToken)
	s.Step(`^I got "([^"]*)" response$`, iGotResponse)
}
func stopBlogApp() (err error) {
	fmt.Println("stopping editor app")
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
