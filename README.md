# Turnkey CLI

## Building the CLI

We use [GoReleaser](https://goreleaser.com/) to build and release our binaries.

To build locally:
```
$ goreleaser release --snapshot --rm-dist
```

To release:
```
# Generate a Github token with "write:packages"
# ==> https://github.com/settings/tokens
$ export GITHUB_TOKEN=<your token>
$ git tag -a vx.y.z -m "New Release: x.y.z"
$ git push origin vx.y.z
$ goreleaser release --rm-dist
```

## Installing the CLI

```
brew tap tkhq/taps
brew install tk
```

## Usage

Create a new API key:

```sh
./build/tk gen --name rno
Creating /Users/rno/.tk/rno.public
Creating /Users/rno/.tk/rno.private
```

Sign a request:

```sh
./build/tk approve-request --method POST --path /api/v1/sign --body '{"payload": "hello from TKHQ"}' --key=rno
Raw signature: 3046022100a99781a6b1d7ff7c4ce3951ded09a7757c74f1c6d7c7e1a2e617ac2921d74674022100f75d167abe426eb8f89884afe5e864cb965c6370611566f50b46690209b3a95b
Approval header: X-Approved-By-035acbc8b7751b7703736ae16cb22112451372f7b77717bbecdfa8300d4038432: 3046022100a99781a6b1d7ff7c4ce3951ded09a7757c74f1c6d7c7e1a2e617ac2921d74674022100f75d167abe426eb8f89884afe5e864cb965c6370611566f50b46690209b3a95b
--------
To make this request with curl:
        curl -X POST -d {"payload": "hello from TKHQ"} -H'X-Approved-By-035acbc8b7751b7703736ae16cb22112451372f7b77717bbecdfa8300d4038432: 3046022100a99781a6b1d7ff7c4ce3951ded09a7757c74f1c6d7c7e1a2e617ac2921d74674022100f75d167abe426eb8f89884afe5e864cb965c6370611566f50b46690209b3a95b' -v 'https://api.turnkey.io/api/v1/sign'
--------
```
