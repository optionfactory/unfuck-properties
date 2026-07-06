package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func process(config Config) error {
	if config.InPlace {
		return processInPlace(config.InPath)
	}
	return processToNewFile(config.InPath, config.OutPath)
}

func processInPlace(filePath string) error {
	resolvedPath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	filePath = resolvedPath

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	inFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}

	dir := filepath.Dir(filePath)
	tmpFile, err := os.CreateTemp(dir, "unfuck-*.tmp")
	if err != nil {
		inFile.Close()
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()

	defer os.Remove(tmpName)

	processErr := processStream(bufio.NewReader(inFile), bufio.NewWriter(tmpFile))

	if chmodErr := tmpFile.Chmod(info.Mode()); chmodErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not set permissions on temp file: %v\n", chmodErr)
	}

	inFile.Close()

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to safely flush and close temp file: %w", err)
	}

	if processErr != nil {
		return processErr
	}

	if err := os.Rename(tmpName, filePath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}
	return nil
}

func processToNewFile(inPath, outPath string) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	processErr := processStream(bufio.NewReader(inFile), bufio.NewWriter(outFile))

	if closeErr := outFile.Close(); closeErr != nil {
		return fmt.Errorf("failed to safely close output file: %w", closeErr)
	}

	return processErr
}

func processStream(reader *bufio.Reader, writer *bufio.Writer) error {
	for {
		line, err := reader.ReadBytes('\n')

		if len(line) > 0 {
			fixedLine := make([]byte, 0, len(line))

			for i := 0; i < len(line); {
				if i+5 < len(line) && line[i] == '\\' && line[i+1] == 'u' && line[i+2] == '0' && line[i+3] == '0' {
					h1 := hex2byte(line[i+4])
					h2 := hex2byte(line[i+5])

					if h1 != 255 && h2 != 255 {
						fixedLine = append(fixedLine, (h1<<4)|h2)
						i += 6
						continue
					}
				}

				fixedLine = append(fixedLine, line[i])
				i++
			}

			if _, wErr := writer.Write(fixedLine); wErr != nil {
				return fmt.Errorf("write error: %w", wErr)
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read error: %w", err)
		}
	}

	return writer.Flush()
}

func hex2byte(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 255
}
