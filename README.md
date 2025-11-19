# secretenv

A command-line tool that fetches secrets from secret management services (such as AWS Secrets Manager) and executes commands with those secrets as environment variables.

## Overview

`secretenv` retrieves secrets stored in **dotenv (`.env`) format** from secret management services and runs specified commands with those secrets as environment variables.
This enables secure secret management following the [12 Factor App](https://12factor.net/config) methodology without changing your existing dotenv-based configuration.

## Features

- **Secure Secret Management**: Retrieve secrets from centralized secret management services instead of storing them in files
- **Simple Integration**: Works as a wrapper command like [`godotenv`](https://github.com/joho/godotenv), easily integrated into existing applications
- **Multiple Provider Support**: Extensible architecture supporting various secret management services

## Installation

### Using Go install

```bash
go install github.com/hakadoriya/secretenv/cmd/secretenv@latest
```

### Download Binary

Download the latest binary from the [Releases](https://github.com/hakadoriya/secretenv/releases) page.

## Usage

### Basic Usage

```bash
secretenv --provider <provider> --secret <secret-name> -- <command> [args...]
```

Or using environment variables:

```bash
export SECRETENV_PROVIDER=<provider>
export SECRETENV_SECRET=<secret-name>
secretenv -- <command> [args...]
```

### Options

| Option | Environment Variable | Description | Required |
|--------|---------------------|-------------|----------|
| `--provider` | `SECRETENV_PROVIDER` | Secret management service provider (e.g., `aws`) | Yes |
| `--secret` | `SECRETENV_SECRET` | Secret name containing the .env file | Yes |
| `--secret-version` | `SECRETENV_SECRET_VERSION` | Secret version to retrieve (default: provider-specific latest version) | No |

### Examples

#### Running with AWS Secrets Manager

```bash
# Using command-line options
secretenv --provider aws --secret my-app-secrets -- ./myapp

# Using environment variables
export SECRETENV_PROVIDER=aws
export SECRETENV_SECRET=my-app-secrets
secretenv -- ./myapp arg1 arg2

# Specifying a version
secretenv --provider aws --secret my-app-secrets --secret-version AWSCURRENT -- ./myapp
```

#### Dockerfile Example

```dockerfile
FROM alpine:latest

# Install secretenv
COPY secretenv /usr/local/bin/secretenv

# Set environment variables
ENV SECRETENV_PROVIDER=aws
ENV SECRETENV_SECRET=my-app-secrets

# Run application with secretenv
ENTRYPOINT ["secretenv", "--"]
CMD ["./myapp"]
```

## Supported Providers

### `aws` provider: AWS Secrets Manager

**Prerequisites:**
- AWS credentials configured (via environment variables, IAM role, or AWS credentials file)
- Appropriate IAM permissions to access Secrets Manager

**Required IAM Permissions:**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": "arn:aws:secretsmanager:region:account-id:secret:secret-name"
    }
  ]
}
```

**Default Version:**
- If `--secret-version` is not specified, `AWSCURRENT` is used

## .env File Format

The secret value should be stored in `.env` format.

### Example Secret Content

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=admin
DB_PASSWORD="p@ssw0rd"

# API Keys
API_KEY=abc123xyz
SECRET_TOKEN='secret-token-value'

# Feature Flags
FEATURE_X_ENABLED=true
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
