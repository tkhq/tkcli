package pkg

import (
	"encoding/hex"
	"fmt"

	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"

	"github.com/tkhq/go-sdk/pkg/api/client/wallets"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/encryptionkey"
	"github.com/tkhq/go-sdk/pkg/util"
)

var (
	walletNameOrID             string
	walletAccountAddressFormat string
	walletAccountCurve         string
	walletAccountPathFormat    string
	walletAccountPath          string
	walletAccountAddress       string
)

func init() {
	walletCreateCmd.Flags().StringVar(&walletNameOrID, "name", "", "name to be applied to the wallet.")

	walletExportCmd.Flags().StringVar(&walletNameOrID, "name", "", "name or ID of wallet to export.")
	walletExportCmd.Flags().StringVar(&exportBundlePath, "export-bundle-output", "", "filepath to write the export bundle to.")

	walletAccountExportCmd.Flags().StringVar(&walletAccountAddress, "address", "", "address of wallet account to export.")
	walletAccountExportCmd.Flags().StringVar(&exportBundlePath, "export-bundle-output", "", "filepath to write the export bundle to.")

	walletInitImportCmd.Flags().StringVar(&User, "user", "", "ID of user to importing the wallet")
	walletInitImportCmd.Flags().StringVar(&importBundlePath, "import-bundle-output", "", "filepath to write the import bundle to.")

	walletImportCmd.Flags().StringVar(&User, "user", "", "ID of user to importing the wallet")
	walletImportCmd.Flags().StringVar(&walletNameOrID, "name", "", "name to be applied to the wallet.")
	walletImportCmd.Flags().StringVar(&encryptedBundlePath, "encrypted-bundle-input", "", "filepath to read the encrypted bundle from.")

	walletAccountsListCmd.Flags().StringVar(&walletNameOrID, "wallet", "", "name or ID of wallet used to fetch accounts.")

	walletAccountCreateCmd.Flags().StringVar(&walletNameOrID, "wallet", "", "name or ID of wallet used for account creation.")
	walletAccountCreateCmd.Flags().StringVar(&walletAccountAddressFormat, "address-format", "", "address format for account. For a list of formats, use 'turnkey address-formats list'.")
	walletAccountCreateCmd.Flags().StringVar(&walletAccountCurve, "curve", "", "curve for account. For a list of curves, use 'turnkey curves list'. If unset, will predict based on address format.")
	walletAccountCreateCmd.Flags().StringVar(&walletAccountPathFormat, "path-format", string(models.PathFormatBip32), "the derivation path format for account.")
	walletAccountCreateCmd.Flags().StringVar(&walletAccountPath, "path", "", "the derivation path for account. If unset, will predict next path.")

	walletAccountsCmd.AddCommand(walletAccountsListCmd)
	walletAccountsCmd.AddCommand(walletAccountCreateCmd)
	walletsCmd.AddCommand(walletCreateCmd)
	walletsCmd.AddCommand(walletsListCmd)
	walletsCmd.AddCommand(walletExportCmd)
	walletsCmd.AddCommand(walletInitImportCmd)
	walletsCmd.AddCommand(walletImportCmd)
	walletsCmd.AddCommand(walletAccountsCmd)

	rootCmd.AddCommand(walletsCmd)
}

