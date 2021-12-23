#jwtToken Example

##主要技术点：
+ 实现JWT token 无状态机制
+ 有状态的 JWT token 保存到 redis
+ mariadb golang 基本操作
+ refresh token机制
+ bearer 验证
+ 统一配置文件yaml
+ 实现拦截器
+ 模拟signup/signin/signout/refresh/welcome handler处理
## 第三方库使用
+ github.com/dgrijalva/jwt-go
+ github.com/go-redis/redis
+ github.com/go-sql-driver/mysql
+ github.com/google/uuid 暂时未用
+ gopkg.in/yaml.v3 支持解析yaml/json等
+ "golang.org/x/crypto/bcrypt" bcrypt hash算法，保存密码hash较为合适