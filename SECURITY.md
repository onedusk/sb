# Security Policy

## Supported Versions

We actively support the latest minor version of SB. Security updates will be backported to the previous minor version for 60 days after a new minor version is released.

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow responsible disclosure practices.

### How to Report

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security issues to:
- **Email**: security@dusk-labs.com (if available)
- **GitHub Security Advisory**: Use the [GitHub Security Advisory](https://github.com/dusk-labs/sb/security/advisories/new) feature

### What to Include

Please include as much information as possible:

1. **Description** of the vulnerability
2. **Steps to reproduce** the issue
3. **Affected versions** of SB
4. **Potential impact** assessment
5. **Suggested fix** (if you have one)
6. **Your contact information** for follow-up

### Response Timeline

- **Initial Response**: Within 48 hours of report
- **Status Update**: Within 7 days with assessment
- **Fix Timeline**: Depends on severity
  - **Critical**: 1-7 days
  - **High**: 7-14 days
  - **Medium**: 14-30 days
  - **Low**: 30-90 days

### Disclosure Policy

- We follow a **90-day disclosure timeline**
- Security fixes will be released in a patch version
- CVE identifiers will be requested for significant vulnerabilities
- Public disclosure will be coordinated with the reporter
- Credit will be given to the reporter (unless anonymity is requested)

## Security Considerations

### External Dependencies

SB relies on **ffmpeg** as an external dependency. Security considerations:

1. **FFmpeg Vulnerabilities**
   - Keep ffmpeg updated to the latest stable version
   - Monitor [FFmpeg security advisories](https://ffmpeg.org/security.html)
   - SB does not include ffmpeg; users must install separately

2. **Input Validation**
   - SB validates file extensions before processing
   - File paths are sanitized to prevent directory traversal
   - User-provided options are validated before passing to ffmpeg

3. **Command Injection**
   - All ffmpeg arguments are properly escaped
   - No shell interpretation of user input
   - Arguments are passed as arrays, not strings

### Safe Usage Guidelines

1. **File Permissions**
   ```bash
   # Run with minimal permissions
   # Avoid running as root/administrator
   ```

2. **Untrusted Input**
   ```bash
   # Validate files before processing
   sb info suspicious_file.mov

   # Use dry-run mode first
   sb mp4 -n suspicious_file.mov
   ```

3. **Output Directory**
   ```bash
   # Specify output directory explicitly
   sb mp4 -o /safe/output/path input.mov

   # Avoid overwriting critical files
   ```

4. **Configuration Files**
   ```bash
   # Protect configuration files
   chmod 600 ~/.sb.yaml

   # Review config before use
   cat ~/.sb.yaml
   ```

### Known Limitations

1. **FFmpeg Security**
   - SB inherits all security characteristics of ffmpeg
   - Malformed media files could exploit ffmpeg vulnerabilities
   - Always use the latest stable ffmpeg version

2. **File System Access**
   - SB has full read access to input files
   - SB has write access to output directories
   - No sandboxing is implemented

3. **Resource Consumption**
   - Large batch operations can consume significant CPU/memory
   - No rate limiting on worker pools
   - Users should monitor system resources

## Security Features

### Implemented

- ✓ Input validation (file extensions, paths)
- ✓ Argument sanitization for ffmpeg
- ✓ Context-based cancellation for operations
- ✓ Error handling for invalid inputs
- ✓ No elevated privileges required

### Future Enhancements

- [ ] File size limits and validation
- [ ] Media file format validation (beyond extension)
- [ ] Sandboxing options for ffmpeg execution
- [ ] Resource usage limits (CPU, memory, disk)
- [ ] Audit logging for operations
- [ ] SBOM (Software Bill of Materials) generation

## Vulnerability History

No security vulnerabilities have been reported or addressed in SB to date.

## Security Best Practices for Contributors

If you're contributing code:

1. **Never trust user input**
   - Validate all input parameters
   - Sanitize file paths
   - Escape shell arguments

2. **Handle errors securely**
   - Don't expose sensitive information in error messages
   - Log security-relevant events
   - Fail securely (deny by default)

3. **Review dependencies**
   - Check go.mod for vulnerable dependencies
   - Run `go mod tidy` regularly
   - Use `govulncheck` for vulnerability scanning

4. **Code review**
   - All PRs require review
   - Security-sensitive changes require extra scrutiny
   - Use static analysis tools (golangci-lint)

## Contact

For security-related questions or concerns:
- Create a private security advisory on GitHub
- Check existing issues/discussions (for non-sensitive topics)
- Refer to CONTRIBUTING.md for general contribution guidelines

---

**Last Updated**: 2025-10-17
