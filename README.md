# DirectoryListing Fuzzer

A powerful tool for discovering directory listings on websites. This tool helps security researchers and system administrators identify potentially exposed directory listings that could pose security risks.

## Features

- Fast directory listing discovery
- Multiple output formats (text, JSON, CSV)
- Custom timeout settings (2 seconds default)
- Automatic directory creation for output files
- Support for both HTTP and HTTPS
- Recursive path discovery
- Configurable verbosity

## Installation

```bash
# Clone the repository
git clone https://github.com/begininvoke/DirectoryListenFinder.git

# Change to project directory
cd DirectoryListenFinder

# Build the project
go build -o DirListerFuzzer
```

## Usage

Basic usage:
```bash
./DirListerFuzzer -url https://example.com
```

Advanced usage:
```bash
# Save results to JSON file
./DirListerFuzzer -url https://example.com -o results/scan.json -f json

# Show only successful results
./DirListerFuzzer -url https://example.com -v

# Save as CSV format
./DirListerFuzzer -url https://example.com -o results/scan.csv -f csv
```

### Command Line Options

```bash
Usage of ./DirListerFuzzer:
  -url string
        URL address (e.g., https://example.com)
  -v    
        Show success results only
  -o string
        Output file path (e.g., results/output.json)
  -f string
        Output format (text, json, csv) (default "text")
```

## Output Formats

### Text (Default)
Simple list of discovered URLs, one per line.

### JSON
Structured output including:
- URL
- HTTP Status
- Content Type

### CSV
Comma-separated values with headers:
- URL
- Status
- Content-Type

## Performance

- 2-second timeout per request
- Efficient duplicate URL filtering
- Concurrent scanning capabilities

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security Considerations

This tool should only be used on systems you have permission to test. Unauthorized scanning may be illegal in your jurisdiction.