# Hyper

Hyper is a lightweight and easy-to-use framework. The hyper framework works out of the box. It has built-in commonly used components, you only need to turn it on according to your needs.

## Feature

- Based on [gin](https://github.com/gin-gonic/gin) framework, Lightweight and easy to use
- Config
    - Based on [spf13/viper](https://github.com/spf13/viper) package
    - Support local files (yaml format) and consul
- Logger
    - Base on [sirupsen/logrus](https://github.com/sirupsen/logrus) package
    - Support multiple output sources (stdout/stderr/file)
    - Support text or json format log output
- Queue
    - Support redis and kafka two types of queue storage backend
    - Ability to recover from abnormal consumption tasks
    - Support consumer task failure retry function
    - Built-in prometheus exporter, you can review your queue
- Mysql
    - Base on [go-gorm/gorm](https://github.com/go-gorm/gorm) package
	- The maximum number of connections, enable sql log and other settings only need to be configured
- Discover
    - Support Consul-based service registration and discovery
- Selector
    - Support getting nodes from Discover for load balancing
    - Support RoundRobin strategy
- Cache
    - Support redis as a storage backend
- Built-in performance analysis tool pprof

## Install and Use

### Install

```shell
go get -u github.com/cyub/hyper/cmd/hyper
```

### Create Application

```shell
cd www/
hyper new your_project_name
```

### Run Application

```
cd your_project_name
make run
```

### Access Application

```
curl localhost:8000/welcome
```