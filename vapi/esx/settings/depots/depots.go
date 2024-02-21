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

// SettingsDepotsOfflineSummary is a type mapping for
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/data-structures/Settings/Depots/Offline/Summary/
type SettingsDepotsOfflineSummary struct {
	Description string `json:"description"`
	SourceType  string `json:"source_type"`
	FileId      string `json:"file_id,omitempty"`
	Location    string `json:"location,omitempty"`
	Owner       string `json:"owner,omitempty"`
	OwnerData   string `json:"owner_data,omitempty"`
}

// SettingsDepotsOfflineInfo is a type mapping for
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/data-structures/Settings/Depots/Offline/Info/
type SettingsDepotsOfflineInfo struct {
	CreateTime  string `json:"create_time"`
	Description string `json:"description"`
	SourceType  string `json:"source_type"`
	FileId      string `json:"file_id,omitempty"`
	Location    string `json:"location,omitempty"`
	Owner       string `json:"owner,omitempty"`
	OwnerData   string `json:"owner_data,omitempty"`
}

// SettingsDepotsOfflineCreateSpec is a type mapping for
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/data-structures/Settings/Depots/Offline/CreateSpec/
type SettingsDepotsOfflineCreateSpec struct {
	Description string `json:"description,omitempty"`
	SourceType  string `json:"source_type"`
	FileId      string `json:"file_id,omitempty"`
	Location    string `json:"location,omitempty"`
	OwnerData   string `json:"owner_data,omitempty"`
}

// SettingsDepotContentComponentsSummary is a type mapping for
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/data-structures/Settings/DepotContent/Components/Summary/
type SettingsDepotContentComponentsSummary struct {
	DisplayName string           `json:"display_name"`
	Name        string           `json:"name"`
	Vendor      string           `json:"vendor"`
	Versions    []VersionSummary `json:"versions"`
}

// VersionSummary is a type mapping for
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/data-structures/Settings/DepotContent/Components/ComponentVersionSummary/
type VersionSummary struct {
	Category       string `json:"category"`
	DisplayVersion string `json:"display_version"`
	Kb             string `json:"kb"`
	ReleaseDate    string `json:"release_date"`
	Summary        string `json:"summary"`
	Urgency        string `json:"urgency"`
	Version        string `json:"version"`
}

// GetOfflineDepot retrieves an offline depot by its identifier
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offline/depot/get/
func (c *Manager) GetOfflineDepot(depotId string) (SettingsDepotsOfflineSummary, error) {
	path := c.Resource(DepotsOfflinePath).WithSubpath(depotId)
	req := path.Request(http.MethodGet)
	var res SettingsDepotsOfflineSummary
	return res, c.Do(context.Background(), req, &res)
}

// GetOfflineDepots retrieves all offline depots
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offline/get/
func (c *Manager) GetOfflineDepots() (map[string]SettingsDepotsOfflineInfo, error) {
	path := c.Resource(DepotsOfflinePath)
	req := path.Request(http.MethodGet)
	var res map[string]SettingsDepotsOfflineInfo
	return res, c.Do(context.Background(), req, &res)
}

// DeleteOfflineDepot triggers a task to delete an offline depot by its identifier
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offline/depotvmw-tasktrue/delete/
func (c *Manager) DeleteOfflineDepot(depotId string) (string, error) {
	path := c.Resource(DepotsOfflinePath).WithSubpath(depotId).WithParam("vmw-task", "true")
	req := path.Request(http.MethodDelete)
	var res string
	return res, c.Do(context.Background(), req, &res)
}

// CreateOfflineDepot triggers a task to create an offline depot
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depots/offlinevmw-tasktrue/post/
func (c *Manager) CreateOfflineDepot(spec SettingsDepotsOfflineCreateSpec) (string, error) {
	path := c.Resource(DepotsOfflinePath).WithParam("vmw-task", "true")
	req := path.Request(http.MethodPost, spec)
	var res string
	return res, c.Do(context.Background(), req, &res)
}

// GetDepotContentComponents retrieves the components in a depot
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/depot-content/components/get/
func (c *Manager) GetDepotContentComponents(bundleTypes, names, vendors, versions *[]string, minVersion *string) ([]SettingsDepotContentComponentsSummary, error) {
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
	var res []SettingsDepotContentComponentsSummary
	return res, c.Do(context.Background(), req, &res)
}
