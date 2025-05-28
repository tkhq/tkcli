## Turnkey CLI (tkcli) Command Guide

`tkcli` is a command-line interface tool for interacting with the Turnkey API. It allows you to manage activities, keys, wallets, and more directly from your terminal.

### Persistent Flags

The following flags can be used with most `tkcli` commands:

- `--keys-folder`, `-d <directory>`: Specifies the directory in which to locate API keys (defaults to a standard location).
- `--encryption-keys-folder <directory>`: Specifies the directory in which to locate encryption keys (defaults to a standard location).
- `--key-name`, `-k <name>`: Specifies the name of the API key to use for interacting with the Turnkey API service (defaults to "default").
- `--encryption-key-name <name>`: Specifies the name of the encryption key to use (defaults to "default").
- `--host <hostname>`: Specifies the hostname of the API server (defaults to "api.turnkey.com").
- `--organization <ID>`: Specifies the organization ID to interact with. This is often required.
- `--output`, `-o <format>`: Specifies the output format. Available formats are `json` (default) and `yaml`.

### Command Reference

#### `activities`

Interact with the API activities.

- __`list`__: Return the set of activities for an organization.

  - Flags:
    - `--status <status1,status2,...>`: Only include activities whose status is declared in this set. Available statuses: `created`, `pending`, `completed`, `failed`, `consensus`, `consensus_needed`, `rejected`, `all`.
  - Example: `tkcli activities list --status completed,failed`

- __`get <activity-id>`__: Return the details and status of a particular activity.

  - Arguments:
    - `<activity-id>`: The ID of the activity to retrieve.
  - Example: `tkcli activities get <your-activity-id>`

#### `curves`

Interact with the available curves.

- __`list`__: Return the available curve types.
  - Example: `tkcli curves list`

#### `decrypt`

Decrypt a ciphertext from a bundle exported from a Turnkey secure enclave.

- Flags:

  - `--export-bundle-input <filepath>`: Filepath to read the export bundle from (Required).
  - `--plaintext-output <filepath>`: Optional filepath to write the plaintext to.
  - `--signer-quorum-key <hex>`: Optional override for the signer quorum key (for testing only).
  - `--solana-address <address>`: Optional Solana address, for use when exporting Solana private keys in the proper format.

- Example: `tkcli decrypt --export-bundle-input ./export-bundle.txt --plaintext-output ./decrypted-key.txt`

#### `encrypt`

Encrypt a plaintext into a bundle to be imported to a Turnkey secure enclave.

- Flags:

  - `--import-bundle-input <filepath>`: Filepath to read the import bundle from (result of `init-import`) (Required).
  - `--encrypted-bundle-output <filepath>`: Filepath to write the encrypted bundle to (Required).
  - `--plaintext-input <filepath>`: Filepath to read the plaintext from that will be encrypted (Required).
  - `--key-format <format>`: Optional formatting to apply to the plaintext before it is encrypted. Available formats: `mnemonic` (default), `hexadecimal`, `solana`.
  - `--user <ID>`: ID of the user encrypting the plaintext (Required).
  - `--signer-quorum-key <hex>`: Optional override for the signer quorum key (for testing only).

- Example: `tkcli encrypt --import-bundle-input ./import-bundle.txt --encrypted-bundle-output ./encrypted-bundle.txt --plaintext-input ./plaintext-key.txt --user <your-user-id>`

#### `ethereum` (`eth`)

Perform actions related to Ethereum.

- __`transaction` (`tx`)__: Perform signing and other actions for a transaction.

  - Flags:

    - `--signer`, `-s <signer>`: Wallet account address, private key address, or private key ID (Required).
    - `--payload <payload>`: Payload of the transaction (Required).

  - Example: `tkcli eth tx --signer <your-signer-id> --payload <your-transaction-payload>`

#### `generate`

Generate keys.

- __`api-key`__: Generate a Turnkey API key.

  - Flags:

    - `--curve <type>`: Curve type for API key. Supported types: `p256` (default), `secp256k1`, and `ed25519`.
    - `--organization <ID>`: Organization ID (Required).

  - Example: `tkcli generate api-key --organization <your-organization-id> --curve secp256k1`

- __`encryption-key`__: Generate a Turnkey encryption key.

  - Flags:

    - `--user <ID>`: ID of user generating the encryption key (Required).
    - `--organization <ID>`: Organization ID (Required).

  - Example: `tkcli generate encryption-key --user <your-user-id> --organization <your-organization-id>`

#### `organizations` (`o`, `org`, `organization`)

Interact with organizations stored in Turnkey.

- __`create`__: Create a new organization.

  - Flags:
    - `--name <name>`: Name of the organization (Required).
  - *Note: Based on the code, this command appears to be implemented to create private keys instead of an organization.*
  - Example: `tkcli organizations create --name "My New Org"`

#### `private-keys` (`pk`)

Interact with private keys.

- __`create`__: Create a new private key.

  - Flags:

    - `--address-format <format1,format2,...>`: Address format(s) for private key (Required). Use `tkcli address-formats list` for available formats.
    - `--curve <curve>`: Curve to use for the generation of the private key (Required). Use `tkcli curves list` for available curves.
    - `--name <name>`: Name to be applied to the private key (Required).
    - `--tag <tag1,tag2,...>`: Tag(s) to be applied to the private key.

  - Example: `tkcli private-keys create --address-format ethereum --curve secp256k1 --name "My Eth Key"`

- __`list`__: Return private keys for the organization.

  - Example: `tkcli private-keys list`

