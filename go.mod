module github.com/mailbadger/app

go 1.17

require (
	github.com/aws/aws-sdk-go v1.42.11
	github.com/aws/aws-sdk-go-v2 v1.11.2
	github.com/aws/aws-sdk-go-v2/config v1.10.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.11.0
	github.com/cbroglie/mustache v1.3.0
	github.com/didip/tollbooth v4.0.2+incompatible
	github.com/didip/tollbooth_gin v0.0.0-20170928041415-5752492be505
	github.com/gavv/httpexpect/v2 v2.2.0
	github.com/gin-contrib/sessions v0.0.4
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.10.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/go-github/v25 v25.1.3
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/gorilla/csrf v1.7.1
	github.com/gwatts/gin-adapter v0.0.0-20170508204228-c44433c485ad
	github.com/huandu/facebook v2.3.1+incompatible
	github.com/jinzhu/now v1.1.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/open-policy-agent/opa v0.36.0
	github.com/rakyll/statik v0.1.7
	github.com/robbiet480/go.sns v0.0.0-20181124163742-ca087b49e1da
	github.com/rubenv/sql-migrate v0.0.0-20200616145509-8d140a17f351
	github.com/segmentio/ksuid v1.0.4
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/unrolled/secure v1.0.9
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	golang.org/x/net v0.0.0-20211111083644-e5c967477495
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/api v0.67.0
	gopkg.in/ezzarghili/recaptcha-go.v3 v3.0.1
	gorm.io/driver/mysql v1.2.2
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.4
)

require (
	cloud.google.com/go/compute v0.1.0 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.5.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.9.0 // indirect
	github.com/aws/smithy-go v1.9.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/structs v1.0.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.13.5 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/ugorji/go/codec v1.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.9.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/yashtewari/glob-intersection v0.0.0-20180916065949-5c77d914dd0b // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220126215142-9970aeb2e350 // indirect
	google.golang.org/grpc v1.43.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	moul.io/http2curl v1.0.1-0.20190925090545-5cd742060b0e // indirect
)
