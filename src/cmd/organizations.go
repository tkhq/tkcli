package cmd

import (
	"github.com/pkg/errors"
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
	Use:     "organizations interacts with organizations stored in Turnkey",
	Short:   "organizations interacts with organizations stored in Turnkey",
	Aliases: []string{"o", "org", "organization"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadClient()
	},
}

var organizationsCreateCmd = &cobra.Command{
	Use:   "create a new organization",
	Short: "create a new organization",
	PreRun: func(cmd *cobra.Command, args []string) {
		if organizationsCreateName == "" {
			OutputError(errors.New("name for private key must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		curve := models.Immutableactivityv1Curve(privateKeysCreateCurve)

		addressFormats := make([]models.Immutableactivityv1AddressFormat, len(privateKeysCreateAddressFormats))

		for n, f := range privateKeysCreateAddressFormats {
			addressFormats[n] = models.Immutableactivityv1AddressFormat(f)
		}

		params := private_keys.NewPublicAPIServiceCreatePrivateKeysParams()
		params.SetBody(&models.V1CreatePrivateKeysRequest{
			OrganizationID: &privateKeysOrgID,
			Parameters: &models.V1CreatePrivateKeysIntent{
				PrivateKeys: []*models.V1PrivateKeyParams{
					{
						AddressFormats: addressFormats,
						Curve:          &curve,
						PrivateKeyName: &privateKeysCreateName,
						PrivateKeyTags: privateKeysCreateTags,
					},
				},
			},
		})

		resp, err := APIClient.V0().PrivateKeys.PublicAPIServiceCreatePrivateKeys(params, APIClient.Authenticator)
		if err != nil {
			OutputError(errors.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(errors.Errorf("failed to create private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
