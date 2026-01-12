# Turnkey CLI

[![Go Build Status](https://github.com/tkhq/tkcli/actions/workflows/go-build.yml/badge.svg)](https://github.com/tkhq/tkcli/actions/workflows/go-build.yml)

## Installation

We have multiple ways to install the CLI depending on your threat model.

Please check our work to whatever extent appropriate for your use case.

### Prerequisites

The Makefile assumes the presence of a few basic tools:

- `make`
- `bash`
- `Docker`

### Blind Trust

> :warning: Before you copy/paste, note that these are /low/ security options

If you are on an untrusted machine and are only evaluating our tools, we offer
easy low security install paths common in the industry.

Do note that any time you run an unverified binary off the internet you are
giving a third party full permission to execute any code they want on your
system. Github accounts, CDNs, and package repository accounts get compromised
all the time.

#### Download

| Version | OS    | Architecture | Download                                                                                       |
| ------- | ----- | ------------ | ---------------------------------------------------------------------------------------------- |
| v1.1.5  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.5/dist/turnkey.linux-x86_64)     |
| v1.1.5  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.5/dist/turnkey.linux-aarch64)   |
| v1.1.5  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.5/dist/turnkey.darwin-x86_64)   |
| v1.1.5  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.5/dist/turnkey.darwin-aarch64) |
| v1.1.4  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.4/dist/turnkey.linux-x86_64)     |
| v1.1.4  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.4/dist/turnkey.linux-aarch64)   |
| v1.1.4  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.4/dist/turnkey.darwin-x86_64)   |
| v1.1.4  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.4/dist/turnkey.darwin-aarch64) |
| v1.1.3  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.3/dist/turnkey.linux-x86_64)     |
| v1.1.3  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.3/dist/turnkey.linux-aarch64)   |
| v1.1.3  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.3/dist/turnkey.darwin-x86_64)   |
| v1.1.3  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.3/dist/turnkey.darwin-aarch64) |
| v1.1.2  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.linux-x86_64)     |
| v1.1.2  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.linux-aarch64)   |
| v1.1.2  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.darwin-x86_64)   |
| v1.1.2  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.darwin-aarch64) |
| v1.1.1  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.1/dist/turnkey.linux-x86_64)     |
| v1.1.1  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.1/dist/turnkey.linux-aarch64)   |
| v1.1.1  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.1.1/dist/turnkey.darwin-x86_64)   |
| v1.1.1  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.1.1/dist/turnkey.darwin-aarch64) |
| v1.0.5  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.5/dist/turnkey.linux-x86_64)     |
| v1.0.5  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.5/dist/turnkey.linux-aarch64)   |
| v1.0.5  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.5/dist/turnkey.darwin-x86_64)   |
| v1.0.5  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.5/dist/turnkey.darwin-aarch64) |
| v1.0.4  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.linux-x86_64)     |
| v1.0.4  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.linux-aarch64)   |
| v1.0.4  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.darwin-x86_64)   |
| v1.0.4  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.darwin-aarch64) |
| v1.0.3  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.3/dist/turnkey.linux-x86_64)     |
| v1.0.3  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.3/dist/turnkey.linux-aarch64)   |
| v1.0.3  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.3/dist/turnkey.darwin-x86_64)   |
| v1.0.3  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.3/dist/turnkey.darwin-aarch64) |
| v1.0.2  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.2/dist/turnkey.linux-x86_64)     |
| v1.0.2  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.2/dist/turnkey.linux-aarch64)   |
| v1.0.2  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.2/dist/turnkey.darwin-x86_64)   |
| v1.0.2  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.2/dist/turnkey.darwin-aarch64) |
| v1.0.1  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.1/dist/turnkey.linux-x86_64)     |
| v1.0.1  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.1/dist/turnkey.linux-aarch64)   |
| v1.0.1  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.1/dist/turnkey.darwin-x86_64)   |
| v1.0.1  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.1/dist/turnkey.darwin-aarch64) |
| v1.0.0  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.0/dist/turnkey.linux-x86_64)     |
| v1.0.0  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.0/dist/turnkey.linux-aarch64)   |
| v1.0.0  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v1.0.0/dist/turnkey.darwin-x86_64)   |
| v1.0.0  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v1.0.0/dist/turnkey.darwin-aarch64) |
| v0.3.4  | Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/v0.3.4/dist/turnkey.linux-x86_64)     |
| v0.3.4  | Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/v0.3.4/dist/turnkey.linux-aarch64)   |
| v0.3.4  | MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/v0.3.4/dist/turnkey.darwin-x86_64)   |
| v0.3.4  | MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/v0.3.4/dist/turnkey.darwin-aarch64) |

