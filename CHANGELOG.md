# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2025-12-13

### Added

#### South African Payments Module
- **PayFast**: Payment creation, webhook verification, refund processing
- **Ozow**: Instant EFT payment integration with hash verification
- **Yoco**: Checkout API, webhook handling, refund functionality
- **PayShap**: Real-time payments, ShapID lookup, QR code generation

#### Communications Module
- **SMS**: Multi-provider support (Clickatell, Africa's Talking, Twilio)
- **WhatsApp**: Business API integration (text, templates, media messages)
- **USSD**: Interactive menu support for SA networks (MTN, Vodacom, Cell C, Telkom)

#### AI/ML Expansion
- Multi-provider LLM support (OpenAI, Anthropic, Google Gemini, Ollama, Azure)
- Text embeddings with cosine similarity utilities
- Image analysis and vision capabilities
- South African language translation (isiZulu, isiXhosa, Afrikaans, Sepedi, Sesotho, Setswana, Xitsonga, Tshivenda, siSwati, isiNdebele)

#### South African Business Tools
- **CIPC**: Company registration lookup and verification
- **B-BBEE**: Compliance verification, level calculation, EME/QSE classification
- **ID Validation**: SA ID number parsing with Luhn checksum, DOB extraction, age verification
- **Address Utilities**: Postal code lookup, province detection, address validation

### Changed
- Updated `LunaClient` with new resource modules: `payments`, `messaging`, `zaTools`
- Expanded `AiResource` with embeddings, vision, and translation methods
- Enhanced type exports in main `index.ts`


## [1.0.2] - 2025-12-09

### Changed
#### Documentation
- **UI Overhaul**: Redesigned landing page with modern dark theme and glassmorphism
- Added animated components using `framer-motion`
- Added feature icons using `react-icons`
- Improved typography (Inter font) and spacing

## [1.0.1] - 2025-12-09

### Added
#### Documentation
- Launched official documentation site at [docs.eclipse-softworks.com](https://docs.eclipse-softworks.com/)
- Integrated Docusaurus with Cloudflare Pages for continuous deployment

#### CI/CD
- Added `docs.yml` workflow for automated documentation builds and deployments
- Configured OIDC Trusted Publishers for PyPI releases
- Configured npm provenance for `@eclipse-softworks/luna-sdk` releases

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

