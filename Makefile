.PHONY: init
init:
	go install github.com/google/wire/cmd/wire@latest

.PHONY: wire
wire:
	wire gen -output_file_prefix app_ ./internal/app
