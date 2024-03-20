# snmp monitor  

> 用于 snmp 监控   
> go version 1.20 about

# cfg 说明  

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
