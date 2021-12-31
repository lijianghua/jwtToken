# jwt API Example

## 主要技术点：
+ 实现JWT token 无状态机制
+ 有状态的 JWT token 保存到 redis,登出之后无法再使用refreshtoken和accesstoken
+ mariadb golang 基本操作
+ refresh token机制
+ bearer token 机制
+ 统一配置文件yaml
+ 实现拦截器
+ 模拟signup/signin/signout/refresh/welcome handler处理
## 流程
![avatar](https://github.com/lijianghua/jwtToken/blob/master/tokenFlow.png)
## 第三方库使用
+ github.com/dgrijalva/jwt-go v3.2.0+incompatible
+ github.com/go-redis/redis v6.15.9+incompatible
+ github.com/go-sql-driver/mysql v1.6.0
+ github.com/twinj/uuid v1.0.0
+ golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
+ gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
## bug
未按照oauth标准响应格式返回（包括成功和错误）
