# go-port-scanner
golang 网络端口扫描器，支持 tcp/udp 协议

# 使用教程
```
 go run main.go [options] address  
```

## 例子
```
 go run main.go -c 4 127.0.0.1 //使用 4 核心扫描 127.0.0.1 挂起的端口
```

## options：  
- -p   
使用指定协议，如：tcp、udp，默认为 tcp    
- -c    
使用指定 cpu 数量，默认为 1  
- -t    
使用指定线程数，默认为 1（使用多线程扫描，速度会更快）  
- -d    
扫描单个端口超时控制，单位：秒，默认为 5 秒
