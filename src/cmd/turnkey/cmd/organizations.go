package cmd

import (
	"github.com/tkhq/tkcli/src/api/client"
	"github.com/tkhq/tkcli/src/api/client/private_keys"
	"github.com/tkhq/tkcli/src/api/models"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	organizationsCreateName string
)

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
		LoadKeypair("")
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

		resp, err := client.Default.PrivateKeys.PublicAPIServiceCreatePrivateKeys(params, new(Authenticator))
		if err != nil {
			OutputError(errors.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(errors.Errorf("failed to create private key: %d: %s", resp.Code(), resp.Error()))
		}

		Output(resp.Payload)
	},
}
