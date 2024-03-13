import requests

class FishBase:
    def search_spec(self, spec):
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