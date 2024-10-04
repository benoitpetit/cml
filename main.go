package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hpcloud/tail"
)

// Constants for pagination and default settings
const (
	defaultPageSize  = 15
	tailPollInterval = 1 * time.Second
)

// Struct to represent the application's state
type model struct {
	file        *os.File
	logs        []string
	filterLevel string
	searchTerm  string
	currentPage int
	pageSize    int
	totalPages  int
	live        bool
	filePath    string
	tailer      *tail.Tail
	err         error
	lineOffsets []int64

	// UI Components
	spinner     spinner.Model
	searchInput textinput.Model
	searching   bool

	// Width of the terminal window
	windowWidth int
}

// Message types for Bubble Tea
type (
	newLogMsg struct {
		log string
	}
	errMsg error
)

// Define styles using lipgloss
var (
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("32")) // Green
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("33")) // Yellow
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("31")) // Red
	defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("37")) // White
)

// Init initializes the application
func (m *model) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.live {
		cmds = append(cmds, m.startTailing())
	}

	// Initialize spinner
	m.spinner = spinner.New()
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Orange
	cmds = append(cmds, m.spinner.Tick)

	// Initialize search input
	m.searchInput = textinput.New()
	m.searchInput.Placeholder = "Search..."
	m.searchInput.CharLimit = 100
	m.searchInput.Width = 30

	return tea.Batch(cmds...)
}

// Update handles messages and updates the model
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "enter", "esc":
				m.searching = false
				m.searchTerm = m.searchInput.Value()
				m.currentPage = 0
				// Re-index with the new search term
				lineOffsets, err := indexFilteredLineOffsets(m.filePath, m.filterLevel, m.searchTerm)
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				m.lineOffsets = lineOffsets
				m.totalPages = (len(m.lineOffsets) + m.pageSize - 1) / m.pageSize
				logs, err := m.loadPage(m.currentPage)
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				m.logs = logs
			}
		} else {
			switch strings.ToLower(msg.String()) {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "ctrl+s":
				m.searching = true
				m.searchInput.Focus()
				cmds = append(cmds, textinput.Blink)
			case "enter", "down":
				if m.currentPage < m.totalPages-1 {
					m.currentPage++
					logs, err := m.loadPage(m.currentPage)
					if err != nil {
						m.err = err
						return m, tea.Quit
					}
					m.logs = logs
				}
			case "up":
				if m.currentPage > 0 {
					m.currentPage--
					logs, err := m.loadPage(m.currentPage)
					if err != nil {
						m.err = err
						return m, tea.Quit
					}
					m.logs = logs
				}
			case "home":
				m.currentPage = 0
				logs, err := m.loadPage(m.currentPage)
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				m.logs = logs
			case "end":
				if m.totalPages > 0 {
					m.currentPage = m.totalPages - 1
					logs, err := m.loadPage(m.currentPage)
					if err != nil {
						m.err = err
						return m, tea.Quit
					}
					m.logs = logs
				}
			}
		}

	case newLogMsg:
		if m.live && matchesFilters(msg.log, m.filterLevel, m.searchTerm) {
			formatted := formatLogLine(msg.log)
			m.logs = append(m.logs, formatted)

			// Limit the number of logs displayed in live mode
			if len(m.logs) > m.pageSize {
				// Remove the oldest log
				m.logs = m.logs[1:]
			}

			// Update line offsets if necessary
			newOffset, err := m.getCurrentFileOffset()
			if err == nil {
				m.lineOffsets = append(m.lineOffsets, newOffset)
				m.totalPages = (len(m.lineOffsets) + m.pageSize - 1) / m.pageSize
				// No need to reload logs from file in live mode
			}
		}
		if m.live {
			return m, watchTail(m.tailer)
		}

	case errMsg:
		m.err = fmt.Errorf("%v", msg)
		return m, tea.Quit

	case tea.WindowSizeMsg:
		// Adjust the size if necessary
		m.windowWidth = msg.Width // Capture the current window width
	}

	// Update the spinner and add the cmd to the slice cmds
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	// Update the search field if in search mode
	if m.searching {
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func wrapLog(log string, width int) string {
	if width <= 0 {
		return log // No wrapping if width is incorrect
	}

	var wrappedLog strings.Builder
	words := strings.Fields(log)
	currentLineLength := 0

	for _, word := range words {
		wordLength := len(word)

		// If the word can't be added to the current line, move to the next line
		if currentLineLength+wordLength+1 > width {
			wrappedLog.WriteString("\n") // Move to the next line
			currentLineLength = 0
		}

		// Add a space between words unless it's the start of a new line
		if currentLineLength > 0 {
			wrappedLog.WriteString(" ")
			currentLineLength++
		}

		wrappedLog.WriteString(word)
		currentLineLength += wordLength
	}

	return wrappedLog.String()
}

// Function to generate the header view
func headerView(m *model) string {
	var b strings.Builder

	// Enhanced header with modern and innovative lipgloss styles
	// Primary Header Style
	headerTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("62")).
		Bold(true).
		Padding(1, 2).
		MarginBottom(1)

	// Status Bar Style
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("60"))

	// Navigation Instructions Style
	navStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Background(lipgloss.Color("238")).
		Italic(true).
		Padding(0, 1).
		MarginBottom(1)

	// Construct the header title
	headerTitle := headerTitleStyle.Render("Check My Logs - Monitoring Logs Tool")

	// Construct the status bar
	status := fmt.Sprintf("%s Page %d/%d | Total Logs: %d | Live Mode: %v",
		func() string {
			if m.live {
				return m.spinner.View()
			} else {
				return ""
			}
		}(), m.currentPage+1, m.totalPages, len(m.lineOffsets), m.live)
	statusBar := statusBarStyle.Render(status)

	// Construct the navigation instructions
	var navInstructions string
	if m.searching {
		navInstructions = "Press [Enter] to apply search, [Esc] to cancel."
	} else if m.live {
		navInstructions = "Press [Ctrl+S] to search, [Q] or [Ctrl+C] to quit."
	} else if m.totalPages > 1 {
		navInstructions = "Press [Ctrl+S] to search, [Enter]/[Down] for next page, [Up] for previous page, [Home] for first page, [End] for last page, [Q] to quit."
	} else {
		navInstructions = "Press [Ctrl+S] to search, [Q] or [Ctrl+C] to quit."
	}
	instructions := navStyle.Render(navInstructions)

	// Write the enhanced header
	b.WriteString(headerTitle + "\n")
	b.WriteString(statusBar + "\n")
	b.WriteString(instructions + "\n")

	// Separator line with a modern look
	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true).
		Render(strings.Repeat("━", m.windowWidth))
	b.WriteString(separatorStyle + "\n")

	return b.String()
}

