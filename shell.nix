{ pkgs ? import <nixpkgs> {} }:

let

	run = pkgs.writeShellScriptBin "run" ''
		#!/usr/bin/env bash

		${pkgs.templ}/bin/templ generate
		${pkgs.tailwindcss}/bin/tailwindcss -i css/tw.css -o css/style.css

		echo "Starting server"

		${pkgs.go}/bin/go run .
	'';


in
pkgs.mkShell {
	# Fix for delve
	hardeningDisable = [ "fortify" ];
	packages = with pkgs; [
		go
		templ
		tailwindcss

		golangci-lint

		run
	];
}

