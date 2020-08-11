# httpeek

With `httpeek` you can quickly do HTTP lookups during your recon phase. 
You can pipe your list with URLs into `httpeek` and it will output the status code
and HTML `title`. 

## Usage

```

Usage of ./httpeek:
  -f, --file string    Input file to use instead of stdin
  -q, --query string   XPath query to lookup in HTML output (default "//title")
  -s, --silent         Only output actual results

```

## Examples

I like to use [gobuster](https://github.com/OJ/gobuster) and [httprobe](https://github.com/tomnomnom/httprobe) and combine it with `httpeek` during my recon.

```
# First enumerate all DNS subdomains
$ gobuster dns -d github.com -w <some wordlist> > subdomains.txt

# Use httprobe to find out which are exposing HTTP(S) services
$ cat subdomains.txt | httprobe > https.txt

# Check what the status code and title is with httpeek
$ cat https.txt | httpeek > httpeek.json
```

Or just as one liner combined with [jq](https://github.com/stedolan/jq) to get the title of all 200 responses:

```
$ cat subdomains.txt | httprobe | httpeek | jq '. | select(.status_code == 200) | .result'
```
