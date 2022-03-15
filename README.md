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
   -u, -url string         url to scan
   -wf, -word-file string  wordlist file

WORD OPTIONS:
   -we, -word-ext string     word extension
   -wp, -word-prefix string  word prefix
   -ws, -word-suffix string  word suffix

OUTPUT OPTIONS:
   -fc, -filter-code int   filter by status code
   -ec, -exclude-code int  exclude by status code

OTHER OPTIONS:
   -concurrent int  concurrent goroutines (default 10)
   -proxy string    proxy
   -timeout int     timeout (default 5)
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
- [x] custom word extension
- [x] custom word prefix, word suffix
- [x] custom URL replacement location
- [x] support proxy
- [x] concurrency settings
- [x] filter by status code
- [x] exclude by status code
- [x] custom timeout
- [x] QPS limit
- [x] basic auth
- [ ] custom headers
- [ ] progress save
