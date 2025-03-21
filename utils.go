package winimg

import (
	"path/filepath"
	"strings"
)

const (
	LongPathPrefix = `\\?\`
)

// DefaultCaptureExclude returns default capture exclusion list
// which can be used for adding custom exclusion.
func DefaultCaptureExclude() []string {
	return []string{
		"\\$ntfs.log",
		"\\hiberfil.sys",
		"\\pagefile.sys",
		"\\swapfile.sys",
		"\\System Volume Information",
		"\\$Recycle.Bin\\*",
		"\\Recycler",
		"\\Recycled",
		"\\Windows\\CSC",
		"\\winpepge.sys",
		"\\$windows.~ls",
		"\\$windows.~bt",
	}
}

// DefaultCompressionExclusion returns default compression exclusion
// list which can be used for adding custom exclusion.
func DefaultCompressionExclusion() []string {
	return []string{
		"*.mp3",
		"*.zip",
		"*.cab",
		"*.wmv",
		"*.wma",
		"*.wim",
		"*.swm",
		"*.dvr-ms",
		"\\Windows\\inf\\*.pnf",
	}
}

// ShouldExclude checks if fileName matches the pattern used for DISM exclusion.
//
// https://learn.microsoft.com/en-us/windows-hardware/manufacture/desktop/dism-configuration-list-and-wimscriptini-files-winnext#exclusion-list-guidelines
func ShouldExclude(fileName string, pattern string) bool {
	// make case-insensitive for comparsion
	fileName = strings.ToLower(fileName)
	pattern = strings.ToLower(pattern)

	sepStr := string(filepath.Separator)

	// root-relative patterns
	if strings.HasPrefix(pattern, sepStr) {
		if strings.ContainsRune(pattern, '*') {
			patternDir, patternBase := filepath.Split(pattern)
			patternDir = strings.TrimSuffix(patternDir, sepStr)

			// for root patterns with wildcards:
			// 1. The directory prefix must match exactly
			// 2. The file part must match the wildcard pattern

			fileDir, fileBase := filepath.Split(fileName)
			fileDir = strings.TrimSuffix(fileDir, sepStr)

			// check if directory matches and file matches wildcard
			if fileDir == patternDir {
				matched, _ := filepath.Match(patternBase, fileBase)
				return matched
			}
			return false
		}

		// non-wildcard root-relative patterns
		return fileName == pattern || strings.HasPrefix(fileName, pattern+sepStr)
	}

	// non-root patterns with wildcards
	if strings.ContainsRune(pattern, '*') {
		// split the path into components
		segments := strings.Split(fileName, sepStr)
		if len(segments) > 0 && segments[0] == "" {
			segments = segments[1:] // skip empty segment from leading slash
		}

		dirPart, filePart := filepath.Split(pattern)
		dirPart = strings.TrimSuffix(dirPart, sepStr)

		fileDir, fileBase := filepath.Split(fileName)
		fileDir = strings.TrimSuffix(fileDir, sepStr)

		// compare directories and match filename pattern
		if fileDir == dirPart {
			matched, _ := filepath.Match(filePart, fileBase)
			return matched
		}

		// check if any part of the path matches
		for i := 0; i < len(segments); i++ {
			checkPath := strings.Join(segments[i:], sepStr)

			checkDir, checkBase := filepath.Split(checkPath)
			checkDir = strings.TrimSuffix(checkDir, sepStr)

			if checkDir == dirPart {
				matched, _ := filepath.Match(filePart, checkBase)
				return matched
			}
		}
		return false
	}

	// for regular patterns without wildcards or root indicator,
	// match as substring anywhere in the path
	return strings.Contains(fileName, pattern)
}
