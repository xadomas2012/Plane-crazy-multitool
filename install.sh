#!/bin/sh
set -eu

APP_NAME="PC-Multitool"
BIN_NAME="PC-Gear-Calculator"

INSTALL_DIR="$HOME/.local/opt/$APP_NAME"
BIN_DIR="$HOME/.local/bin"
DESKTOP_DIR="$HOME/.local/share/applications"

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
SOURCE_BIN="$SCRIPT_DIR/$BIN_NAME"

if [ ! -f "$SOURCE_BIN" ]; then
    echo "Error: $BIN_NAME not found next to install.sh"
    exit 1
fi

echo "Installing $APP_NAME..."

mkdir -p "$INSTALL_DIR"
mkdir -p "$BIN_DIR"
mkdir -p "$DESKTOP_DIR"

cp "$SOURCE_BIN" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

ln -sf "$INSTALL_DIR/$BIN_NAME" "$BIN_DIR/pc-multitool"

cat > "$DESKTOP_DIR/pc-multitool.desktop" <<DESKTOP
[Desktop Entry]
Type=Application
Name=PC Multitool
Comment=Plane Crazy engineering calculator and multitool
Exec=$BIN_DIR/pc-multitool
Terminal=true
Categories=Utility;
DESKTOP

echo
echo "PC Multitool installed successfully."
echo
echo "Installed to:"
echo "  $INSTALL_DIR/$BIN_NAME"
echo
echo "Launcher:"
echo "  $BIN_DIR/pc-multitool"
echo

case ":${PATH:-}:" in
    *:"$BIN_DIR":*)
        ;;
    *)
        echo "Note: $BIN_DIR is not currently in PATH."
        echo "Add this to your shell config:"
        echo
        echo 'export PATH="$HOME/.local/bin:$PATH"'
        echo
        ;;
esac
