package ls

import (
	"agent-dev-environment/src/api/v1/filesystem/ls"
	"agent-dev-environment/src/library/api"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Handler(req ls.Request) (*ls.Response, error) {
	// First verify the path exists
	fileInfo, err := os.Stat(req.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, api.NewError(api.NotFound, "Path not found")
		}
		return nil, err
	}

	var entries []ls.FileInfo

	// Handle single file case
	if !fileInfo.IsDir() {
		entries = append(entries, fileInfoToAPIFileInfo(fileInfo, req.Long))
		return &ls.Response{Entries: entries}, nil
	}

	// Use Linux ls command for directory listing
	entries, err = executeLinuxLS(req.Path, req.Recursive, req.Long)
	if err != nil {
		return nil, err
	}

	// Sort entries alphabetically by name
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return &ls.Response{Entries: entries}, nil
}

func executeLinuxLS(path string, recursive bool, long bool) ([]ls.FileInfo, error) {
	var args []string

	// Build ls command arguments
	if long {
		args = append(args, "-l")
	}
	if recursive {
		args = append(args, "-R")
	}
	args = append(args, "-a") // Include all files
	args = append(args, path)

	cmd := exec.Command("ls", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ls command failed: %s", stderr.String())
	}

	return parseLinuxLSOutput(stdout.String(), path, recursive, long)
}

func parseLinuxLSOutput(output string, basePath string, recursive bool, long bool) ([]ls.FileInfo, error) {
	var entries []ls.FileInfo
	lines := strings.Split(output, "\n")

	currentDir := basePath

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle recursive directory headers (e.g., "/path/to/dir:")
		if recursive && strings.HasSuffix(line, ":") {
			currentDir = strings.TrimSuffix(line, ":")
			continue
		}

		// Skip total line from ls -l
		if strings.HasPrefix(line, "total ") {
			continue
		}

		var entry ls.FileInfo

		if long {
			// Parse ls -l format
			entry, err := parseLongFormatLine(line, currentDir)
			if err != nil {
				// Skip lines that can't be parsed
				continue
			}
			entries = append(entries, entry)
		} else {
			// Simple format: just filenames
			name := line
			// Skip . and ..
			if name == "." || name == ".." {
				continue
			}

			// Get file info for size and directory check
			fullPath := filepath.Join(currentDir, name)
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}

			entry = ls.FileInfo{
				Name:        name,
				IsDirectory: info.IsDir(),
				Size:        info.Size(),
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func parseLongFormatLine(line string, currentDir string) (ls.FileInfo, error) {
	// ls -l output format:
	// -rw-r--r-- 1 user group 1234 Jan 02 15:04 filename
	// drwxr-xr-x 2 user group 4096 Jan 02 15:04 dirname

	fields := strings.Fields(line)
	if len(fields) < 9 {
		return ls.FileInfo{}, fmt.Errorf("invalid ls -l format")
	}

	permissions := fields[0]
	sizeStr := fields[4]
	month := fields[5]
	day := fields[6]
	timeOrYear := fields[7]
	name := strings.Join(fields[8:], " ")

	// Skip . and ..
	if name == "." || name == ".." {
		return ls.FileInfo{}, fmt.Errorf("skip special directories")
	}

	// Parse size
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 0
	}

	// Parse modification time
	modTime := parseModTime(month, day, timeOrYear)

	// Determine if directory
	isDir := strings.HasPrefix(permissions, "d")

	return ls.FileInfo{
		Name:        name,
		IsDirectory: isDir,
		Size:        size,
		Mode:        permissions,
		ModTime:     modTime,
	}, nil
}

func parseModTime(month string, day string, timeOrYear string) time.Time {
	now := time.Now()
	currentYear := now.Year()

	// Parse month
	monthMap := map[string]time.Month{
		"Jan": time.January, "Feb": time.February, "Mar": time.March,
		"Apr": time.April, "May": time.May, "Jun": time.June,
		"Jul": time.July, "Aug": time.August, "Sep": time.September,
		"Oct": time.October, "Nov": time.November, "Dec": time.December,
	}
	monthNum, ok := monthMap[month]
	if !ok {
		return time.Time{}
	}

	// Parse day
	dayNum, err := strconv.Atoi(day)
	if err != nil {
		return time.Time{}
	}

	// Determine if timeOrYear is time (HH:MM) or year
	var year int
	var hour, minute int

	if strings.Contains(timeOrYear, ":") {
		// It's a time, assume current year
		timeParts := strings.Split(timeOrYear, ":")
		if len(timeParts) == 2 {
			hour, _ = strconv.Atoi(timeParts[0])
			minute, _ = strconv.Atoi(timeParts[1])
		}
		year = currentYear
		// If the date is in the future, it's probably from last year
		testDate := time.Date(year, monthNum, dayNum, hour, minute, 0, 0, time.UTC)
		if testDate.After(now) {
			year--
		}
	} else {
		// It's a year
		year, _ = strconv.Atoi(timeOrYear)
		hour = 0
		minute = 0
	}

	return time.Date(year, monthNum, dayNum, hour, minute, 0, 0, time.UTC)
}

func fileInfoToAPIFileInfo(info os.FileInfo, includeLongFormat bool) ls.FileInfo {
	fileInfo := ls.FileInfo{
		Name:        info.Name(),
		IsDirectory: info.IsDir(),
		Size:        info.Size(),
	}

	if includeLongFormat {
		fileInfo.Mode = info.Mode().String()
		fileInfo.ModTime = info.ModTime()
	}

	return fileInfo
}
