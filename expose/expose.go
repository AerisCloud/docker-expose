package expose

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/gosimple/slug"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/ini.v1"
)

var composeRe = regexp.MustCompile("^([a-z-]+)_([a-z-]+)_([0-9]+)$")

type Expose struct {
	url      string
	email    string
	username string
	token    string
}

type ExposeHost struct {
	Mapping  string `json:"mapping"`
	Hostname string `json:"hostname"`
	Status   int    `json:"status"`
}

type ExposeList []ExposeHost

type Service struct {
	Service string `json:"service"`
	Port    int    `json:"port"`
}

// the current aeriscloud folder
func dataPath() string {
	if runtime.GOOS == "linux" {
		return path.Join(os.Getenv("HOME"), ".config", "AerisCloud")
	} else if runtime.GOOS == "darwin" {
		return path.Join(os.Getenv("HOME"), "Library", "Application Support", "AerisCloud")
	}
	return ""
}

// load the aeris.cd config from aeriscloud
func NewExposeFromConf() (Expose, error) {
	var expose = Expose{}
	// try aeriscloud
	acConfFileName := path.Join(dataPath(), "config.ini")
	acConf, err := ini.Load(acConfFileName)
	if err == nil && acConf.Section("aeris").Key("url").MustString("") != "" {
		expose.url = acConf.Section("aeris").Key("url").MustString("aeris.cd")
		expose.email = acConf.Section("aeris").Key("email").MustString("")
		expose.token = acConf.Section("aeris").Key("token").MustString("")
		expose.username = strings.Split(expose.email, "@")[0]
		return expose, nil
	}

	return expose, errors.New("No valid aeris.cd configuration found")
}

// slugify the container's name and append docker
func (expose Expose) ContainerHost(name string) string {
	if composeRe.MatchString(name) {
		parts := composeRe.FindStringSubmatch(name)
		name = slug.Make(parts[2]) + parts[3] + "." + slug.Make(parts[1])
	} else {
		name = slug.Make(name)
	}
	return fmt.Sprintf("%s.docker", name)
}

// create the full url for the container
func (expose Expose) UserURL(name string) string {
	return fmt.Sprintf("http://%s.%s.%s", expose.ContainerHost(name), expose.username, expose.url)
}

// expose a container on the service
func (expose Expose) Add(name string, port int) error {
	var err error
	var query struct {
		LocalIP  string    `json:"localip"`
		Services []Service `json:"services"`
	}

	query.LocalIP, err = local_ip()
	if err != nil {
		return err
	}

	query.Services = make([]Service, 1)
	query.Services[0].Service = expose.ContainerHost(name)
	query.Services[0].Port = port

	request := gorequest.New()
	res, _, errs := request.Post("http://"+expose.url+"/api/service").Set("Auth-Username", expose.username).Set("Auth-Token", expose.token).Send(query).End()

	if len(errs) > 0 {
		return errs[0]
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Invalid return code from server: %d", res.StatusCode)
	}

	return nil
}

// un-expose (is that a word?) the container from the service
func (expose Expose) Delete(name string) error {
	var query struct {
		Services []Service `json:"services"`
	}

	query.Services = make([]Service, 1)
	query.Services[0].Service = expose.ContainerHost(name)

	request := gorequest.New()
	res, _, errs := request.Delete("http://"+expose.url+"/api/service").Set("Auth-Username", expose.username).Set("Auth-Token", expose.token).Send(query).End()

	if len(errs) > 0 {
		return errs[0]
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Invalid return code from server: %d", res.StatusCode)
	}

	return nil
}

// list currently exposed services
func (expose Expose) List(owned bool) (ExposeList, error) {
	request := gorequest.New().SetBasicAuth(expose.username, expose.token)
	res, body, errs := request.Get("http://" + expose.url + "/api/vms").End()
	if len(errs) > 0 {
		return ExposeList{}, errs[0]
	}

	if res.StatusCode != 200 {
		return ExposeList{}, fmt.Errorf("Invalid return code from server: %d", res.StatusCode)
	}

	el := ExposeList{}
	err := json.Unmarshal([]byte(body), &el)
	if err != nil {
		return ExposeList{}, err
	}

	if owned {
		res := ExposeList{}
		for _, eh := range el {
			components := strings.Split(eh.Hostname, ".")
			user := components[len(components)-1]
			project := components[len(components)-2]
			if user == expose.username && project == "docker" {
				res = append(res, eh)
			}
		}
		return res, nil
	}

	return el, nil
}

// find a service in the list of exposes
func (el ExposeList) Find(name string) (ExposeHost, bool) {
	for _, eh := range el {
		components := strings.Split(eh.Hostname, ".")
		project_name := strings.Join(components[:len(components)-2], ".")
		if project_name == name {
			return eh, true
		}
	}

	return ExposeHost{}, false
}
