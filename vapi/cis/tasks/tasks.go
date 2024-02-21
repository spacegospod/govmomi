/*
Copyright (c) 2024 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"context"
	"github.com/vmware/govmomi/vapi/rest"
	"net/http"
	"time"
)

const (
	// TasksPath The endpoint for retrieving tasks
	TasksPath = "/api/cis/tasks"
)

// Manager extends rest.Client, adding task related methods.
type Manager struct {
	*rest.Client
}

// NewManager creates a new Manager instance with the given client.
func NewManager(client *rest.Client) *Manager {
	return &Manager{
		Client: client,
	}
}

func (c *Manager) WaitForCompletion(taskId string) (string, error) {
	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			taskInfo, err := c.getTaskInfo(taskId)
			status := taskInfo["status"].(string)
			if err != nil {
				return status, err
			}

			if status != "RUNNING" {
				return status, nil
			}
		}
	}
}

func (c *Manager) getTaskInfo(taskId string) (map[string]interface{}, error) {
	path := c.Resource(TasksPath).WithSubpath(taskId)
	req := path.Request(http.MethodGet)
	// TODO create bindings
	var res map[string]interface{}
	return res, c.Do(context.Background(), req, &res)
}
