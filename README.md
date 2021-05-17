# probabDrill

钻孔数据是人工录入的，因此不可避免的会出现问题。一般而言，出现的问题包括且不限于如下问题：

1. 钻孔基本坐标异常（钻孔坐标数据录入错误）
2. 钻孔标准编号与地层层序编号不匹配（标准编号不能覆盖所有的钻孔分层信息）
3. 原始数据可能存在零厚度层
4. 钻孔层位信息中的地层高程不单调（是递减或递增）

```
go mod init  # 初始化go.mod
go mod tidy  # 更新依赖文件
go mod download  # 下载依赖文件
go mod vendor  # 将依赖转移至本地的vendor文件
go mod edit  # 手动修改依赖文件
go mod graph  # 打印依赖图
go mod verify  # 校验依赖
```
