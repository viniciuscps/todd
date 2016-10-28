/*
    ToDD comms functions

    This file holds the infrastructure for agent-server communication abstractions in ToDD.

	Copyright 2016 Matt Oswalt. Use or modification of this
	source code is governed by the license provided here:
	https://github.com/toddproject/todd/blob/master/LICENSE
*/

package comms

import (
	"errors"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/toddproject/todd/agent/defs"
	"github.com/toddproject/todd/agent/responses"
	"github.com/toddproject/todd/agent/tasks"
	"github.com/toddproject/todd/config"
)

// CommsPackage will ensure that whatever specific comms struct is loaded at compile time will support
// all of the necessary features/functions that we need to make ToDD work. In short, this interface
// represents a list of things that the server and agents do on the message queue.
type CommsPackage interface {

	// TODO(mierdin) best way to document interface or function args? I've tried to document
	// them minimally below, but would like a better way to document the meaning behind
	// the arguments defined here.

	// (agent advertisement to advertise)
	AdvertiseAgent(defs.AgentAdvert) error

	// (map of assets:hashes, lock for asset map)
	ListenForAgent(*map[string]map[string]string, *sync.Mutex) error

	// (uuid)
	ListenForTasks(string) error

	// (queuename, task)
	SendTask(string, tasks.Task) error

	// watches for new group membership instructions in the cache and reregisters
	WatchForGroup()

	ListenForGroupTasks(string, chan bool) error

	ListenForResponses(*chan bool) error
	SendResponse(responses.Response) error
}

// toddComms is a struct to hold anything that satisfies the CommsPackage interface
type toddComms struct {
	CommsPackage
}

// NewToDDComms will create a new instance of toddComms, and load the desired
// CommsPackage-compatible comms package into it.
func NewToDDComms(cfg config.Config) (*toddComms, error) {

	var tc toddComms

	// Load the appropriate comms package based on config file
	switch cfg.Comms.Plugin {
	case "rabbitmq":
		tc.CommsPackage = newRabbitMQComms(cfg)
	default:
		log.Error("Invalid comms plugin in config file")
		return nil, errors.New("Invalid comms plugin in config file")
	}

	return &tc, nil

}
