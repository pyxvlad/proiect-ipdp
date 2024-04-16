{ pkgs ? import <nixpkgs> {} }:

let

	run = pkgs.writeShellScriptBin "run" ''
		#!/usr/bin/env bash

		${pkgs.templ}/bin/templ generate
		${pkgs.tailwindcss}/bin/tailwindcss -i css/tw.css -o css/style.css

		$(cd database; ${pkgs.sqlc}/bin/sqlc generate)

		echo "Starting server"

		${pkgs.go}/bin/go run .
	'';


	ctpTW = pkgs.buildNpmPackage rec {
		pname = "catppuccin-tailwindcss";
		version = "0.1.6";
		src = pkgs.fetchFromGitHub {
			owner = "catppuccin";
			repo = "tailwindcss";
			rev = "v${version}";
			hash = "sha256-ae5v9cB21Rs6O1+Y9QgoNr2e3Qio5MEXuh84dkgcg0A=";
		};
		npmDepsHash = "sha256-Q4b0WLUBLPfArieC7D0h/KfX/wpB/MgZIk0VzYIO1IQ=";
		npmPackFlags = [ "--ignore-scripts" ];
		NODE_OPTIONS = "--openssl-legacy-provider";

		meta = with pkgs.lib; {
			description = "A modern web UI for various torrent clients with a Node.js backend and React frontend";
			homepage = "https://flood.js.org";
			license = licenses.gpl3Only;
			maintainers = with maintainers; [ winter ];
		};	
	};

in
pkgs.mkShell {
	# Fix for delve
	hardeningDisable = [ "fortify" ];
	packages = with pkgs; [
		go
		templ
		tailwindcss

		golangci-lint
		sqlc

		run

		ctpTW
	];

	NODE_PATH = "${ctpTW}/lib/node_modules";
}

