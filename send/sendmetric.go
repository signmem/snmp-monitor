package send

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/signmem/snmpmonitor/g"
	"net/http"
	"errors"
)

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
