#!/bin/bash

###############################################################################
# TramaTex - Automated Installation Script for Linux
# 
# Usage: ./scripts/install.sh [dev|prod|test]
# 
# Supported distros: Ubuntu, Debian, CentOS, RHEL, Fedora
###############################################################################

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ENVIRONMENT=${1:-dev}  # Default: dev
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INSTALL_DIR="/opt/tramatex"

###############################################################################
# Functions
###############################################################################

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}TramaTex - Linux Installation Script${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        DISTRO=$ID
        VERSION=$VERSION_ID
    elif [ -f /etc/redhat-release ]; then
        DISTRO="rhel"
        VERSION=$(cat /etc/redhat-release | grep -oE '[0-9]+')
    else
        print_error "Cannot detect Linux distribution"
        exit 1
    fi
    
    print_info "Detected: $DISTRO $VERSION"
}

# Check if running as root or with sudo
check_sudo() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This script must be run as root or with sudo"
        exit 1
    fi
    print_success "Running with appropriate privileges"
}

# Check if Docker is installed
check_docker() {
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version | grep -oE '[0-9]+\.[0-9]+')
        print_success "Docker is already installed (version $DOCKER_VERSION)"
        return 0
    else
        print_info "Docker not found, will install"
        return 1
    fi
}

# Check if Docker Compose is installed
check_docker_compose() {
    if command -v docker-compose &> /dev/null; then
        DC_VERSION=$(docker-compose --version | grep -oE '[0-9]+\.[0-9]+')
        print_success "Docker Compose is already installed (version $DC_VERSION)"
        return 0
    else
        print_info "Docker Compose not found, will install"
        return 1
    fi
}

# Install Docker on Ubuntu/Debian
install_docker_ubuntu() {
    print_info "Installing Docker on Ubuntu/Debian..."
    
    # Add Docker repository
    apt-get update
    apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Install Docker
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    # Start Docker
    systemctl start docker
    systemctl enable docker
    
    print_success "Docker installed successfully"
}

# Install Docker on CentOS/RHEL/Fedora
install_docker_redhat() {
    print_info "Installing Docker on CentOS/RHEL/Fedora..."
    
    # Add Docker repository
    yum install -y yum-utils
    yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    
    # Install Docker
    yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    # Start Docker
    systemctl start docker
    systemctl enable docker
    
    print_success "Docker installed successfully"
}

# Install Docker based on distribution
install_docker() {
    case $DISTRO in
        ubuntu|debian)
            install_docker_ubuntu
            ;;
        centos|rhel|fedora)
            install_docker_redhat
            ;;
        *)
            print_error "Unsupported distribution: $DISTRO"
            exit 1
            ;;
    esac
}

# Install required utilities
install_utilities() {
    print_info "Installing required utilities..."
    
    case $DISTRO in
        ubuntu|debian)
            apt-get install -y git curl wget openssl
            ;;
        centos|rhel|fedora)
            yum install -y git curl wget openssl
            ;;
    esac
    
    print_success "Utilities installed"
}

# Setup Docker user permissions
setup_docker_user() {
    if [ -n "$SUDO_USER" ]; then
        print_info "Adding $SUDO_USER to docker group..."
        usermod -aG docker $SUDO_USER
        print_success "User $SUDO_USER added to docker group"
        print_warning "Please log out and log back in for changes to take effect"
    fi
}

# Clone or update repository
setup_repository() {
    print_info "Setting up TramaTex repository..."
    
    if [ ! -d "$INSTALL_DIR" ]; then
        print_info "Cloning repository to $INSTALL_DIR..."
        mkdir -p $(dirname $INSTALL_DIR)
        cd $(dirname $INSTALL_DIR)
        git clone https://github.com/tu-usuario/tramatex.git
    else
        print_info "Updating existing repository..."
        cd $INSTALL_DIR
        git pull origin main
    fi
    
    cd $INSTALL_DIR
    print_success "Repository ready at $INSTALL_DIR"
}

# Generate secure JWT secret
generate_jwt_secret() {
    openssl rand -base64 32
}

