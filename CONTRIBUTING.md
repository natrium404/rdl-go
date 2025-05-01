# Contributing

First off, thanks for taking the time to contribute! We appreciate your help in improving RDL.
Whether you're fixing bugs, suggesting improvements, or adding new features, this guide will help you get started.

## Table of Contents

- [How to Contribute](<CONTRIBUTING#How to Contribute>)
- [Development Setup](<CONTRIBUTING#Development Setup>)
- [Code Guidelines](<CONTRIBUTING#Code Guidelines>)
- [Submitting a Pull Request](<CONTRIBUTING#Submitting a Pull Request>)
- [Feature Requests & Bug Reports](<CONTRIBUTING#Feature Requests & Bug Reports>)

## How to Contribute

- Fork the repository
- Create a feature branch:

```sh
git checkout -b feature/your-feature-name
```

- Make your changes
- Push to your fork:

```sh
git push origin feature/your-feature-name
```

Open a Pull Request from your branch to main

## Development Setup

To run RDL locally:

### Prerequisites

- Go v1.20+
- Node.js
- Wails CLI

### Steps

```sh
git clone https://github.com/yourusername/rdl.git
cd rdl
cd frontend && npm install && cd ..
wails dev
```

To build a production-ready binary:

```sh
wails build
```

## Code Guidelines

- Follow idiomatic Go formatting (go fmt)
- Keep UI code clean and modular (Svelte components)
- Add comments where necessary, especially for complex logic
- Test your changes before pushing

## Submitting a Pull Request

- Ensure your PR has a clear title and description
- Reference related issues (if any)
- Stick to the scope of a single feature or fix per PR

## Feature Requests & Bug Reports

- Use the Issues tab to:
- Report bugs
- Suggest new features
- Ask questions or request help
- Please include clear, detailed information when opening an issue.
