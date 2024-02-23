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

package component

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

type lsResult map[string]clusters.SettingsComponentInfo

func (r lsResult) Write(w io.Writer) error {
	for k, v := range r {
		var obj []byte
		var err error
		if obj, err = json.MarshalIndent(v, "", "  "); err != nil {
			return err
		}
		if _, err = fmt.Fprintf(w, "%s -> %s\n", k, string(obj)); err != nil {
			return err
		}
	}
	return nil
}

type ls struct {
	*flags.ClientFlag
	*flags.OutputFlag

	clusterId string
	draftId   string
}

func init() {
	cli.Register("cluster.draft.component.ls", &ls{})
}

func (cmd *ls) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)

	f.StringVar(&cmd.clusterId, "cluster-id", "", "The identifier of the cluster.")
	f.StringVar(&cmd.draftId, "draft-id", "", "The identifier of the software draft.")

}

func (cmd *ls) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.OutputFlag.Process(ctx); err != nil {
		return err
	}

	return nil
}

func (cmd *ls) Usage() string {
	return "CLUSTER"
}

func (cmd *ls) Description() string {
	return `Displays the list of components in a software draft.

Examples:
  govc cluster.draft.component.ls -cluster-id=domain-c21 -draft-id=13`
}

func (cmd *ls) Run(ctx context.Context, f *flag.FlagSet) error {
	rc, err := cmd.RestClient()

	dm := clusters.NewManager(rc)

	var d map[string]clusters.SettingsComponentInfo
	if d, err = dm.ListSoftwareDraftComponents(cmd.clusterId, cmd.draftId); err != nil {
		return err
	}

	return cmd.WriteResult(lsResult(d))
}
