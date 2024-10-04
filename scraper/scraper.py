import requests
from bs4 import BeautifulSoup
from json import dumps

baseUrl = "https://www.libgen.is"
possibleMirrors = ["Libgen & IPFS & Tor", "Libgen.li"]

def scrape(subject):
    searchUrl = f"{baseUrl}/search.php?req={subject}"
    res = requests.get(searchUrl)
    parser = BeautifulSoup(res.content, "html")
    urls = [x.get("href") for x in parser.find_all("a") if x.get("title") in possibleMirrors and x.getText() == "[1]"]


    s = requests.Session()
    links = []
    for url in urls:
        res = s.get(url)
        parser = BeautifulSoup(res.content, "html")
        links += [link.get("href") for link in parser.find_all("a") if link.getText() == "GET"]

    for link in links:
        requests.post("http://localhost:8001/add", data = dumps({
            "url": link,
            "subject": subject
        }))
