package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/cjburchell/uatu-go"

	dockerCompose "github.com/cjburchell/docker-compose"
	"github.com/pkg/errors"
)

type server struct {
	name       string
	configFile string
	endpoints  map[string]*endpoint
	log        log.ILog
}

// IServer interface
type IServer interface {
	Endpoint(name string, method string, path string) IEndpoint
}

type request struct {
	Path        string `json:"path"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
	Endpoint    string `json:"endpoint"`
}

// CreateServer creates the server
func CreateServer(name string, log log.ILog) IServer {
	return &server{name: name, log: log}
}

func (server *server) Write(p []byte) (n int, err error) {
	requestStart := "Request:{"

	lines := strings.Split(string(p), string([]byte{27}))

	for _, data := range lines {
		if strings.Contains(data, requestStart) {
			index := strings.Index(data, requestStart)
			pos := index + len(requestStart) - 1

			server.log.Debugf("%d", pos)
			requestJSONString := data[pos:]
			server.log.Debugf("%s", requestJSONString)

			var requestObject = request{}
			err := json.Unmarshal([]byte(requestJSONString), &requestObject)
			if err != nil {
				server.log.Error(err)
				return len(p), nil
			}

			for _, endpoint := range server.endpoints {
				if endpoint.config.Name == requestObject.Endpoint {
					if endpoint.customRequestHandler != nil {
						endpoint.customRequestHandler(requestObject.Path, requestObject.ContentType, requestObject.Body)
					}
				}
			}
		}
	}

	return len(p), nil
}

func (server *server) BuildComposeService(path string) (dockerCompose.Service, error) {
	var filename, err = server.SaveMockFile(path)
	return dockerCompose.Service{
		Image:   "cjburchell/loki:latest",
		Volumes: []string{fmt.Sprintf("./%s:/mock/%s", filename, filename)},
		Environment: []string{
			fmt.Sprintf("ConfigFile=/mock/%s", filename),
			fmt.Sprintf("ServerName=%s", server.name),
		},
	}, err
}

func (server *server) Stop() error {
	return os.Remove(server.configFile)
}

func (server *server) SaveMockFile(path string) (string, error) {
	var fileObject = map[string]endpointConfig{}

	for name, endpoint := range server.endpoints {
		fileObject[name] = endpoint.config
	}

	configJSON, err := json.Marshal(&fileObject)
	if err != nil {
		return "", errors.WithStack(err)
	}

	filename := fmt.Sprintf("%s_mock_test.json", server.name)
	server.configFile = fmt.Sprintf("%s/%s", path, filename)

	return filename, ioutil.WriteFile(server.configFile, configJSON, 0644)
}

func (server *server) Endpoint(name string, method string, path string) IEndpoint {
	server.log.Printf("%s: Registering Endpoint %s %s", name, method, path)
	newEndpoint := &endpoint{config: endpointConfig{Name: name, Path: path, Method: method}}
	if server.endpoints == nil {
		server.endpoints = map[string]*endpoint{}
	}

	server.endpoints[name] = newEndpoint
	return newEndpoint
}

func (server *server) AttachToLogs(containers dockerCompose.IContainers) error {
	return containers.LogServiceWithHandler(server.name, server)
}
