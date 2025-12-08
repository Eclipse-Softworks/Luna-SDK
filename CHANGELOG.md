# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-12-08

### Added

#### SDK Core
- TypeScript, Python, and Go client libraries
- API key authentication (`lk_<env>_<key>` format)
- OAuth token authentication with automatic refresh
- Users and Projects resource support (CRUD operations)
- Pagination tooling with async iterators (`iterate()` helpers)
- Secure storage integration (`keytar`, `keyring`, `zalando/go-keyring`)
- Standardized error handling with specific error types (`AuthenticationError`, etc.)
- Automatic retry with exponential backoff and jitter
- Configurable timeouts and max retries
- Environment variable configuration loading
- Structured JSON logging with sensitive data redaction

#### Production Readiness
- Go SDK: Refactored `NewClient`, `NewAPIKeyAuth`, `NewTokenAuth` to return errors (no panics)
- Python SDK: Dynamic User-Agent generation with system info
- Python SDK: Async context managers for proper resource cleanup
- TypeScript SDK: Dynamic User-Agent implementation
- CLI: Localhost-only OAuth callback binding for security
- All SDKs: Comprehensive test suites (Go: 26, Python: 80, TypeScript: 84 tests)

#### Service Modules
- **ResMate**: Student Residences management (listings, search, filtering)
- **Identity**: Group management and access control
- **Storage**: Buckets and file management (upload/download URLs)
- **AI Tools**: Chat completions integration
- **Automation**: Workflow management and execution triggers

#### TypeScript SDK
- ESM/CJS dual output with full TypeScript types
- Vitest test suite with coverage thresholds
- ESLint + Prettier for code quality
- npm publishing configuration

#### Python SDK
- Async/await support with httpx
- Pydantic models for request/response validation
- pytest test suite with async support
- ruff + mypy for linting and type checking
- PyPI publishing configuration

#### Go SDK
- Idiomatic Go with functional options pattern
- Thread-safe token refresh
- golangci-lint for code quality
- Makefile for build/test/lint commands

#### CLI
- Luna CLI tool built with Cobra
- Commands: `auth`, `users`, `projects`, `config`
- Multiple output formats: table, JSON, YAML
- Profile-based configuration (`~/.luna/config.yaml`)
- Browser-based OAuth login flow

#### Documentation
- OpenAPI 3.1.0 specification
- Canonical error codes schema
- README with quick start examples for all languages

