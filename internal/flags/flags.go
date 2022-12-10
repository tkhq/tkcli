package flags

import "github.com/urfave/cli/v2"

func KeysFolder() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "keys-folder",
		Usage:    "Folder in which to write API public and private key files. Defaults to `~/.config/turnkey/keys`, or $XDG_CONFIG_HOME/turnkey/keys if XDG_CONFIG_HOME is set",
		Required: false,
	}
}

func CreateKeyName() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "name",
		Usage:    "Name of the API key. Will be used to create <name>.public and <name>.private. To write to stdout, use \"--name -\".",
		Required: true,
	}
}

func Key() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "key",
		Aliases:  []string{"k"},
		Usage:    "Private key to sign with. Provide a name to lookup the private key in your turnkey/keys directory (e.g. \"my_api_key\" will use \"~/.config/turnkey/keys/my_api_key.private\"), or a full path to a valid private key (e.g. \"/path/to/key.private\")",
		Required: true,
	}
}

func Message() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "message",
		Usage:    "Message to sign",
		Required: true,
	}
}

func Host() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "host",
		Usage:    "HTTP host _without_ protocol, e.g. api.domain.tld",
		Required: true,
	}
}

func Method() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "method",
		Usage:    "HTTP Method. Should be \"GET\" or \"POST\"",
		Required: true,
	}
}

func Path() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "path",
		Usage:    "Path, including the leading \"/\" and query string if any. For example: /api/v1/keys?curve=ed25519",
		Required: true,
	}
}

func Body() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "body",
		Usage:    "HTTP body, only relevant for POST requests. For example: {\"message\": \"Hello, world!\"}",
		Required: true,
	}
}
