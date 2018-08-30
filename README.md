# gommy

## Overview

The Command line tool is time series dummy data generator.

## Installation

#### Using go get,

```bash
$ go get github.com/yatakemi/gommy
```

## Usage

```bash
$ gommy -h
NAME:
   gommy - This app create to the dummy data files.

USAGE:
   gommy [global options] command [command options] [arguments...]

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

## e.g.

```bash
$ gommy -c config.toml -o dummyData.csv
```

### [config.toml](config.toml)

```toml
[datetime]
datetimeFormat = "2006-01-02 15:04:05"
start = 2018-11-11T00:00:00Z
end = 2018-12-11T23:00:00Z
column = 1
  [datetime.sampling]
  num = 1
  unit = "min"


[data]
min = 0.0
max = 0.001
pointtype = "float"

  [[data.datetime]]
  column = 7
  datetimeFormat = "15:04:05"
  add = 2000 # ms

  [[data.tag]]
  column = 4
  value = ["x", "y"]

  [[data.tag]]
  column = 5
  value = ["a", "b", "c"]

  [[data.tag]]
  column = 6
  value = ["XX", "YY"]

  [[data.abnormal]]
  column = 2
  min = 100.0
  max = 200.0
  pointtype = "int"
  start = 2018-11-12T00:00:00Z
  end = 2018-11-12T23:00:00Z

    [data.abnormal.transition]
    num = 1
    unit = "hour"

  [[data.abnormal]]
  column = 3
  min = 20.0
  max = 30.0
  pointtype = "float"
  start = 2018-11-11T00:00:00Z
  end = 2018-11-11T10:00:00Z

    [data.abnormal.transition]
    num = 5
    unit = "min"

[[header]]
row = ["datetime", "aaaa", "bbbb", "cccc", "dddd", "eeee", "ffff", "gggg"]

[[header]]
row = ["timestamp", "number", "number", "number", "string", "string", "string", "string"]
```

### dummyData.csv

```csv
datetime,aaaa,bbbb,cccc,dddd,eeee,ffff,gggg
timestamp,number,number,number,string,string,string,string
2018-11-11 00:00:00,0.0007355972757972727,5.231728471813334,x,a,XX,00:00:02,0.0006532558526877338
2018-11-11 00:00:00,1.3177514878274794e-05,5.271597817616059,x,a,YY,00:00:02,0.0004023891559321828
2018-11-11 00:00:00,0.0006436030332137436,5.499780626736413,x,b,XX,00:00:02,0.0008558502004125847
2018-11-11 00:00:00,0.0006224964414509832,5.679696630024756,x,b,YY,00:00:02,0.0005695700323944159
2018-11-11 00:00:00,0.0004698830439772164,5.207167146360693,x,c,XX,00:00:02,0.00018046295275153692
:
:
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
