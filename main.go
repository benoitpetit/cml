package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hpcloud/tail"
)

// Constants for pagination and default settings
const (
	defaultPageSize  = 15              // Default number of log entries per page
	tailPollInterval = 1 * time.Second // Polling interval if necessary
)

// Struct to represent the application's state
type model struct {
	logs        []string   // Logs to display
	filterLevel string     // Log level to filter
	searchTerm  string     // Search term in logs
	currentPage int        // Current log page
	pageSize    int        // Number of logs per page
	totalPages  int        // Total number of pages
	live        bool       // Indicator for live mode
	filePath    string     // Path to the log file
	tailer      *tail.Tail // Tail instance
	err         error      // Potential error
}

// Message types for Bubble Tea
type (
	newLogMsg struct {
		log string
	}
	errMsg error
)

// Init initializes the application
func (m *model) Init() tea.Cmd {
	if m.live {
		return m.startTailing()
	}
	return nil
}

// Update handles messages and updates the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Handle keyboard key messages
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "down":
			if m.currentPage < m.totalPages-1 {
				m.currentPage++
			}
		case "up":
			if m.currentPage > 0 {
				m.currentPage--
			}
		case "home":
			m.currentPage = 0
		case "end":
			if m.totalPages > 0 {
				m.currentPage = m.totalPages - 1
			}
		}

	// Handle new logs in live mode
	case newLogMsg:
		if matchesFilters(msg.log, m.filterLevel, m.searchTerm) {
			formatted := formatLogLine(msg.log)
			m.logs = append(m.logs, formatted)
			m.totalPages = (len(m.logs) + m.pageSize - 1) / m.pageSize
			// If in live mode, automatically move to the last page
			if m.live {
				m.currentPage = m.totalPages - 1
			}
		}
		// If in live mode, queue the next log
		if m.live {
			return m, watchTail(m.tailer)
		}

	// Handle errors
	case errMsg:
		m.err = fmt.Errorf("%v", msg)
		return m, tea.Quit
	}

	return m, nil
}

// View displays the current state of the application
func (m *model) View() string {
	var b strings.Builder

	// Display paginated logs
	start := m.currentPage * m.pageSize
	end := min(start+m.pageSize, len(m.logs))

	if start < len(m.logs) {
		for _, log := range m.logs[start:end] {
			b.WriteString(log + "\n")
		}
	} else {
		b.WriteString("No logs found.\n")
	}

	// Add a separator and instructions if necessary
	b.WriteString("\n--------------------------------------------------------------------------------------------------------\n")
	if m.live {
		b.WriteString("Live mode enabled. Waiting for new logs...\n")
	} else if m.totalPages > 1 {
		b.WriteString("Press [Enter]/[Down] to view more logs, [Up] for previous, [Home] for first, [End] for last, or [Q] to quit.\n")
	} else {
		b.WriteString("Press [Q] or [Ctrl+C] to quit.\n")
	}

	// Display any potential errors
	if m.err != nil {
		b.WriteString(fmt.Sprintf("\nError: %v\n", m.err))
	}

	return b.String()
}

// Function to start live tailing
func (m *model) startTailing() tea.Cmd {
	t, err := tail.TailFile(m.filePath, tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Poll:      true,
		Logger:    tail.DiscardingLogger,
	})
	if err != nil {
		return func() tea.Msg {
			return errMsg(fmt.Errorf("error opening log file: %v", err))
		}
	}

	m.tailer = t

	return watchTail(t)
}

// Function to read a log line and send a message
func watchTail(t *tail.Tail) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-t.Lines
		if !ok {
			return errMsg(fmt.Errorf("log file tailing was interrupted"))
		}
		if line.Err != nil {
			return errMsg(fmt.Errorf("error reading log file: %v", line.Err))
		}
		return newLogMsg{log: line.Text}
	}
}

// Function to format a log line
func formatLogLine(line string) string {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		return line // Return as-is if format is not respected
	}

	dateStr := parts[0] + " " + parts[1]
	level := parts[2]
	message := parts[3]

	// Format the date if parsing succeeds
	if date, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
		dateStr = date.Format("2006-01-02 15:04:05")
	}

	// Format and return the log line
	return fmt.Sprintf("%s %s: %s", dateStr, level, message)
}

