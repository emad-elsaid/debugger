# Feature: Open Files Panel

## Overview
Show the currently opened files by the process being debugged in a scrollable view.

## User Capabilities

### 1. View File list
**User Action** Clicking on the tab
**Expected Behavior**
- A list of open files by the debugged process
- It should show the following information:
  - File name
  - File path
  - File size
  - File type (socket, file, directory, ...etc)
  - Open mode (read, write, append...etc)
  - A button to open the directory of the file

### 2. Updates when the program pauses
When the program pauses at a breakpoint the file list should be updated to
reflect currently open

### 3. Files should be sortable
A drop down at the top right should allow the user to sort by any of the attributes shown (name, path, size, type, mode...)


### 4. Files should be filterable
A text field at the top of the panel should allow the user to filter the file by
the path/name. the search should be fuzzy search not exact value.
