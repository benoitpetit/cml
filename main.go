package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

// Constants for pagination and default settings
const (
	defaultPageSize = 15 // Default number of log entries per page
)

// Model structure to hold application state
type model struct {
	logs        []string // Logs to display
	filterLevel string   // Log level to filter
	searchTerm  string   // Term to search in logs
	currentPage int      // Current page of logs
	pageSize    int      // Number of logs per page
	totalPages  int      // Total number of pages
}

// Init initializes the application
func (m model) Init() tea.Cmd {
	return nil
}

// Update processes user input and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": // Exit the application
			return m, tea.Quit
		case "enter": // Move to the next page
			if m.currentPage < m.totalPages-1 {
				m.currentPage++
			}
		case "up": // Move to the previous page
			if m.currentPage > 0 {
				m.currentPage--
			}
		case "down": // Move to the next page
			if m.currentPage < m.totalPages-1 {
				m.currentPage++
			}
		case "home": // Go to the first page
			m.currentPage = 0
		case "end": // Go to the last page
			m.currentPage = m.totalPages - 1
		}
	}
	return m, nil
}

// View renders the current state of the application
func (m model) View() string {
	var b strings.Builder
	start := m.currentPage * m.pageSize
	end := min(start+m.pageSize, len(m.logs))

	// Display logs or a message if no logs are found
	if start < len(m.logs) {
		for _, log := range m.logs[start:end] {
			b.WriteString(log + "\n")
		}
		if end < len(m.logs) {
			b.WriteString("--------------------------------------------------------------------------------------------------------\n")
			b.WriteString("Press [Enter] to see more logs, [Up] for previous, [Down] for next, [Home] for first, or [End] for last.\n")
		}
	} else {
		b.WriteString("No logs found.\n")
	}

	return b.String()
}

// clearScreen clears the terminal screen using ANSI escape codes
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// formatLogLine formats a log line for display, including color coding.
func formatLogLine(line string) string {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		return line // Return as-is if format is not respected
	}

	dateStr := parts[0] + " " + parts[1]
	level := parts[2]
	message := parts[3]

	// Format the date if parsing is successful
	if date, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
		dateStr = date.Format("2006-01-02 15:04:05")
	}

	// Get color function based on log level
	colorFunc := logColors[level]
	if colorFunc == nil {
		colorFunc = color.New(color.Reset) // Default color if level is unknown
	}

	// Format and return the log line with color
	return fmt.Sprintf("%s %s: %s", colorFunc.Sprint(dateStr), level, message)
}

// loadLogs loads logs from the specified file, applying filters as necessary.
func loadLogs(filePath, filterLevel, searchTerm string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		color.Red("Error: The specified file does not exist.")
		return nil
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
		color.Red("Error reading file: %v", err)
	}

	return filteredLogs
}

// matchesFilters checks if the log line meets the filtering criteria
func matchesFilters(line, filterLevel, searchTerm string) bool {
	if filterLevel != "" && !strings.Contains(line, filterLevel) {
		return false
	}
	if searchTerm != "" && !strings.Contains(line, searchTerm) {
		return false
	}
	return true
}

// exportLogs exports the filtered logs to the specified file.
func exportLogs(logs []string, exportPath string) error {
	file, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("could not create export file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, log := range logs {
		_, err := writer.WriteString(log + "\n")
		if err != nil {
			return fmt.Errorf("could not write log to file: %v", err)
		}
	}

	return writer.Flush()
}

// Helper function to get the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	clearScreen() // Clear the terminal screen

	if len(os.Args) < 2 {
		fmt.Println("Usage: cml <file_path> [--filter <level>] [--search <term>] [--pagesize <size>] [--export <export_path>]")
		return
	}

	filePath := os.Args[1]
	var filterLevel, searchTerm, exportPath string
	pageSize := defaultPageSize // Default page size

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
					fmt.Println("Invalid page size. Using default value of 10.")
					pageSize = defaultPageSize // Default value if invalid
				}
				i++
			}
		case "--export":
			if i+1 < len(os.Args) {
				exportPath = os.Args[i+1]
				i++
			}
		}
	}

	logs := loadLogs(filePath, filterLevel, searchTerm)

	// Calculate total pages
	totalPages := (len(logs) + pageSize - 1) / pageSize

	m := model{
		logs:        logs,
		filterLevel: filterLevel,
		searchTerm:  searchTerm,
		pageSize:    pageSize,
		totalPages:  totalPages,
	}

	p := tea.NewProgram(m)
	_, err := p.Run() // Handle both return values but ignore the model
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Export logs if an export path is provided
	if exportPath != "" {
		if err := exportLogs(logs, exportPath); err != nil {
			color.Red("Error exporting logs: %v", err)
		} else {
			fmt.Printf("Logs successfully exported to %s\n", exportPath)
		}
	}
}

// logColors maps log levels to corresponding colors.
var logColors = map[string]*color.Color{
	"INFO":    color.New(color.FgGreen),
	"WARNING": color.New(color.FgYellow),
	"ERROR":   color.New(color.FgRed),
	"DEBUG":   color.New(color.FgCyan),
}
