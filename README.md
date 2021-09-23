# poker for poker planning

As I was annoyed that most poker planning arent't free for a simple functionnality, here is a solution. Free poker planning.

![badge](https://github.com/worming004/poker/workflows/builddeploy/badge.svg)

## installation

with go version 1.17 installed, simply do

```bash
go install github.com/worming004/poker@v1.0.0
poker
```

show help with `poker -h`

## docker run

```bash
git clone https://www.github.com/worming004/poker
cd poker
docker build -t poker .
docker run --rm -p 8000:8000 poker
```

It doesn't support for now password management.

## password

in case you set a password, you'll have to reach the website with query parameter `<endpoint>?password=<your password>`
