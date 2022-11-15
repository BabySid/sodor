# Sodor
Sodor is a distributed and extensible scheduler system, with lower operating expenses and high performance.

# Usage
```shell
# fat_ctrl
./fat_ctrl --metastore mysql://user:pass@tcp(1.2.3.4:3306)/charset=utf8 --listen_addr=:9527 --log.path=../logs

# thomas
./thomas --data.path=../data --listen_addr=:9528 --log.path=../logs
```