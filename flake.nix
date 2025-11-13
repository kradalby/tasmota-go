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

        packages = {
          tasmota-cli = pkgs.buildGoModule {
            pname = "tasmota";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-qaS2PLxDCttfzJZrnz1d2Mk5oZ4a4uTZnQFjJe2CKek=";
            subPackages = [ "cmd/tasmota" ];
          };

          default = self.packages.${system}.tasmota-cli;
        };

        apps = {
          tasmota = {
            type = "app";
            program = "${self.packages.${system}.tasmota-cli}/bin/tasmota";
          };

          default = self.apps.${system}.tasmota;
        };
      }
    );
}
