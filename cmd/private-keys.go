package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkhq/tkcli/api/client"
	"github.com/tkhq/tkcli/api/client/private_keys"
	"github.com/tkhq/tkcli/api/models"
)

var (
	privateKeysOrgID string

	privateKeysCreateAddressFormats []string
	privateKeysCreateCurve          string
	privateKeysCreateName           string
	privateKeysCreateTags           []string
)

func init() {
	privateKeysCmd.PersistentFlags().StringVar(&privateKeysOrgID, "organization", "", "organization ID to be used")

	privateKeysCreateCmd.Flags().StringSliceVar(&privateKeysCreateAddressFormats, "address-format", nil, "address format(s) for private key.  For a list of formats, use 'turnkey private-keys formats list'.")
	privateKeysCreateCmd.Flags().StringVar(&privateKeysCreateCurve, "curve", "", "curve to use for the generation of the private key.  For a list of available curves, use 'turnkey private-keys curves list'.")
	privateKeysCreateCmd.Flags().StringVar(&privateKeysCreateName, "name", "", "name to be applied to the private key")
	privateKeysCreateCmd.Flags().StringSliceVar(&privateKeysCreateTags, "tag", nil, "tag(s) to be applied to the private key")

   privateKeysCmd.AddCommand(privateKeysCreateCmd)

   rootCmd.AddCommand(privateKeysCmd)
}

var privateKeysCmd = &cobra.Command{
	Use:   "private-keys interacts with private keys stored in Turnkey",
	Short: "private-keys interacts with private keys",
	PreRun: func(cmd *cobra.Command, args []string) {
      if err := LoadKeypair(""); err != nil {
         OutputError(errors.Wrap(err, "failed to load API key"))
      }
	},
	Aliases: []string{"pk"},
}

var privateKeysCreateCmd = &cobra.Command{
   Use: "create a new private key",
   Short: "create a new private key",
	PreRun: func(cmd *cobra.Command, args []string) {
		if privateKeysOrgID == "" {
         OutputError(errors.New("organization ID must be set"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		apiClient := client.New(nil, nil)

		curve := models.Immutableactivityv1Curve(privateKeysCreateCurve)

		addressFormats := make([]models.Immutableactivityv1AddressFormat, len(privateKeysCreateAddressFormats))

		for n, f := range privateKeysCreateAddressFormats {
			addressFormats[n] = models.Immutableactivityv1AddressFormat(f)
		}

		params := &private_keys.PublicAPIServiceCreatePrivateKeysParams{
			Body: &models.V1CreatePrivateKeysRequest{
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
			},
		}

		resp, err := apiClient.PrivateKeys.PublicAPIServiceCreatePrivateKeys(params, new(Authenticator))
		if err != nil {
			OutputError(errors.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(errors.Errorf("failed to create private key: %d: %s", resp.Code(), resp.Error()))
		}

		Output(resp.Payload)
	},
}
