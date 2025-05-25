# nani

[![Go Reference](https://pkg.go.dev/badge/github.com/asaidimu/nani.svg)](https://pkg.go.dev/github.com/asaidimu/nani)
[![Build Status](https://github.com/asaidimu/nani/workflows/Test%20Workflow/badge.svg)](https://github.com/asaidimu/nani/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple terminal ui

## Installation

```bash
go get github.com/asaidimu/nani
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/asaidimu/nani/pkg"
)

func main() {
	fmt.Println(pkg.Greeting("World"))
}
```

## Development

This project follows conventional commits for automated semantic versioning.

### Commit Message Format

- **fix:** a commit that fixes a bug (corresponds to PATCH in SemVer)
- **feat:** a commit that adds new functionality (corresponds to MINOR in SemVer)
- **feat!:** or **fix!:** or **refactor!:** etc., a commit with a footer `BREAKING CHANGE:` introduces a breaking API change (corresponds to MAJOR in SemVer)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
