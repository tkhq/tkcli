{
  description = "tkcli devshell";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        gci = pkgs.buildGoModule rec {
           name = "gci";
           src = pkgs.fetchFromGitHub {
              owner = "daixiang0";
              repo = "gci";
              rev = "v0.10.1";
              sha256 = "sha256-/YR61lovuYw+GEeXIgvyPbesz2epmQVmSLWjWwKT4Ag=";
           };

           # Switch to fake vendor sha for upgrades:
           #vendorSha256 = pkgs.lib.fakeSha256;
           vendorSha256 = "sha256-g7htGfU6C2rzfu8hAn6SGr0ZRwB8ZzSf9CgHYmdupE8=";
        };

        go = pkgs.go_1_21;

        tkbuild = pkgs.writeScriptBin "build" ''
          #!/bin/sh
          pushd $(git rev-parse --show-toplevel)/src
          ${go}/bin/go build -o $(${go}/bin/go env GOPATH)/bin/turnkey ./cmd/turnkey
          ${go}/bin/go build -o ../out/turnkey.linux-x86_64 ./cmd/turnkey # hack for local CLI go test
        '';

        tklint = pkgs.writeScriptBin "lint" ''
          #!/bin/sh
          pushd $(git rev-parse --show-toplevel)/src
          ${go}/bin/go mod tidy
          ${pkgs.gofumpt}/bin/gofumpt -w *.go ./cmd/*
          ${gci}/bin/gci write --skip-generated -s standard -s default -s "Prefix(github.com/tkhq)" .
          ${pkgs.golangci-lint}/bin/golangci-lint run ./...
          ${go}/bin/go build -o ../out/turnkey.linux-x86_64 ./cmd/turnkey # hack for local CLI go test
          ${go}/bin/go test -v ./...
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.bashInteractive
            pkgs.envsubst
            gci
            pkgs.gofumpt
            pkgs.golangci-lint
            go
            pkgs.go-swagger
            pkgs.go-tools
            tkbuild
            tklint
          ];
        };
      });
}
