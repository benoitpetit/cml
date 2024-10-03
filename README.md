# Check My Logs (CML)

## Overview

**CML** is a command-line tool for viewing, filtering, and exporting log files. It provides a user-friendly interface for navigating log entries and applying various filters to quickly find relevant information.

## Features

- **Pagination**: Navigate through logs with customizable page sizes.
- **Filtering**: Filter logs by log level (e.g., INFO, WARNING, ERROR) and search terms.
- **Exporting**: Export filtered logs to a specified file for further analysis.
- **Live Mode**: Follow log files in real-time as new entries are added.

## Installation

### Method 1: Using `go install`

Ensure you have Go installed on your machine. You can verify this by running:

```bash
go version
```

If Go is not installed, download it from [golang.org](https://golang.org/dl/) and follow the installation instructions.

Run the following command to install CML:

```bash
go install github.com/benoitpetit/cml@latest
```

This command will install CML and make it available in your `$GOPATH/bin`.

### Method 2: Using the Installation Script

You can also install CML by downloading and executing the installation script directly from your terminal. The scripts are now located in the `/hack` directory of the repository.

- Using **curl**:

  ```bash
  bash <(curl -s https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/install.sh)
  ```

- Using **wget**:

  ```bash
  bash <(wget -qO - https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/install.sh)
  ```

**Note:** Ensure that the URLs correctly point to the `/hack/install.sh` script in your repository. Replace `master` with the appropriate branch name if necessary.

## Uninstallation

To uninstall CML, you can run the uninstallation script provided in the `/hack` directory of the repository.

### Using the Uninstallation Script

Execute the uninstallation script directly from your terminal:

- Using **curl**:

  ```bash
  bash <(curl -s https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/uninstall.sh)
  ```

- Using **wget**:

  ```bash
  bash <(wget -qO - https://raw.githubusercontent.com/benoitpetit/cml/refs/heads/master/hack/uninstall.sh)
  ```

**Note:** Ensure that the URLs correctly point to the `/hack/uninstall.sh` script in your repository. Replace `master` with the appropriate branch name if necessary.

### Manual Uninstallation

If you prefer manual uninstallation, you can remove the binary directly:

1. **Check if CML is installed:**

   ```bash
   ls /usr/local/bin/cml
   ```

   If the binary exists, you can remove it.

2. **Remove the Binary:**

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
- `--live`: Enable live mode to follow the log file in real-time.
- `--help`: Display the help message.

### Examples

1. **Filter and Export:**

   To filter logs from `logs.txt`, search for "timeout", set the page size to 10, and export the results:

   ```bash
   cml logs.txt --filter ERROR --search "timeout" --pagesize 10 --export filtered_logs.txt
   ```

2. **Live Mode:**

   To monitor `logs.txt` in live mode, displaying new logs as they are added:

   ```bash
   cml logs.txt --live
   ```

3. **Combination of Options:**

   You can also combine multiple options:

   ```bash
   cml logs.txt --filter WARNING --live --pagesize 20
   ```

## Controls

While viewing logs, you can navigate using the following controls:

- **Enter or Down Arrow**: Go to the next log page.
- **Up Arrow**: Go to the previous log page.
- **Home**: Go to the first page.
- **End**: Go to the last page.
- **Q or Ctrl+C**: Quit the application.

## Troubleshooting

If you encounter issues during installation or uninstallation:

1. **Binary Not Found**: Ensure the download URL is correct and the binary exists.
2. **Permissions**: If you receive permission errors, check the executable permissions with `ls -l /usr/local/bin/cml` and adjust as necessary.
3. **Missing Dependencies**: If commands fail due to missing tools (e.g., `zip`), install them using your package manager.

## License

This project is licensed under the MIT License.