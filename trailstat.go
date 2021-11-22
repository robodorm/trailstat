package main

import (
	"encoding/xml"
	"fmt"
	"github.com/concourse/concourse/atc/api/ccserver"
	"github.com/concourse/concourse/fly/rc"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/robodorm/tripoli"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var config = make(map[string]string)
var result []string
var client = &http.Client{Timeout: 10 * time.Second}

func httpGet(addr string) (string, error) {
	var r string
	rq, err := http.NewRequest(http.MethodGet, addr, nil)

	if err != nil {
		return "", err
	}

	rq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config[addr]))
	rq.Header.Set("Content-Type", "application/xml")

	ret, err := client.Do(rq)
	if err != nil {
		return "", err
	}

	defer ret.Body.Close()
	var projects []ccserver.Project

	channel := ccserver.ProjectsContainer{Projects: projects}

	if err = xml.NewDecoder(ret.Body).Decode(&channel); err != nil {
		return "", err
	} else if len(channel.Projects) != 0 {
		for _, item := range channel.Projects {
			r += fmt.Sprintf("Status: %s, Name: %s \n", item.LastBuildStatus, item.Name)
		}
	}

	return r, nil
}

func getTeamStatusesLinks() error {
	targets, err := rc.LoadTargets()

	if err != nil {
		return err
	}

	format := "%s/api/v1/teams/%s/cc.xml"

	for _, c := range targets {
		config[fmt.Sprintf(format, c.API, c.TeamName)] = fmt.Sprintf("%s", c.Token.Value)
	}

	return nil
}

func getConcourseData() ([]interface{}, error) {
	var data []interface{}

	for k := range config {
		data = append(data, k)
	}

	return tripoli.Run(httpGet, runtime.NumCPU(), data), nil
}

func main() {
	for {
		err := getTeamStatusesLinks()

		if err != nil {
			return
		}

		r2, _ := getConcourseData()

		for _, res := range r2 {
			v := fmt.Sprintf("%s", reflect.ValueOf(res).Index(0))
			for _, s := range strings.Split(v, "\n") {
				result = append(result, s)
			}
		}

		systray.Run(onReady, onExit)
		time.Sleep(1000)
	}
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Pipelines")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.SetIcon(icon.Data)

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	for _, i := range result {
		systray.AddMenuItem(i, "")
	}
}

func onExit() {}
