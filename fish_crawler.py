import requests
import time
import re
import pandas as pd
from datetime import datetime
from bs4 import BeautifulSoup

df_trait = pd.read_excel("2018-north-species.xlsx")
species = df_trait.iloc[:, 0]
print(species)
tmp_dict = {}

def search_spec(spec):
    """exp_url = 'https://www.fishbase.se/summary/Acanthurus-dussumieri.html'
    """
    url = f'https://www.fishbase.se/summary/{spec}.html'
    res = requests.get(url)
    if res.status_code == 200:
        try: 
            html_content = res.content
        except Exception as e:
            print(e)
    return html_content

def parse_max_len(text):
    pattern = r"(\d+\.\d+) (.+?) (.*?);([^;]+)"
    match = re.search(pattern, text)
    if match:
        max_length = float(match.group(1))
        unit = match.group(2).strip()
        category = match.group(3).strip()
        additional_info = match.group(4).strip()
        print("Match found:", match.group())
        return [max_length, unit, category, additional_info]
    else:
        print("No match found.", "\n")
        return [None] * 4

def parse_bayesian_len_wei(text):
    pattern = r"Bayesian length-weight: a=([\d.]+) \(([\d.]+) - ([\d.]+)\), b=([\d.]+) \(([\d.]+) - ([\d.]+)\), (in \w+ total length)"
    matches = re.findall(pattern, text)
    if matches:
        print("Match found.", "\n")
        return list(matches[0])
    else:
        print("No match found.", "\n")
        return [None] * 7
    
def parse_html(html_content):
    soup = BeautifulSoup(html_content, "html.parser")
    main_content = soup.find('div', id="ss-main")
    columns = main_content.find_all('h1', class_="slabel bottomBorder")

    for col in columns:
        if "Size / Weight / Age" in col.text:
            max_length = col.find_next('div', class_="smallSpace").text
            res_max_length = parse_max_len(max_length)
        elif "Estimates based on models" in col.text:
            bay_len_wei = col.find_next('div', class_="smallSpace").text
            res_bay_len_wei = parse_bayesian_len_wei(bay_len_wei)
    tot_res = res_max_length + res_bay_len_wei

    return tot_res

for s in species:
    s_prased = s.replace(" ", "-")
    print("\n", s_prased)
    try:
        html_content = search_spec(s_prased)
        res = parse_html(html_content)
    except Exception as e:
        print(e)
    tmp_dict[s] = res
    time.sleep(1.5)

df_res = pd.DataFrame.from_dict(tmp_dict, 
                                orient='index',
                                columns=['max_len', 'unit', 'info', 'ref', 'bayesian_a', 'a_min', 'a_max', 'bayesian_b', 'b_min', 'b_max', 'info'])
df_res.to_csv(f"data/result_{datetime.now().strftime('%Y-%m-%d')}.csv")