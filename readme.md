项目结构概览：

``` go
├── config	配置文件
├── domain	业务逻辑
├── dto	传输数据结构
│   ├── req	请求数据
│   └── resp	响应数据
├── gerror	错误信息
├── model	逻辑数据结构模型定义
├── persist	持久化
│   └── mysql	mysql持久化
└── router	api路由endpoint
```

* First of all,you should copy config.toml.example as config.toml,then give the correct configuration.



* build:
``` bash
make build
```

* run:
``` bash
./build/bin/LightServer
```