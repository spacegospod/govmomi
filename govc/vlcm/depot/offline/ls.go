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

package offline

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/vapi/esx/settings/depots"
	"io"
)

type lsResult map[string]depots.SettingsDepotsOfflineInfo

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
}

func init() {
	cli.Register("vlcm.depot.offline.ls", &ls{})
}

func (cmd *ls) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
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
	return "VLCM"
}

func (cmd *ls) Description() string {
	return `Displays the list of offline image depots. 

Examples:
  govc vlcm.depot.offline.ls`
}

func (cmd *ls) Run(ctx context.Context, f *flag.FlagSet) error {
	rc, err := cmd.RestClient()

	dm := depots.NewManager(rc)

	var d map[string]depots.SettingsDepotsOfflineInfo
	if d, err = dm.GetOfflineDepots(); err != nil {
		return err
	}

	return cmd.WriteResult(lsResult(d))
}
