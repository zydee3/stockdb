#!/bin/bash

# Database configuration
DATABASE_USERNAME="postgres"
DATABASE_PASSWORD="password"
DATABASE_NAME="stockdb"
DATABASE_PORT=5432
DATABASE_DATA_DIR="/var/lib/postgresql/data"

# Docker configuration
# 0 means no limit
DOCKER_CPU_CORE_LIMIT=0
DOCKER_CPU_MEMORY_LIMIT=0
DOCKER_IMAGE="timescale/timescaledb:latest-pg16"

# Supported Ubuntu versions
declare -A supported_versions
supported_versions=(
    ["20.04"]="focal"
    ["22.04"]="jammy"
    ["24.04"]="noble"
    ["24.10"]="oracular"
)

info() {
    echo -e "[INFO] $1"
}

warn() {
    echo -e "[WARN] $1"
}

error() {
    echo -e "[ERRO] $1"
    exit 1
}

assert_os_valid() {
    # Check if the OS is acceptable by checking the distribution info
    if [ ! -f /etc/os-release ]; then
        error "Cannot detect OS type. This script is for Ubuntu 22.04 only."
    fi

    # Load the OS release information
    . /etc/os-release

    # Check if the OS is Ubuntu
    if [ "$NAME" != "Ubuntu" ]; then
        error "This script is for Ubuntu only. Detected OS: $NAME"
    fi

    # Check if the OS version is supported
    version_supported=false
    for version in "${!supported_versions[@]}"; do
        if [ "$VERSION_ID" == "$version" ] && [ "$VERSION_CODENAME" == "${supported_versions[$version]}" ]; then
            version_supported=true
            break
        fi
    done

    # If the version is not supported, display supported versions and exit
    if [ "$version_supported" != true ]; then
        warn "Detected OS: $NAME $VERSION_ID ($VERSION_CODENAME)"
        warn "This script supports these versions:"
        for version in $(echo "${!supported_versions[@]}" | tr ' ' '\n' | sort); do
            warn "\t- Ubuntu $version (${supported_versions[$version]})"
        done

        error "Unsupported Ubuntu version."
    fi
}

run_command() {
    local command="$1"

    info "Running Command: $command"
    # Use eval to properly execute complex commands
    eval "$command" >/tmp/command_output.log 2>&1
    local status=$?

    if [ $status -ne 0 ]; then
        local output=$(cat /tmp/command_output.log)
        info "Command Output: $output"
        error "Command failed with status $status"
    fi
}

is_docker_installed() {
    if command -v docker &>/dev/null; then
        return 0
    else
        return 1
    fi
}

enable_docker_service() {
    info "Enabling Docker service"

    # Enable and start the Docker service
    run_command "sudo systemctl enable docker"
    run_command "sudo systemctl start docker"
}

install_docker() {
    # install_docker logic based on this guide:
    # https://www.postgresql.org/download/linux/ubuntu/

    # Check if Docker is already installed
    if is_docker_installed; then
        info "Docker is already installed. Version: $(docker --version)"
        return
    fi

    info "Installing Docker"

    # Run each installation step
    run_command "sudo apt-get update"
    run_command "sudo apt-get install -y ca-certificates curl"
    run_command "sudo install -m 0755 -d /etc/apt/keyrings"
    run_command "sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc"
    run_command "sudo chmod a+r /etc/apt/keyrings/docker.asc"

    # Add the repository to Apt sources
    local repo_command="echo \"deb [arch=\$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \$(. /etc/os-release && echo \"\${UBUNTU_CODENAME:-\$VERSION_CODENAME}\") stable\" | sudo tee /etc/apt/sources.list.d/docker.list"
    run_command "$repo_command"

    # Update the package index again to install Docker Engine, Docker CLI,
    # containerd, Buildx, and Compose plugins
    run_command "sudo apt-get update"
    run_command "sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin"

    # Add the current user to the docker group
    if is_docker_installed; then
        run_command "sudo usermod -aG docker $(whoami)"
    else
        error "Failed to install Docker (post-install verification check)."
    fi

    enable_docker_service
}

create_database_docker_volume() {
    info "Creating Docker volume for PostgreSQL"

    # Check if the Docker volume already exists
    if docker volume inspect "$DATABASE_NAME" &>/dev/null; then
        return
    fi

    # Create the Docker volume
    run_command "docker volume create $DATABASE_NAME"
}

create_database_docker_image() {
    info "Creating Docker image for PostgreSQL"

    # Check if the Docker image already exists, delete it
    if docker image inspect "$DOCKER_IMAGE" &>/dev/null; then
        error "Docker image $DOCKER_IMAGE already exists. Please remove it."
    fi

    # Pull and run the Docker image
    NAME_ARG="--name $DATABASE_NAME"
    CPUS_ARG="--cpus $DOCKER_CPU_CORE_LIMIT"
    MEMORY_ARG="--memory $DOCKER_CPU_MEMORY_LIMIT"
    PORT_ARG="-p $DATABASE_PORT:$DATABASE_PORT"
    VOLUME_ARG="-v $DATABASE_NAME:$DATABASE_DATA_DIR"
    PASSWORD_ARG="-e POSTGRES_PASSWORD=$DATABASE_PASSWORD"
    run_command "docker run -d $NAME_ARG $CPUS_ARG $MEMORY_ARG $PORT_ARG $VOLUME_ARG $PASSWORD_ARG $DOCKER_IMAGE"
}

allow_remote_connections() {
    info "Allowing remote connections to PostgreSQL"

    # Check if the firewall is installed
    if ! command -v ufw &>/dev/null; then
        info "Installing UFW"
        run_command "sudo apt-get install -y ufw"
        run_command "sudo ufw allow OpenSSH"
        run_command "sudo ufw enable"
    fi

    # Open up the PostgreSQL port in the firewall
    run_command "sudo ufw allow 5432/tcp"

    # Modify PostgreSQL configuration
    run_command "docker exec -u postgres $DATABASE_NAME bash -c \"sed -i \\\"s/#listen_addresses = 'localhost'/listen_addresses = '*'/\\\" /var/lib/postgresql/data/postgresql.conf\""
    run_command "docker exec -u postgres $DATABASE_NAME bash -c \"echo \\\"host all all 0.0.0.0/0 md5\\\" >> /var/lib/postgresql/data/pg_hba.conf\""

    # Restart the container to apply changes
    run_command "docker restart $DATABASE_NAME"

    info "PostgreSQL configured for remote connections"
}

start_docker_container() {
    info "Starting Docker container for PostgreSQL"

    # Check if the Docker container is already running
    if docker ps -q --filter "name=$DATABASE_NAME" &>/dev/null; then
        return
    fi

    # Start the Docker container
    run_command "docker start $DATABASE_NAME"
}

assert_os_valid
install_docker
create_database_docker_volume
create_database_docker_image
allow_remote_connections
start_docker_container
