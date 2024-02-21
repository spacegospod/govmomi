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

package clusters

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi/vapi/rest"
	"net/http"
	"strings"
)

const (
	// SoftwareDraftsPath The endpoint for the software drafts API
	SoftwareDraftsPath     = "/api/esx/settings/clusters/%s/software/drafts"
	SoftwareComponentsPath = SoftwareDraftsPath + "/%s/software/components"
)

// Manager extends rest.Client, adding Software Drafts related methods.
type Manager struct {
	*rest.Client
}

// NewManager creates a new Manager instance with the given client.
func NewManager(client *rest.Client) *Manager {
	return &Manager{
		Client: client,
	}
}

// ListSoftwareDrafts retrieves the software drafts for a cluster
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/clusters/cluster/software/drafts/get/
func (c *Manager) ListSoftwareDrafts(clusterId string, owners *[]string) (map[string]interface{}, error) {
	path := c.Resource(fmt.Sprintf(SoftwareDraftsPath, clusterId))

	if owners != nil && len(*owners) > 0 {
		path = path.WithParam("owners", strings.Join(*owners, ","))
	}

	req := path.Request(http.MethodGet)
	// TODO create bindings
	var res map[string]interface{}
	return res, c.Do(context.Background(), req, &res)
}

// CreateSoftwareDraft creates a software draft on the provided cluster
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/clusters/cluster/software/drafts/post/
func (c *Manager) CreateSoftwareDraft(clusterId string) (string, error) {
	path := c.Resource(fmt.Sprintf(SoftwareDraftsPath, clusterId))
	req := path.Request(http.MethodPost)
	// TODO create bindings
	var res string
	return res, c.Do(context.Background(), req, &res)
}

// DeleteSoftwareDraft removes the specified draft
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/clusters/cluster/software/drafts/draft/delete/
func (c *Manager) DeleteSoftwareDraft(clusterId, draftId string) error {
	path := c.Resource(fmt.Sprintf(SoftwareDraftsPath, clusterId)).WithSubpath(draftId)
	req := path.Request(http.MethodDelete)
	return c.Do(context.Background(), req, nil)
}

// GetSoftwareDraft returns the set of components in the specified draft
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/clusters/cluster/software/drafts/draft/get/
func (c *Manager) GetSoftwareDraft(clusterId, draftId string) (map[string]interface{}, error) {
	path := c.Resource(fmt.Sprintf(SoftwareDraftsPath, clusterId)).WithSubpath(draftId)
	req := path.Request(http.MethodGet)
	// TODO create bindings
	var res map[string]interface{}
	return res, c.Do(context.Background(), req, &res)
}

// GetSoftwareDraft returns the set of components in the specified draft
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/clusters/cluster/software/drafts/draft/get/
func (c *Manager) CommitSoftwareDraft(clusterId, draftId string, spec map[string]interface{}) (string, error) {
	path := c.Resource(fmt.Sprintf(SoftwareDraftsPath, clusterId)).WithSubpath(draftId).WithParam("action", "commit").WithParam("vmw-task", "true")
	req := path.Request(http.MethodPost, spec)
	var res string
	return res, c.Do(context.Background(), req, &res)
}

// UpdateSoftwareDraftComponents updates the set of components in the specified draft
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/clusters/cluster/software/drafts/draft/software/components/patch/
func (c *Manager) UpdateSoftwareDraftComponents(clusterId, draftId string, spec map[string]interface{}) error {
	path := c.Resource(fmt.Sprintf(SoftwareComponentsPath, clusterId, draftId))
	req := path.Request(http.MethodPatch, spec)
	return c.Do(context.Background(), req, nil)
}

// RemoveSoftwareDraftComponents removes a component from the specified draft
// https://developer.vmware.com/apis/vsphere-automation/latest/esx/api/esx/settings/clusters/cluster/software/drafts/draft/software/components/component/delete/
func (c *Manager) RemoveSoftwareDraftComponents(clusterId, draftId, component string) error {
	path := c.Resource(fmt.Sprintf(SoftwareComponentsPath, clusterId, draftId)).WithSubpath(component)
	req := path.Request(http.MethodDelete)
	return c.Do(context.Background(), req, nil)
}
