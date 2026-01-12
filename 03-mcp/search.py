import os
import zipfile
import requests
from pathlib import Path
from minsearch import Index

def download_fastmcp_repo():
    """Download the fastmcp repository zip file"""
    url = "https://github.com/jlowin/fastmcp/archive/refs/heads/main.zip"
    zip_path = "fastmcp-main.zip"
    
    if os.path.exists(zip_path):
        print(f"File {zip_path} already exists, skipping download")
        return zip_path
    
    print(f"Downloading {url}...")
    response = requests.get(url)
    response.raise_for_status()
    
    with open(zip_path, 'wb') as f:
        f.write(response.content)
    
    print(f"Downloaded {zip_path}")
    return zip_path

def extract_and_index():
    """Extract zip, read .md and .mdx files, and index with minsearch"""
    zip_path = download_fastmcp_repo()
    extract_dir = "fastmcp-main-extracted"
    
    # Extract zip file
    if os.path.exists(extract_dir):
        print(f"Directory {extract_dir} already exists, skipping extraction")
    else:
        print(f"Extracting {zip_path}...")
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(extract_dir)
        print(f"Extracted to {extract_dir}")
    
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
    
    print(f"Found {len(documents)} markdown files")
    
    # Create MinSearch index
    index = Index(text_fields=["content", "filename"])
    index.fit(documents)
    
    return index

def search_docs(index, query: str, limit: int = 5):
    """Search the indexed documents and return top results"""
    results = index.search(query, num_results=limit)
    return results

if __name__ == "__main__":
    # Create index
    print("Creating search index...")
    index = extract_and_index()
    
    # Test search with "demo"
    print("\nSearching for 'demo'...")
    results = search_docs(index, "demo", limit=5)
    
    print(f"\nFound {len(results)} results:")
    for i, result in enumerate(results, 1):
        print(f"{i}. {result.get('filename', 'N/A')}")
        print(f"   Score: {result.get('score', 'N/A')}")
        print(f"   Preview: {result.get('content', '')[:100]}...")
        print()
    
    if results:
        print(f"First file returned: {results[0].get('filename', 'N/A')}")
