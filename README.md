Core
====

## Install libgit2
We use v0.23.4 because v0.24 sucks

```shell
$ wget https://github.com/libgit2/libgit2/archive/v0.23.4.zip
$ unzip v0.23.4.zip
```
Install dependencies: `libssh2`, `http-parser`, `cmake`, `libcurl`.

```shell
On Mac OSX:
$ brew install libssh2 http-parser cmake libcurl

On CentOS:
$ yum install libssh2-devel http-parser cmake libcurl-devel
```
Then install libgit2.

```shell
$ cmake .  -DCMAKE_VERBOSE_MAKEFILE=ON -Wno-dev -DUSE_SSH=YES
$ make && make install
```
Note on Mac OSX may need to set CFLAGS="-std=c99".
Now libgit2 is installed under `/usr/local/lib` as default. We still need to set pkg-config and link dynamic libraries.

```shell
On Mac OSX:
$ cd /usr/local/lib/pkgconfig
$ ln -s /path/to/libgit2/pkgconfig/libgit2.pc  libgit2.pc
$ cd /usr/local/lib
$ ln -s /path/to/libgit2/libgit2.23.dylib libgit2.23.dylib

On CentOS:
$ cd /usr/lib64/pkgconfig/
$ ln -s /usr/local/lib/pkgconfig/libgit2.pc libgit2.pc
$ cd /usr/lib64
$ ln -s /usr/local/lib/libgit2.so.23 libgit2.so.23
```

## setup dev environment

```shell
$ git config --global url."git@gitlab.ricebook.net:".insteadOf "https://gitlab.ricebook.net/"
$ go get gitlab.ricebook.net/platform/core.git
$ mv $GOPATH/src/gitlab.ricebook.net/platform/core.git $GOPATH/src/gitlab.ricebook.net/platform/core
$ cd $GOPATH/src/gitlab.ricebook.net/platform/core && go install
$ ln -s $GOPATH/src/gitlab.ricebook.net/platform/core $MY_WORK_SPACE/eru-core2
$ make deps
```

### GRPC

Generate golang & python code

```shell
$ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
$ make golang
$ make python
```

Current version of dependencies are:

* google.golang.org/grpc: v1.0.1-GA
* github.com/golang/protobuf: f592bd283e

do not forget first command...

### deploy core on local environment

* create `core.yaml` like this

```yaml
bind: ":5000" # gRPC server 监听地址
agent_port: "12345" # agent 的 HTTP API 端口, 暂时没有用到
permdir: "/mnt/mfs/permdirs" # 宿主机的 permdir 的路径
etcd: # etcd 集群的地址
    - "http://127.0.0.1:2379"
etcd_lock_prefix: "/eru-core/_lock" # etcd 分布式锁的前缀, 一般会用隐藏文件夹

git:
    public_key: "[path_to_pub_key]" # git clone 使用的公钥
    private_key: "[path_to_pri_key]" # git clone 使用的私钥
    gitlab_token: "[token]" # gitlab API token

docker:
    log_driver: "json-file" # 日志驱动, 线上会需要用 none
    network_mode: "bridge" # 默认网络模式, 用 bridge
    cert_path: "[cert_file_dir]" # docker tls 证书目录
    hub: "hub.ricebook.net" # docker hub 地址
    hub_prefix: "namespace/test" # 存放镜像的命名空间, 两边的/会被去掉, 中间的会保留. 镜像名字会是$hub/$hub_prefix/appname:version

scheduler:
    lock_key: "_scheduler_lock" # scheduler 用的锁的 key, 会在 etcd_lock_prefix 里面
    lock_ttl: 10 # scheduler 超时时间
    type: "complex" # scheduler 类型, complex 或者 simple
```

* start eru core

```shell
$ core --config /path/to/core.yaml --log-level debug
```

or

```shell
$ export ERU_CONFIG_PATH=/path/to/core.yaml
$ export ERU_LOG_LEVEL=DEBUG
$ core
```


### Use client.py

```
$ devtools/client.py --grpc-host core-grpc.intra.ricebook.net node:get intra c2-docker-7
```
