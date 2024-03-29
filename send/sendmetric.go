package send

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/signmem/snmp-monitor/g"
	"net/http"
	"time"
)

func GenSnmpMetricAlive(snmpServerDict g.SnmpServers, value float64) (err error) {

	var metrics MetricValue
	var sendMetrics []MetricValue

	metrics.Type = "GAUGE"
	metrics.Step = g.Config().Step
	metrics.Timestamp = time.Now().Unix()
	metrics.Metric = "snmpd.alive"
	metrics.Endpoint = snmpServerDict.HostName
	metrics.Value = value

	sendMetrics = append(sendMetrics, metrics)

	metrics.Metric = "agent.alive"
	sendMetrics = append(sendMetrics, metrics)

	err = UploadMetric(sendMetrics)

	if err != nil {
		g.Logger.Errorf("upload err:%s", err)
	}
	return nil
}

func UploadMetric(metrics []MetricValue) (err error) {

	url := g.Config().UploadUrl

	jdata, err := json.Marshal(metrics)

	if err != nil {
		g.Logger.Errorf("UploadMetric() err:%s", err)
		return err
	}

	header := "application/json"
	resp, err := http.Post(url, header, bytes.NewBuffer(jdata))

	if err != nil {
		g.Logger.Errorf("UploadMetric() err:%s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("UploadMetric() http resp.status err %d", resp.StatusCode)
		g.Logger.Errorf(msg)
		return errors.New(msg)
	}

	return nil
}
