# Documentation Update Guide

This guide outlines how to keep the A-Drive documentation up-to-date with code changes.

## Documentation Structure

```
docs/
├── TECHNICAL_DOCUMENTATION.md    # Architecture, code structure, dependencies
├── FUNCTIONAL_DOCUMENTATION.md   # Features, API endpoints, user functionality
├── CHANGELOG.md                   # Version history and changes
└── DOC_UPDATE_GUIDE.md           # This guide
```

## When to Update Documentation

### Required Updates
Documentation MUST be updated when:

1. **API Changes**
   - New endpoints added
   - Endpoint parameters modified
   - Response format changes
   - Authentication requirements change

2. **Database Schema Changes**
   - New models added
   - Existing model fields modified
   - Relationships between models change
   - Migration requirements

3. **Configuration Changes**
   - New environment variables
   - Changed default values
   - New configuration options
   - Security settings modifications

4. **Architecture Changes**
   - New dependencies added
   - Framework upgrades
   - Middleware modifications
   - Route structure changes

5. **Feature Additions/Removals**
   - New functionality implemented
   - Existing features removed or deprecated
   - User interface changes
   - Admin functionality modifications

### Optional Updates
Documentation SHOULD be updated when:
- Performance optimizations
- Bug fixes that affect behavior
- Code refactoring that changes structure
- Security improvements

## Update Process

### 1. Identify Documentation Impact
Before making code changes, determine which documentation files need updates:

- **Technical Documentation**: Code structure, dependencies, models
- **Functional Documentation**: User-facing features, API endpoints
- **Changelog**: All significant changes

### 2. Update Documentation Files

#### Technical Documentation Updates
When updating `TECHNICAL_DOCUMENTATION.md`:

```markdown
## What to Update

### Dependencies Section
- Add new dependencies to the dependency list
- Update version numbers for upgraded packages
- Remove deprecated dependencies

### Architecture Section
- Update directory structure if files are added/moved
- Modify data models when schema changes
- Update configuration examples

### Code Examples
- Reflect actual code structure
- Update import statements
- Ensure examples compile and run
```

#### Functional Documentation Updates
When updating `FUNCTIONAL_DOCUMENTATION.md`:

```markdown
## What to Update

### API Endpoints
- Document new endpoints with full details
- Update parameter lists and examples
- Include response format changes
- Update authentication requirements

### Features Section
- Add descriptions of new functionality
- Update user workflows
- Include new configuration options
- Document changed behaviors
```

#### Changelog Updates
When updating `CHANGELOG.md`:

```markdown
## Format for New Entries

### [Version] - YYYY-MM-DD

#### Added
- New features and functionality

#### Changed
- Modifications to existing features

#### Deprecated
- Features marked for removal

#### Removed
- Deleted features

#### Fixed
- Bug fixes

#### Security
- Security-related changes
```

### 3. Version Control Integration

#### Pre-commit Checklist
Before committing code changes:

- [ ] Technical documentation reflects code changes
- [ ] Functional documentation covers new/changed features
- [ ] Changelog includes all significant changes
- [ ] Documentation version numbers updated
- [ ] Examples tested and verified

#### Commit Message Convention
When documentation is updated with code:
```
feat: add user bulk operations

- Implement bulk user creation endpoint
- Add admin middleware for bulk operations
- Update API documentation with new endpoints
- Add changelog entry for v1.1.0

Docs-Updated: TECHNICAL_DOCUMENTATION.md, FUNCTIONAL_DOCUMENTATION.md, CHANGELOG.md
```

## Automated Documentation Updates

### Git Hooks (Recommended)
Set up a pre-commit hook to remind about documentation:

```bash
#!/bin/sh
# .git/hooks/pre-commit

# Check if any Go files are modified
if git diff --cached --name-only | grep -q "\.go$"; then
    echo "Go files modified. Please ensure documentation is updated:"
    echo "- docs/TECHNICAL_DOCUMENTATION.md (for code structure changes)"
    echo "- docs/FUNCTIONAL_DOCUMENTATION.md (for feature changes)"
    echo "- docs/CHANGELOG.md (for all significant changes)"
    echo ""
    echo "Continue with commit? (y/n)"
    read answer
    if [ "$answer" != "y" ]; then
        exit 1
    fi
fi
```

### Documentation Review Process

#### Pull Request Template
Include documentation checklist in PR template:

```markdown
## Documentation Checklist
- [ ] Technical documentation updated for code changes
- [ ] Functional documentation covers new features
- [ ] Changelog includes all changes
- [ ] API documentation reflects endpoint changes
- [ ] Configuration examples are current
```

## Documentation Maintenance

### Regular Reviews
Schedule monthly documentation reviews to:
- Verify accuracy of all documentation
- Update outdated examples
- Improve clarity and completeness
- Check for broken links or references

### Version Synchronization
- Keep documentation version aligned with code version
- Update "Last Updated" dates in all documentation files
- Ensure changelog accurately reflects release history

### Quality Standards
Documentation should be:
- **Accurate**: Reflects actual code behavior
- **Complete**: Covers all features and functionality
- **Clear**: Easy to understand for intended audience
- **Current**: Updated with each significant change
- **Consistent**: Follows established format and style

## Tools and Automation

### Recommended Tools
- **markdownlint**: Ensure consistent markdown formatting
- **alex**: Check for inclusive language
- **textlint**: Grammar and style checking
- **doctoc**: Auto-generate table of contents

### CI/CD Integration
Consider adding documentation checks to CI pipeline:
```yaml
# Example GitHub Action step
- name: Check Documentation
  run: |
    # Check if documentation files exist
    test -f docs/TECHNICAL_DOCUMENTATION.md
    test -f docs/FUNCTIONAL_DOCUMENTATION.md
    test -f docs/CHANGELOG.md
    
    # Validate markdown format
    markdownlint docs/
```

## Contact and Support

For questions about documentation updates:
- Review this guide first
- Check existing documentation for examples
- Consult with team lead for major changes
- Update this guide if process changes

---

**Remember**: Good documentation is a gift to your future self and your teammates. Keep it current, clear, and comprehensive.