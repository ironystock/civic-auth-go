# Contributing to Civic Auth Go SDK

Thank you for your interest in contributing to the Civic Auth Go SDK! We welcome contributions from the community.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/civic-auth-go.git
   cd civic-auth-go
   ```
3. Install dependencies:
   ```bash
   make install-deps
   ```

## Development Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and add tests if applicable

3. Run the test suite:
   ```bash
   make test
   ```

4. Run linting:
   ```bash
   make lint
   ```

5. Build examples to ensure they work:
   ```bash
   make build
   ```

6. Commit your changes:
   ```bash
   git add .
   git commit -m "Add your meaningful commit message"
   ```

7. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

8. Create a Pull Request on GitHub

## Code Style

- Follow Go conventions and best practices
- Use `gofmt` for code formatting (run `make fmt`)
- Add tests for new functionality
- Update documentation as needed
- Keep commits focused and atomic

## Testing

- Write unit tests for new functions and methods
- Ensure all tests pass before submitting PR
- Include integration tests for major features
- Test coverage should not decrease

## Documentation

- Update README.md if adding new features
- Add godoc comments for public functions
- Include examples for complex functionality
- Update CHANGELOG.md for notable changes

## Pull Request Guidelines

- Provide a clear description of the problem and solution
- Reference any related issues
- Include tests that cover the changes
- Update documentation as needed
- Ensure CI passes

## Issues

- Use the issue tracker for bugs and feature requests
- Provide detailed reproduction steps for bugs
- Include Go version and OS information
- Check for existing issues before creating new ones

## Questions

If you have questions about contributing, please open an issue or reach out to the maintainers.

Thank you for contributing!
