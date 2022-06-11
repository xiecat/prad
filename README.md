# prad

Web directory and file discovery.

## Usage

```shell
➜ prad  .\prad.exe -h
╱╱╱╱╱╱╱╱╱╱╭╮
╱╱╱╱╱╱╱╱╱╱┃┃
╭━━┳━┳━━┳━╯┃
┃╭╮┃╭┫╭╮┃╭╮┃
┃╰╯┃┃┃╭╮┃╰╯┃
┃╭━┻╯╰╯╰┻━━╯
┃┃
╰╯ v0.0.1

web directory and file discovery.

Usage:
  C:\prad.exe [flags]

Flags:
INPUT OPTIONS:
   -u, -url string           url to scan
   -wf, -word-file string    wordlist file
   -wl, -word-list string[]  wordlist

OUTPUT OPTIONS:
   -nc, -no-color  disable color in output

OTHER OPTIONS:
   -concurrent int  concurrent goroutines (default 10)
   -proxy string    proxy
   -timeout int     timeout (default 5)
   -qps int         QPS (default 10)
```

```shell
➜ prad  .\prad.exe -u http://127.0.0.1:8000
╱╱╱╱╱╱╱╱╱╱╭╮
╱╱╱╱╱╱╱╱╱╱┃┃
╭━━┳━┳━━┳━╯┃
┃╭╮┃╭┫╭╮┃╭╮┃
┃╰╯┃┃┃╭╮┃╰╯┃
┃╭━┻╯╰╯╰┻━━╯
┃┃
╰╯ v0.0.1

404 - http://127.0.0.1:8000/.svn
404 - http://127.0.0.1:8000/admin
404 - http://127.0.0.1:8000/login
404 - http://127.0.0.1:8000/.git
404 - http://127.0.0.1:8000/backup
404 - http://127.0.0.1:8000/manager
200 - http://127.0.0.1:8000/.idea
```

```shell
➜ prad  .\prad.exe -u 'http://127.0.0.1:8000/{{path}}/admin'
╱╱╱╱╱╱╱╱╱╱╭╮
╱╱╱╱╱╱╱╱╱╱┃┃
╭━━┳━┳━━┳━╯┃
┃╭╮┃╭┫╭╮┃╭╮┃
┃╰╯┃┃┃╭╮┃╰╯┃
┃╭━┻╯╰╯╰┻━━╯
┃┃
╰╯ v0.0.1

404 - http://127.0.0.1:8000/backup/admin
404 - http://127.0.0.1:8000/login/admin
404 - http://127.0.0.1:8000/admin/admin
404 - http://127.0.0.1:8000/manager/admin
404 - http://127.0.0.1:8000/.svn/admin
404 - http://127.0.0.1:8000/.idea/admin
404 - http://127.0.0.1:8000/.git/admin
```

## Features

- [x] custom wordlist file
- [x] custom URL replacement location
- [x] support proxy
- [x] concurrency settings
- [x] custom timeout
- [x] QPS limit
- [ ] custom word extension
- [ ] custom word prefix, word suffix
- [ ] filter by status code
- [ ] exclude by status code
- [ ] progress save
- [ ] WAF detection
