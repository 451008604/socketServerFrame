# socketServerFrame
```
│  go.mod
│  go.sum
│  main.go                              //入口文件
│
├─api                                   //API注册与逻辑实现
├─client                                //客户端服务（示例代码，可直接删除）
│  │  main.go
│  │
│  ├─api
│  │      base.go
│  │      ping.go
│  │
│  ├─base
│  │      connection.go
│  │
│  └─conf
├─config                                //项目配置文件
│      conf.go
│      config.json
│
├─iface                                 //抽象接口定义
│      iconnection.go                   //连接对象（一个客户端对应一个连接）
│      iconnmanager.go                  //连接管理器对象
│      idatapack.go                     //处理消息封包/拆包
│      imessage.go                      //消息对象（包含消息ID、消息长度、消息内容）
│      imsghandler.go                   //消息处理对象（接收消息处理中间层）
│      irequest.go                      //请求对象（包含请求对应的连接、请求数据）
│      irouter.go                       //业务处理对象（可以处理每条协议逻辑前置、后置）
│      iserver.go                       //服务器对象（负责管理对外服务的生命周期）
│
├─logic                                 //逻辑层（负责处理公共逻辑）
├─logFiles                              //日志文件存放目录
├─logs                                  //打印日志管理器（异步打印）
│      printlog2console.go              //输出到控制台
│      printlog2file.go                 //输出到文件
│      printlogManager.go               //对外提供打印接口（可以控制打印模式）
│
├─proto                                 //proto文件管理
│  │  generate_pb.sh                    //生成 pb 协议和 grpc 协议文件
│  │
│  ├─bin                                //protoc编译后的输出目录
│  │      message.pb.go
│  │      message_grpc.pb.go
│  │
│  └─src                                //proto源文件目录
│          message.proto
│
└─znet                                  //抽象接口实现
        connection.go
        connmanager.go
        datapack.go
        message.go
        msghandler.go
        request.go
        router.go
        server.go
```

## gRPC配置

- 安装protoc编译器

> https://github.com/protocolbuffers/protobuf/releases/  
> 下载后解压到任意目录把`bin`里面的`protoc.exe`复制到`%GOPATH%/bin`里面，并配置`PATH`环境变量，确保 protoc 可以正常执行

- 安装相关模块

> go install google.golang.org/protobuf/proto  
> go install google.golang.org/protobuf/cmd/protoc-gen-go@latest  
> go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  