# Turnkey CLI

![Go Build Status](https://github.com/tkhq/tkcli/actions/workflows/go-build.yml/badge.svg)

## Installation

We have multiple ways to install the CLI depending on your threat model.

Please check our work to whatever extent appropriate for your use case.

### Blind Trust

> :warning: Before you copy/paste, note that these are /low/ security options

If you are on an untrusted machine and are only evaluating our tools, we offer
easy low security install paths common in the industry.

Do note that any time you run an unverified binary off the internet you are
giving a third party full permission to execute any code they want on your
system. Github accounts, CDNs, and package repository accounts get compromised
all the time.

#### Download

| OS    | Architecture | Download                                             |
|-------|--------------|------------------------------------------------------|
| Linux | x86_64       | [turnkey.linux-x86_64](https://github.com/tkhq/tkcli/raw/main/dist/turnkey.linux-x86_64)    |
| Linux | aarch64      | [turnkey.linux-aarch64](https://github.com/tkhq/tkcli/raw/main/dist/turnkey.linux-aarch64)  |
| MacOS | x86_64       | [turnkey.darwin-x86_64](https://github.com/tkhq/tkcli/raw/main/dist/turnkey.darwin-x86_64)  |
| MacOS | aarch64      | [turnkey.darwin-aarch64](https://github.com/tkhq/tkcli/raw/main/dist/turnkey.darwin-aarch64)|

#### Git

```sh
git clone https://github.com/thkq/tkcli
cd tkcli
# This installs in  ~/.local/bin; make sure this is in your $PATH!
make install
```

#### Brew

```sh
brew tap tkhq/tkcli
brew install turnkey
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

    ```
    git clone https://github.com/tkhq/tkcli
    cd tkcli
    ```

2. Review binary signatures

    ```
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

    ```
    git clone https://github.com/tkhq/tkcli
    cd tkcli
    ```

2. Review source

    * Ideal: Review the entire supply chain is recommended for high risk uses
    * Minimal: review the "attest" "sign" and "verify" targets in the Makefile

3. Reproduce binaries

    ```
    make reproduce
    ```

    Note: See Trust section below for expected keys/signers

4. Install binaries

    ```
    make install
    ```

5. Upload signature

    While this step is totally optional, if you took the time to verify our
    binaries we would welcome you signing them and submitting your signature so
    we have public evidence third parties are checking our work.

    ```
    gh repo fork
    git add dist/*
    git commit -m "add signature"
    git sig add
    git push origin main
    gh pr create
    ```

## Usage

Create a new API key:

```sh
$ turnkey gen --name my-test-key
{
    "privateKeyFile": "/Users/rno/.config/turnkey/keys/my-test-key.private",
    "publicKeyFile": "/Users/rno/.config/turnkey/keys/my-test-key.public"
}
```

Make an API request:

```sh
$ turnkey request --host coordinator-beta.turnkey.io --path /api/v1/sign --body '{"payload": "hello from TKHQ"}' --key=my-test-key
{
    "result": "I am a teapot"
}
```

Create, but do not _post_ a request:

```sh
$ turnkey request --no-post --host coordinator-beta.turnkey.io --path /api/v1/sign --body '{"payload": "hello from TKHQ"}' --key=my-test-key
{
    "curlCommand": "curl -X POST -d'{\"payload\": \"hello from TKHQ\"}' -H'X-Stamp: eyJwdWJsaWNLZXkiOiIwM2JmMTYyNTc2ZWI4ZGZlY2YzM2Q5Mjc1ZDA5NTk1Mjg0ZjZjNGRmMGRiNjE1NmMzYzU4Mjc3Nzg4NmEwZWUwYWMiLCJzaWduYXR1cmUiOiIzMDQ0MDIyMDZiMmRlYmIwYjA3YmYwMDJlMjI1ZmQ4NTgzZjZmNGUxNGE5YTUxYWRiYWJjNDAyYzY5YTZlN2Q4N2ViNWNjMDgwMjIwMjE0ZTdkMGJlODFjMGYyNDEyOWE0MmNkZGFlOTUxYTBmZTViMGM1Mzc3YjM2NzZiOTUyNDgyNmYwODdhMWU4ZiIsInNjaGVtZSI6IlNJR05BVFVSRV9TQ0hFTUVfVEtfQVBJX1AyNTYifQ' -v 'https://coordinator-beta.turnkey.io/api/v1/sign'",
    "message": "{\"payload\": \"hello from TKHQ\"}",
    "stamp": "eyJwdWJsaWNLZXkiOiIwM2JmMTYyNTc2ZWI4ZGZlY2YzM2Q5Mjc1ZDA5NTk1Mjg0ZjZjNGRmMGRiNjE1NmMzYzU4Mjc3Nzg4NmEwZWUwYWMiLCJzaWduYXR1cmUiOiIzMDQ0MDIyMDZiMmRlYmIwYjA3YmYwMDJlMjI1ZmQ4NTgzZjZmNGUxNGE5YTUxYWRiYWJjNDAyYzY5YTZlN2Q4N2ViNWNjMDgwMjIwMjE0ZTdkMGJlODFjMGYyNDEyOWE0MmNkZGFlOTUxYTBmZTViMGM1Mzc3YjM2NzZiOTUyNDgyNmYwODdhMWU4ZiIsInNjaGVtZSI6IlNJR05BVFVSRV9TQ0hFTUVfVEtfQVBJX1AyNTYifQ"
}
```

## Building

### Build for all platforms
```
make
```

### Build for one platform

```
make out/turnkey.linux-amd64
```

## Release

To release a new version of the CLI:

```sh
$ make VERSION=vX.Y.Z dist
```

This will produce a new set of artifact in the `dist/` directory, along with a new manifest.

Open a pull request, and once you have enough approvals, tag the release:

```sh
$ git tag -sa vX.Y.Z -m "New release: X.Y.Z"
```

Finally, update the download table above, with links pointing to the new binaries.

Once the pull request is merged, ask your reviewer(s) to attest with `git sig`:

```sh
$ make reproduce

# If the reproduce command succeeds:
$ git sig add
```

Once enough signatures have been collected, the following command should succeed:

```sh
$ git sig verify --threshold 2
```

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

```git sig verify --threshold 2```

Commits will be signed by at least one of the keys under the signers section
below.

Released binaries should be signed by at least two of them signifying
successful reproducible builds.

We encourage you to review the below keyoxide links and any available
web-of-trust for each key to ensure it is really owned by the person it claims
to be owned by.

### Signers

| Name             | PGP Fingerprint                                                                          |
|------------------|------------------------------------------------------------------------------------------|
| Andrew Min       |[DE05 0A45 1E6F AF94 C677 B58B 9361 DEC6 47A0 87BD](https://keyoxide.org/9361DEC647A087BD)|
| Arnaud Brousseau |[6870 5ACF 41E8 ECDE E292 5A42 4AAB 800C FFA3 065A](https://keyoxide.org/4AAB800CFFA3065A)|
| Keyan Zhang      |[0211 6F38 FB32 9E98 65A1 D08B 5880 CFD7 A7D9 5342](https://keyoxide.org/5880CFD7A7D95342)|
| Lance Vick       |[6B61 ECD7 6088 748C 7059 0D55 E90A 4013 36C8 AAA9](https://keyoxide.org/E90A401336C8AAA9)|
| Seán C McCord    |[D9C8 71CC DEBD 2C7A 23D8 3041 798D DA52 5182 4DEA](https://keyoxide.org/798DDA5251824DEA)|
