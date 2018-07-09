package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/cmds/options"
	"github.com/spf13/cobra"
)

func NewCmdMount() *cobra.Command {
	cfg := options.NewConfig()
	cmd := &cobra.Command{
		Use:               "mount",
		Short:             "Mount the volume at the mount dir",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				Error(cloud.ErrIncorrectArgNumber).Print()
			}
			cloud, err := cloud.GetCloudManager(cfg.Provider, context.Background())
			if err != nil {
				Error(err).Print()
			}

			dir := args[0]
			opt := cloud.NewOptions()
			if err := json.Unmarshal([]byte(args[1]), opt); err != nil {
				Error(fmt.Errorf("could not parse options for attach; got %v", os.Args[1])).Print()
			}

			if err = cloud.Initialize(); err != nil {
				Error(err).Print()
			}

			if err := cloud.Mount(dir, opt); err != nil {
				Error(err).Print()
			}
			Success().Print()
		},
	}
	return cmd
}
