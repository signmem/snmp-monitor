package g

type GlobalConfig struct {
	Debug			bool		`json:"debug"`
	LogFile			string		`json:"logfile"`
	LogMaxAge		int			`json:"logmaxage"`
	LogRotateAge	int			`json:"logrotateage"`
	Oids	 		[]OIDMAP	`json:"oids"`
	OidWalks		[]OidWalk	`json:"oidwalks"`
	SnmpServer 		[]SnmpServers 	`json:"snmpserver"`
	Step 			int64 		`json:"step"`
	SkipTime		int			`json:"skiptime"`
	UploadUrl		string		`json:"uploadurl"`
	Upload			bool		`json:"upload"`
	Listen			string		`json:"listen"`
}

type OIDMAP struct {
	OID 	string 			`json:"oid"`
	Alias 	string 			`json:"alias"`
	Type	string			`json:"type"`
}

type OidWalk struct {
	TagName 			string 		`json:"tagname"`
	TagOid				string		`json:"tagoid"`
	Check 				[]OIDMAP	`json:"check"`
}

type SnmpServers struct {
	IPAddr		string		`json:"ipaddr"`
	HostName	string		`json:"hostname"`
}