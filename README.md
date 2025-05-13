# Fileshear

Fileshear is a Go-based project for efficient file sharing and management.

## Features

- Fast and secure file transfers
- Simple CLI interface
- Cross-platform support

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) 1.18 or higher

### Installation

```bash
git clone https://github.com/yourusername/fileshear.git
cd fileshear
go build
```

### Usage

```bash
./fileshear [options]
```

Example:

```bash
./fileshear --send /path/to/file --to 192.168.1.2
```

## Project Structure

```
fileshear/
├── cmd/         # CLI commands
├── internal/    # Internal packages
├── pkg/         # Public packages
├── README.md
└── main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/foo`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/foo`)
5. Open a pull request

## License

This project is licensed under the MIT License.
