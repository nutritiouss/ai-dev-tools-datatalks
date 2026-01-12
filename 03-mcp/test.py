import requests

def scrape_web(url: str) -> str:
    """Scrape web content using Jina Reader API and return markdown"""
    jina_url = f"https://r.jina.ai/{url}"
    headers = {
        "Accept": "text/markdown",
        "X-Return-Format": "markdown"
    }
    response = requests.get(jina_url, headers=headers)
    response.raise_for_status()
    return response.text

if __name__ == "__main__":
    # Test with the specified URL
    url = "https://github.com/alexeygrigorev/minsearch"
    try:
        content = scrape_web(url)
        char_count = len(content)
        print(f"Character count: {char_count}")
        print(f"First 200 characters: {content[:200]}")
    except Exception as e:
        print(f"Error: {e}")
        print(f"Response status: {e.response.status_code if hasattr(e, 'response') else 'N/A'}")
