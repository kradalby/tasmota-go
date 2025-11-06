{
  description = "Tasmota Go library";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_25
            golangci-lint
            gopls
            gotools
            go-tools
            delve
          ];
        };

        packages.default = pkgs.buildGoModule {
          pname = "tasmota-go";
          version = "0.1.0";
          src = ./.;
          vendorHash = null;
        };
      }
    );
}