# Setup environment configuration
setup_environment() {
    print_info "Setting up environment configuration..."
    
    if [ ! -f "$INSTALL_DIR/.env" ]; then
        cp "$INSTALL_DIR/.env.example" "$INSTALL_DIR/.env"
        print_success "Created .env file from template"
    else
        print_warning ".env file already exists, skipping"
    fi
    
    # Generate JWT secret if not set
    if ! grep -q "JWT_SECRET=" "$INSTALL_DIR/.env" || grep "JWT_SECRET=your-secret" "$INSTALL_DIR/.env"; then
        JWT_SECRET=$(generate_jwt_secret)
        sed -i "s|JWT_SECRET=.*|JWT_SECRET=$JWT_SECRET|" "$INSTALL_DIR/.env"
        print_success "Generated JWT_SECRET"
    fi
    
    # Configure based on environment
    case $ENVIRONMENT in
        prod)
            sed -i 's|LOG_LEVEL=.*|LOG_LEVEL=warn|' "$INSTALL_DIR/.env"
            sed -i "s|DB_PASSWORD=.*|DB_PASSWORD=$(openssl rand -base64 16)|" "$INSTALL_DIR/.env"
            print_warning "IMPORTANT: Review .env file and set production values"
            ;;
        test)
            sed -i 's|LOG_LEVEL=.*|LOG_LEVEL=debug|' "$INSTALL_DIR/.env"
            ;;
        *)  # dev
            sed -i 's|LOG_LEVEL=.*|LOG_LEVEL=info|' "$INSTALL_DIR/.env"
            ;;
    esac
    
    print_success "Environment configured for: $ENVIRONMENT"
}

# Start Docker services
start_services() {
    print_info "Starting Docker services..."
    
    cd $INSTALL_DIR
    
    docker-compose down 2>/dev/null || true
    docker-compose up -d
    
    print_success "Services started"
    
    # Wait for services to be ready
    print_info "Waiting for services to be ready..."
    sleep 5
    
    # Check health
    if docker-compose ps | grep -q "Up"; then
        print_success "Services are running"
    else
        print_error "Some services failed to start"
        docker-compose logs
        exit 1
    fi
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    # Check Docker is running
    if ! docker ps &> /dev/null; then
        print_error "Docker daemon is not running"
        exit 1
    fi
    print_success "Docker daemon running"
    
    # Check services are running
    cd $INSTALL_DIR
    if docker-compose ps | grep -q "Up"; then
        print_success "Docker services running"
    else
        print_error "Some Docker services not running"
        exit 1
    fi
    
    # Check API health endpoint
    sleep 3  # Give API time to start
    if curl -sf http://localhost:8080/api/health > /dev/null 2>&1; then
        print_success "API health check passed"
    else
        print_warning "API not yet responding (might still be starting)"
    fi
    
    # Check PostgreSQL
    if docker-compose exec -T tramatex_db psql -U tramatex -d tramatex -c "SELECT 1;" > /dev/null 2>&1; then
        print_success "PostgreSQL connection verified"
    else
        print_warning "PostgreSQL not yet ready"
    fi
}

# Print final status
print_final_status() {
    echo ""
    print_success "Installation complete!"
    echo ""
    echo "Next steps:"
    echo "1. Review .env file: $INSTALL_DIR/.env"
    echo "2. Access API: http://localhost:8080"
    echo "3. Access apps/frontend: http://localhost:5173 (if installed)"
    echo "4. View logs: cd $INSTALL_DIR && docker-compose logs -f"
    echo ""
    echo "Useful commands:"
    echo "  docker-compose ps           # View service status"
    echo "  docker-compose logs -f      # View live logs"
    echo "  docker-compose stop         # Stop services"
    echo "  docker-compose down         # Stop and remove containers"
    echo ""
}

###############################################################################
# Main Execution
###############################################################################

main() {
    print_header
    
    # Validate environment parameter
    case $ENVIRONMENT in
        dev|prod|test)
            print_info "Installation environment: $ENVIRONMENT"
            ;;
        *)
            print_error "Invalid environment: $ENVIRONMENT"
            print_info "Usage: ./scripts/install.sh [dev|prod|test]"
            exit 1
            ;;
    esac
    
    echo ""
    
    # Execution steps
    check_sudo
    detect_distro
    install_utilities
    
    if ! check_docker; then
        install_docker
    fi
    
    if ! check_docker_compose; then
        print_info "Docker Compose plugin will be installed with Docker"
    fi
    
    setup_docker_user
    setup_repository
    setup_environment
    start_services
    verify_installation
    print_final_status
}

# Run main function
main "$@"
