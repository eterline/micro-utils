
# microutils

Project mini with utils. Written in ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

Use for:
![Proxmox](https://img.shields.io/badge/proxmox-proxmox?style=for-the-badge&logo=proxmox&logoColor=%23E57000&labelColor=%232b2a33&color=%232b2a33)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)

## Build

To build project in full prod variant
```
make build-prod
```

for testing build
```
make build
```

All binariaes will be in './build' folder at repository


## Usage per util:

Mini usage CLI apps with argument input functional.

### seeip:
Tool for searching about IP of domain name (or list of them).

#### supported dns provide:
- local DNS         - `local`
- Google DoH        - `google`
- Cloudflare DoH    - `cloudflare`
- remote DNS server - `10.192.0.1:53` (DNS server addr for example. Port must be exists)

```
Usage: seeip [--addr ADDR] [--reslov RESLOV] [--json] [--format] [--workers WORKERS]

Options:
  --addr ADDR, -a        Search ip address or domain. Can be list or single value. [default: []]
  --reslov RESLOV, -r    Resolver service name | DNS server address. [default: local]
  --json, -j             JSON object output.
  --format, -f           JSON formatted object output.
  --workers WORKERS, -w  Process worker count.
  --help, -h             display this help and exit
```

Exampled output:
```
user@host~#: seeip -a google.com -r google

google.com:
    resolve_duration_ms: 504
    resumes:
        - request_ip: 142.250.74.78
          resume:
            status: success
            continent: Europe
            continentCode: EU
            country: Sweden
            countryCode: SE
            region: AB
            regionName: Stockholm County
            city: Stockholm
            districresumet: ""
            zip: 100 04
            lat: 59.3293
            lon: 18.0686
            timezone: Europe/Stockholm
            offset: 7200
            currency: SEK
            isp: Google LLC
            org: Google LLC
            as: AS15169 Google LLC
            asname: GOOGLE
            reverse: arn09s23-in-f14.1e100.net
            mobile: false
            proxy: false
            hosting: true
        - request_ip: 2a00:1450:400f:802::200e
          resume:
            status: success
            continent: Europe
            continentCode: EU
            country: Sweden
            countryCode: SE
            region: AB
            regionName: Stockholm County
            city: Stockholm
            districresumet: ""
            zip: 100 04
            lat: 59.3327
            lon: 18.0656
            timezone: Europe/Stockholm
            offset: 7200
            currency: SEK
            isp: Google LLC
            org: GOOGLE 2a
            as: AS15169 Google LLC
            asname: GOOGLE
            reverse: arn09s23-in-x0e.1e100.net
            mobile: false
            proxy: false
            hosting: true
    ns:
        - ns1.google.com.
        - ns3.google.com.
        - ns2.google.com.
        - ns4.google.com.
```

### filehash:
Later...

### gpufo:
Later...

### stubit:
Later...

### uuid:
Later...

## License

[MIT](https://choosealicense.com/licenses/mit/)