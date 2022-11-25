package codespace

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newCodeCmd(app *App) *cobra.Command {
	var (
		filterOptions getOrChooseCodespaceFilterOptions
		useInsiders   bool
		useWeb        bool
	)

	codeCmd := &cobra.Command{
		Use:   "code",
		Short: "Open a codespace in Visual Studio Code",
		Args:  noArgsConstraint,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return validateGetOrChooseCodespaceCommandArgs(filterOptions)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.VSCode(cmd.Context(), filterOptions, useInsiders, useWeb)
		},
	}

	addGetOrChooseCodespaceCommandArgs(codeCmd, &filterOptions)

	codeCmd.Flags().BoolVar(&useInsiders, "insiders", false, "Use the insiders version of Visual Studio Code")
	codeCmd.Flags().BoolVarP(&useWeb, "web", "w", false, "Use the web version of Visual Studio Code")

	return codeCmd
}

// VSCode opens a codespace in the local VS VSCode application.
func (a *App) VSCode(ctx context.Context, filterOptions getOrChooseCodespaceFilterOptions, useInsiders bool, useWeb bool) error {
	codespace, err := getOrChooseCodespace(ctx, a.apiClient, filterOptions)
	if err != nil {
		return err
	}

	browseURL := vscodeProtocolURL(codespace.Name, useInsiders)
	if useWeb {
		browseURL = codespace.WebURL
		if useInsiders {
			u, err := url.Parse(browseURL)
			if err != nil {
				return err
			}
			q := u.Query()
			q.Set("vscodeChannel", "insiders")
			u.RawQuery = q.Encode()
			browseURL = u.String()
		}
	}

	if err := a.browser.Browse(browseURL); err != nil {
		return fmt.Errorf("error opening Visual Studio Code: %w", err)
	}

	return nil
}

func vscodeProtocolURL(codespaceName string, useInsiders bool) string {
	application := "vscode"
	if useInsiders {
		application = "vscode-insiders"
	}
	return fmt.Sprintf("%s://github.codespaces/connect?name=%s", application, url.QueryEscape(codespaceName))
}
