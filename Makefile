BINARY_NAME := lamp
BUILD_FLAGS := -ldflags="-s -w"

.PHONY: build install-linux install-macos clean

build:
	echo "Building $BINARY_NAME..."
	go build -ldflags="-s -w" -o "$BINARY_NAME" .

install-macos: build
	bash install-macos.sh
	codesign --force --deep --sign - Lamp.app
	sudo cp -r Lamp.app /Applications/
	sudo xattr -cr /Applications/Lamp.app
	@echo "Clearing app cache..."
	/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -f /Applications/Lamp.app

clean:
	rm -f $(BINARY_NAME)
	rm -rf Lamp.app