#### Git

```sh
git clone https://github.com/tkhq/tkcli
cd tkcli
# This installs in  ~/.local/bin; make sure this is in your $PATH!
make install
```

#### Brew

```sh
brew install tkhq/tap/turnkey
```

### Moderate Trust

These steps will allow you to prove that at least two Turnkey engineers
signed off on the produced binaries, signaling that they reproduced them from
source code and got identical results, in addition to our usual two-party code
review processes.

This minimizes a single point of trust (and failure) in our binary release
process.

See the [Reproducible Builds](https://reproducible-builds.org/) project for
more information on these practices.

We use git for all development, releases, and signing. Unfortunately git has no
native method for large file storage or multi-signature workflows so some git
add-ons are required.

To follow these steps please install [git-lfs][gl] and [git-sig][gs].

[gs]: https://codeberg.org/distrust/git-sig
[gl]: https://git-lfs.com

1. Clone repo

   ```sh
   git clone https://github.com/tkhq/tkcli
   cd tkcli
   ```

2. Review binary signatures

   ```sh
   git sig verify
   ```

   Note: See Trust section below for expected keys/signers

3. Install binary

   ```
   make install
   ```

### Zero Trust

If you intend to use the Turnkey CLI on a system you need to be able to trust
or for a high risk use case, we strongly recommend taking the time to hold us
accountable to the maximum degree you have resources and time for.

This protects not only you, but also protects our team. If many people are
checking our work for tampering it removes the incentive for someone malicious
to attempt to force one or more of us to tamper with the software.

1. Clone repo

   ```sh
   git clone https://github.com/tkhq/tkcli
   cd tkcli
   ```

2. Review source
   - Ideal: Review the entire supply chain is recommended for high risk uses
   - Minimal: review the "attest" "sign" and "verify" targets in the Makefile

3. Reproduce binaries

   ```sh
   make reproduce
   ```

   Note: See Trust section below for expected keys/signers

4. Install binaries

   ```sh
   make install
   ```

5. Upload signature

   While this step is totally optional, if you took the time to verify our
   binaries we would welcome you signing them and submitting your signature so
   we have public evidence third parties are checking our work.

   **NOTE**: this additionally uses Github's official CLI tool, [gh](https://github.com/cli/cli).

   ```sh
   gh repo fork
   git add dist/*
   git commit -m "add signature"
   git sig add
   git push origin main
   gh pr create
   ```

## Usage

### Generate a new API key

```sh
$ turnkey generate api-key --organization $ORGANIZATION_ID --key-name default
{
   "privateKeyFile": "/Users/andrew/Library/Application Support/turnkey/keys/default.private",
   "publicKey": "0236f17892a4649d97b2e4a4ad3c22d815e4e77848a0b8e4a5b0956ae4d6be382e",
   "publicKeyFile": "/Users/andrew/Library/Application Support/turnkey/keys/default.public"
}
```

### Add your public API key

As an authenticated user on the Turnkey dashboard, navigate to your user page by clicking on "User Details" in the user dropdown menu.

Click on "Create API keys and follow the prompts to add the generated public API key. You'll be required to authenticate with the same authenticator used during onboarding. After this succeeds, you should be all set to interact with our API.

#### Notes

- If you would like to manually copy your locally-stored public/private API key files (e.g. `default.public`, `default.private`), you will have to save the files without newlines (which occupy extra bytes). For example, for VIM, use `:set binary noeol` or `:set binary noendofline` before writing.
- Only P-256 keys (`API_KEY_CURVE_P256`) are currently supported.

### Create a Wallet

Wallets are collections of cryptographic key pairs typically used for sending and receiving digital assets. To create on, we need to provide a name:

```sh
turnkey wallets create --name default --key-name default
```

### Create an Ethereum account

To create a cryptographic key pair on our new Wallet, we neet to pass our desired address format:

```sh
turnkey wallets accounts create --wallet default --address-format ADDRESS_FORMAT_ETHEREUM --key-name default
```

This command will produce an Ethereum address (e.g. `0x08cb1216C95149DF66978b574E484869512CE2bF`) that we'll need to sign a transaction. You can see your new Wallet account with:

```sh
turnkey wallets accounts list --wallet default --key-name default
```

### Sign a Transaction

Now you can sign an Ethereum transaction with this new address with our [`sign_transaction` endpoint](https://docs.turnkey.com/api-reference/signing/sign-transaction). Make sure to replace the `unsignedTransaction` below with your own. You can use our [simple transaction generator](https://build.tx.xyz/) if you need a quick transaction for testing:

```json
turnkey request --path /public/v1/submit/sign_transaction --body '{
    "timestampMs": "'"$(date +%s)"'000",
    "type": "ACTIVITY_TYPE_SIGN_TRANSACTION_V2",
    "organizationId": "'"$ORGANIZATION_ID"'",
    "parameters": {
      "type": "TRANSACTION_TYPE_ETHEREUM",
      "signWith": "<Your Ethereum address>",
      "unsignedTransaction": "<Your Transaction>"
    }
}' --key-name default
```

### Next Steps

See the [official docs](https://docs.turnkey.com/sdks/cli#next-steps) for additional usage information and examples.

## Building

### Build for all platforms

```sh
make
```

### Build for one platform

```sh
make out/turnkey.linux-amd64
```

### Local build (for development only)

The following will drop a binary in `build/turnkey`:

```sh
make build-local
```

Note that you may need to do the following:

- Install `git-lfs`: https://git-lfs.com
- Setup: `git lfs install`

## Release

To release a new version of the CLI:

Determine the next version:

```sh
git tag | sort -n | tail -n5
```

Export your new version:

```sh
export VERSION=vX.Y.Z
```

Build the release artifacts:

```sh
make VERSION=$VERSION dist
```

Cut a new release branch:

```sh
git checkout -b release-$VERSION
```

Open a pull request, and once you have enough approvals, tag the release:

```sh
git tag -sa $VERSION -m "New release: $VERSION"
```

Finally, update the download table above, with links pointing to the new binaries.

Once the pull request is merged, ask your reviewer(s) to attest with `git sig`:

```sh
make reproduce

# If the reproduce command succeeds:
git sig add
```

Once enough signatures have been collected, the following command should succeed:

```sh
git sig verify --threshold 2
```

Finally, post the new release on Github with a changelog and update the Homebrew tap.

## Trust

### Process

You should never trust random binaries or code you find on the internet. Even
if it is from a reputable git identity, developers are phished all the time.

Supply chain attacks are becoming increasingly common in our industry and it
takes strong accountability to prevent them from happening.

The only way to be reasonably confident code was actually authored by the
people we think it was, is if that software is cryptographically signed by a
key only those individuals have access to.

Similarly if a company releases binaries, you have no idea if the machine that
compiled it is compromised or not, and no idea if the code in that binary
corresponds to the actual code in the repo that you or someone you trust
authored or reviewed.

To address both problems we take the following steps:

1. All commits are signed with keys that only exist on hardware security
   modules held by each engineer
2. All binaries are signed by the engineer that compiled them
3. Attesting engineers compile and sign binaries if they get the same hashes

### Signature Verification

To learn who signed the current release run:

`git sig verify --threshold 2`

Commits will be signed by at least one of the keys under the signers section
below.

Released binaries should be signed by at least two of them signifying
successful reproducible builds.

We encourage you to review the below keyoxide links and any available
web-of-trust for each key to ensure it is really owned by the person it claims
to be owned by.

### Signers

| Name             | PGP Fingerprint                                                                            |
| ---------------- | ------------------------------------------------------------------------------------------ |
| Andrew Min       | [DE05 0A45 1E6F AF94 C677 B58B 9361 DEC6 47A0 87BD](https://keyoxide.org/9361DEC647A087BD) |
| Arnaud Brousseau | [6870 5ACF 41E8 ECDE E292 5A42 4AAB 800C FFA3 065A](https://keyoxide.org/4AAB800CFFA3065A) |
| Keyan Zhang      | [0211 6F38 FB32 9E98 65A1 D08B 5880 CFD7 A7D9 5342](https://keyoxide.org/5880CFD7A7D95342) |
| Lance Vick       | [6B61 ECD7 6088 748C 7059 0D55 E90A 4013 36C8 AAA9](https://keyoxide.org/E90A401336C8AAA9) |
| Seán C McCord    | [39B2 095B 61DD 23EE E1BF 883A 8A1F 0484 90D2 3AFD](https://keyoxide.org/8A1F048490D23AFD) |
