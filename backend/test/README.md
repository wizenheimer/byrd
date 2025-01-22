# API Test Suite

This repository contains Hurl tests for validating the API endpoints. The test suite covers Workspaces, Members, Competitors, Pages, Users, and Token validation.

## Prerequisites

- [Hurl](https://hurl.dev) installed (version 6.0.0 or higher)
- Access to the API (development/staging environment recommended)
- Valid authentication token

## Configuration

Create a `variables.env` file with your configuration:

```env
# Base configuration
host=http://localhost:8080
auth_token=YOUR_AUTH_TOKEN

# Optional: Pre-existing IDs for independent test runs
workspace_id=YOUR_WORKSPACE_ID
competitor_id=YOUR_COMPETITOR_ID
```

## Usage

```bash
hurl --variables-file variables.env workspace_tests.hurl
```

```bash
# Using workspace_id from variables file
hurl --variables-file variables.env member_tests.hurl

# Or using command line variables
hurl --variable host=http://localhost:8080 \
     --variable auth_token=YOUR_TOKEN \
     --variable workspace_id=YOUR_WORKSPACE_ID \
     member_tests.hurl
```

```bash
hurl --variables-file variables.env competitor_tests.hurl
```

```bash
hurl --variables-file variables.env page_tests.hurl
```

```bash
hurl --variables-file variables.env user_and_token_tests.hurl
```

```bash
#!/bin/bash

# Create a sequential_run.sh script
hurl --variables-file variables.env workspace_tests.hurl && \
hurl --variables-file variables.env member_tests.hurl && \
hurl --variables-file variables.env competitor_tests.hurl && \
hurl --variables-file variables.env page_tests.hurl && \
hurl --variables-file variables.env user_and_token_tests.hurl
```

```bash
# Verbose output
hurl --verbose --variables-file variables.env workspace_tests.hurl

# Debug output
hurl --error-format long --variables-file variables.env workspace_tests.hurl

# Test specific entry
hurl --test --from-entry 2 --to-entry 4 workspace_tests.hurl
```
