# MCP Homework - Custom Documentation Search Server

This project implements a custom MCP (Model Context Protocol) server that provides web scraping and documentation search capabilities, similar to Context7.

## Features

- **Web Scraping**: Scrape web content using Jina Reader API
- **Documentation Search**: Search FastMCP documentation using minsearch
- **MCP Tools**: All functionality exposed as MCP tools for AI assistants

## Setup

### Prerequisites

- Python 3.10+
- `uv` package manager (install with `pip install uv`)

### Installation

1. Initialize the project (already done):
```bash
uv init
```

2. Install dependencies (already done):
```bash
uv add fastmcp requests minsearch
```

## Running the Server

To run the MCP server:

```bash
uv run python main.py
```

Or with explicit directory:

```bash
uv --directory /opt/test_feature/datatalks/ai-dev-tools-datatalks/03-mcp run python main.py
```

## Available Tools

### 1. `add(a: int, b: int) -> int`
Add two numbers together.

### 2. `scrape_web(url: str) -> str`
Scrape web content using Jina Reader API and return markdown.

**Example:**
```python
scrape_web("https://datatalks.club/")
```

**Note:** Jina Reader API may have rate limits or restrictions on certain URLs.

### 3. `search_docs(query: str, limit: int = 5) -> list`
Search the FastMCP documentation for relevant content.

**Example:**
```python
search_docs("demo", limit=5)
```

Returns a list of documents matching the query, each containing:
- `filename`: Path to the document
- `content`: Full text content of the document

## Project Structure

```
03-mcp/
├── main.py              # Main MCP server with all tools
├── test.py              # Test file for web scraping tool
├── search.py            # Test file for search functionality
├── pyproject.toml       # Project configuration
├── uv.lock              # Dependency lock file
├── HOMEWORK_ANSWERS.md  # Answers to homework questions
└── README.md            # This file
```

## Testing

### Test Web Scraping

```bash
uv run python test.py
```

### Test Search Functionality

```bash
uv run python search.py
```

This will:
1. Download the FastMCP repository
2. Extract and index all markdown files
3. Test search with query "demo"
4. Display the top 5 results

## Integration with AI Assistants

To use this MCP server with an AI assistant (like Cursor, Claude, or GitHub Copilot):

1. Configure your MCP client to point to this server
2. Use the command: `uv --directory <full-path-to-this-directory> run python main.py`
3. The AI assistant will be able to call the tools:
   - `scrape_web` for web content
   - `search_docs` for documentation search

## Dependencies

- `fastmcp`: MCP server framework
- `requests`: HTTP library for web requests
- `minsearch`: Full-text search library

## Notes

- The search index is built on first use and cached for performance
- The FastMCP repository is downloaded and extracted automatically
- Only `.md` and `.mdx` files are indexed
- Jina Reader API may have service restrictions (451 errors observed)

## Homework Submission

All homework answers are documented in `HOMEWORK_ANSWERS.md`.
