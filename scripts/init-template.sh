#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in a git repository
if [ ! -d .git ]; then
    print_error "This script must be run from the root of a git repository"
    exit 1
fi

# Check if this is the template repository
if [ -d "scripts" ] && [ -f "scripts/init-template.sh" ]; then
    print_info "Detected template repository structure"
else
    print_error "This doesn't appear to be the go-cli-template repository"
    exit 1
fi

print_info "Welcome to the Go CLI Template Initialization Script!"
echo ""

# Get project information
read -p "Enter your project name (e.g., mycli): " PROJECT_NAME
read -p "Enter your GitHub username: " GITHUB_USERNAME
read -p "Enter a short description of your project: " SHORT_DESC
read -p "Enter a longer description of your project: " LONG_DESC

# Convert project name to uppercase for env vars
PROJECT_NAME_UPPER=$(echo "$PROJECT_NAME" | tr '[:lower:]' '[:upper:]')

print_info "Initializing project with the following settings:"
echo "  Project Name: $PROJECT_NAME"
echo "  GitHub Username: $GITHUB_USERNAME"
echo "  Short Description: $SHORT_DESC"
echo ""

# Create temporary directory for processing
TEMP_DIR=$(mktemp -d)
print_info "Working in temporary directory: $TEMP_DIR"

# Function to replace placeholders in a file
replace_placeholders() {
    local file=$1
    sed -i.bak \
        -e "s/{{\.ProjectName}}/$PROJECT_NAME/g" \
        -e "s/{{\.GitHubUsername}}/$GITHUB_USERNAME/g" \
        -e "s/{{\.ShortDescription}}/$SHORT_DESC/g" \
        -e "s/{{\.LongDescription}}/$LONG_DESC/g" \
        -e "s/{{\.PROJECTNAME}}/$PROJECT_NAME_UPPER/g" \
        "$file" && rm "${file}.bak"
}

# Process all files
print_info "Processing template files..."

# Find all files (excluding .git directory)
find . -type f -not -path "./.git/*" -not -path "./scripts/*" | while read -r file; do
    # Skip binary files
    if file "$file" | grep -q "text"; then
        print_info "Processing: $file"
        replace_placeholders "$file"
    fi
done

# Rename directories and files with template placeholders
print_info "Renaming template directories and files..."

# Rename the main binary directory
if [ -d "cmd/{{.ProjectName}}" ]; then
    mv "cmd/{{.ProjectName}}" "cmd/$PROJECT_NAME"
    print_info "Renamed cmd/{{.ProjectName}} to cmd/$PROJECT_NAME"
fi

# Update .gitignore entries
sed -i.bak "s/{{\.ProjectName}}/$PROJECT_NAME/g" .gitignore && rm .gitignore.bak

# Initialize go module with the correct name
print_info "Reinitializing Go module..."
rm go.mod go.sum 2>/dev/null || true
go mod init "github.com/$GITHUB_USERNAME/$PROJECT_NAME"
go mod tidy

# Remove the initialization script and update README
print_info "Cleaning up template files..."
rm -rf scripts/init-template.sh

# Update README to remove template section
sed -i.bak '/## Using This Template/,/- `{{\.PROJECTNAME}}`/d' README.md && rm README.md.bak

# Create initial commit
print_info "Creating initial commit..."
git add -A
git commit -m "Initialize $PROJECT_NAME from go-cli-template"

print_info "✅ Template initialization complete!"
echo ""
print_info "Next steps:"
echo "  1. Update the LICENSE file with your preferred license"
echo "  2. Update the README.md with project-specific information"
echo "  3. Run 'task dev' to start development with hot reload"
echo "  4. Run 'task test' to run the example tests"
echo "  5. Start building your CLI by modifying internal/cli/example.go"
echo ""
print_info "Happy coding! 🚀"