package owner

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/internal/text"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type iprompter interface {
	Select(string, string, []string) (int, error)
}

type OwnerOptions struct {
	Config     func() (gh.Config, error)
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)
	Prompter   iprompter

	Owner       string
	List        bool
	ListFilter  string
	SelectOwner bool
}

func NewCmdOwner(f *cmdutil.Factory) *cobra.Command {
	opts := &OwnerOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		HttpClient: f.HttpClient,
		Prompter:   f.Prompter,
	}

	cmd := &cobra.Command{
		Use:   "owner [OWNER] |",
		Short: "Manage default owner for GitHub CLI commands",
		Long: `The owner command allows you to manage the default owner
			for GitHub CLI commands.`,
		Example: heredoc.Doc(`
			$ gh owner
			$ gh owner GITHUB_USERNAME
			$ gh owner --list
			$ gh owner --select
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return cmdutil.FlagErrorf("accepts at most 1 arg(s), received %d", len(args))
			}

			if len(args) == 1 && (opts.List || opts.SelectOwner) {
				return cmdutil.FlagErrorf("cannot use OWNER argument with --list flag")
			}

			if opts.SelectOwner && opts.List {
				return cmdutil.FlagErrorf("cannot use --select and --list flags together")
			}

			if len(args) == 1 && !opts.List && !opts.SelectOwner {
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
			ownersList, err := getOwners(opts)
			if err != nil {
				return err
			}

			if opts.List {
				// List organizations
				err = listRun(opts, ownersList)
				if err != nil {
					return err
				}

				return nil
			}

			if opts.SelectOwner {
				// Select default owner
				opts.Owner, err = selectOwnerPrompt(opts.Prompter, ownersList.User, ownersList.Organizations)
				if err != nil {
					return err
				}
			}

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
				err := setDefaultOwner(*opts, ownersList)
				if err != nil {
					return err
				}
				return nil
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.List, "list", "l", false, "List organizations")
	cmd.Flags().BoolVarP(&opts.SelectOwner, "select", "s", false, "Interactively select a default owner")

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

func setDefaultOwner(opts OwnerOptions, ownerList *OrganizationList) error {
	// Set default owner
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	// Check if owner is in the list of organizations
	found := false
	for _, org := range ownerList.Organizations {
		if org.Login == opts.Owner {
			found = true
			break
		}
	}

	if !found {
		fmt.Fprintf(opts.IO.Out, "Owner %s not found\n", opts.Owner)
	} else {
		cfg.Set("", "gh-owner", opts.Owner)
		err = cfg.Write()
		if err != nil {
			return err
		}
		fmt.Fprintf(opts.IO.Out, "Default owner set to %s\n", opts.Owner)
	}

	return nil
}

func getOwners(opts *OwnerOptions) (*OrganizationList, error) {
	httpClient, err := opts.HttpClient()
	if err != nil {
		return nil, err
	}

	cfg, err := opts.Config()
	if err != nil {
		return nil, err
	}

	host, _ := cfg.Authentication().DefaultHost()

	ownersList, err := listOrgs(httpClient, host, 100)
	if err != nil {
		return nil, err
	}

	return ownersList, nil
}

func listRun(opts *OwnerOptions, ownersList *OrganizationList) error {
	if err := opts.IO.StartPager(); err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "error starting pager: %v\n", err)
	}
	defer opts.IO.StopPager()

	if opts.IO.IsStdoutTTY() {
		header := listHeader(ownersList.User, len(ownersList.Organizations), ownersList.TotalCount)
		fmt.Fprintf(opts.IO.Out, "\n%s\n\n", header)
	}

	for _, org := range ownersList.Organizations {
		fmt.Fprintln(opts.IO.Out, org.Login)
	}

	return nil
}

func listHeader(user string, resultCount, totalCount int) string {
	if totalCount == 0 {
		return "There are no organizations associated with @" + user
	}

	return fmt.Sprintf("Showing %d of %s", resultCount, text.Pluralize(totalCount, "organization"))
}

func selectOwnerPrompt(prompter iprompter, user string, orgs []Organization) (string, error) {
	selectedOwner, err := prompter.Select("Select a default owner", user, organizationsToList(orgs))
	if err != nil {
		return "", err
	}
	return orgs[selectedOwner].Login, nil
}

func organizationsToList(orgs []Organization) []string {
	var list []string
	for _, org := range orgs {
		list = append(list, org.Login)
	}
	return list
}