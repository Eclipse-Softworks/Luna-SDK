# Luna CLI

Command-line interface for the Eclipse Softworks Platform.

## Installation

```bash
go install github.com/eclipse-softworks/luna-sdk/cli@latest
```

## Quick Start

```bash
# Configure API key
luna config set api_key lk_live_xxxx

# Or use environment variable
export LUNA_API_KEY=lk_live_xxxx

# Check authentication status
luna auth status

# List users
luna users list

# Get a specific user
luna users get usr_123

# Create a project
luna projects create --name "My Project"
```

## Commands

### Authentication

```bash
luna auth login     # Interactive browser-based login
luna auth logout    # Clear stored credentials
luna auth status    # Show authentication status
```

### Users

```bash
luna users list                       # List all users
luna users get usr_123                # Get a user by ID
luna users create --name "..." --email "..."  # Create a user
luna users delete usr_123             # Delete a user
```

### Projects

```bash
luna projects list                    # List all projects
luna projects get prj_123             # Get a project by ID
luna projects create --name "..."     # Create a project
luna projects delete prj_123          # Delete a project
```

### Configuration

```bash
luna config list                      # List all configuration
luna config get api_key               # Get a config value
luna config set api_key lk_live_xxxx  # Set a config value
```

## Global Flags

```
--api-key string     API key (overrides config)
--profile string     Configuration profile (default: "default")
--format string      Output format: table, json, yaml (default: "table")
--no-color           Disable colored output
--verbose            Enable verbose output
--debug              Enable debug mode
```

## Configuration File

Located at `~/.luna/config.yaml`:

```yaml
default_profile: production

profiles:
  production:
    api_key: lk_live_xxxx
    base_url: https://api.eclipse.dev
    
  staging:
    api_key: lk_test_yyyy
    base_url: https://api.staging.eclipse.dev

settings:
  output_format: table
  color: true
```

## License

MIT © Eclipse Softworks
