package owner

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type OwnerOptions struct {
	Config     func() (gh.Config, error)
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Owner string
}

func NewCmdOwner(f *cmdutil.Factory) *cobra.Command {
	opts := &OwnerOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "owner [OWNER]",
		Short: "Manage default owner for GitHub CLI commands",
		Long: `The owner command allows you to manage the default owner
			for GitHub CLI commands.`,
		Example: heredoc.Doc(`
			$ gh owner
			$ gh owner GITHUB_USERNAME
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return cmdutil.FlagErrorf("accepts at most 1 arg(s), received %d", len(args))
			}

			return nil
		},
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
				A owner can be supplied as an argument in any of the following formats:
				- "OWNER"
			`),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Owner == "" {
				// List default owner
				// TODO: Implement list default owner
				return nil
			}

			if opts.Owner != "" {
				// Set default owner
				// TODO: Implement set default owner
				return nil
			}

			return nil
		},
	}

	return cmd
}
