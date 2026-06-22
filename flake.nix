{
  description = "Tasmota Go library";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-checks.url = "github:kradalby/flake-checks";
    flake-checks.inputs.nixpkgs.follows = "nixpkgs";
    flake-checks.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, flake-checks }:
    {
      overlays.default = final: prev: {
        tasmota = self.packages.${prev.system}.default;
      };
    } // flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        fc = flake-checks.lib;
        common = {
          inherit pkgs;
          root = ./.;
          pname = "tasmota-go";
          version = "0.0.1";
          vendorHash = "sha256-qaS2PLxDCttfzJZrnz1d2Mk5oZ4a4uTZnQFjJe2CKek=";
          goPkg = pkgs.go_1_26;
          goRace = true;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_26
            golangci-lint
            gopls
            gofumpt
            delve
            prek
            nixpkgs-fmt
          ];
        };

        packages = {
          tasmota-cli = fc.goBuild (common // { subPackages = [ "cmd/tasmota" ]; });

          default = self.packages.${system}.tasmota-cli;
        };

        apps = {
          tasmota = {
            type = "app";
            program = "${self.packages.${system}.tasmota-cli}/bin/tasmota";
          };

          default = self.apps.${system}.tasmota;
        };

        formatter = fc.formatter common;

        checks = {
          build = fc.goBuild common;
          gotest = fc.goTest (common // { goRace = false; });
          gotest-race = fc.goTest (common // { goRace = true; });
          golangci-lint = fc.goLint common;
          formatting = fc.goFormat common;
        };
      }
    );
}
