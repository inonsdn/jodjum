#!/usr/bin/env bash
set -euo pipefail

# Create a new module in server/internal following the user/ template.
# Usage: ./scripts/new-module.sh <module-name>
# Example: ./scripts/new-module.sh product

if [ $# -ne 1 ]; then
	echo "Usage: $0 <module-name>" >&2
	echo "Example: $0 product" >&2
	exit 1
fi

RAW="$1"

# Validate: letters/digits only, must start with a letter. Case-insensitive so
# "Product", "PRODUCT", and "product" are all accepted.
if ! [[ "$RAW" =~ ^[A-Za-z][A-Za-z0-9]*$ ]]; then
	echo "Error: module name must be letters/digits and start with a letter (got '$RAW')" >&2
	exit 1
fi

# Normalize regardless of the input casing:
#   MODULE (package + folder name) -> all lowercase   (Product -> product)
#   TITLE  (exported type prefix)  -> capitalized      (product -> Product)
MODULE="$(tr '[:upper:]' '[:lower:]' <<< "$RAW")"
TITLE="$(tr '[:lower:]' '[:upper:]' <<< "${MODULE:0:1}")${MODULE:1}"

# Resolve internal dir relative to this script (server/internal)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INTERNAL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)/internal"
MODULE_DIR="$INTERNAL_DIR/$MODULE"

if [ -d "$MODULE_DIR" ]; then
	echo "Error: module '$MODULE' already exists at $MODULE_DIR" >&2
	exit 1
fi

mkdir -p "$MODULE_DIR"

cat > "$MODULE_DIR/repo.go" <<EOF
package $MODULE

import "github.com/jackc/pgx/v5/pgxpool"

type ${TITLE}Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *${TITLE}Repo {
	return &${TITLE}Repo{
		db: db,
	}
}
EOF

cat > "$MODULE_DIR/services.go" <<EOF
package $MODULE

type ${TITLE}Service struct {
	${MODULE}Repo *${TITLE}Repo
}

func NewService(${MODULE}Repo *${TITLE}Repo) *${TITLE}Service {
	return &${TITLE}Service{
		${MODULE}Repo: ${MODULE}Repo,
	}
}
EOF

cat > "$MODULE_DIR/handler.go" <<EOF
package $MODULE

type ${TITLE}Handler struct {
	service *${TITLE}Service
}

func NewHandler(service *${TITLE}Service) *${TITLE}Handler {
	return &${TITLE}Handler{
		service: service,
	}
}
EOF

cat > "$MODULE_DIR/routes.go" <<EOF
package $MODULE

import (
	"net/http"
	"server/internal/types"
)

func RegisterRoutes(router *http.ServeMux, middleware types.Middleware, handler *${TITLE}Handler) {

}
EOF

echo "Created module '$MODULE' at $MODULE_DIR:"
echo "  - repo.go     (${TITLE}Repo)"
echo "  - services.go (${TITLE}Service)"
echo "  - handler.go  (${TITLE}Handler)"
echo "  - routes.go   (RegisterRoutes)"
