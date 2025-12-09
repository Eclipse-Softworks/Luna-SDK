# Versioning Policy

Luna SDK follows [Semantic Versioning 2.0.0](https://semver.org/).

## Version Format: `MAJOR.MINOR.PATCH`

- **MAJOR**: Incompatible API changes.
- **MINOR**: Additions in a backward-compatible manner.
- **PATCH**: Backward-compatible bug fixes.

## Release Cadence

- **Stable**: Released monthly or upon critical bug fixes.
- **Beta**: Released bi-weekly for testing new features.

## Deprecation

When a feature is deprecated:
1. It is marked as deprecated in the `CHANGELOG`.
2. A warning is logged at runtime when used.
3. It will remain available for at least one major version cycle before removal.

## Changelog

See [CHANGELOG.md](https://github.com/Eclipse-Softworks/Luna-SDK/blob/master/CHANGELOG.md) for the full history of changes.
