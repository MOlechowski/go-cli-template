#!/bin/bash

# Self-Hosted GitHub Runner Setup Script
# For macOS ARM64 (Apple Silicon)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running on macOS ARM64
if [[ "$(uname)" != "Darwin" ]]; then
    print_error "This script is designed for macOS only"
    exit 1
fi

if [[ "$(uname -m)" != "arm64" ]]; then
    print_error "This script is designed for Apple Silicon (ARM64) Macs only"
    exit 1
fi

print_status "Setting up GitHub Actions self-hosted runner..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    print_warning "Homebrew not found. Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    
    # Add Homebrew to PATH
    echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
    eval "$(/opt/homebrew/bin/brew shellenv)"
else
    print_status "Homebrew found"
fi

# Install required tools
print_status "Installing required tools..."

# Install Go
if ! command -v go &> /dev/null; then
    print_status "Installing Go..."
    brew install go
else
    print_status "Go already installed: $(go version)"
fi

# Install Task
if ! command -v task &> /dev/null; then
    print_status "Installing Task..."
    brew install go-task/tap/go-task
else
    print_status "Task already installed: $(task --version)"
fi

# Install golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    print_status "Installing golangci-lint..."
    brew install golangci-lint
else
    print_status "golangci-lint already installed: $(golangci-lint --version | head -1)"
fi

# Create actions-runner directory
RUNNER_DIR="$HOME/actions-runner"
if [[ -d "$RUNNER_DIR" ]]; then
    print_warning "Runner directory already exists at $RUNNER_DIR"
    read -p "Do you want to remove it and start fresh? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf "$RUNNER_DIR"
        print_status "Removed existing runner directory"
    else
        print_error "Aborting setup. Please remove $RUNNER_DIR manually if needed."
        exit 1
    fi
fi

print_status "Creating runner directory..."
mkdir -p "$RUNNER_DIR"
cd "$RUNNER_DIR"

# Download the latest runner
print_status "Downloading GitHub Actions runner..."
RUNNER_VERSION="2.311.0"
RUNNER_URL="https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-osx-arm64-${RUNNER_VERSION}.tar.gz"

curl -o "actions-runner-osx-arm64-${RUNNER_VERSION}.tar.gz" -L "$RUNNER_URL"

print_status "Extracting runner..."
tar xzf "./actions-runner-osx-arm64-${RUNNER_VERSION}.tar.gz"

print_status "Setting permissions..."
chmod +x ./config.sh ./run.sh

# Print next steps
print_status "Setup complete! Next steps:"
echo ""
echo "1. Go to your repository settings:"
echo "   https://github.com/MOlechowski/go-cli-template/settings/actions/runners"
echo ""
echo "2. Click 'New self-hosted runner' and select 'macOS' and 'ARM64'"
echo ""
echo "3. Run the configuration command they provide (it will look like):"
echo "   ./config.sh --url https://github.com/MOlechowski/go-cli-template --token XXXXXX"
echo ""
echo "   When prompted for labels, use: macos,arm64,self-hosted,local"
echo ""
echo "4. Start the runner:"
echo "   ./run.sh"
echo ""
echo "5. Update repository variables in GitHub:"
echo "   - Go to Settings → Secrets and variables → Actions → Variables"
echo "   - Set MACOS_RUNNER = \"[self-hosted, macos, arm64, local]\""
echo ""
echo "6. Test with a small commit to trigger CI"
echo ""

print_status "Your runner directory is ready at: $RUNNER_DIR"
print_status "System specs detected:"
echo "  - CPU: $(sysctl -n machdep.cpu.brand_string) ($(sysctl -n hw.ncpu) cores)"
echo "  - Memory: $(echo "scale=1; $(sysctl -n hw.memsize) / 1024 / 1024 / 1024" | bc)GB"
echo "  - Architecture: $(uname -m)"

print_warning "Remember to install the runner as a service for production use:"
echo "  sudo ./svc.sh install"
echo "  sudo ./svc.sh start"