var walletsCmd = &cobra.Command{
	Use:   "wallets",
	Short: "Interact with wallets",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		basicSetup(cmd)
		LoadKeypair("")
		LoadClient()
		LoadEncryptionKeypair("")
	},
	Aliases: []string{},
}

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	PreRun: func(cmd *cobra.Command, args []string) {
		if walletNameOrID == "" {
			OutputError(eris.New("--name must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		activity := string(models.ActivityTypeCreateWallet)

		params := wallets.NewCreateWalletParams()
		params.SetBody(&models.CreateWalletRequest{
			OrganizationID: &Organization,
			Parameters: &models.CreateWalletIntent{
				WalletName: &walletNameOrID,
				Accounts:   []*models.WalletAccountParams{},
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.CreateWallet(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to create wallet: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

var walletsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return wallets for the organization",
	Run: func(cmd *cobra.Command, args []string) {
		params := wallets.NewGetWalletsParams()

		params.SetBody(&models.GetWalletsRequest{
			OrganizationID: &Organization,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.GetWallets(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to list wallets: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

var walletExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a wallet",
	PreRun: func(cmd *cobra.Command, args []string) {
		if walletNameOrID == "" {
			OutputError(eris.New("--name must be specified"))
		}

		if EncryptionKeyName == "" {
			OutputError(eris.New("--encryption-key-name must be specified"))
		}

		if exportBundlePath == "" {
			OutputError(eris.New("--export-bundle-output must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		wallet, err := lookupWallet(walletNameOrID)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to lookup wallet"))
		}

		tkPublicKey := EncryptionKeypair.GetPublicKey()
		kemPublicKey, err := encryptionkey.DecodeTurnkeyPublicKey(tkPublicKey)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to decode encryption public key"))
		}
		kemPublicKeyBytes, err := (*kemPublicKey).MarshalBinary()
		if err != nil {
			OutputError(eris.Wrap(err, "failed to marshal encryption public key"))
		}
		targetPublicKey := hex.EncodeToString(kemPublicKeyBytes)

		activity := string(models.ActivityTypeExportWallet)

		params := wallets.NewExportWalletParams()
		params.SetBody(&models.ExportWalletRequest{
			OrganizationID: &Organization,
			Parameters: &models.ExportWalletIntent{
				WalletID:        wallet.WalletID,
				TargetPublicKey: &targetPublicKey,
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.ExportWallet(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to export wallet: %s", resp.Error()))
		}

		exportBundle := resp.Payload.Activity.Result.ExportWalletResult.ExportBundle
		err = writeFile(*exportBundle, exportBundlePath)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to write export bundle to file"))
		}

		exportedWalletID := resp.Payload.Activity.Result.ExportWalletResult.WalletID
		Output(exportedWalletID)
	},
}

var walletAccountExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a wallet account",
	PreRun: func(cmd *cobra.Command, args []string) {
		if walletAccountAddress == "" {
			OutputError(eris.New("--address must be specified"))
		}

		if EncryptionKeyName == "" {
			OutputError(eris.New("--encryption-key-name must be specified"))
		}

		if exportBundlePath == "" {
			OutputError(eris.New("--export-bundle-output must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		tkPublicKey := EncryptionKeypair.GetPublicKey()
		kemPublicKey, err := encryptionkey.DecodeTurnkeyPublicKey(tkPublicKey)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to decode encryption public key"))
		}
		kemPublicKeyBytes, err := (*kemPublicKey).MarshalBinary()
		if err != nil {
			OutputError(eris.Wrap(err, "failed to marshal encryption public key"))
		}
		targetPublicKey := hex.EncodeToString(kemPublicKeyBytes)

		activity := string(models.ActivityTypeExportWalletAccount)

		params := wallets.NewExportWalletAccountParams()
		params.SetBody(&models.ExportWalletAccountRequest{
			OrganizationID: &Organization,
			Parameters: &models.ExportWalletAccountIntent{
				Address:         &walletAccountAddress,
				TargetPublicKey: &targetPublicKey,
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.ExportWalletAccount(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to export wallet account: %s", resp.Error()))
		}

		exportBundle := resp.Payload.Activity.Result.ExportWalletAccountResult.ExportBundle
		if err := writeFile(*exportBundle, exportBundlePath); err != nil {
			OutputError(eris.Wrap(err, "failed to write export bundle to file"))
		}

		exportedWalletAccountAddress := resp.Payload.Activity.Result.ExportWalletAccountResult.Address
		Output(exportedWalletAccountAddress)
	},
}

var walletInitImportCmd = &cobra.Command{
	Use:   "init-import",
	Short: "Initialize wallet import",
	PreRun: func(cmd *cobra.Command, args []string) {
		if User == "" {
			OutputError(eris.New("--user must be specified"))
		}

		if importBundlePath == "" {
			OutputError(eris.New("--import-bundle-output must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		activity := string(models.ActivityTypeInitImportWallet)

		params := wallets.NewInitImportWalletParams()
		params.SetBody(&models.InitImportWalletRequest{
			OrganizationID: &Organization,
			Parameters: &models.InitImportWalletIntent{
				UserID: &User,
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.InitImportWallet(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to initialize wallet import: %s", resp.Error()))
		}

		importBundle := resp.Payload.Activity.Result.InitImportWalletResult.ImportBundle
		err = writeFile(*importBundle, importBundlePath)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to write import bundle to file"))
		}
	},
}

var walletImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a wallet",
	PreRun: func(cmd *cobra.Command, args []string) {
		if User == "" {
			OutputError(eris.New("--user must be specified"))
		}

		if walletNameOrID == "" {
			OutputError(eris.New("--name must be specified"))
		}

		if encryptedBundlePath == "" {
			OutputError(eris.New("--encrypted-bundle-input must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		encryptedBundle, err := readFile(encryptedBundlePath)
		if err != nil {
			OutputError(err)
		}

		activity := string(models.ActivityTypeImportWallet)

		params := wallets.NewImportWalletParams()
		params.SetBody(&models.ImportWalletRequest{
			OrganizationID: &Organization,
			Parameters: &models.ImportWalletIntent{
				UserID:          &User,
				WalletName:      &walletNameOrID,
				EncryptedBundle: &encryptedBundle,
				Accounts:        []*models.WalletAccountParams{},
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.ImportWallet(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to import wallet: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

var walletAccountsCmd = &cobra.Command{
	Use:     "accounts",
	Short:   "Interact with wallet accounts",
	Aliases: []string{"acc"},
}

var walletAccountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new account for a wallet",
	PreRun: func(cmd *cobra.Command, args []string) {
		if walletNameOrID == "" {
			OutputError(eris.New("--name must be specified"))
		}

		if walletAccountAddressFormat == "" {
			OutputError(eris.New("--address-format must not be empty"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		wallet, err := lookupWallet(walletNameOrID)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to lookup wallet"))
		}

		addressFormat := models.AddressFormat(walletAccountAddressFormat)
		curve := models.Curve(walletAccountCurve)
		pathFormat := models.PathFormat(walletAccountPathFormat)
		path := walletAccountPath

		// set standard curve, if we can, if no override
		if curve == "" {
			if standardCurve := getCurveForAddressFormat(addressFormat); standardCurve != nil {
				curve = *standardCurve
			}
		}

		// set standard path, if we can, if no override
		if path == "" {
			accounts, err := listAccountsForWallet(wallet.WalletID)
			if err != nil {
				OutputError(eris.Wrap(err, "failed to lookup wallet accounts"))
			}

			// build path map to avoid conflicts
			paths := make(map[string]struct{})
			for _, account := range accounts {
				// we only need to care about accounts w/ this same address format
				if *account.AddressFormat != addressFormat {
					continue
				}

				paths[*account.Path] = struct{}{}
			}

			// find the next unused standard path
			for i := 0; i < len(paths)+1; i++ {
				if standardPath := getStandardPath(pathFormat, addressFormat, i); standardPath != nil {
					// we've found an unused path!
					if _, ok := paths[*standardPath]; !ok {
						path = *standardPath
						break
					}
				}
			}
		}

		activity := string(models.ActivityTypeCreateWalletAccounts)

		params := wallets.NewCreateWalletAccountsParams()
		params.SetBody(&models.CreateWalletAccountsRequest{
			OrganizationID: &Organization,
			Parameters: &models.CreateWalletAccountsIntent{
				WalletID: wallet.WalletID,
				Accounts: []*models.WalletAccountParams{
					{
						AddressFormat: &addressFormat,
						Curve:         &curve,
						PathFormat:    &pathFormat,
						Path:          &path,
					},
				},
			},
			TimestampMs: util.RequestTimestamp(),
			Type:        &activity,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.CreateWalletAccounts(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to create wallet account: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

var walletAccountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return accounts for the wallet",
	PreRun: func(cmd *cobra.Command, args []string) {
		if walletNameOrID == "" {
			OutputError(eris.New("--name must be specified"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		wallet, err := lookupWallet(walletNameOrID)
		if err != nil {
			OutputError(eris.Wrap(err, "failed to lookup wallet"))
		}

		params := wallets.NewGetWalletAccountsParams()

		params.SetBody(&models.GetWalletAccountsRequest{
			OrganizationID: &Organization,
			WalletID:       wallet.WalletID,
		})

		if err := params.Body.Validate(nil); err != nil {
			OutputError(eris.Wrap(err, "request validation failed"))
		}

		resp, err := APIClient.V0().Wallets.GetWalletAccounts(params, APIClient.Authenticator)
		if err != nil {
			OutputError(eris.Wrap(err, "request failed"))
		}

		if !resp.IsSuccess() {
			OutputError(eris.Errorf("failed to list wallets: %s", resp.Error()))
		}

		Output(resp.Payload)
	},
}

func lookupWallet(nameOrID string) (*models.Wallet, error) {
	params := wallets.NewGetWalletsParams()

	params.SetBody(&models.GetWalletsRequest{
		OrganizationID: &Organization,
	})

	if err := params.Body.Validate(nil); err != nil {
		OutputError(eris.Wrap(err, "request validation failed"))
	}

	resp, err := APIClient.V0().Wallets.GetWallets(params, APIClient.Authenticator)
	if err != nil {
		OutputError(eris.Wrap(err, "request failed"))
	}

	if !resp.IsSuccess() {
		OutputError(eris.Errorf("failed to list wallets: %s", resp.Error()))
	}

	for _, wallet := range resp.Payload.Wallets {
		if *wallet.WalletName == nameOrID || *wallet.WalletID == nameOrID {
			return wallet, nil
		}
	}

	return nil, eris.Errorf("wallet %q not found in list of wallets", nameOrID)
}

func listAccountsForWallet(walletID *string) ([]*models.WalletAccount, error) {
	params := wallets.NewGetWalletAccountsParams()

	params.SetBody(&models.GetWalletAccountsRequest{
		OrganizationID: &Organization,
		WalletID:       walletID,
	})

	if err := params.Body.Validate(nil); err != nil {
		OutputError(eris.Wrap(err, "request validation failed"))
	}

	resp, err := APIClient.V0().Wallets.GetWalletAccounts(params, APIClient.Authenticator)
	if err != nil {
		OutputError(eris.Wrap(err, "request failed"))
	}

	if !resp.IsSuccess() {
		OutputError(eris.Errorf("failed to list wallets: %s", resp.Error()))
	}

	return resp.Payload.Accounts, nil
}

func getCurveForAddressFormat(addressFormat models.AddressFormat) *models.Curve {
	switch addressFormat {
	case models.AddressFormatEthereum, models.AddressFormatCosmos, models.AddressFormatUncompressed:
		return models.NewCurve(models.CurveSecp256k1)
	case models.AddressFormatSolana:
		return models.NewCurve(models.CurveEd25519)
	default:
		// we're here because either we haven't updated this switch statement after adding new
		// address formats OR we've hit an address format that supports multiple curves so we'll
		// make no assumptions on the expected curve
		return nil
	}
}

func getStandardPath(pathFormat models.PathFormat, addressFormat models.AddressFormat, accountIndex int) *string {
	// we currently only support BIP-32 so we'll make no assumptions about the path if given a different path format
	if pathFormat != models.PathFormatBip32 {
		return nil
	}

	var path string

	switch addressFormat {
	case models.AddressFormatEthereum:
		path = fmt.Sprintf(`m/44'/60'/%d'/0/0`, accountIndex)
	case models.AddressFormatCosmos:
		path = fmt.Sprintf(`m/44'/118'/%d'/0/0`, accountIndex)
	case models.AddressFormatSolana:
		path = fmt.Sprintf(`m/44'/501'/%d'/0'`, accountIndex)
	case models.AddressFormatUncompressed, models.AddressFormatCompressed:
		path = fmt.Sprintf(`m/%d'`, accountIndex)
	default:
		// we're here because we haven't updated this switch statement after adding new
		// address formats so we'll make no assumptions on the expected path
		return nil
	}

	return &path
}
