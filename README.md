# gommy

## Overview

The Command line tool is time series dummy data generator.

## Installation

#### Using go get,

```
$ go get github.com/yatakemi/gommy
```

## Usage

#### help

```
$ gommy -h
NAME:
   DummyGenerator - This app create to the dummy data files.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  a config file path (default: "./config.toml")
   --output value, -o value  a creating dummy data file path (default: "./dummyData.csv")
   --help, -h                show help
   --version, -v             print the version
```

## Future Perspectives

- update using by go lang way
- add testing
- add CI
- improve performance
- add newline characters declaration and character encoding declaration

## Contributing

Always welcome for contributing

## License & Authors

- Author:: @yatakemi
- License:: MIT
