import time
import pandas as pd
from datetime import datetime

import parser
import web_data

html_parser = parser.Parser()
fish_base = web_data.FishBase()

df_trait = pd.read_excel("data/2018-north-species.xlsx")
species = df_trait.iloc[:, 0]
print(species)
tmp_dict = {}

for s in species:
    s_prased = s.replace(" ", "-")
    print("\n", s_prased)
    try:
        html_content = fish_base.search_spec(s_prased)
        res = html_parser.parse_html(html_content)
    except Exception as e:
        print(e)
    tmp_dict[s] = res
    time.sleep(1.5)

df_res = pd.DataFrame.from_dict(tmp_dict, 
                                orient='index',
                                columns=['max_len', 'unit', 'info', 'ref', 'bayesian_a', 'a_min', 'a_max', 'bayesian_b', 'b_min', 'b_max', 'info'])
df_res.to_csv(f"data/result_{datetime.now().strftime('%Y-%m-%d')}.csv")