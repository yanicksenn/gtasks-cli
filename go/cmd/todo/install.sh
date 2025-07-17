#!/bin/bash

# This script builds and installs the 'todo' binary.

set -e

BINARY_NAME="todo"
INSTALL_DIR="/usr/local/bin"
INSTALL_PATH="$INSTALL_DIR/$BINARY_NAME"
BAZEL_TARGET="//go/cmd/todo:todo"

echo "Building $BINARY_NAME..."
bazel build $BAZEL_TARGET >/dev/null

EXECUTION_ROOT=$(bazel info execution_root)
BINARY_REL_PATH=$(bazel cquery --output=files ${BAZEL_TARGET})
BUILT_BINARY_PATH="$EXECUTION_ROOT/$BINARY_REL_PATH"

if [ ! -f "$BUILT_BINARY_PATH" ]; then
    echo "Error: Built binary not found at '$BUILT_BINARY_PATH'"
    exit 1
fi

echo "Installing $BUILT_BINARY_PATH in $INSTALL_DIR..."
sudo cp -f "$BUILT_BINARY_PATH" "$INSTALL_DIR"

echo "Installation successful. You can now use '$BINARY_NAME'."
