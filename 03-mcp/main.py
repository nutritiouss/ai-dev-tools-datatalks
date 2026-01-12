from fastmcp import FastMCP
import requests
import os
import zipfile
from pathlib import Path
from minsearch import Index

mcp = FastMCP("Demo 🚀")

# Global index cache
_search_index = None

@mcp.tool
def add(a: int, b: int) -> int:
    """Add two numbers"""
    return a + b

@mcp.tool
def scrape_web(url: str) -> str:
    """Scrape web content using Jina Reader API and return markdown.
    
    Args:
        url: The URL to scrape (e.g., https://datatalks.club/)
    
    Returns:
        The markdown content of the webpage
    """
    jina_url = f"https://r.jina.ai/{url}"
    headers = {
        "Accept": "text/markdown",
        "X-Return-Format": "markdown"
    }
    response = requests.get(jina_url, headers=headers)
    response.raise_for_status()
    return response.text

def _get_or_create_index():
    """Get or create the search index for fastmcp documentation"""
    global _search_index
    if _search_index is not None:
        return _search_index
    
    zip_path = "fastmcp-main.zip"
    extract_dir = "fastmcp-main-extracted"
    
    # Download if needed
    if not os.path.exists(zip_path):
        url = "https://github.com/jlowin/fastmcp/archive/refs/heads/main.zip"
        response = requests.get(url)
        response.raise_for_status()
        with open(zip_path, 'wb') as f:
            f.write(response.content)
    
    # Extract if needed
    if not os.path.exists(extract_dir):
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(extract_dir)
    
    # Collect all .md and .mdx files
    documents = []
    base_path = Path(extract_dir) / "fastmcp-main"
    
    for file_path in base_path.rglob("*.md"):
        relative_path = file_path.relative_to(base_path)
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        documents.append({
            "filename": str(relative_path),
            "content": content
        })
    
    for file_path in base_path.rglob("*.mdx"):
        relative_path = file_path.relative_to(base_path)
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        documents.append({
            "filename": str(relative_path),
            "content": content
        })
    
    # Create and fit index
    _search_index = Index(text_fields=["content", "filename"])
    _search_index.fit(documents)
    
    return _search_index

@mcp.tool
def search_docs(query: str, limit: int = 5) -> list:
    """Search the FastMCP documentation for relevant content.
    
    Args:
        query: The search query string
        limit: Maximum number of results to return (default: 5)
    
    Returns:
        List of documents matching the query, each containing 'filename' and 'content'
    """
    index = _get_or_create_index()
    results = index.search(query, num_results=limit)
    return results

if __name__ == "__main__":
    mcp.run()
