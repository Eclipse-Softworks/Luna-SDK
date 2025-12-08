---
description: Create a new release version
---

# Create Release

Bump version, update changelog, and push release tag.

## Steps

1. Determine current version:
```bash
node -p "require('./packages/typescript/package.json').version"
```

2. Update TypeScript version:
```bash
cd packages/typescript && npm version <patch|minor|major> --no-git-tag-version
```

3. Update Python version in `pyproject.toml`:
- Edit `packages/python/pyproject.toml`
- Change `version = "X.Y.Z"` to new version

4. Update Go SDK version constant:
- Edit `packages/go/luna/client.go`
- Change `const Version = "X.Y.Z"` to new version

5. Update CLI version:
- Edit `cli/cmd/root.go`
- Change `const Version = "X.Y.Z"` to new version

6. Update CHANGELOG.md with release notes

7. Commit changes:
```bash
git add -A && git commit -m "chore: bump version to vX.Y.Z"
```

8. Create and push tag:
```bash
git tag vX.Y.Z && git push Luna-SDK master && git push Luna-SDK vX.Y.Z
```

## Notes
- The release workflow will automatically publish to npm, PyPI, and trigger Go proxy
- Check GitHub Actions for release status
