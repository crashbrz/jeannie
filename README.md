# jeannie (Opsgenie API Key Validator)

jeannie is a powerful tool for validating Opsgenie API keys, ensuring they are functional and properly authenticated. This tool supports single API key validation or batch processing using files, with configurable options for color-coded output, debugging, and parallel processing.

## Features

- Validate single or multiple API keys efficiently.
- Use color-coded output to distinguish valid and invalid keys.
- Debug mode to display HTTP response bodies.
- Multi-threaded validation with configurable goroutines.

## Installation

```bash
$ git clone https://github.com/crashbrz/jeannie.git
$ cd jeannie
$ go build -o jeannie main.go
```

## Usage

```bash
$ ./jeannie [options]
```

### Options

- `-k`: Validate a single API key.
- `-f`: Provide a file containing API keys (one per line) for batch validation.
- `-e`: Specify the Opsgenie endpoint. Default: `https://api.opsgenie.com/v1/json/cloudwatch?apiKey=`.
- `-remove-color`: Disable color output.
- `-d`: Display invalid API keys.
- `-debug`: Display the HTTP response body for each API key.
- `-t`: Set the number of goroutines to use for parallel validation. Default: 1.

### Examples

**Validate a single API key:**

```bash
$ ./jeannie -k your_api_key_here
```

**Validate API keys from a file with 4 goroutines:**

```bash
$ ./jeannie -f api_keys.txt -t 4
```

**Debug invalid API keys:**

```bash
$ ./jeannie -f api_keys.txt -d -debug
```

### Output

- `Valid`: Green output indicates valid API keys.
- `Invalid`: Red output indicates invalid API keys.
- Debug output displays the full HTTP response body when `-debug` is enabled.

### License

jeannie is licensed under the SushiWare license. For more information, check [docs/license.txt](docs/license.txt).
