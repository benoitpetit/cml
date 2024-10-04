
# Check My Logs (CML)

## Overview

**Check My Logs (CML)** is a command-line tool designed for viewing, filtering, and exporting log files. It offers an intuitive interface with pagination, real-time log monitoring, filtering, and exporting capabilities, making it ideal for developers and system administrators who need to efficiently manage and analyze log files.

## Key Features

- **Pagination**: Navigate through logs with customizable page sizes.
- **Filtering**: Filter logs by level (e.g., INFO, WARNING, ERROR) and search for specific terms.
- **Exporting**: Export filtered logs to a file for further analysis or reporting.
- **Live Mode**: Follow log files in real-time as new entries are added.
- **Search & Filter Combination**: Combine log level filtering with keyword searches to narrow down relevant log entries.
- **Interactive Interface**: Easy-to-use interface with keyboard navigation and search functionality.

## Installation

### Method 1: Using `go install`

To install **CML** using Go, ensure you have Go installed:

```bash
go version
```

If Go is not installed, download it from [golang.org](https://golang.org/dl/) and follow the installation instructions.

Install **CML** with:

```bash
go install github.com/benoitpetit/cml@latest
```

This command installs CML in your `$GOPATH/bin`.

### Method 2: Using the Installation Script

You can install **CML** by running the installation script directly from the repository:

- Using **curl**:

  ```bash
  bash <(curl -s https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/install.sh)
  ```

- Using **wget**:

  ```bash
  bash <(wget -qO - https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/install.sh)
  ```

**Note**: Ensure that the URLs point correctly to the `/hack/install.sh` script in your repository. Replace `master` with the appropriate branch if necessary.

## Uninstallation

### Using the Uninstallation Script

To uninstall **CML**, run the uninstallation script from the repository:

- Using **curl**:

  ```bash
  bash <(curl -s https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/uninstall.sh)
  ```

- Using **wget**:

  ```bash
  bash <(wget -qO - https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/uninstall.sh)
  ```

### Manual Uninstallation

To manually remove **CML**, follow these steps:

1. Check if **CML** is installed:

   ```bash
   ls /usr/local/bin/cml
   ```

2. If the binary exists, remove it with:

   ```bash
   sudo rm /usr/local/bin/cml
   ```

## Usage

Run **CML** with the following command:

```bash
cml <file_path> [options]
```

### Options

- `--filter <level>`: Filter logs by level (e.g., INFO, WARNING, ERROR).
- `--search <term>`: Search for a specific term in the logs.
- `--pagesize <size>`: Set the number of log entries to display per page (default is 15).
- `--export <export_path>`: Export filtered logs to the specified file.
- `--live`: Enable live mode to follow the log file in real-time as new entries are added.
- `--help`: Display the help message.

### Examples

1. **Filter and Export**:

   To filter logs from `logs.txt`, search for "timeout", set the page size to 10, and export the results:

   ```bash
   cml logs.txt --filter ERROR --search "timeout" --pagesize 10 --export filtered_logs.txt
   ```

2. **Live Mode**:

   To monitor `logs.txt` in live mode, displaying new logs as they are added:

   ```bash
   cml logs.txt --live
   ```

3. **Combination of Options**:

   You can combine multiple options to tailor your log viewing experience:

   ```bash
   cml logs.txt --filter WARNING --live --pagesize 20
   ```

## Interactive Controls

When using **CML**, you can navigate logs with the following keyboard controls:

- **Enter or Down Arrow**: Go to the next page of logs.
- **Up Arrow**: Go to the previous page of logs.
- **Home**: Jump to the first page.
- **End**: Jump to the last page.
- **Ctrl+S**: Start a search within the logs.
- **Esc**: Cancel the search.
- **Q or Ctrl+C**: Quit the application.

## Exporting Logs

To export filtered logs to a file, use the `--export` option. For example, after applying filters or searches, export the results:

```bash
cml logs.txt --filter INFO --export exported_logs.txt
```

This will save the filtered logs to `exported_logs.txt`.

## Troubleshooting

If you encounter issues during installation or use, here are some common solutions:

1. **Binary Not Found**: Check if the download URL is correct and if the binary exists in the expected directory.
2. **Permission Issues**: Ensure **CML** has executable permissions. If needed, adjust permissions using `chmod +x /usr/local/bin/cml`.
3. **Missing Dependencies**: If certain commands fail due to missing dependencies (e.g., `zip`, `curl`), install the necessary tools using your package manager (e.g., `apt`, `brew`, `dnf`).

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Clone your fork locally:
3. Make your changes, and push to your fork.
4. Create a pull request from your fork to the main repository.

Feel free to report issues or suggest improvements.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
