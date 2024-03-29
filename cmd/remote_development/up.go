package remote_development

import (
	"bunnyshell.com/cli/pkg/config"
	"bunnyshell.com/cli/pkg/environment"
	remoteDevPkg "bunnyshell.com/cli/pkg/remote_development"
	remoteDevMutagenConfig "bunnyshell.com/dev/pkg/mutagen/config"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

type SyncMode enumflag.Flag

const (
	None SyncMode = iota
	TwoWaySafe
	TwoWayResolved
	OneWaySafe
	OneWayReplica
)

var SyncModeToMutagenMode = map[SyncMode]remoteDevMutagenConfig.Mode{
	None:           remoteDevMutagenConfig.None,
	TwoWaySafe:     remoteDevMutagenConfig.TwoWaySafe,
	TwoWayResolved: remoteDevMutagenConfig.TwoWayResolved,
	OneWaySafe:     remoteDevMutagenConfig.OneWaySafe,
	OneWayReplica:  remoteDevMutagenConfig.OneWayReplica,
}

var SyncModeIds = map[SyncMode][]string{
	None:           {string(remoteDevMutagenConfig.None)},
	TwoWaySafe:     {string(remoteDevMutagenConfig.TwoWaySafe)},
	TwoWayResolved: {string(remoteDevMutagenConfig.TwoWayResolved)},
	OneWaySafe:     {string(remoteDevMutagenConfig.OneWaySafe)},
	OneWayReplica:  {string(remoteDevMutagenConfig.OneWayReplica)},
}

func init() {
	options := config.GetOptions()
	settings := config.GetSettings()

	var (
		syncMode       SyncMode = TwoWayResolved
		localSyncPath  string
		remoteSyncPath string
		resourcePath   string

		portMappings []string

		waitTimeout int64
		noTTY       bool
	)

	command := &cobra.Command{
		Use: "up",

		ValidArgsFunction: cobra.NoFileCompletions,

		RunE: func(cmd *cobra.Command, args []string) error {
			remoteDevelopment := remoteDevPkg.NewRemoteDevelopment()

			if localSyncPath != "" {
				remoteDevelopment.WithLocalSyncPath(localSyncPath)
			}

			if remoteSyncPath != "" {
				remoteDevelopment.WithRemoteSyncPath(remoteSyncPath)
			}

			if len(portMappings) > 0 {
				remoteDevelopment.WithPortMappings(portMappings)
			}

			environmentResource, err := environment.NewFromWizard(&settings.Profile.Context, resourcePath)
			if err != nil {
				return err
			}

			remoteDevelopment.
				WithEnvironmentResource(environmentResource).
				WithWaitTimeout(waitTimeout).
				WithSyncMode(SyncModeToMutagenMode[syncMode])

			// init
			if err = remoteDevelopment.Up(); err != nil {
				return err
			}

			// start
			if !noTTY {
				if err = remoteDevelopment.StartSSHTerminal(); err != nil {
					return err
				}
			}

			return remoteDevelopment.Wait()
		},
	}

	flags := command.Flags()

	flags.AddFlag(options.Organization.GetFlag("organization"))
	flags.AddFlag(options.Project.GetFlag("project"))
	flags.AddFlag(options.Environment.GetFlag("environment"))
	flags.AddFlag(options.ServiceComponent.GetFlag("component"))

	flags.StringVarP(&localSyncPath, "local-sync-path", "l", "", "Local folder path to sync")
	flags.StringVarP(&remoteSyncPath, "remote-sync-path", "r", "", "Remote folder path to sync")
	flags.StringVarP(&resourcePath, "resource", "s", "", "The cluster resource to use (namespace/kind/name format).")
	flags.StringSliceVarP(
		&portMappings,
		"port-forward",
		"f",
		portMappings,
		"Port forward: '8080>3000'\nReverse port forward: '9003<9003'\nComma separated: '8080>3000,9003<9003'",
	)

	flags.BoolVar(&noTTY, "no-tty", false, "Start remote development with no SSH terminal")
	flags.Int64VarP(&waitTimeout, "wait-timeout", "w", 120, "Time to wait for the pod to be ready")

	flags.Var(
		enumflag.New(&syncMode, "sync-mode", SyncModeIds, enumflag.EnumCaseSensitive),
		"sync-mode",
		"Mutagen sync mode.\nAvailable sync modes: none, two-way-safe, two-way-resolved, one-way-safe, one-way-replica.\n\"none\" sync mode disables mutagen.",
	)

	mainCmd.AddCommand(command)
}
