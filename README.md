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

### Prerequisites

- **Operating System**: Linux
- **Required Tools**:
  - `wget` (download tool)
  - `tar` (extraction tool)
- **Administrator Access**: You must have `sudo` privileges to install software and move files into system directories.

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

This command installs CML in your `$GOPATH/bin`. Ensure that `$GOPATH/bin` is included in your system's `$PATH` to execute `cml` from anywhere.

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

**Note**: Ensure that the URLs correctly point to the `/hack/install.sh` script in the repository. Replace `master` with the appropriate branch if necessary.

### Method 3: Manual Installation

Follow the steps below to manually install **CML** on your Linux system.

#### Installation Steps

1. **Download the CML Binary:**

   Use `wget` to download the CML tarball from the official GitHub repository.

   ```bash
   wget -O /tmp/cml.tar.gz https://github.com/benoitpetit/cml/releases/download/v1.2.0/cml-linux-amd64.tar.gz
   ```

   *Replace `v1.2.0` with the desired version if a newer release is available.*

2. **Extract the Binary and Move to `/usr/local/bin`:**

   Use `tar` to extract the downloaded file and move the `cml` binary to a directory included in your `$PATH` for easy execution.

   ```bash
   sudo tar -xzvf /tmp/cml.tar.gz -C /usr/local/bin/
   ```

3. **Make the Binary Executable:**

   Ensure that the `cml` binary has the necessary execution permissions.

   ```bash
   sudo chmod +x /usr/local/bin/cml
   ```

4. **Clean Up Temporary Files:**

   Remove the downloaded tarball to free up space.

   ```bash
   rm /tmp/cml.tar.gz
   ```

5. **Verify the Installation:**

   Confirm that **CML** has been installed correctly by checking its version.

   ```bash
   cml --version
   ```

   You should see output similar to:

   ```
   cml version v1.2.0
   ```

---

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

**Note**: Ensure that the URLs correctly point to the `/hack/uninstall.sh` script in the repository. Replace `master` with the appropriate branch if necessary.

### Manual Uninstallation

To manually remove **CML**, follow these steps:

1. **Check if CML is Installed:**

   Verify the presence of the `cml` binary in `/usr/local/bin`.

   ```bash
   ls /usr/local/bin/cml
   ```

2. **Remove the Binary if It Exists:**

   If the `cml` binary is found, remove it using the following command:

   ```bash
   sudo rm /usr/local/bin/cml
   ```

---

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

1. **Binary Not Found**: Check if the download URL is correct and if the binary exists in the expected directory (`/usr/local/bin`).

2. **Permission Issues**: Ensure **CML** has executable permissions. If needed, adjust permissions using:

   ```bash
   sudo chmod +x /usr/local/bin/cml
   ```

3. **Missing Dependencies**: If certain commands fail due to missing dependencies (e.g., `wget`, `tar`), install the necessary tools using your package manager. For example:

   - **Debian/Ubuntu**:

     ```bash
     sudo apt update
     sudo apt install -y wget tar
     ```

   - **Red Hat/Fedora**:

     ```bash
     sudo dnf install -y wget tar
     ```

   - **openSUSE**:

     ```bash
     sudo zypper install -y wget tar
     ```

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Clone your fork locally:
   
   ```bash
   git clone https://github.com/your-username/cml.git
   ```
   
3. Make your changes, and push to your fork.
4. Create a pull request from your fork to the main repository.

Feel free to report issues or suggest improvements.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
