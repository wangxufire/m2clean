# m2clean

### Install
```
go get github.com/wangxufire/m2clean@latest
```

### OR Download [Click here](https://github.com/wangxufire/m2clean/releases)

### Usage
```
$ m2clean -h    
Usage:
  m2clean [OPTIONS]

Application Options:
  -p, --path=             Path to m2 directory, if using a custom path, default is homedir/.m2/repository
  -a, --accessed-before=  Delete all libraries (even if latest version) last accessed on or before this date (2006-01-02).Default 3 month ago.
  -f, --ignore-artifacts= artifactIds (full or part) to be ignored.
  -g, --ignore-groups=    groupIds (full or part) to be ignored.
  -d, --dryrun            Do not delete files, just simulate and print result.

Help Options:
  -h, --help              Show this help message
```