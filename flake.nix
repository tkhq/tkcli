{
  description = "tkcli devshell";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.05";
    nixpkgs-go.url = "github:nixos/nixpkgs/nixos-23.11";
    #nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; # localstack is broken right now (2023-01-25) in unstable, due to a missing dependency
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, nixpkgs-go, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        goPkgs = nixpkgs-go.legacyPackages.${system};
        go = goPkgs.go_1_21;
        golangci-lint = goPkgs.golangci-lint;

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

        tkbuild = pkgs.writeScriptBin "build" ''
          #!/bin/sh
          pushd $(git rev-parse --show-toplevel)/src
          ${go}/bin/go build -o $(go env GOPATH)/bin/turnkey ./cmd/turnkey
          ${go}/bin/go build -o ../out/turnkey.linux-x86_64 ./cmd/turnkey # hack for local CLI go test
        '';

        tklint = pkgs.writeScriptBin "lint" ''
          #!/bin/sh
          pushd $(git rev-parse --show-toplevel)/src
          ${go}/bin/go mod tidy
          ${pkgs.gofumpt}/bin/gofumpt -w ./cmd/*
          ${gci}/bin/gci write --skip-generated -s standard -s default -s "Prefix(github.com/tkhq)" .
          ${golangci-lint}/bin/golangci-lint run ./...
          ${go}/bin/go build -o ../out/turnkey.linux-x86_64 ./cmd/turnkey # hack for local CLI go test
          ${go}/bin/go test -v ./...
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            bashInteractive
            envsubst
            gci
            gofumpt
            golangci-lint
            go
            go-swagger
            go-tools
            tkbuild
            tklint
          ];
        };
      });
}
