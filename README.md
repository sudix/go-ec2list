# go-ec2list
Retrieve EC2 instances information from all profiles and all regions concurrently.
This tool is intended to use with interactive selection tools like [peco](https://github.com/peco/peco) or [percol](https://github.com/mooz/percol).


INSTALLATION
=============

## Download binary file

Dwonload the binary file for your platform from [release page](https://github.com/sudix/go-ec2list/releases) and move it where you want like below.

```
$ mv go-ec2list_linux_amd64 /usr/local/bin/go-ec2list
```

## By go get

```
$ go get github.com/sudix/go-ec2list
```

USAGE
=============

## Setup your AWS credential file.

See [Configuring the AWS Command Line Interface - AWS Command Line Interface](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

## Run command

```
$ go-ec2list
```

## Result

You will get EC2 instances list.  
Each line is tab separeted values.

__Output Values__

- Instance Name
- Instance ID
- Public IP
- Private IP
- Profile Name
- Availability zone
- Instance type
- Instance state

## Options

### Cache

If you set `-cachemin` option, go-ec2list caches results and uses that until it expires.
Value means expire minutes.

__USAGE__

```
$ go-ec2list -cachemin 10
```

The cache is stored to `$HOME/.go-ec2list/cache`.

### Delete cache

When `-refresh` is specified, the cache will always be deleted and re-create cache regardless of `-cachemin` option.

__USAGE__

```
$ go-ec2list -cachemin 10 -refresh
```

EXAMPLE USAGE WITH INTERACTIVE SELECTION TOOLS.
=============

Set alias like below, and you can ssh login to the selected instance.

```sh
alias ec2="go-ec2list | peco | cut -f3 | xargs -I{} sh -c 'ssh "ec2-user@{}" </dev/tty' ssh"
```
