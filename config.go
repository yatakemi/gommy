package main

import "time"

// Config for toml
type Config struct {
	Datetime DatetimeConfig
	Data     DataConfig
	Header   []HeaderConfig
}

// DatetimeConfig for toml
type DatetimeConfig struct {
	DatetimeFormat string           `toml:"datetimeFormat"`
	Start          time.Time        `toml:"start"`
	End            time.Time        `toml:"end"`
	Column         int              `toml:"column"`
	Sampling       SamplingDatetime `toml:"sampling"`
}

// SamplingDatetime for toml
type SamplingDatetime struct {
	Num  int    `toml:"num"`
	Unit string `toml:"unit"`
}

// DataConfig for toml
type DataConfig struct {
	Min       float64        `toml:"min"`
	Max       float64        `toml:"max"`
	Pointtype string         `toml:"pointtype"`
	Abnormal  []AbnormalData `toml:"abnormal"`
	Tag       []TagData      `toml:"tag"`
	Datetime  []DatetimeData `toml:"datetime"`
}

// AbnormalData for toml
type AbnormalData struct {
	Min        float64                `toml:"min"`
	Max        float64                `toml:"max"`
	Pointtype  string                 `toml:"pointtype"`
	Column     int                    `toml:"column"`
	Start      time.Time              `toml:"start"`
	End        time.Time              `toml:"end"`
	Transition TransitionAbnormalData `toml:"transition"`
}

// TransitionAbnormalData for toml
type TransitionAbnormalData struct {
	Num  int    `toml:"num"`
	Unit string `toml:"unit"`
}

// TagData for toml
type TagData struct {
	Column int      `toml:"column"`
	Rate   int      `toml:"rate"` // TODO only rate=100 now
	Value  []string `toml:"value"`
}

// DatetimeData for toml
type DatetimeData struct {
	Column         int    `toml:"column"`
	DatetimeFormat string `toml:"datetimeFormat"`
	Add            int    `toml:"add"` // [ms]
}

// HeaderConfig for toml
type HeaderConfig struct {
	Row []string `toml:"row"`
}
