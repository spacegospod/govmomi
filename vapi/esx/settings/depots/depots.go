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

package depots

import (
	"context"
	"github.com/vmware/govmomi/vapi/rest"
	"net/http"
	"strings"
)

const (
	// DepotsOfflinePath The endpoint for the offline depots API
	DepotsOfflinePath = "/api/esx/settings/depots/offline"
	// DepotContentComponentsPath The endpoint for retrieving the components in a depot
	DepotContentComponentsPath = "/api/esx/settings/depot-content/components"
)

// Manager extends rest.Client, adding vLCM related methods.
type Manager struct {
	*rest.Client
}

// NewManager creates a new Manager instance with the given client.
func NewManager(client *rest.Client) *Manager {
	return &Manager{
		Client: client,
	}
}

// GetOfflineDepot retrieves an offline depot by its identifier
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offline/depot/get/
func (c *Manager) GetOfflineDepot(depotId string) (map[string]interface{}, error) {
	path := c.Resource(DepotsOfflinePath).WithSubpath(depotId)
	req := path.Request(http.MethodGet)
	// TODO create bindings
	var res map[string]interface{}
	return res, c.Do(context.Background(), req, &res)
}

// GetOfflineDepots retrieves all offline depots
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offline/get/
func (c *Manager) GetOfflineDepots() (map[string]interface{}, error) {
	path := c.Resource(DepotsOfflinePath)
	req := path.Request(http.MethodGet)
	// TODO create bindings
	var res map[string]interface{}
	return res, c.Do(context.Background(), req, &res)
}

// DeleteOfflineDepot triggers a task to delete an offline depot by its identifier
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offline/depotvmw-tasktrue/delete/
func (c *Manager) DeleteOfflineDepot(depotId string) (string, error) {
	path := c.Resource(DepotsOfflinePath).WithSubpath(depotId).WithParam("vmw-task", "true")
	req := path.Request(http.MethodDelete)
	// TODO create bindings
	var res string
	return res, c.Do(context.Background(), req, &res)
}

// CreateOfflineDepot triggers a task to create an offline depot
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offlinevmw-tasktrue/post/
func (c *Manager) CreateOfflineDepot(spec map[string]interface{}) (string, error) {
	path := c.Resource(DepotsOfflinePath).WithParam("vmw-task", "true")
	req := path.Request(http.MethodPost, spec)
	// TODO create bindings
	var res string
	return res, c.Do(context.Background(), req, &res)
}

// GetDepotContentComponents retrieves the components in a depot
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offlinevmw-tasktrue/post/
func (c *Manager) GetDepotContentComponents(bundleTypes, names, vendors, versions *[]string, minVersion *string) ([]map[string]interface{}, error) {
	addArrayParam := func(path *rest.Resource, name string, value *[]string) {
		if value != nil && len(*value) > 0 {
			path = path.WithParam(name, strings.Join(*value, ","))
		}
	}
	path := c.Resource(DepotContentComponentsPath)
	addArrayParam(path, "bundle_types", bundleTypes)
	addArrayParam(path, "names", names)
	addArrayParam(path, "vendors", vendors)
	addArrayParam(path, "versions", versions)

	if minVersion != nil {
		path = path.WithParam("min_version", *minVersion)
	}
	req := path.Request(http.MethodGet)
	// TODO create bindings
	var res []map[string]interface{}
	return res, c.Do(context.Background(), req, &res)
}