- __`export`__: Export a private key.

  - Flags:

    - `--id <name-or-id>`: Name or ID of private key to export (Required).
    - `--encryption-key-name <name>`: Name of the encryption key to use (Required).
    - `--export-bundle-output <filepath>`: Filepath to write the export bundle to (Required).

  - Example: `tkcli private-keys export --id "My Eth Key" --encryption-key-name default --export-bundle-output ./private-key-export.txt`

- __`init-import`__: Initialize private key import.

  - Flags:

    - `--user <ID>`: ID of user importing the private key (Required).
    - `--import-bundle-output <filepath>`: Filepath to write the import bundle to (Required).

  - Example: `tkcli private-keys init-import --user <your-user-id> --import-bundle-output ./private-key-import-bundle.txt`

- __`import`__: Import a private key.

  - Flags:

    - `--user <ID>`: ID of user importing the private key (Required).
    - `--encrypted-bundle-input <filepath>`: Filepath to read the encrypted bundle from (Required).
    - `--address-format <format1,format2,...>`: Address format(s) for private key (Required).
    - `--curve <curve>`: Curve to use for the generation of the private key (Required).
    - `--name <name>`: Name to be applied to the private key (Required).

  - Example: `tkcli private-keys import --user <your-user-id> --encrypted-bundle-input ./encrypted-private-key.txt --address-format ethereum --curve secp256k1 --name "Imported Eth Key"`

#### `raw`

Send low-level (raw) requests to the Turnkey API.

- __`sign`__: Sign a raw payload.

  - Flags:

    - `--signer`, `-s <signer>`: Wallet account address, private key address, or private key ID (Required).
    - `--payload <payload>`: Payload to be signed (Required).
    - `--payload-encoding <encoding>`: Encoding of payload (defaults to `text/utf8`).
    - `--hash-function <function>`: Hash function (defaults to `sha256`).

  - Example: `tkcli raw sign --signer <your-signer-id> --payload "Hello, Turnkey!"`

#### `request` (`req`, `r`)

Given a request body, generate a stamp for it, and send it to the Turnkey API server.

- Flags:

  - `--path <path>`: Path for the API request (Required).
  - `--body <body-or-filepath>`: Body of the request. Can be `-` for stdin or prefixed with `@` for a filename (defaults to `-`).
  - `--no-post`: Generates the stamp and displays the cURL command but does NOT post the request.

- Example: `tkcli request --path /v1/my-endpoint --body '{"data": "some data"}'`

#### `transaction-types`

Interact with the available transaction types.

- __`list`__: Return the available transaction types.
  - Example: `tkcli transaction-types list`

#### `version`

Display build and version information.

- Example: `tkcli version`

#### `wallets`

Interact with wallets.

- __`create`__: Create a new wallet.

  - Flags:
    - `--name <name>`: Name to be applied to the wallet (Required).
  - Example: `tkcli wallets create --name "My New Wallet"`

- __`list`__: Return wallets for the organization.

  - Example: `tkcli wallets list`

- __`export`__: Export a wallet.

  - Flags:

    - `--name <name-or-id>`: Name or ID of wallet to export (Required).
    - `--encryption-key-name <name>`: Name of the encryption key to use (Required).
    - `--export-bundle-output <filepath>`: Filepath to write the export bundle to (Required).

  - Example: `tkcli wallets export --name "My New Wallet" --encryption-key-name default --export-bundle-output ./wallet-export.txt`

- __`init-import`__: Initialize wallet import.

  - Flags:

    - `--user <ID>`: ID of user importing the wallet (Required).
    - `--import-bundle-output <filepath>`: Filepath to write the import bundle to (Required).

  - Example: `tkcli wallets init-import --user <your-user-id> --import-bundle-output ./wallet-import-bundle.txt`

- __`import`__: Import a wallet.

  - Flags:

    - `--user <ID>`: ID of user importing the wallet (Required).
    - `--name <name>`: Name to be applied to the wallet (Required).
    - `--encrypted-bundle-input <filepath>`: Filepath to read the encrypted bundle from (Required).

  - Example: `tkcli wallets import --user <your-user-id> --name "Imported Wallet" --encrypted-bundle-input ./encrypted-wallet.txt`

- __`accounts` (`acc`)__: Interact with wallet accounts.

  - __`list`__: Return accounts for the wallet.

    - Flags:
      - `--wallet <name-or-id>`: Name or ID of wallet used to fetch accounts (Required).
    - Example: `tkcli wallets accounts list --wallet "My New Wallet"`

  - __`create`__: Create a new account for a wallet.

    - Flags:

      - `--wallet <name-or-id>`: Name or ID of wallet used for account creation (Required).
      - `--address-format <format>`: Address format for account (Required). Use `tkcli address-formats list` for available formats.
      - `--curve <curve>`: Curve for account. If unset, will predict based on address format. Use `tkcli curves list` for available curves.
      - `--path-format <format>`: The derivation path format for account (defaults to `bip32`).
      - `--path <path>`: The derivation path for account. If unset, will predict next path.

    - Example: `tkcli wallets accounts create --wallet "My New Wallet" --address-format ethereum`

  - __`export`__: Export a wallet account.

    - Flags:

      - `--address <address>`: Address of wallet account to export (Required).
      - `--encryption-key-name <name>`: Name of the encryption key to use (Required).
      - `--export-bundle-output <filepath>`: Filepath to write the export bundle to (Required).

    - Example: `tkcli wallets accounts export --address <your-account-address> --encryption-key-name default --export-bundle-output ./wallet-account-export.txt`
