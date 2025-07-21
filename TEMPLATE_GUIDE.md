# GitHub Template Usage Guide

This repository is set up as a **GitHub Template** with enhanced AI/LLM integration capabilities.

## 🚀 Using This Template

### Option 1: GitHub Web Interface
1. Click the **"Use this template"** button on the repository page
2. Choose **"Create a new repository"**
3. Configure your new repository settings
4. Click **"Create repository from template"**

### Option 2: GitHub CLI
```bash
gh repo create my-new-cli --template MOlechowski/go-cli-template
cd my-new-cli
```

## 🤖 AI Integration Features

### Claude Code Integration
This template includes:
- **Automatic Claude responses** to @claude mentions in issues/PRs
- **AI-powered code review** and suggestions
- **Intelligent bug analysis** and debugging help
- **Code implementation assistance** following Go best practices

### Issue Templates
- **Bug Report** - Enhanced with AI debugging guidance
- **Feature Request** - Includes AI implementation assistance 
- **Claude Task** - Dedicated template for AI coding tasks

### Pull Request Template
- **Comprehensive checklist** including linting and standards
- **AI code review** integration prompts
- **Go-specific quality checks**

## 🛠️ Template Customization

### After Creating from Template

⚠️ **Security Note**: This template uses Probot Settings app for automated repository management. The `.github/CODEOWNERS` file protects sensitive configuration files, but ensure you review all PRs that modify `.github/settings.yml` carefully.

1. **Update Repository Information**
   ```bash
   # Update go.mod with your module path
   go mod edit -module github.com/yourusername/your-cli-name
   
   # Update import paths in code
   find . -name "*.go" -exec sed -i 's|github.com/go-cli-template/hello-world-cli|github.com/yourusername/your-cli-name|g' {} \;
   ```

2. **Customize Application Details**
   - Update `cmd/hello-world-cli/main.go` with your CLI name
   - Modify `internal/cli/root.go` with your application description
   - Update `pkg/version/version.go` with your version info

3. **Configure Claude Integration** (Optional)
   - Add your `ANTHROPIC_API_KEY` to repository secrets
   - Customize `.github/workflows/claude.yml` if needed
   - Modify issue templates for your specific use cases

4. **Update Documentation**
   - Customize `README.md` with your project details
   - Update `CLAUDE.md` with project-specific context
   - Modify `RULES-GO.md` for your coding standards

## 📁 Template Structure

```
your-new-cli/
├── cmd/                     # Application entrypoints
│   └── hello-world-cli/     # Main CLI application
├── internal/               # Private packages
│   ├── cli/               # Command implementations
│   ├── errors/            # Error handling
│   ├── greeting/          # Business logic example
│   └── logger/            # Logging utilities
├── pkg/                   # Public packages
│   └── version/           # Version information
├── .github/               # GitHub templates & workflows
│   ├── ISSUE_TEMPLATE/    # Issue templates with AI integration
│   ├── workflows/         # CI/CD + Claude integration
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── dependabot.yml     # Dependency management
├── docs/                  # Documentation
├── .golangci.yml          # Comprehensive linting config
├── CLAUDE.md             # AI assistant context
├── RULES-GO.md           # Go coding standards
└── Taskfile.yml          # Build automation
```

## 🤖 Working with Claude

### Getting Help
- **Issue Creation**: Use the "Claude Code Task" template
- **Code Review**: Mention `@claude review this PR`
- **Bug Analysis**: Mention `@claude` in bug reports
- **Feature Implementation**: Use `@claude implement [feature description]`

### Best Practices for AI Integration
1. **Be Specific**: Provide detailed requirements and context
2. **Reference Standards**: Claude follows your RULES-GO.md automatically
3. **Iterative Refinement**: Use follow-up @claude mentions for improvements
4. **Review AI Suggestions**: Always review and test AI-generated code

## 🔧 Advanced Configuration

### Linting Configuration
- **golangci-lint v2.2.2** with 19+ comprehensive linters
- **Security scanning** with gosec
- **Performance optimization** checks
- **Go 1.24** compatibility

### CI/CD Pipeline
- **Multi-platform testing** (Linux, macOS, Windows)
- **Automated linting** with quality gates
- **Dependency management** with Dependabot
- **Release automation** with semantic versioning

### Development Tools
- **Task runner** with comprehensive build commands
- **Structured logging** with slog
- **Error handling** patterns and utilities
- **Testing** utilities and examples

## 📚 Learning Resources

- [Go CLI Best Practices](./RULES-GO.md)
- [Project Architecture](./CLAUDE.md)
- [Contributing Guidelines](./docs/)
- [Claude Code Documentation](https://github.com/anthropics/claude-code-action)

## 🆘 Support

- **Create an issue** using our templates
- **Mention @claude** for AI assistance
- **Check existing issues** for similar problems
- **Review documentation** in the docs/ directory

---

**Happy coding with AI assistance! 🚀🤖**