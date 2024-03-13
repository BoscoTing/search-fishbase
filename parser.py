import re
from bs4 import BeautifulSoup

class Parser: 
    def parse_max_len(self, text):
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

    def parse_bayesian_len_wei(self, text):
        pattern = r"Bayesian length-weight: a=([\d.]+) \(([\d.]+) - ([\d.]+)\), b=([\d.]+) \(([\d.]+) - ([\d.]+)\), (in \w+ total length)"
        matches = re.findall(pattern, text)
        if matches:
            print("Match found.", "\n")
            return list(matches[0])
        else:
            print("No match found.", "\n")
            return [None] * 7
        
    def parse_html(self, html_content):
        soup = BeautifulSoup(html_content, "html.parser")
        main_content = soup.find('div', id="ss-main")
        columns = main_content.find_all('h1', class_="slabel bottomBorder")

        for col in columns:
            if "Size / Weight / Age" in col.text:
                max_length = col.find_next('div', class_="smallSpace").text
                res_max_length = self.parse_max_len(max_length)
            elif "Estimates based on models" in col.text:
                bay_len_wei = col.find_next('div', class_="smallSpace").text
                res_bay_len_wei = self.parse_bayesian_len_wei(bay_len_wei)
        tot_res = res_max_length + res_bay_len_wei

        return tot_res