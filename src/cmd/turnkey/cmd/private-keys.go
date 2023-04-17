package cmd

import (
	"github.com/tkhq/tkcli/src/api/client"
	"github.com/tkhq/tkcli/src/api/client/private_keys"
	"github.com/tkhq/tkcli/src/api/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	signingKeyID string

	privateKeysOrgID string

	privateKeysCreateAddressFormats []string
	privateKeysCreateCurve          string
	privateKeysCreateName           string
	privateKeysCreateTags           []string
)

func init() {
	privateKeysCreateCmd.Flags().StringSliceVar(&privateKeysCreateAddressFormats, "address-format", nil, "address format(s) for private key.  For a list of formats, use 'turnkey address-formats list'.")
	privateKeysCreateCmd.Flags().StringVar(&privateKeysCreateCurve, "curve", "", "curve to use for the generation of the private key.  For a list of available curves, use 'turnkey curves list'.")
	privateKeysCreateCmd.Flags().StringVar(&privateKeysCreateName, "name", "", "name to be applied to the private key")
	privateKeysCreateCmd.Flags().StringSliceVar(&privateKeysCreateTags, "tag", make([]string, 0), "tag(s) to be applied to the private key")

	privateKeysCmd.AddCommand(privateKeysCreateCmd)
	privateKeysCmd.AddCommand(privateKeysListCmd)

	rootCmd.AddCommand(privateKeysCmd)
}

var privateKeysCmd = &cobra.Command{
	Use:   "private-keys interacts with private keys stored in Turnkey",
	Short: "private-keys interacts with private keys",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		LoadKeypair("")
	},
	Aliases: []string{"pk"},
}

var privateKeysCreateCmd = &cobra.Command{
	Use:   "create a new private key",
	Short: "create a new private key",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(privateKeysCreateAddressFormats) < 1 {
			OutputError(errors.New("must specify at least one address format"))
		}

		if privateKeysCreateCurve == "" {
			OutputError(errors.New("curve cannot be empty"))
		}

		if privateKeysCreateName == "" {
			OutputError(errors.New("name for private key must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		curve := models.Immutableactivityv1Curve(privateKeysCreateCurve)

		addressFormats := make([]models.Immutableactivityv1AddressFormat, len(privateKeysCreateAddressFormats))

		for n, f := range privateKeysCreateAddressFormats {
			addressFormats[n] = models.Immutableactivityv1AddressFormat(f)
		}

		activity := string(models.V1ActivityTypeACTIVITYTYPECREATEPRIVATEKEYS)

		// if privateKeysCreateTags == nil {
		// 	privateKeysCreateTags = make([]string, 0)
		// }

		params := private_keys.NewPublicAPIServiceCreatePrivateKeysParams()
		params.SetBody(&models.V1CreatePrivateKeysRequest{
			OrganizationID: &Organization,
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
			TimestampMs: RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(errors.Wrap(err, "request validation failed"))
		}

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

var privateKeysListCmd = &cobra.Command{
	Use:   "list private keys",
	Short: "list private keys for the organization",
	Run: func(cmd *cobra.Command, args []string) {
		params := private_keys.NewPublicAPIServiceGetPrivateKeysParams()

		params.SetBody(&models.V1GetPrivateKeysRequest{
			OrganizationID: &Organization,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(errors.Wrap(err, "request validation failed"))
		}

		resp, err := client.Default.PrivateKeys.PublicAPIServiceGetPrivateKeys(params, new(Authenticator))
		if err != nil {
			OutputError(errors.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(errors.Errorf("failed to list private keys: %d: %s", resp.Code(), resp.Error()))
		}

		Output(resp.Payload)
	},
}

// LoadSigningKey require-loads a signing key
func LoadSigningKey(name string) {
	if name != "" {
		signingKeyID = name
	}

	if signingKeyID == "" {
		OutputError(errors.New("no private key provided"))
	}

	if _, err := uuid.Parse(signingKeyID); err != nil {
		signingKeyID, err = lookupPrivateKeyByName(signingKeyID)
		if err != nil {
			OutputError(errors.Wrap(err, "provided private key was not a UUID and lookup by name failed"))
		}
	}
}

func lookupPrivateKeyByName(name string) (string, error) {
	params := private_keys.NewPublicAPIServiceGetPrivateKeysParams()

	params.SetBody(&models.V1GetPrivateKeysRequest{
		OrganizationID: &Organization,
	})

	if err := params.Body.Validate(nil); err != nil {
		return "", errors.Wrap(err, "formulation of a lookup by name request failed")
	}

	resp, err := client.Default.PrivateKeys.PublicAPIServiceGetPrivateKeys(params, new(Authenticator))
	if err != nil {
		return "", errors.Wrap(err, "lookup by name failed")
	}

	if !resp.IsSuccess() {
		return "", errors.Errorf("lookup by name failed: %d: %s", resp.Code(), resp.Error())
	}

	for _, k := range resp.Payload.PrivateKeys {
		if *k.PrivateKeyName == name {
			return *k.PrivateKeyID, nil
		}
	}

	return "", errors.Errorf("private key name %q not found in list of private keys", name)
}
