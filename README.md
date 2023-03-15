# Turnkey CLI

![Go Build Status](https://github.com/tkhq/tkcli/actions/workflows/go-build.yml/badge.svg)

## Building the CLI

If you simply want to build a binary locally, run
```sh
$ make build/turnkey
```

We use [GoReleaser](https://goreleaser.com/) to build and release our binaries.

To build a release locally:
```
$ goreleaser release --snapshot --rm-dist
```

To release:
```
# Generate a Github token with "write:packages"
# ==> https://github.com/settings/tokens
$ export GITHUB_TOKEN=<your token>
$ git tag -s -a vx.y.z -m "New release: x.y.z" # create a signed tag
$ git tag -v vx.y.z # verify the tag
$ git push origin vx.y.z
$ goreleaser release --rm-dist
```

## Installing the CLI

```
brew tap tkhq/tap
brew install turnkey
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
