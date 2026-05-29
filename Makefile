ODIR=bin

help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

linux: # Build the arm64 and x86/x64 binaries for Linux systems.
	@echo "Building arm64 binary"
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o $(ODIR)/soloterm_Linux_arm64 .

	@echo "Building x86/x64 binary"
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o $(ODIR)/soloterm_Linux_x86_64 .

mac: # Build the arm64 and x86/x64 binaries for Mac systems.
	@echo "Building arm64 binary"
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o $(ODIR)/soloterm_Darwin_arm64 .

	@echo "Building x86/x64 binary"
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o $(ODIR)/soloterm_Darwin_x86_64 .

windows: # Build the arm64 and x86/x64 executables for Windows systems.
	@echo "Building arm64 binary"
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o $(ODIR)/soloterm_Windows_arm64.exe .

	@echo "Building x86/x64 binary"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(ODIR)/soloterm_Windows_x86_64.exe .


all: # Builds soloterm for all systems.
	@echo "Building soloterm for all platforms..."

	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o $(ODIR)/soloterm_Darwin_arm64 .
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o $(ODIR)/soloterm_Darwin_x86_64 .
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o $(ODIR)/soloterm_Linux_arm64 .
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o $(ODIR)/soloterm_Linux_x86_64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o $(ODIR)/soloterm_Windows_arm64.exe .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(ODIR)/soloterm_Windows_x86_64.exe .

	@echo "Done."

.PHONY: clean

clean: # Removes any build binaries / executables from the bin folder.
	rm -f $(ODIR)/*