// View displays the current state of the application
func (m *model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	var b strings.Builder

	// Display the header first
	b.WriteString(headerView(m))

	// If searching, display the search field
	if m.searching {
		b.WriteString("Search: " + m.searchInput.View() + "\n\n")
	} else {
		// Display logs with wrapping if necessary
		for _, log := range m.logs {
			wrappedLog := wrapLog(log, m.windowWidth) // Apply wrapping
			b.WriteString(wrappedLog + "\n")
		}
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

// Function to format a log line with colors based on log level
func formatLogLine(line string) string {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		return defaultStyle.Render(line)
	}

	dateStr := parts[0] + " " + parts[1]
	level := parts[2]
	message := parts[3]

	// Format the date if parsing succeeds
	if date, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
		dateStr = date.Format("2006-01-02 15:04:05")
	}

	// Apply style based on level
	var styledLevel string
	switch strings.ToUpper(level) {
	case "INFO":
		styledLevel = infoStyle.Render(level)
	case "WARNING", "WARN":
		styledLevel = warningStyle.Render(level)
	case "ERROR", "ERR":
		styledLevel = errorStyle.Render(level)
	default:
		styledLevel = defaultStyle.Render(level)
	}

	return fmt.Sprintf("%s %s: %s", dateStr, styledLevel, message)
}

// Function to check if a log line matches the filters
func matchesFilters(line, filterLevel, searchTerm string) bool {
	if filterLevel != "" && !strings.Contains(strings.ToUpper(line), strings.ToUpper(filterLevel)) {
		return false
	}
	if searchTerm != "" && !strings.Contains(strings.ToUpper(line), strings.ToUpper(searchTerm)) {
		return false
	}
	return true
}

