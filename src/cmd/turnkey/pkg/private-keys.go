package pkg

import (
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/private_keys"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/util"
)

var (
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

	privateKeyInitImportCmd.Flags().StringVar(&user, "user", "", "ID of user to importing the private key")
	privateKeyInitImportCmd.Flags().StringVar(&importBundlePath, "import-bundle-output", "", "filepath to write the import bundle to.")

	privateKeyImportCmd.Flags().StringVar(&user, "user", "", "ID of user to importing the private key")
	privateKeyImportCmd.Flags().StringVar(&privateKeysCreateName, "name", "", "name to be applied to the private key.")
	privateKeyImportCmd.Flags().StringVar(&encryptedBundlePath, "encrypted-bundle-input", "", "filepath to read the encrypted bundle from.")
	privateKeyImportCmd.Flags().StringSliceVar(&privateKeysCreateAddressFormats, "address-format", nil, "address format(s) for private key.  For a list of formats, use 'turnkey address-formats list'.")
	privateKeyImportCmd.Flags().StringVar(&privateKeysCreateCurve, "curve", "", "curve to use for the generation of the private key.  For a list of available curves, use 'turnkey curves list'.")

	privateKeysCmd.AddCommand(privateKeysCreateCmd)
	privateKeysCmd.AddCommand(privateKeysListCmd)
	privateKeysCmd.AddCommand(privateKeyInitImportCmd)
	privateKeysCmd.AddCommand(privateKeyImportCmd)

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
			OutputError(eris.New("--address-format must not be empty"))
		}

		if privateKeysCreateCurve == "" {
			OutputError(eris.New("--curve must be specified"))
		}

		if privateKeysCreateName == "" {
			OutputError(eris.New("--name must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		curve := models.Curve(privateKeysCreateCurve)

		if curve == Help {
			Output(models.CurveEnum)
			return
		}

		addressFormats := make([]models.AddressFormat, len(privateKeysCreateAddressFormats))

		for n, f := range privateKeysCreateAddressFormats {
			if f == Help {
				Output(models.AddressFormatEnum)
				return
			}

			addressFormats[n] = models.AddressFormat(f)
		}

		activity := string(models.ActivityTypeCreatePrivateKeysV2)

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
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

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

var privateKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return private keys for the organization",
	Run: func(cmd *cobra.Command, args []string) {
		params := private_keys.NewGetPrivateKeysParams()

		params.SetBody(&models.GetPrivateKeysRequest{
			OrganizationID: &Organization,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().PrivateKeys.GetPrivateKeys(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to list private keys: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

var privateKeyInitImportCmd = &cobra.Command{
	Use:   "init-import",
	Short: "Initialize private key import",
	PreRun: func(cmd *cobra.Command, args []string) {
		if user == "" {
			OutputError(eris.New("--user must be specified"))
		}

		if importBundlePath == "" {
			OutputError(eris.New("--import-bundle-output must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		activity := string(models.ActivityTypeInitImportPrivateKey)

		params := private_keys.NewInitImportPrivateKeyParams()
		params.SetBody(&models.InitImportPrivateKeyRequest{
			OrganizationID: &Organization,
			Parameters: &models.InitImportPrivateKeyIntent{
				UserID: &user,
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().PrivateKeys.InitImportPrivateKey(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to initialize private key import: %s", resp.Error()))
		}

		importBundle := resp.Payload.Activity.Result.InitImportPrivateKeyResult.ImportBundle
		err = writeFile(*importBundle, importBundlePath)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to write import bundle to file"))
		}
	},
}

var privateKeyImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a private key",
	PreRun: func(cmd *cobra.Command, args []string) {
		if user == "" {
			OutputError(eris.New("--user must be specified"))
		}

		if encryptedBundlePath == "" {
			OutputError(eris.New("--encrypted-bundle-input must be specified"))
		}

		if len(privateKeysCreateAddressFormats) < 1 {
			OutputError(eris.New("--address-format must not be empty"))
		}

		if privateKeysCreateCurve == "" {
			OutputError(eris.New("--curve must be specified"))
		}

		if privateKeysCreateName == "" {
			OutputError(eris.New("--name must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		encryptedBundle, err := readFile(encryptedBundlePath)
		if err != nil {
			OutputError(err)
		}

		curve := models.Curve(privateKeysCreateCurve)

		if curve == Help {
			Output(models.CurveEnum)
			return
		}

		addressFormats := make([]models.AddressFormat, len(privateKeysCreateAddressFormats))

		for n, f := range privateKeysCreateAddressFormats {
			if f == Help {
				Output(models.AddressFormatEnum)
				return
			}

			addressFormats[n] = models.AddressFormat(f)
		}

		activity := string(models.ActivityTypeImportPrivateKey)

		params := private_keys.NewImportPrivateKeyParams()
		params.SetBody(&models.ImportPrivateKeyRequest{
			OrganizationID: &Organization,
			Parameters: &models.ImportPrivateKeyIntent{
				UserID:          &user,
				PrivateKeyName:  &privateKeysCreateName,
				EncryptedBundle: &encryptedBundle,
				Curve:           &curve,
				AddressFormats:  addressFormats,
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().PrivateKeys.ImportPrivateKey(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to import private key: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}
