package monitor

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/pkg/errors"
	"github.com/signmem/snmp-monitor/g"
	"github.com/signmem/snmp-monitor/send"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)


func GlobalStart() {

	g.SnmpServerDict = g.Config().SnmpServer
	// 没有配置服务器则退出
	if len(g.SnmpServerDict) == 0 {
		log.Fatalf("can not find any server in config file.")
		os.Exit(1)
	}

	var servers []string
	for _, snmpList := range g.SnmpServerDict {
		servers = append(servers, snmpList.IPAddr)
	}

	// 支持最大 10 并发
	task_chan := make(chan bool, 10)
	wg := sync.WaitGroup{}
	defer close(task_chan)

	// do snmp check program

	for {
		if time.Now().Unix() % g.Config().Step == 0 {

			for _, addr := range servers {

				wg.Add(1)
				task_chan <- true
				go func(addr string) {
					<-task_chan
					// 主入口
					runSnmapCheck(addr)
					defer wg.Done()
				}(addr)
			}
			wg.Wait()
		}

		time.Sleep(time.Second * time.Duration(1))
	}

}

func runSnmapCheck(addr string) {

	var metricsValue float64
	err := snmapProgram(addr)

	if err != nil {
		return
	}

	metricsValue = 1
	// 监控成功，则上报 falcon agent.alive = 1
	send.GenSnmpMetricAlive(addr, metricsValue)
	return
}

func snmapProgram(address string) (err error){

	// totalSnmpMetric 私有变量
	var totalSnmpMetric []send.MetricValue

	// range oids in cfg file (snmpget 用法)
	oidsmap := g.Config().Oids

	for _, idsGroup := range oidsmap {

		// fix get metrics from global variable to private variable
		snmpGetMetrics, err := snmpGet(address, "", idsGroup)

		if err != nil {
			continue
		}

		totalSnmpMetric = append(totalSnmpMetric, snmpGetMetrics...)
	}

	// snmpwalk 用法

	oidwalk := g.Config().OidWalks

	for _, oidWalkMap := range oidwalk {

		// 通过 walk 获取当前监控个数并创建新 tag
		//var oidWalkTag []OidStruct
		oidWalkTag, err := walkMakeTag(address, oidWalkMap.TagOid)

		if err != nil {
			g.Logger.Errorf("spip check %s, oid %s", address, oidWalkMap.TagOid)
			continue
		}

		// 读取 cfg check 配置 (后面用于遍历的 oid)
		oidWalkCheck := oidWalkMap.Check

		tagLeft :=  oidWalkMap.TagName

		if len(oidWalkTag) > 0 {
			for _, getWalkTags := range oidWalkTag {
				num := getWalkTags.WalkNum

				// 创建 tag
				tagRight := getWalkTags.WalkReturn
				tagFull := tagLeft + tagRight

				for _, walkCheck := range  oidWalkCheck {

					// 重组 oidmap 用于 snmp get 
					var walkOidMap  g.OIDMAP
					// regroup walkOidMap
					walkCheckOid := walkCheck.OID + "." + num

					walkOidMap.OID =  walkCheckOid
					walkOidMap.Alias = walkCheck.Alias
					walkOidMap.Type = walkCheck.Type

					// 利用重组的 oid map 进行 snmpget 操作
					snmpMetrics, _ := snmpGet(address, tagFull, walkOidMap)
					totalSnmpMetric = append(totalSnmpMetric, snmpMetrics...)
				}

			}
		}

		if len(totalSnmpMetric) == 0 {
			errmsg := errors.New("server:%s, metric total len is 0")
			return errmsg
		}

		g.Logger.Debugf("server: %s, metric total len: %d", address, len(totalSnmpMetric))

		if g.Config().Debug {
			for _, metric := range totalSnmpMetric {
				g.Logger.Debug(metric.String())
			}
		}

		if len(totalSnmpMetric) > 0 && g.Config().Upload {
			send.UploadMetric(totalSnmpMetric)
		}

	}
	return nil
}


func snmpGet(address string, tagName string, idsGroup g.OIDMAP) (snmpMetrics []send.MetricValue, err error) {

	gosnmp.Default.Target = address
	gosnmp.Default.Timeout = time.Duration(3 * time.Second)
	err = gosnmp.Default.Connect()

	if err != nil {
		g.Logger.Errorf("snmp connect error:%s", err)
		return snmpMetrics, err
	}

	defer gosnmp.Default.Conn.Close()

	var metrics send.MetricValue

	metrics.Endpoint = g.GetHostname(address)

	if len(tagName) != 0 {
		metrics.Tags = tagName
	}

	metrics.Timestamp = time.Now().Unix()
	metrics.Step = g.Config().Step

	var idsmap []string
	idsmap = append(idsmap, idsGroup.OID)

	result, err := gosnmp.Default.Get(idsmap)

	if err != nil {
		g.Logger.Errorf("snmp Get() ids:%s, alias:%s, err:%s",
			idsGroup.OID, idsGroup.Alias, err)
		return snmpMetrics, err
	}

	for _, variables := range result.Variables {

		// g.Logger.Infof("oid: %s, nmae: %s", variables.Name, idsDetail.Name)

		metrics.Metric = idsGroup.Alias
		metrics.Type = idsGroup.Type

		switch variables.Type {

		case gosnmp.OctetString:
			// g.Logger.Infof("string: %v", string(variables.Value.([]byte)))
			valueString := string(variables.Value.([]byte))
			metrics.Value, _ = strconv.ParseFloat(valueString, 64)
		default:
			// g.Logger.Infof("num: %d", gosnmp.ToBigInt(variables.Value))
			value := gosnmp.ToBigInt(variables.Value).Int64()
			metrics.Value =  float64(value)
		}

		snmpMetrics = append(snmpMetrics, metrics)
	}

	return snmpMetrics,nil
}




func walkMakeTag(address string, oids string) (oidwalktag []OidStruct, err error) {

	gosnmp.Default.Target = address
	gosnmp.Default.Timeout = time.Duration(3 * time.Second)
	err = gosnmp.Default.Connect()

	if err != nil {
		g.Logger.Errorf("snmp connect error:%s", err)
		return oidwalktag, err
	}

	defer gosnmp.Default.Conn.Close()

	// err = gosnmp.Default.BulkWalk(oids, walkValue)
	allPduReport, err := gosnmp.Default.BulkWalkAll(oids)

	if err != nil {
		g.Logger.Errorf("snmp BulkWalkAll get error:%s", err)
		return oidwalktag, err
	}

	oidwalktag, err = walkValueToTag(allPduReport)

	if err != nil {
		g.Logger.Errorf("snmp %s walk %s error:%s", address, oids, err)
		return oidwalktag, err
	}

	return oidwalktag, nil
}

func walkValueToTag(report []gosnmp.SnmpPDU) (totalOidTag []OidStruct, err error) {

	// g.Logger.Infof("walk name: %s", pdu.Name)
	for _, pdu := range report {

		var oidTag OidStruct

		lastOID := pdu.Name[strings.LastIndex(pdu.Name, ".")+1:]
		oidTag.WalkNum = lastOID

		switch pdu.Type {
		case gosnmp.OctetString:
			b := pdu.Value.([]byte)
			oidTag.WalkReturn = string(b)
		}

		if len(oidTag.WalkNum) == 0 || len(oidTag.WalkReturn) == 0 {
			msg := fmt.Sprintf("snmpWalk get Tag error %s", pdu.Name)
			g.Logger.Errorf(msg)
			continue
		}

		totalOidTag = append(totalOidTag, oidTag)
	}

	return totalOidTag,nil
}
