package monitor

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/pkg/errors"
	"github.com/signmem/snmp-monitor/g"
	"github.com/signmem/snmp-monitor/send"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ServerTotalList map[string]int
	GRAP  = int64(rand.Intn(100))
)

func GlobalStart() {

	g.SnmpServerDict = g.Config().SnmpServer
	ServerTotalList = make(map[string]int)

	if len(g.SnmpServerDict) == 0 {
		log.Fatalf("can not find any server in config file.")
		os.Exit(1)
	}

	var servers []string
	for _, snmpList := range g.SnmpServerDict {
		servers = append(servers, snmpList.IPAddr)
	}

	for _, server := range servers {
		ServerTotalList[server] = 0
	}

	for {
		if time.Now().Unix() % g.Config().Step == 0 {

			for addr, time := range ServerTotalList {
				var metricsValue int64

				if time == 0 {
					err := snmapProgram(addr)

					if err != nil {
						metricsValue = 0
						ServerTotalList[addr] = g.Config().SkipTime
					} else {
						metricsValue = 1
					}

				} else {
					ServerTotalList[addr] = time - 1
				}

				send.GenSnmpMetricAlive(addr, metricsValue)
			}
		}

		time.Sleep(time.Second * time.Duration(1))
	}

}



func snmapProgram(address string) (err error){

	// use to clean up TOTALMETRICS
	TOTALMETRICS = TOTALMETRICS[:0]

	// range oids in cfg file
	// TERRY
	oidsmap := g.Config().Oids

	for _, idsGroup := range oidsmap {
		err := snmpGet(address, "", idsGroup)
		if err != nil {
			return err
		}
	}

	// runing walk

	oidwalk := g.Config().OidWalks

	for _, oidWalkMap := range  oidwalk {

		OIDWALKTAG = OIDWALKTAG[:0]

		// 通过 walk 获取当前监控个数并创建新 tag
		walkMakeTag(address, oidWalkMap.TagOid)
		// get OIDWALKTAG

		// 读取 cfg check 配置 (后面用于遍历的 oid)
		oidWalkCheck := oidWalkMap.Check

		tagLeft :=  oidWalkMap.TagName

		if len(OIDWALKTAG) > 0 {
			for _, getWalkTags := range OIDWALKTAG {
				num := getWalkTags.WalkNum

				// 创建 tag
				tagRight := getWalkTags.WalkReturn
				tagFull := tagLeft + tagRight

				for _, walkCheck := range  oidWalkCheck {

					// 重组 oidmap 用于 snmp get 
					var walkOidMap  g.OIDMAP
					// regroup walkOidMap
					walkCheckOid := walkCheck.OID + "." + num

					// TERRY
					walkOidMap.OID =  walkCheckOid
					walkOidMap.Alias = walkCheck.Alias
					walkOidMap.Type = walkCheck.Type

					// 利用重组的 oid map 进行 snmpget 操作
					snmpGet(address, tagFull, walkOidMap)
				}

			}
		}

		g.Logger.Debugf("server: %s, metric total len: %d", address, len(TOTALMETRICS))

		if g.Config().Debug {

			for _, metric := range TOTALMETRICS {
				g.Logger.Debug(metric.String())
			}
		}

		if len(TOTALMETRICS) > 0 && g.Config().Upload {
			send.UploadMetric(TOTALMETRICS)
		}

		// 清空变量
		OIDWALKTAG = OIDWALKTAG[:0]
		TOTALMETRICS = TOTALMETRICS[:0]
	}
	return nil
}


func snmpGet(address string, tagName string, idsGroup g.OIDMAP) (err error) {

	gosnmp.Default.Target = address
	gosnmp.Default.Timeout = time.Duration(10 * time.Second)
	err = gosnmp.Default.Connect()

	if err != nil {
		g.Logger.Errorf("snmp connect error:%s", err)
		return err
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
		return err
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

		TOTALMETRICS = append(TOTALMETRICS, metrics)
	}

	return nil
}




func walkMakeTag(address string, oids string) (err error) {

	gosnmp.Default.Target = address
	gosnmp.Default.Timeout = time.Duration(10 * time.Second)
	err = gosnmp.Default.Connect()

	if err != nil {
		g.Logger.Errorf("snmp connect error:%s", err)
		return err
	}

	defer gosnmp.Default.Conn.Close()

	err = gosnmp.Default.BulkWalk(oids, walkValue)

	if err != nil {
		g.Logger.Errorf("snmp %s walk %s error:%s", address, oids, err)
		return err
	}

	return nil
}


func walkValue(pdu gosnmp.SnmpPDU) (err error) {

	// g.Logger.Infof("walk name: %s", pdu.Name)

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
		return errors.New(msg)
	}

	OIDWALKTAG = append(OIDWALKTAG, oidTag)

	return nil
}
