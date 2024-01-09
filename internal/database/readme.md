# gorm gen 自动生成

> 参考文档: [https://gorm.io/zh_CN/gen/](https://gorm.io/zh_CN/gen/)

使用

```bash
# 在项目根目录执行
go run . gorm gen
```

## gorm 生成说明

### 目录结构

```
├── gorm_gen                 ---
│   ├── base.go              ---
│   └── [示例model].go        --- xxx.go 自行创建的生成规则文件，可选
├── model                    ---
│   ├── [model_01].gen.go    --- gorm_gen 自动生成的文件，不要修改
│   ├── [model_0x].gen.go    --- gorm_gen 自动生成的文件，不要修改
│   └── [custom].go          --- 自定义model文件
└── query
    ├── [model_01].gen.go    --- gorm_gen 自动生成的文件，不要修改
    ├── [model_0x].gen.go    --- gorm_gen 自动生成的文件，不要修改
    ├── [custom].go          --- 自定义query文件
    └── gen.go
```

### 步骤

- 根据需要确定是否需要自定义 `gorm_gen/xxx.go` 文件 (可以参考现有的 `gorm_gen/demo.go`)
- 根据需要确定是否需要修改 `grom_gen.go` 中的 db 连接配置
- 在项目根目录执行 `go run . gorm gen` 开始生成相应文件
