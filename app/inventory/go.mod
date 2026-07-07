module github.com/MoScenix/mes/app/inventory

go 1.25.5

replace (
	github.com/MoScenix/mes/common => ../../common
	github.com/MoScenix/mes/rpc_gen => ../../rpc_gen
	github.com/apache/thrift => github.com/apache/thrift v0.13.0
)

require (
	github.com/MoScenix/mes/rpc_gen v0.0.0-00010101000000-000000000000
	github.com/cloudwego/kitex v0.15.4
	github.com/kitex-contrib/obs-opentelemetry/logging/logrus v0.0.0-20251121033812-f6c3e41f13e9
	github.com/kr/pretty v0.3.1
	go.uber.org/zap v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/validator.v2 v2.0.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)