// Function to load logs for a specific page
func (m *model) loadPage(page int) ([]string, error) {
	if m.totalPages == 0 {
		return []string{}, nil // Return an empty slice without error
	}

	if page < 0 || page >= m.totalPages {
		return nil, fmt.Errorf("invalid page number")
	}

	startLine := page * m.pageSize
	endLine := min(startLine+m.pageSize, len(m.lineOffsets))

	var logs []string
	for i := startLine; i < endLine; i++ {
		line, err := readLineAt(m.file, m.lineOffsets[i])
		if err != nil {
			return nil, err
		}
		formatted := formatLogLine(line)
		logs = append(logs, formatted)
	}

	return logs, nil
}

// Function to read a specific line at a given offset
func readLineAt(file *os.File, offset int64) (string, error) {
	_, err := file.Seek(offset, 0)
	if err != nil {
		return "", fmt.Errorf("error seeking in file: %v", err)
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	return strings.TrimRight(line, "\n"), nil
}

// Function to index the positions of lines that match filters
func indexFilteredLineOffsets(filePath, filterLevel, searchTerm string) ([]int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	var offsets []int64
	var offset int64 = 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if matchesFilters(line, filterLevel, searchTerm) {
			offsets = append(offsets, offset)
		}
		offset += int64(len(line)) + 1 // +1 for the newline character
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return offsets, nil
}

// Function to get the current file offset
func (m *model) getCurrentFileOffset() (int64, error) {
	offset, err := m.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, fmt.Errorf("error retrieving offset: %v", err)
	}
	return offset, nil
}

// Function to export logs
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
  Ctrl+S                : Start searching
  Enter or Down Arrow   : Go to the next log page.
  Up Arrow              : Go to the previous log page.
  Home                  : Go to the first page.
  End                   : Go to the last page.
  Q or Ctrl+C           : Quit the application.
  Esc                   : Cancel search
`
	fmt.Println(help)
}

func main() {
	clearScreen()

	// Define command-line flags
	filterLevel := flag.String("filter", "", "Filter logs by level (e.g., INFO, WARNING, ERROR).")
	searchTerm := flag.String("search", "", "Search for a specific term in the logs.")
	pageSize := flag.Int("pagesize", defaultPageSize, "Set the number of log entries to display per page (default 15).")
	exportPath := flag.String("export", "", "Export filtered logs to the specified file.")
	live := flag.Bool("live", false, "Enable live mode to follow the log file in real-time.")
	helpFlag := flag.Bool("help", false, "Display this help message.")

	// Parse the flags
	flag.Parse()

	// Display help if requested
	if *helpFlag {
		printHelp()
		return
	}

	// Ensure that the file path is provided
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: File path is required.")
		printHelp()
		return
	}

	filePath := args[0]

	// Load filtered logs
	lineOffsets, err := indexFilteredLineOffsets(filePath, *filterLevel, *searchTerm)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Calculate total number of pages
	totalPages := 0
	if len(lineOffsets) > 0 {
		totalPages = (len(lineOffsets) + *pageSize - 1) / *pageSize
	}

	// Open the log file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	// Initialize the model
	m := &model{
		file:        file,
		lineOffsets: lineOffsets,
		filterLevel: *filterLevel,
		searchTerm:  *searchTerm,
		pageSize:    *pageSize,
		totalPages:  totalPages,
		live:        *live,
		filePath:    filePath,
	}

	// Load the first page if possible
	if totalPages > 0 {
		logs, err := m.loadPage(0)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		m.logs = logs
		m.currentPage = 0
	}

	// Initialize the spinner and search input in Init
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Export logs if an export path is provided
	if *exportPath != "" {
		finalLogs := finalModel.(*model).logs
		if err := exportLogs(finalLogs, *exportPath); err != nil {
			fmt.Printf("Error exporting logs: %v\n", err)
		} else {
			fmt.Printf("Logs successfully exported to %s\n", *exportPath)
		}
	}
}

// Function to clear the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
