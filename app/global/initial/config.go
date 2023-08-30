package initial

const ConfigToml = `
# 数据库配置 以下 两者任选其一 不要同时使用
# sqlite
[db]
DBType = "sqlite"
DBName = "pear-admin"

# mysql
# [db]
# DBType = "mysql"
# DBName = "pear-admin"
# DBUser = "root"
# DBPwd = "123456"
# DBHost = "127.0.0.1:3306"

# redis配置
[redis]
RedisAddr = "127.0.0.1:6379"
RedisPWD = ""
RedisDB = 0

# 日志配置
[zaplog]
director = 'runtime/log'

# 其他杂项配置
[app]
HttpPort = 8009
PageSize = 20
RunMode = "debug"
JwtSecret = "0102$%#&*^*&150405"
ImgSavePath = "static/upload"
ImgUrlPath = "runtime/upload/images"
`
