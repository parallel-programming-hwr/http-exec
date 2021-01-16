# http-exec

Start remote command via http request.

## How it works

This is only a wrapper for some command, that can be started via an http request.
It can only run __one__ concurrent job.

## Usage

### Install

Go get it with:
```bash
go get github.com/parallel-programming-hwr/http-exec
```

### Run

Start http server with:
```bash
./http-exec --command <command> --args="<arg1,arg2,...>"
```

This will start a server and listen to port 8080.

## API

### /start

Start executing command and write stdout back.

### /

Write status of worker. Either `running` or `ready`.
