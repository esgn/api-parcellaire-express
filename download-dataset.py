#!/usr/bin/env python3 
#-*- coding: utf-8 -*-

import requests
from bs4 import BeautifulSoup
import os
import shutil
from multiprocessing.pool import ThreadPool
import re

def extract_all_links(site):
    html = requests.get(site).text
    soup = BeautifulSoup(html, 'html.parser').find_all('a')
    links = [link.get('href') for link in soup]
    return links

def download_url(url):
    file_name_start_pos = url.rfind("/") + 1
    filename = url[file_name_start_pos:]
    out_file = os.path.join(out_dir, filename)
    r = requests.get(url, stream=True)
    if r.status_code == requests.codes.ok:
        with open(out_file, 'wb') as f:
            for data in r:
                f.write(data)
    else:
        return filename + " téléchargement en échec"
    return filename + " téléchargement réussi"


# Téléchargement des archives

url = "https://files.opendatarchives.fr/professionnels.ign.fr/parcellaire-express/PCI-par-DEPT_2021-02/"
out_dir = "parcellaire-express"

regex = re.compile(r'.*\.7z$')

all_links = extract_all_links(url)
all_links = [i for i in all_links if regex.match(i)]
all_links = [url+x for x in all_links]

if os.path.exists(out_dir):
    shutil.rmtree(out_dir)
os.mkdir(out_dir)

results = ThreadPool(5).imap_unordered(download_url, all_links)
for r in results:
    print(r)
