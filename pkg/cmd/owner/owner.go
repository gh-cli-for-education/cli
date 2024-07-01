package owner

import (
	"fmt"
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

			if len(args) == 1 {
				opts.Owner = args[0]
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
				owner, err := getDefaultOwner(*opts)
				if err != nil {
					return err
				}

				if owner == "" {
					fmt.Fprintf(opts.IO.Out, "No default owner set\n")
				} else {
					fmt.Fprintf(opts.IO.Out, "Default owner: %s\n", owner)
				}

				return nil
			}

			if opts.Owner != "" {
				// Set default owner
				err := setDefaultOwner(*opts, opts.Owner)
				if err != nil {
					return err
				}
				return nil
			}

			return nil
		},
	}

	return cmd
}

func getDefaultOwner(opts OwnerOptions) (string, error) {
	// Get default owner
	cfg, err := opts.Config()
	if err != nil {
		return "", err
	}

	optValue := cfg.GetOrDefault("", "gh-owner")
	if optValue.IsSome() {
		return optValue.Unwrap().Value, nil
	}

	return "", nil
}

func setDefaultOwner(opts OwnerOptions, owner string) error {
	// Set default owner
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	cfg.Set("", "gh-owner", owner)
	err = cfg.Write()
	if err != nil {
		return err
	}

	return nil
}
