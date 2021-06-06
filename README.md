# probabDrill

### api_url如下

http://171.16.1.119:4399/swagger/index.html

更换模型数据，要相应修改配置文件 可直接使用run.sh重新编译和启动程序

更改ip要更改gin router和 注释

### 数据格式

basic_info.dat：钻孔名，钻孔坐标x，钻孔坐标y，钻孔坐标z

boundary_info.dat：坐标x，坐标y

layer_info.dat：钻孔名，地层名，地层深度（距离孔口）

std_layer_info.dat：地层代码，地层描述，地层序号

### 主要文件目录说明

assets：存放数据文件
项目目录下：config.yaml是配置文件
app中：存放api接口代码

### 钻孔数据结构的验证

钻孔数据是人工录入的，因此不可避免的会出现问题。一般而言，出现的问题包括且不限于如下问题：

1. 钻孔基本坐标异常（钻孔坐标数据录入错误）
2. 钻孔标准编号与地层层序编号不匹配（标准编号不能覆盖所有的钻孔分层信息）
3. 原始数据可能存在零厚度层
4. 钻孔层位信息中的地层高程不单调（是递减或递增）

### 安装go swag

https://github.com/swaggo/swag
go get -u github.com/swaggo/swag/cmd/swag


### 如果 gomod 无法使用，需要设置代理
[go mod 代理设置](https://www.cnblogs.com/tomtellyou/p/13053825.html)
```
# gomod常用命令
go mod init     # 初始化go.mod
go mod tidy     # 更新依赖文件
go mod download # 下载依赖文件
go mod vendor   # 将依赖转移至本地的vendor文件
go mod edit     # 手动修改依赖文件
go mod graph    # 打印依赖图
go mod verify   # 校验依赖
```

编译软件需要执行的命令
``` bash
# linux下使用
 go mod tidy
 go mod vendor
 swag init
 make build
 ./pd conf --path .

 # windows use
  go mod tidy                 # 配置依赖1
  go mod vendor               # 配置依赖2
  swag init                   # 重新生成api页面
  go build -o pd.exe main.go  # 编译服务 
  start pd.exe conf --path .  # 用户启动服务
```
