
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
user@host~# seeip -a google.com -r google

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
Tool for file hash calc. (Multi-thread working)
#### Hash types
- SHA256
- SHA1
- MD5

```
Usage: filehash [--file FILE] [--json] [--pretty]

Options:
  --file FILE, -f FILE   Target file
  --json, -j             JSON output format
  --pretty, -p           JSON output pretty style
  --help, -h             display this help and exit
```

Exampled output:
```
user@host~# filehash -f .bashrc

Target file: .bashrc
File size: 185.0000 B
===========================
 SHA256 = cd107f934b472a5d7efaca5af25cca5b64c4da59b49b66ed34400e8982f077ff
 SHA1   = d767be47298f160a41cfb5585a4c22f66b5e8776
 MD5    = d2f9f9f9bea5746979fba64a43797f2e
```

### gpufo:
Later...

### stubit:
Later...

### uuid:
Later...

```
Usage: uuid [--count COUNT] [--payload PAYLOAD] [--domain DOMAIN] [--version VERSION]

Options:
  --count COUNT, -c     UUID v4 count to generate [default: 1]
  --payload PAYLOAD, -p UUID hashed payload
  --domain DOMAIN, -d   UUID domain namespace [default: dns]
  --version VERSION, -v UUID version 3 or 5 [default: 3]
  --help, -h            display this help and exit
```

Exampled output:
```
# generates 2 UUID
user@host~#  uuid -c 2
- 457cbded-2561-4b17-a96f-d3fc98b46408
- 336b31e8-2b7a-4f8c-b8fd-194388c5f100

# generates 1 UUID from file
user@host~# uuid -p .bashrc
- 05045efe-5934-3ef3-b4f7-a7457c5908c
```
  !Don't support pipeline now

## License

[MIT](https://choosealicense.com/licenses/mit/)