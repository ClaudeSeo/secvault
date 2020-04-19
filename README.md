# secvault
simple cli tool to easily manage sensitive environment variable using AWS Secrets Manager 

## Prerequisites

- AWS Credentials

## Installation

```bash
$ go get github.com/claudeseo/secvault
```

## Usage

Retrieving the list 

```bash
$ secvault list
```

Get Environment variable 

```bash
$ secvault get --secret-name {SECRET_NAME} --type {dotenv | kubernetes | json}
```


Put Environment variable

```bash
$ secvault put --secret-name {SECRET_NAME} --file {FILE_PATH} --type {dotenv | json}
```
