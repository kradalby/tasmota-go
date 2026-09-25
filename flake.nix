{
  description = "Tasmota Go library";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-checks.url = "github:kradalby/flake-checks";
    flake-checks.inputs.nixpkgs.follows = "nixpkgs";
    flake-checks.inputs.flake-utils.follows = "flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      flake-checks,
    }:
    {
      overlays.default = final: prev: {
        tasmota = self.packages.${prev.stdenv.hostPlatform.system}.default;
      };
    }
    // flake-utils.lib.eachDefaultSystem (
      system:
      let
        # Go 1.27. Note that bare `pkgs.go` / `pkgs.buildGoModule` still resolve
        # to 1.26 in nixpkgs-unstable, so the latest toolchain has to be named
        # explicitly. `go_latest` / `buildGoLatestModule` are used rather than
        # `go_1_27` / `buildGo127Module` so this keeps tracking the newest Go
        # without another edit here.
        #
        # Rebuild the dev tools nixpkgs still builds against its default Go.
        # This is load-bearing, not cosmetic: goimports built with 1.26 rejects
        # the generic method in status.go outright --
        # `method must have no type parameters` -- which would fail the
        # `formatting` check on valid 1.27 source. golangci-lint and gopls
        # already track go_latest upstream, so they need no override.
        goToolsOverlay = _: prev: {
          gofumpt = prev.gofumpt.override { buildGoModule = prev.buildGoLatestModule; };
          delve = prev.delve.override { buildGoModule = prev.buildGoLatestModule; };
          # goimports ships wrapped with a `go` on PATH. That `go` must be at
          # least the go.mod directive, or GOTOOLCHAIN=auto tries to fetch a
          # toolchain from inside the network-less treefmt sandbox.
          gotools = prev.gotools.override {
            buildGoModule = prev.buildGoLatestModule;
            go = prev.go_latest;
          };
        };

        pkgs = import nixpkgs {
          inherit system;
          overlays = [ goToolsOverlay ];
        };
        fc = flake-checks.lib;
        common = {
          inherit pkgs;
          root = ./.;
          pname = "tasmota-go";
          version = "0.0.1";
          vendorHash = "sha256-qaS2PLxDCttfzJZrnz1d2Mk5oZ4a4uTZnQFjJe2CKek=";
          goPkg = pkgs.go_latest;
          goRace = true;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_latest
            golangci-lint
            gopls
            gofumpt
            delve
            prek
            nixfmt
          ];

          shellHook = ''
            # Never fetch a toolchain. A go.mod ahead of nixpkgs' Go must be a
            # clear error, not a silent download from go.dev outside the store.
            #
            # CGO_ENABLED is deliberately left alone: the CI coverage job runs
            # `go test -race` through `nix develop`, and -race requires cgo.
            export GOTOOLCHAIN=local
          '';
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
