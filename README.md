# Turnkey CLI

![Go Build Status](https://github.com/tkhq/tkcli/actions/workflows/go-build.yml/badge.svg)

## Building the CLI

If you simply want to build binaries locally, run
```
$ make
```

To release:
```
make dist
```

## Installing the CLI

1. Clone repo
    ```
    git clone https://github.com/tkhq/tkcli
    cd tkcli
    ```

2. Review source
    * Ideal: Review of the entire supply chain is recommended for high risk uses
    * Minimal: review the "attest" "sign" and "verify" targets in the Makefile

3. Review binaries
    * Ideal: ```make attest sign```
        * Reproduce and sign binaries yourself
	* We welcome PRs with external verification signatures
    * Recommended: ```make attest```
        * Prove published source code matches pubished binaries
    * Minimal: ```make verify```
        * Prove multiple people signed binaries
        * Ensure signatures are by people whose reputations you trust

4. Install binary

    Replace "linux" and "amd64" with your preferred architecture:

    ```
    mkdir -p ~/.local/bin
    cp tkcli/dist/turnkey.linux.amd64 ~/.local/bin/turnkey
    chmod +x ~/.local/bin/turnkey
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

Sign a request:

```sh
$ turnkey approve-request --host api.turnkey.io --path /api/v1/sign --body '{"payload": "hello from TKHQ"}' --key=my-test-key
{
    "curlCommand": "curl -X POST -d'{\"payload\": \"hello from TKHQ\"}' -H'X-Stamp: eyJwdWJsaWNLZXkiOiIwM2JmMTYyNTc2ZWI4ZGZlY2YzM2Q5Mjc1ZDA5NTk1Mjg0ZjZjNGRmMGRiNjE1NmMzYzU4Mjc3Nzg4NmEwZWUwYWMiLCJzaWduYXR1cmUiOiIzMDQ0MDIyMDZiMmRlYmIwYjA3YmYwMDJlMjI1ZmQ4NTgzZjZmNGUxNGE5YTUxYWRiYWJjNDAyYzY5YTZlN2Q4N2ViNWNjMDgwMjIwMjE0ZTdkMGJlODFjMGYyNDEyOWE0MmNkZGFlOTUxYTBmZTViMGM1Mzc3YjM2NzZiOTUyNDgyNmYwODdhMWU4ZiIsInNjaGVtZSI6IlNJR05BVFVSRV9TQ0hFTUVfVEtfQVBJX1AyNTYifQ' -v 'https://api.turnkey.io/api/v1/sign'",
    "message": "{\"payload\": \"hello from TKHQ\"}",
    "stamp": "eyJwdWJsaWNLZXkiOiIwM2JmMTYyNTc2ZWI4ZGZlY2YzM2Q5Mjc1ZDA5NTk1Mjg0ZjZjNGRmMGRiNjE1NmMzYzU4Mjc3Nzg4NmEwZWUwYWMiLCJzaWduYXR1cmUiOiIzMDQ0MDIyMDZiMmRlYmIwYjA3YmYwMDJlMjI1ZmQ4NTgzZjZmNGUxNGE5YTUxYWRiYWJjNDAyYzY5YTZlN2Q4N2ViNWNjMDgwMjIwMjE0ZTdkMGJlODFjMGYyNDEyOWE0MmNkZGFlOTUxYTBmZTViMGM1Mzc3YjM2NzZiOTUyNDgyNmYwODdhMWU4ZiIsInNjaGVtZSI6IlNJR05BVFVSRV9TQ0hFTUVfVEtfQVBJX1AyNTYifQ"
}
```
