# snmp monitor  

> 用于 snmp 监控   

# cfg 说明  

> 简单说明   

```
    "debug":false,               // 打印 metric 信息到日志  
    "logfile": "/apps/conf/snmp-monitor/snmp-monitor.log",
    "logmaxage": 432000,
    "logrotateage": 86400,
    "step": 60,                 // 监控时间间隔  
    "upload": false,            // true 则 upload 至 falcon agent  
    "listen": "0.0.0.0:22233",  // 用于 http://0.0.0.0:22233/health_check 检测  
    "uploadurl": "http://127.0.0.1:22230/v1/push",  // falcon-agent 上传 metrics 接口  
    "snmpserver": [""],        // 需要执行监控的服务器地址

```


> 独立 get 方法  

```
    "oids": [            <- 使用 get 方式
        {
                "oid": "1.3.6.1.2.1.1.3.0",    // 需要监控的 oid (falcon metric values)  
                "alias": "sysUpTimeInstance",  // falcon metric name 
                "type": "GAUGE"                // 类型定义
        }
    ]
```

> 使用 walk 方法  

```
    "oidwalks": [
        {
        "tagoid" : "1.3.6.1.2.1.2.2.1.2",      <- walk 返回的 String (用于 falcon tag) ex 返回 eth0
        "tagname": "iface=",                   <- falcog tags  iface=eth0
        "check": [
                {
                    "oid": "1.3.6.1.2.1.2.2.1.11",      <- falcon  values
                    "alias": "net.if.in.multicast",     <- falcon  metrics name
                    "type": "COUNTER"                   <- falcon  type
                }
        }
    ]
```
