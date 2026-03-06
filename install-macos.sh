#!/usr/bin/env bash
# install-macos.sh — Package Lamp as a native macOS .app
set -e

APP_NAME="Lamp"
BINARY_NAME="lamp"
BINARY_SRC="./$BINARY_NAME"
APP_DIR="./$APP_NAME.app"
CONTENTS="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS/MacOS"
RESOURCES_DIR="$CONTENTS/Resources"

# Always rebuild
echo "Building $BINARY_NAME..."
go build -ldflags="-s -w" -o "$BINARY_NAME" .

# Generate icns from pixel art
echo "Generating icon..."
sips -z 1024 1024 icon/lamp_pixel_art_32x32.png --out /tmp/lamp_1024.png
mkdir -p /tmp/lamp.iconset
sips -z 16 16     /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_16x16.png
sips -z 32 32     /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_16x16@2x.png
sips -z 32 32     /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_32x32.png
sips -z 64 64     /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_32x32@2x.png
sips -z 128 128   /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_128x128.png
sips -z 256 256   /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_128x128@2x.png
sips -z 256 256   /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_256x256.png
sips -z 512 512   /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_256x256@2x.png
sips -z 512 512   /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_512x512.png
sips -z 1024 1024 /tmp/lamp_1024.png --out /tmp/lamp.iconset/icon_512x512@2x.png
iconutil -c icns /tmp/lamp.iconset -o lamp.icns
rm -rf /tmp/lamp.iconset /tmp/lamp_1024.png

echo "Creating $APP_NAME.app bundle..."
rm -rf "$APP_DIR"
mkdir -p "$MACOS_DIR"
mkdir -p "$RESOURCES_DIR"

cp "$BINARY_SRC" "$MACOS_DIR/$BINARY_NAME"
chmod +x "$MACOS_DIR/$BINARY_NAME"

if [ -f "./lamp.icns" ]; then
  cp "./lamp.icns" "$RESOURCES_DIR/lamp.icns"
  ICON_XML="
  <key>CFBundleIconFile</key>
  <string>lamp</string>"
else
  echo "No lamp.icns found — skipping icon."
  ICON_XML=""
fi

cat > "$CONTENTS/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key>
  <string>${APP_NAME}</string>
  <key>CFBundleDisplayName</key>
  <string>${APP_NAME}</string>
  <key>CFBundleIdentifier</key>
  <string>com.yourname.lamp</string>
  <key>CFBundleVersion</key>
  <string>1.0.0</string>
  <key>CFBundleShortVersionString</key>
  <string>1.0.0</string>
  <key>CFBundleExecutable</key>
  <string>${BINARY_NAME}</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>NSHighResolutionCapable</key>
  <true/>
  <key>NSPrincipalClass</key>
  <string>NSApplication</string>
  <key>NSSupportsAutomaticGraphicsSwitching</key>
  <true/>
  <key>LSMinimumSystemVersion</key>
  <string>10.13</string>
  <key>LSUIElement</key>
  <false/>${ICON_XML}
</dict>
</plist>
PLIST

xattr -cr "$APP_DIR" 2>/dev/null || true

echo "Installing to /Applications..."
rm -rf "/Applications/$APP_NAME.app"
cp -r "$APP_DIR" "/Applications/$APP_NAME.app"

echo ""
echo "Done! $APP_NAME.app installed to /Applications."