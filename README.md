# go-port-scanner
golang网络端口扫描器，支持tcp/udp协议

# 使用教程

命令：go run main.go [options] address  
options：  
-p    使用指定协议，如：tcp、udp，默认为tcp    
-c    使用指定cpu数量，默认为1  
-t     使用指定线程数，默认为1（使用多线程扫描，速度会更快）  
-d    扫描单个端口超时控制，单位：秒，默认为5秒
