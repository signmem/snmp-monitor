package monitor

import (
	"github.com/signmem/snmpmonitor/send"
)

type OidStruct struct {
	WalkNum 		string 			`json:"walknum"`
	WalkReturn		string 			`json:"walkreturn"`
}

var OIDWALKTAG []OidStruct


var TOTALMETRICS []send.MetricValue