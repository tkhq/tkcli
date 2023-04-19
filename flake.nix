{
  description = "tkcli devshell";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-22.11";
    #nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; # localstack is broken right now (2023-01-25) in unstable, due to a missing dependency
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        tkbuild = pkgs.writeScriptBin "build" ''
          #!/bin/sh
          pushd $(git rev-parse --show-toplevel)/src
          ${pkgs.go}/bin/go install
        '';

        tklint = pkgs.writeScriptBin "lint" ''
          #!/bin/sh
          pushd $(git rev-parse --show-toplevel)/src
          ${pkgs.gofumpt}/bin/gofumpt -w *.go ./cmd/*
          ${pkgs.golangci-lint}/bin/golangci-lint run ./...
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            bashInteractive
            envsubst
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
