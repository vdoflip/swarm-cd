# Swarm CD Examples

This directory contains example configurations and deployment files for Swarm CD.

## Quick Start

1. Clone this repository:
   ```bash
   git clone https://github.com/m-adawi/swarm-cd.git
   cd swarm-cd
   ```

2. Copy and modify the example configuration files:
   ```bash
   cp examples/config.yaml .
   cp examples/repos.yaml .
   cp examples/stacks.yaml .
   ```

## Deploying Swarm CD with 1Password Integration

### Prerequisites

1. A Docker Swarm cluster (initialize one if needed):
   ```bash
   docker swarm init
   ```

2. A 1Password account with admin access to create Connect server credentials

### Setting up 1Password Connect

1. Sign in to your 1Password account at [1password.com](https://1password.com)

2. Go to Integrations → Connect server → New Connect server

3. Give your Connect server a name (e.g., "Swarm CD")

4. Select the vaults you want to make accessible

5. Click "Add Connect Server" to create the credentials

6. Download the `1password-credentials.json` file

7. Move the credentials file to your deployment directory:
   ```bash
   mv ~/Downloads/1password-credentials.json .
   ```

### Configuration Files

1. `config.yaml`: Global configuration
   ```bash
   cp examples/config.yaml .
   # Edit as needed, especially:
   # - web.host and web.port
   # - update_interval
   ```

2. `repos.yaml`: Git repository configurations
   ```bash
   cp examples/repos.yaml .
   # Add your repositories and credentials
   ```

3. `stacks.yaml`: Stack configurations
   ```bash
   cp examples/stacks.yaml .
   # Configure your stacks
   ```

### Deployment

1. Build the Swarm CD image:
   ```bash
   docker build -t swarm-cd .
   ```

2. Deploy the stack:
   ```bash
   docker stack deploy -c examples/docker-compose.swarm-cd.yaml swarm-cd
   ```

3. Verify the deployment:
   ```bash
   docker stack ps swarm-cd
   ```

4. Access the Web UI:
   - Open http://localhost:8080 (or your configured host:port)

## Example Stack Configurations

### Using SOPS for Secrets

See `stack-1password-example.yaml` for an example of how to configure a stack to use SOPS for secret management.

Key points:
- Specify SOPS files in `sops_files`
- Enable `sops_secrets_discovery` if needed
- Ensure SOPS is properly configured with your key management system

### Using 1Password for Secrets

See `stack-1password-example.yaml` for an example of how to configure a stack to use 1Password for secret management.

Key points:
- Set `secret_manager: 1password`
- Configure onepassword section:
  ```yaml
  onepassword:
    connect_host: http://op-connect-api:8080
    connect_token: ${OP_CONNECT_TOKEN}
    vault: your-vault-name
  ```

## Troubleshooting

1. Check service logs:
   ```bash
   docker service logs swarm-cd_swarm-cd
   docker service logs swarm-cd_op-connect-api
   docker service logs swarm-cd_op-connect-sync
   ```

2. Common issues:
   - Ensure `1password-credentials.json` is in the correct location
   - Verify network connectivity between services
   - Check that all required secrets and configs are present

## Security Notes

1. Protect your credentials:
   - Keep `1password-credentials.json` secure
   - Never commit credentials to version control
   - Use appropriate file permissions

2. Network security:
   - The 1Password Connect API is only accessible within the swarm-cd-net network
   - Web UI access can be restricted by configuring web.host in config.yaml

## Additional Resources

- [1Password Connect Documentation](https://developer.1password.com/docs/connect)
- [Docker Swarm Documentation](https://docs.docker.com/engine/swarm/)
- [SOPS Documentation](https://github.com/mozilla/sops)
