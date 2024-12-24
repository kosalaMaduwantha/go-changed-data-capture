# Changed Data Capture (CDC) File Processor Documentation

## Overview
This project implements a file-based Changed Data Capture (CDC) system that monitors changes in a text file and publishes the changes to a RabbitMQ queue. Key features include:

- File monitoring with change detection using SHA-256 hashing
- File locking mechanism using [`flock`](go.mod) to handle concurrent access
- RabbitMQ integration for publishing changes
- JSON-based state persistence in hashStore.json
- Structured logging using logrus

## Module Structure
The project is structured into three main packages:

1. fileOps - Handles file operations
2. hashOps - Manages file hashing
3. jsonOps - Handles JSON serialization/deserialization
4. rabitMqAdapter - Manages RabbitMQ message publishing

## How It Works

### Change Detection
- The program monitors input.txt for changes every 5 seconds
- Changes are detected by comparing SHA-256 hashes of the file content
- Changed in the same line is tracked by the line number and length of the last processed line
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
     - Length of the last processed line

3. **Incremental Processing**
   - Only processes new lines and the current line if there are changes in the file.
   - Tracks line count to resume from the correct position

### Implementation Details

The main loop in captcherFileChange.go:

1. Checks for file changes every 5 seconds
2. If changed:
   - Reads new content from last processed line
   - Updates hash store with new state
   - Logs changes using structured JSON logging
   - Push the change to a queue(Rabit-mq)
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