// Function to check if a log line matches the filters
func matchesFilters(line, filterLevel, searchTerm string) bool {
	if filterLevel != "" && !strings.Contains(line, filterLevel) {
		return false
	}
	if searchTerm != "" && !strings.Contains(line, searchTerm) {
		return false
	}
	return true
}

// Function to load logs from the file, applying filters
func loadLogs(filePath, filterLevel, searchTerm string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("the specified file does not exist or cannot be opened")
	}
	defer file.Close()

	var filteredLogs []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if matchesFilters(line, filterLevel, searchTerm) {
			filteredLogs = append(filteredLogs, formatLogLine(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading the file: %v", err)
	}

	return filteredLogs, nil
}

// Function to export filtered logs to a specified file
func exportLogs(logs []string, exportPath string) error {
	file, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("unable to create export file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, log := range logs {
		_, err := writer.WriteString(log + "\n")
		if err != nil {
			return fmt.Errorf("unable to write to export file: %v", err)
		}
	}

	return writer.Flush()
}

// Function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Function to display the help message
func printHelp() {
	help := `
Check My Logs (CML)

Usage:
  cml <file_path> [options]

Options:
  --filter <level>         Filter logs by level (e.g., INFO, WARNING, ERROR).
  --search <term>          Search for a specific term in the logs.
  --pagesize <size>        Set the number of log entries to display per page (default 15).
  --export <export_path>   Export filtered logs to the specified file.
  --live                   Enable live mode to follow the log file in real-time.
  --help                   Display this help message.

Examples:
  1. Filter and Export:
     To filter logs from logs.txt, search for "timeout", set the page size to 10, and export the results:
     cml logs.txt --filter ERROR --search "timeout" --pagesize 10 --export filtered_logs.txt

  2. Live Mode:
     To monitor logs.txt in live mode, displaying new logs as they are added:
     cml logs.txt --live

  3. Combination of Options:
     You can also combine options:
     cml logs.txt --filter WARNING --live --pagesize 20

Controls:
  Enter or Down Arrow : Go to the next log page.
  Up Arrow            : Go to the previous log page.
  Home                : Go to the first page.
  End                 : Go to the last page.
  Q or Ctrl+C         : Quit the application.
`
	fmt.Println(help)
}

func main() {
	clearScreen() // Clear the terminal screen

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	// Check if help is requested
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			printHelp()
			return
		}
	}

	filePath := os.Args[1]
	var filterLevel, searchTerm, exportPath string
	pageSize := defaultPageSize // Default page size
	live := false

	// Parse command line arguments for filters and export path
	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--filter":
			if i+1 < len(os.Args) {
				filterLevel = os.Args[i+1]
				i++
			}
		case "--search":
			if i+1 < len(os.Args) {
				searchTerm = os.Args[i+1]
				i++
			}
		case "--pagesize":
			if i+1 < len(os.Args) {
				var err error
				pageSize, err = strconv.Atoi(os.Args[i+1])
				if err != nil || pageSize <= 0 {
					fmt.Println("Invalid page size. Using default value of 15.")
					pageSize = defaultPageSize // Default value if invalid
				}
				i++
			}
		case "--export":
			if i+1 < len(os.Args) {
				exportPath = os.Args[i+1]
				i++
			}
		case "--live":
			live = true
		}
	}

	// Load logs initially
	logs, err := loadLogs(filePath, filterLevel, searchTerm)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Calculate total number of pages
	totalPages := (len(logs) + pageSize - 1) / pageSize

	m := &model{
		logs:        logs,
		filterLevel: filterLevel,
		searchTerm:  searchTerm,
		pageSize:    pageSize,
		totalPages:  totalPages,
		live:        live,
		filePath:    filePath,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Export logs if an export path is provided
	if exportPath != "" {
		finalLogs := finalModel.(*model).logs
		if err := exportLogs(finalLogs, exportPath); err != nil {
			fmt.Printf("Error exporting logs: %v\n", err)
		} else {
			fmt.Printf("Logs successfully exported to %s\n", exportPath)
		}
	}
}

// Function to clear the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
