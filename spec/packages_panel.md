# Feature: Packages Panel

## Overview
Simple list view of all Go packages included in the debugged binary,
showing package import paths in alphabetical order.

## User Capabilities

### 1. View All Packages
**User Action:** Click "Packages" tab in bottom panel
**Expected Behavior:**
- List displays all packages compiled into binary
- Packages shown as import paths
- Sorted alphabetically
- Read-only list (no interaction beyond viewing)
- Auto-updates on state changes and restarts

**Package Examples:**
- `main` - Main package
- `fmt` - Standard library
- `net/http` - Standard library nested package
- `github.com/user/project/handlers` - Project package
- `github.com/user/project/models` - Another project package
- `github.com/external/library` - Third-party dependency

### 2. Searchable
A text field at the top should allow the user to filter the list by fuzzy text
