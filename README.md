# GO-GIN

## 库

- gin：`github.com/gin-gonic/gin`
- gorm：`gorm.io/gorm`
- mysql：`gorm.io/driver/mysql`
- viper：`github.com/spf13/viper`
- air：`github.com/air-verse/air`

## 目录

```
go-gin
├─config
│  └─config.yaml
├─internal
│  ├─app
│  │  ├─config
│  │  │  └─config.go
│  │  ├─database
│  │  │  └─database.go
│  │  └─logger
│  │     └─logger.go
│  ├─controller
│  │  └─user.go
│  ├─middleware
│  │  ├─auth.go
│  │  └─logger.go
│  ├─model
│  │  └─user.go
│  ├─router
│  │  └─router.go
│  ├─service
│  │  └─user.go
│  └─utils
│     └─validator.go
├─log
├─tmp
└─main.go
```

### 迁移表

```
DB.AutoMigrate(&model.User{}, &model.Account{})
```
