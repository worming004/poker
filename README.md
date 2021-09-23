# poker for poker planning

As I was annoyed that most poker planning arent't free for a simple functionnality, here is a solution. Free to use poker planning.

![badge](https://github.com/worming004/poker/workflows/builddeploy/badge.svg)

## installation

with go version 1.17 installed, simply do

```bash
git clone https://www.github.com/worming004/poker
cd poker
go install
poker --password <insert your password>
```

## docker run

```bash
git clone https://www.github.com/worming004/poker
cd poker
docker build -t poker .
docker run --rm -p 8000:8000 poker
```
