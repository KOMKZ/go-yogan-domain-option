# go-yogan-domain-option

Yogan Framework 的 Option（配置）领域包，提供系统配置的 CRUD 功能。

## 安装

```bash
go get github.com/KOMKZ/go-yogan-domain-option@latest
```

## 功能

- 系统配置管理（增删改查）
- 按 Key 查询配置
- 按分组类型查询配置
- 支持 JSON 类型的组件参数

## 使用示例

```go
import (
    option "github.com/KOMKZ/go-yogan-domain-option"
    "github.com/KOMKZ/go-yogan-framework/logger"
    "gorm.io/gorm"
)

// 初始化
func InitOptionService(db *gorm.DB) *option.Service {
    repo := option.NewGORMRepository(db)
    return option.NewService(repo, logger.GetLogger("option"))
}

// 使用
func Example(svc *option.Service) {
    // 获取所有配置
    options, _ := svc.GetAll(ctx)
    
    // 根据 Key 获取
    opt, _ := svc.GetByKey(ctx, "site_name")
    
    // 创建配置
    svc.Create(ctx, &option.CreateInput{
        Key:       "site_name",
        Value:     "My Site",
        GroupType: "system",
    })
}
```

## 数据模型

```go
type Option struct {
    ID              uint      
    Key             string    // 配置键（唯一）
    Value           string    // 配置值
    GroupType       string    // 分组类型
    Component       string    // 前端组件类型
    ComponentParams JSONValue // 组件参数
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

## 依赖

- [go-yogan-framework](https://github.com/KOMKZ/go-yogan-framework) - 核心框架
- [gorm](https://gorm.io) - ORM

## License

MIT
