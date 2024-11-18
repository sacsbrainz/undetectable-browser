# Undetectable Browser

A Go-based tool for launching an automated browser instance with advanced proxy support and anti-detection features. Built using [Rod](https://go-rod.github.io/), this tool helps create browser instances that are harder to detect as automated browsers.

[![Build and Release](https://github.com/sacsbrainz/undetectable-browser/actions/workflows/release.yml/badge.svg)](https://github.com/sacsbrainz/undetectable-browser/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacsbrainz/undetectable-browser)](https://goreportcard.com/report/github.com/sacsbrainz/undetectable-browser)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- 🌐 **Proxy Support**: Full support for HTTP/HTTPS proxies with authentication
- 🔒 **Profile Management**: Separate user profiles for different proxy configurations
- 🛡️ **Anti-Detection**: Built-in features to avoid automation detection
- 🚀 **Cross-Platform**: Works on Windows, macOS, and Linux
- 📝 **Detailed Logging**: Optional quiet mode for silent operation
- 🔄 **Session Persistence**: Maintains browser state between sessions

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/sacsbrainz/undetectable-browser/releases).

Available builds:

- Windows (32/64-bit)
- macOS (Intel/Apple Silicon)
- Linux (32/64-bit)

### Building from Source

Requirements:

- Go 1.23 or later
- Git

```bash
# Clone the repository
git clone https://github.com/sacsbrainz/undetectable-browser.git
cd undetectable-browser

# Build the project
go build

# (Optional) Install globally
go install
```

## Usage

Basic usage with proxy:

```bash
./undetectable-browser -proxy http://user:pass@proxy.example.com:8080
```

### Command Line Options

```
Usage of undetectable-browser:

Example:
  ./undetectable-browser -proxy http://user:pass@localhost:8080
  ./undetectable-browser -proxy localhost:8080 -user-dir custom_profile

Flags:
  -help
        Show help message
  -proxy string
        Proxy URL in format [protocol://][username:password@]hostname[:port]
  -quiet
        Silence all log output
  -user-dir string
        Custom user directory path (defaults to proxy hostname)
  -version
        Show version information
```

### Examples

1. Using HTTP proxy with authentication:

```bash
./undetectable-browser -proxy http://username:password@proxy.example.com:8080
```

2. Using SOCKS5 proxy:

```bash
./undetectable-browser -proxy socks5://proxy.example.com:1080
```

3. Custom profile directory:

```bash
./undetectable-browser -proxy localhost:8080 -user-dir my_profile
```

4. Silent mode:

```bash
./undetectable-browser -proxy localhost:8080 -quiet
```

## Profile Management

The browser creates separate profiles for each proxy configuration by default. Profiles are stored in the `data/` directory:

```
data/
├── proxy.example.com/       # Default profile name (proxy hostname)
├── custom_profile/         # Custom profile specified with -user-dir
└── ...
```

Each profile maintains its own:

- Cookies
- Local Storage
- Browser Settings
- Extensions (if added)

## Development

### Project Structure

```
.
├── .github/
│   └── workflows/          # GitHub Actions workflows
├── data/                   # Browser profiles directory
├── main.go                 # Main application code
├── README.md              # This file
└── scripts/
    └── release.sh         # Release helper script
```

### Creating a Release

1. Update your code and commit changes
2. Use the release script:

```bash
./scripts/release.sh v1.0.0 "Release message"
```

Or manually:

```bash
git tag v1.0.0
git push && git push --tags
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Rod](https://go-rod.github.io/) - The browser automation framework
- [Chromium](https://www.chromium.org/) - The underlying browser engine

## Disclaimer

This tool is intended for legitimate purposes such as web testing, development, and research. Users are responsible for complying with all applicable laws and website terms of service. The authors are not responsible for any misuse of this software.

## Support

- 🐛 [Report Bug](https://github.com/sacsbrainz/undetectable-browser/issues)
- ✨ [Request Feature](https://github.com/sacsbrainz/undetectable-browser/issues)
