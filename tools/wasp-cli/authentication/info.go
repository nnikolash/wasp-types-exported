package authentication

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/cliclients"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/waspcmd"
)

type AuthInfoOutput struct {
	AuthenticationMethod string
	AuthenticationURL    string
}

var _ log.CLIOutput = &AuthInfoOutput{}

func (l *AuthInfoOutput) AsText() (string, error) {
	template := `Authentication Method: {{ .AuthenticationMethod }}
Authentication URL: {{ .AuthenticationURL }}`
	return log.ParseCLIOutputTemplate(l, template)
}

func initInfoCmd() *cobra.Command {
	var node string
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Receive information about the authentication methods",
		Run: func(cmd *cobra.Command, args []string) {
			// Auth is currently not inside Swagger, so this is a temporary change
			node = waspcmd.DefaultWaspNodeFallback(node)
			client := cliclients.WaspClient(node)
			authInfo, _, err := client.AuthApi.AuthInfo(context.Background()).Execute()

			log.Check(err)

			log.PrintCLIOutput(&AuthInfoOutput{
				AuthenticationMethod: authInfo.Scheme,
				AuthenticationURL:    authInfo.AuthURL,
			})
		},
	}
	waspcmd.WithWaspNodeFlag(cmd, &node)
	return cmd
}
