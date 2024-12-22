# Changed Data Capture (CDC) File Processor Documentation

## Overview
This project implements a file-based Changed Data Capture (CDC) system that monitors changes in a text file and captures new content incrementally. The program continuously watches an input file and records changes in a JSON-based state store.

## Module Structure

The project is structured into three main packages:

1. fileOps - Handles file operations
2. hashOps - Manages file hashing
3. jsonOps - Handles JSON serialization/deserialization

## How It Works

### Change Detection
- The program monitors input.txt for changes every 5 seconds
- Changes are detected by comparing SHA-256 hashes of the file content
- The current state is maintained in hashStore.json
  
### Process Flow
1. **Hash Comparison**
   ```go
   func isFileChanged(inputFilePath string, configFilePath string) bool
   ```
   Compares current file hash with stored hash to detect changes

2. **State Management**
   - Stores file metadata in JSON format:
     - File name
     - File hash
     - Last processed line number

3. **Incremental Processing**
   - Only processes new lines added since last read
   - Tracks line count to resume from the correct position

### Implementation Details

The main loop in captcherFileChange.go:

1. Checks for file changes every 5 seconds
2. If changed:
   - Reads new content from last processed line
   - Updates hash store with new state
   - Logs changes using structured JSON logging
3. If unchanged:
   - Logs status and continues monitoring

## Key Features
- Incremental processing of changes
- Persistent state management
- SHA-256 based change detection
- JSON formatted logging
- File operation error handling

## Dependencies
- `github.com/sirupsen/logrus` for structured logging
- Standard Go libraries for file operations and hashing

#### Note: Current implementation captures changes in a single file only. For multiple files, the program can be extended to support a list of files to monitor.

## Upcoming Development

- Support for multiple files.
- Capture changes and store in a destination (e.g. database).
- Capture changes from remote file servers.