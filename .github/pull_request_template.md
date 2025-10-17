# Pull Request

## Description

<!-- Provide a brief description of your changes -->

## Type of Change

<!-- Mark relevant items with an 'x' -->

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring
- [ ] Performance improvement
- [ ] Test addition/improvement

## Related Issues

<!-- Link related issues using keywords: Fixes #123, Relates to #456 -->

Fixes #
Relates to #

## Changes Made

<!-- List key changes made in this PR -->

-
-
-

## Testing

<!-- Describe the tests you ran and how to reproduce them -->

### Test Commands

```bash
# Commands used to test your changes
make test
./dist/sb mp4 test.mov
```

### Test Results

- [ ] All existing tests pass
- [ ] Added new tests for changes
- [ ] Tested manually with real files
- [ ] Tested on multiple platforms (if applicable)

### Test Environment

- OS:
- Go Version:
- FFmpeg Version:

## Breaking Changes

<!-- If this PR introduces breaking changes, describe them here -->

- [ ] No breaking changes
- [ ] Breaking changes (describe below)

**Breaking Change Description:**


## Checklist

<!-- Mark completed items with an 'x' -->

- [ ] My code follows the project's coding standards
- [ ] I have run `gofmt -s -w .` on my code
- [ ] I have run `golangci-lint run ./...` with no errors
- [ ] I have added tests that prove my fix/feature works
- [ ] I have updated documentation (README, docs/, comments)
- [ ] I have updated CHANGELOG.md
- [ ] My changes generate no new warnings
- [ ] I have checked for potential security issues
- [ ] I have tested with multiple file types/sizes (if applicable)

## Screenshots/Examples

<!-- If applicable, add screenshots or example outputs -->

```bash
# Example usage
sb mp4 --help
```

## Performance Impact

<!-- If applicable, describe performance impact -->

- [ ] No performance impact
- [ ] Performance improvement (describe below)
- [ ] Performance degradation (describe below and justify)

**Performance Notes:**


## Additional Notes

<!-- Any additional information, context, or notes for reviewers -->

## Reviewer Notes

<!-- Optional: specific areas you'd like reviewers to focus on -->

Please pay special attention to:
-
-

---

**For Maintainers:**
- [ ] Code review completed
- [ ] Tests passing in CI
- [ ] Documentation reviewed
- [ ] Ready to merge
