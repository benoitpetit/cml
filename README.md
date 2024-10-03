# Check My Logs (CML)

## Overview

**CML** is a command-line tool for viewing, filtering, and exporting log files. It provides a user-friendly interface for navigating log entries and applying various filters to quickly find relevant information.

## Features

- **Pagination**: Navigate through logs with customizable page sizes.
- **Filtering**: Filter logs by log level (e.g., INFO, WARNING, ERROR) and search terms.
- **Exporting**: Export filtered logs to a specified file for further analysis.

## Installation

### Method 1: Using `go install`

Ensure you have Go installed on your machine. Run the following command to install CML:

```bash
go install github.com/benoitpetit/cml@latest
```

This command will install CML and make it available in your `$GOPATH/bin`.

### Method 2: Using the Installation Script

You can also install CML by downloading and executing the installation script directly from your terminal. Use one of the following commands:

- Using **curl**:

  ```bash
  bash <(curl -s https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/install.sh)
  ```

- Using **wget**:

  ```bash
  bash <(wget -qO - https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/install.sh)
  ```

## Uninstallation

To uninstall CML, you can run the uninstallation script provided in the repository.

### Using the Uninstallation Script

You can uninstall CML by executing the uninstallation script directly from your terminal. Use one of the following commands:

- Using **curl**:

  ```bash
  bash <(curl -s https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/uninstall.sh)
  ```

- Using **wget**:

  ```bash
  bash <(wget -qO - https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/uninstall.sh)
  ```

### Manual Uninstallation

If you prefer manual uninstallation, you can remove the binary directly:

1. **Check if CML is installed**:

   ```bash
   ls /usr/local/bin/cml
   ```

   If the binary exists, you can remove it.

2. **Remove the Binary**:

   ```bash
   sudo rm /usr/local/bin/cml
   ```

## Usage

Run the application with the following command format:

```bash
cml <file_path> [options]
```

### Options

- `--filter <level>`: Filter logs by level (e.g., INFO, WARNING, ERROR).
- `--search <term>`: Search for a specific term in the logs.
- `--pagesize <size>`: Set the number of log entries to display per page (default is 15).
- `--export <export_path>`: Export the filtered logs to the specified file.

### Example

To filter logs from `logs.txt`, searching for "timeout" and exporting results, use:

```bash
cml logs.txt --filter ERROR --search "timeout" --pagesize 10 --export filtered_logs.txt
```

## Controls

While viewing logs, you can navigate using the following controls:

- **Enter**: Move to the next page of logs.
- **Up Arrow**: Move to the previous page of logs.
- **Down Arrow**: Move to the next page of logs.
- **Home**: Go to the first page.
- **End**: Go to the last page.
- **Ctrl+C or Q**: Exit the application.

## Troubleshooting

If you encounter issues during installation or uninstallation:

1. **Binary Not Found**: Ensure the download URL is correct and the binary exists.
2. **Permissions**: If you receive permission errors, check the executable permissions with `ls -l /usr/local/bin/cml` and adjust as necessary.
3. **Missing Dependencies**: If commands fail due to missing tools (e.g., `wget`), install them using your package manager.

## License

This project is licensed under the MIT License.