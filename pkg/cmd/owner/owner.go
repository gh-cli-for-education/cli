package owner

import (
	"github.com/MakeNowJust/heredoc"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdOwner(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner",
		Short: "Manage default owner for GitHub CLI commands",
		Long: `The owner command allows you to manage the default owner
for GitHub CLI commands.`,
		Example: heredoc.Doc(`
			$ gh owner set-default
			$ gh owner view
		`),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
				A owner can be supplied as an argument in any of the following formats:
				- "OWNER"
			`),
		},
		GroupID: "core",
	}

	cmdutil.AddGroup(cmd, "General commands",
		// ownerListCmd.NewCmdList(f, nil),
		// ownerInfoCmd.NewCmdInfo(f, nil),
	)

	cmdutil.AddGroup(cmd, "Targeted commands",
	  // ownerSetDefaultCmd.NewCmdSetDefault(f, nil),
		// ownerReposCmd.NewCmdView(f, nil),
	)

	return cmd
}
