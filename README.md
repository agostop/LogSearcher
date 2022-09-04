# LogSearcher
The LogSearcher use for search log like ELK, but more lighter. 

Use logspout Collect logs.
logspout need use blow format of log: 
```
version: "2"
services:
  logspout:
    image: gliderlabs/logspout:latest
    hostname: logspout
    container_name: logspout
    command: 'raw+tcp://your_remote_logSearcher'
    volumes:
      - /usr/share/zoneinfo/Asia/Shanghai:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - DEBUG=1
      - LOGSPOUT=ignore
      - "RAW_FORMAT={\"container\":\"{{ .Container.Name }}\",\"timestamp\":\"{{ .Time.Format \"2006-01-02 15:04:05\" }}\", \"message\":{{ toJSON .Data }}}\n"
    networks:
      - deployment_default

networks:
  default_network:
    external: true
```

