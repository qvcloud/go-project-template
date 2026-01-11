#!/bin/bash
set -e

# Load .env if exists
if [ -f .env ]; then
  # Use set -a to export variables automatically
  set -a
  source .env
  set +a
fi

# Default values. APP_NAME will read from env if set.
APP_NAME=${1:-${APP_NAME:-"go-project-template"}}
VERSION=${2:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
OUTPUT_DIR="dist"
MAIN_FILE="cmd/main.go"

# Build Flags
COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "none")
BUILD_TIME=$(date +'%Y-%m-%d_%T')
GO_VERSION=$(go env GOVERSION)
VERSION_PKG="github.com/qvcloud/gopkg/version"

LDFLAGS="-s -w"
LDFLAGS+=" -X 'main.appName=${APP_NAME}'"
LDFLAGS+=" -X '${VERSION_PKG}.Version=${VERSION}'"
LDFLAGS+=" -X '${VERSION_PKG}.Commit=${COMMIT}'"
LDFLAGS+=" -X '${VERSION_PKG}.Build=${BUILD_TIME}'"
LDFLAGS+=" -X '${VERSION_PKG}.Go=${GO_VERSION}'"

# Environment Variables with defaults for Linux build (common for containers)
# Allow overriding from outside
export CGO_ENABLED=${CGO_ENABLED:-0}
export GOOS=${GOOS:-linux}
export GOARCH=${GOARCH:-amd64}

echo "Building ${APP_NAME}..."
echo "  Version: ${VERSION}"
echo "  Commit:  ${COMMIT}"
echo "  OS/Arch: ${GOOS}/${GOARCH}"
echo "  Output:  ${OUTPUT_DIR}/${APP_NAME}"

rm -rf "${OUTPUT_DIR:?}/${APP_NAME}"
mkdir -p "${OUTPUT_DIR}"

go build -ldflags "${LDFLAGS}" -o "${OUTPUT_DIR}/${APP_NAME}" "${MAIN_FILE}"

echo "Build success!"
