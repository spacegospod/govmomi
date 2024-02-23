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

package draft

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/vapi/esx/settings/clusters"
	"io"
)

type infoResult clusters.SettingsClustersSoftwareDraftsInfo

func (r infoResult) Write(w io.Writer) error {
	var obj []byte
	var err error
	if obj, err = json.MarshalIndent(r, "", "  "); err != nil {
		return err
	}
	if _, err = fmt.Fprintln(w, string(obj)); err != nil {
		return err
	}
	return nil
}

type info struct {
	*flags.ClientFlag
	*flags.OutputFlag

	clusterId string
	draftId   string
}

func init() {
	cli.Register("cluster.draft.info", &info{})
}

func (cmd *info) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)
	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)

	f.StringVar(&cmd.clusterId, "cluster-id", "", "The identifier of the cluster.")
	f.StringVar(&cmd.draftId, "draft-id", "", "The identifier of the software draft.")
}

func (cmd *info) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.OutputFlag.Process(ctx); err != nil {
		return err
	}

	return nil
}

func (cmd *info) Usage() string {
	return "CLUSTER"
}

func (cmd *info) Description() string {
	return `Displays the details of a software draft.

Examples:
  govc cluster.draft.info -cluster-id=domain-c21 -draft-id=13`
}

func (cmd *info) Run(ctx context.Context, f *flag.FlagSet) error {
	rc, err := cmd.RestClient()

	dm := clusters.NewManager(rc)

	var d clusters.SettingsClustersSoftwareDraftsInfo
	if d, err = dm.GetSoftwareDraft(cmd.clusterId, cmd.draftId); err != nil {
		return err
	}

	return cmd.WriteResult(infoResult(d))
}
