{
  description = "airthings-exporter";

  inputs = {
    nixpkgs = {
      type = "github";
      owner = "NixOS";
      repo = "nixpkgs";
      rev = "f863eeff847e6c462fc78c5605cf8f0dd567c389";
    };
    gomod2nix = {
      type = "github";
      owner = "nix-community";
      repo = "gomod2nix";
      rev = "1c6fd4e862bf2f249c9114ad625c64c6c29a8a08";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    flake-utils = {
      type = "github";
      owner = "numtide";
      repo = "flake-utils";
      rev = "b1d9ab70662946ef0850d488da1c9019f3a9752a";
    };
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      gomod2nix,
      ...
    }:
    let
      utils = flake-utils;
    in
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        mkGoEnv = gomod2nix.legacyPackages.${system}.mkGoEnv;
        goEnv = mkGoEnv { pwd = ./.; };
        buildGoApplication = gomod2nix.legacyPackages.${system}.buildGoApplication;
        pname = "airthings-exporter";
      in
      {
        packages.default = buildGoApplication {
          pname = pname;
          version = "latest";
          pwd = ./.;
          src = ./.;
          modules = ./gomod2nix.toml;
        };

        devShells = {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go_1_22
              golangci-lint
              gomod2nix.packages.${system}.default
              goEnv
            ];
          };
        };
      }
    );
}
