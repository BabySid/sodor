# Sodor
Sodor is a distributed and extensible scheduler system, with lower operating expenses and high performance.

# Usage
```shell
# fat_ctrl
./fat_ctrl --metastore.addr="mysql://user:pass@tcp(1.2.3.4:3306)/dbName?charset=utf8mb4&parseTime=True&loc=Local" --listen_addr=:9527 --log.path=../logs

# thomas
./thomas --data.path=../data --listen_addr=:9528 --log.path=../logs
```

# Todo
* garbage collection task instances of thomas
* multi-task-type. this should be scalable for parsing the content. e.g. http param, shell script
  * parsed_content in task instance
* multi-alert-plugin. this should be scalable
  * unified alarm logging
* system-variables. e.g. latest_instance