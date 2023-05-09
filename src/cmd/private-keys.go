package cmd

import (
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/private_keys"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/util"
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
	Use:   "private-keys",
	Short: "Interact with private keys",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadClient()
	},
	Aliases: []string{"pk"},
}

var privateKeysCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new private key",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(privateKeysCreateAddressFormats) < 1 {
			OutputError(eris.New("must specify at least one address format"))
		}

		if privateKeysCreateCurve == "" {
			OutputError(eris.New("curve cannot be empty"))
		}

		if privateKeysCreateName == "" {
			OutputError(eris.New("name for private key must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		curve := models.Immutableactivityv1Curve(privateKeysCreateCurve)

		if curve == Help {
			Output(models.Curves())

			return
		}

		addressFormats := make([]models.Immutableactivityv1AddressFormat, len(privateKeysCreateAddressFormats))

		for n, f := range privateKeysCreateAddressFormats {
			if f == Help {
				Output(models.AddressFormats())

				return
			}

			addressFormats[n] = models.Immutableactivityv1AddressFormat(f)
		}

		activity := string(models.V1ActivityTypeACTIVITYTYPECREATEPRIVATEKEYS)

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
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().PrivateKeys.PublicAPIServiceCreatePrivateKeys(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to create private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

var privateKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return private keys for the organization",
	Run: func(cmd *cobra.Command, args []string) {
		params := private_keys.NewPublicAPIServiceGetPrivateKeysParams()

		params.SetBody(&models.V1GetPrivateKeysRequest{
			OrganizationID: &Organization,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().PrivateKeys.PublicAPIServiceGetPrivateKeys(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to list private keys: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

// LoadSigningKey require-loads a signing key.
func LoadSigningKey(name string) {
	if name != "" {
		signingKeyID = name
	}

	if signingKeyID == "" {
		OutputError(eris.New("no private key provided"))
	}

	if _, err := uuid.Parse(signingKeyID); err != nil {
		signingKeyID, err = lookupPrivateKeyByName(signingKeyID)
		if err != nil {
			OutputError(eris.Wrap(err, "provided private key was not a UUID and lookup by name failed"))
		}
	}
}

func lookupPrivateKeyByName(name string) (string, error) {
	params := private_keys.NewPublicAPIServiceGetPrivateKeysParams()

	params.SetBody(&models.V1GetPrivateKeysRequest{
		OrganizationID: &Organization,
	})

	if err := params.Body.Validate(nil); err != nil {
		return "", eris.Wrap(err, "formulation of a lookup by name request failed")
	}

	resp, err := APIClient.V0().PrivateKeys.PublicAPIServiceGetPrivateKeys(params, APIClient.Authenticator)
	if err != nil {
		return "", eris.Wrap(err, "lookup by name failed")
	}

	if !resp.IsSuccess() {
		return "", eris.Errorf("lookup by name failed: %s", resp.Error())
	}

	for _, k := range resp.Payload.PrivateKeys {
		if *k.PrivateKeyName == name {
			return *k.PrivateKeyID, nil
		}
	}

	return "", eris.Errorf("private key name %q not found in list of private keys", name)
}
