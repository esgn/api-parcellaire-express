#!/usr/bin/env python3 
#-*- coding: utf-8 -*-

import os
import sys
import re
import shutil
import requests
from bs4 import BeautifulSoup
from multiprocessing.pool import ThreadPool

def extract_all_links(site):
    html = requests.get(site).text
    soup = BeautifulSoup(html, 'html.parser').find_all('a')
    links = [link.get('href') for link in soup]
    return links

def download_url(url):
    file_name_start_pos = url.rfind("/") + 1
    filename = url[file_name_start_pos:]
    out_file = os.path.join(out_dir, filename)
    try:
        with requests.get(url, stream=True) as r:
            r.raise_for_status()
            with open(out_file, 'wb') as f:
                for data in r.iter_content(chunk_size=8192):
                    f.write(data)
    except requests.exceptions.HTTPError as err:
        print(filename + " teléchargement échoué")
        print(err)
        sys.exit(0)
    return filename + " téléchargement réussi"

if __name__ == "__main__":

    idx_url = "https://files.opendatarchives.fr/professionnels.ign.fr/parcellaire-express/PCI-par-DEPT_2021-04/"
    out_dir = "/tmp/parcellaire-express"
    regex = re.compile(r'.*\.7z$')
    max_parallel_dl = 5
    testing=False

    # Utilisation d'une autre URL de téléchargement que l'URL par défaut
    if "DOWNLOAD_URL" in os.environ:
        if os.environ['DOWNLOAD_URL']!="":
            idx_url = os.environ['DOWNLOAD_URL']

    # Utilisation d'un autre maximum de téléchargements en parallèle que celui par défaut
    if "MAX_PARALLEL_DL" in os.environ:
        if int(os.environ['MAX_PARALLEL_DL'])!=0:
            max_parallel_dl = int(os.environ['MAX_PARALLEL_DL'])
        else:
            print("🚧 new value of MAX_PARALLEL_DL is not in integer. Ignored.", file=sys.stderr)

    # Extraction des liens de téléchargement
    all_links = extract_all_links(idx_url)
    all_links = [i for i in all_links if regex.match(i)]
    all_links = [idx_url+x for x in all_links]

    # Pour un simple test on se limite à une seule archive
    if "TEST_IMPORTER" in os.environ:
        if int(os.environ['TEST_IMPORTER'])!=0:
            testing=True
            all_links = all_links[:1]
        else:
            print("🚧 new value of TEST_IMPORTER is not in integer. Ignored.", file=sys.stderr)

    # Réinitialisation du dossier de téléchargement
    if os.path.exists(out_dir):
        shutil.rmtree(out_dir)
    os.mkdir(out_dir)

    # Téléchargement des archives en parallèle
    print("Début du téléchargement du Parcellaire Express")
    print("URL source : " + idx_url)
    print("Nombre de téléchargements en parallèle : " + str(max_parallel_dl))
    if testing:
        print("🟣 EXECUTION EN MODE TEST => UNE SEULE ARCHIVE SERA TELECHARGEE")
    results = ThreadPool(max_parallel_dl).imap_unordered(download_url, all_links)
    for r in results:
        print(r)
