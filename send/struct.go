package send

import "fmt"

type MetaData struct {
	Metric      string            `json:"metric"`
	Endpoint    string            `json:"endpoint"`
	Timestamp   int64             `json:"timestamp"`
	Step        int64             `json:"step"`
	Value       float64           `json:"value"`
	CounterType string            `json:"counterType"`
	Tags        map[string]string `json:"tags"`
}

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

func (this *MetricValue) String() string {
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, Time:%d, Value:%v>",
		this.Endpoint,
		this.Metric,
		this.Type,
		this.Tags,
		this.Step,
		this.Timestamp,
		this.Value,
	)
}