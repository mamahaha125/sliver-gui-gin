module pear-admin-go

go 1.16

require (
	github.com/bishopfox/sliver v1.15.16
	github.com/gin-gonic/gin v1.9.1
	google.golang.org/grpc v1.56.2
)

replace (
	github.com/bishopfox/sliver => ./sliver-master
	github.com/rsteube/carapace v0.39.0 => github.com/reeflective/carapace v0.25.2-0.20230602202234-e8d757e458ca
)

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/bytedance/sonic v1.9.2 // indirect
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f
	github.com/gchaincl/dotsql v1.0.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 v10.14.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/gorm v1.9.16
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/mojocn/base64Captcha v1.3.4
	github.com/moloch--/asciicast v0.1.1 // indirect
	github.com/mssola/user_agent v0.5.2
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/onsi/ginkgo v1.16.1 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	github.com/pkg/sftp v1.13.4
	github.com/rsteube/carapace v0.39.0 // indirect
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.7.1
	github.com/xujiajun/nutsdb v0.6.0
	go.uber.org/zap v1.18.1
	golang.org/x/arch v0.4.0 // indirect
	golang.org/x/crypto v0.11.0
	golang.org/x/exp v0.0.0-20230711153332-06a737ee72cb // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/term v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/protobuf v1.31.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
