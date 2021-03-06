package mock

import (
	"fmt"
	"os"
	"time"

	dockerCompose "github.com/cjburchell/docker-compose"
	log "github.com/cjburchell/uatu-go"
)

type system struct {
	dockerCompose.IContainers
	mockedServers   []*server
	composeFilePath string
	log             log.ILog
}

// ISystem interface
type ISystem interface {
	End()
}

func (system *system) End() {
	err := system.Down()
	if err != nil {
		system.log.Error(err)
	}

	err = os.Remove(system.composeFilePath)
	if err != nil {
		system.log.Error(err)
	}

	for _, server := range system.mockedServers {
		err = server.stop()
		if err != nil {
			system.log.Error(err)
		}
	}
}

// StartSystem starts the mock system
func StartSystem(path, composeFile string, services map[string]dockerCompose.Service, mocked []IServer, log log.ILog, systemLogs bool) (ISystem, error) {
	file := dockerCompose.File{
		Version:  "2.2",
		Services: services,
	}

	var mockedServers []*server
	for _, item := range mocked {
		mockedServer, ok := item.(*server)
		if ok {
			mockedServers = append(mockedServers, mockedServer)
		}
	}

	var composeFilePath = fmt.Sprintf("%s/%s", path, composeFile)

	var err error
	for _, mockedServer := range mockedServers {
		file.Services[mockedServer.name], err = mockedServer.buildComposeService(path)
		if err != nil {
			return nil, err
		}
	}

	if err := file.SaveFile(composeFilePath); err != nil {
		return nil, err
	}

	compose := dockerCompose.CreateFile(composeFilePath)

	log.Printf("Starting up %s", composeFilePath)
	if err := compose.Up(); err != nil {
		err := os.Remove(composeFilePath)
		if err != nil {
			log.Error(err)
		}

		return nil, err
	}

	if systemLogs {
		// Attach to Docker images loggers
		log.Printf("Connecting to logging up %s", composeFilePath)
		for key := range services {
			if err := compose.LogService(key); err != nil {
				err := compose.Down()
				if err != nil {
					log.Error(err)
				}

				err = os.Remove(composeFilePath)
				if err != nil {
					log.Error(err)
				}

				return nil, err
			}
		}
	}

	for _, servers := range mockedServers {
		err := servers.attachToLogs(compose)
		if err != nil {
			log.Error(err)
		}
	}

	<-time.After(10 * time.Second) // give some time for the servers to start up

	return &system{compose, mockedServers, composeFilePath, log}, nil
}
