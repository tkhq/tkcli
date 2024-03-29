package pkg

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/private_keys"
	"github.com/tkhq/go-sdk/pkg/api/models"
)

var organizationsCreateName string

func init() {
	organizationsCreateCmd.Flags().StringVar(&organizationsCreateName, "name", "", "name of the organization")

	organizationsCmd.AddCommand(organizationsCreateCmd)

	rootCmd.AddCommand(organizationsCmd)
}

var organizationsCmd = &cobra.Command{
	Use:     "organizations",
	Short:   "Interact with organizations stored in Turnkey",
	Aliases: []string{"o", "org", "organization"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadClient()
	},
}

var organizationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new organization",
	PreRun: func(cmd *cobra.Command, args []string) {
		if organizationsCreateName == "" {
			OutputError(eris.New("--name must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		curve := models.Curve(privateKeysCreateCurve)

		addressFormats := make([]models.AddressFormat, len(privateKeysCreateAddressFormats))

		for n, f := range privateKeysCreateAddressFormats {
			addressFormats[n] = models.AddressFormat(f)
		}

		params := private_keys.NewCreatePrivateKeysParams()
		params.SetBody(&models.CreatePrivateKeysRequest{
			OrganizationID: &Organization,
			Parameters: &models.CreatePrivateKeysIntentV2{
				PrivateKeys: []*models.PrivateKeyParams{
					{
						AddressFormats: addressFormats,
						Curve:          &curve,
						PrivateKeyName: &privateKeysCreateName,
						PrivateKeyTags: privateKeysCreateTags,
					},
				},
			},
		})

		resp, err := APIClient.V0().PrivateKeys.CreatePrivateKeys(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to create private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
