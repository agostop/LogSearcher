# LogSearcher
The LogSearcher use for search log like ELK, but more lighter. 

Use logspout Collect logs.
logspout need use format of log: 
```
- "RAW_FORMAT={\"container\":\"{{ .Container.Name }}\",\"timestamp\":\"{{ .Time.Format \"2006-01-02 15:04:05\" }}\", \"message\":{{ toJSON .Data }}}\n"
```
