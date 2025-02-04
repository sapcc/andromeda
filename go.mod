module github.com/sapcc/andromeda

go 1.23

toolchain go1.23.5

require (
	github.com/Boostport/migration v1.1.2
	github.com/Boostport/migration/driver/mysql v1.1.2
	github.com/Boostport/migration/driver/postgres v1.1.3
	github.com/Masterminds/squirrel v1.5.4
	github.com/actatum/stormrpc v0.4.2
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v9 v9.1.0
	github.com/apex/log v1.9.0
	github.com/cockroachdb/cockroach-go/v2 v2.3.8
	github.com/didip/tollbooth v4.0.2+incompatible
	github.com/dlmiddlecote/sqlstats v1.0.2
	github.com/dre1080/recovr v1.0.3
	github.com/getsentry/sentry-go v0.25.0
	github.com/go-openapi/errors v0.22.0
	github.com/go-openapi/loads v0.22.0
	github.com/go-openapi/runtime v0.28.0
	github.com/go-openapi/spec v0.21.0
	github.com/go-openapi/strfmt v0.23.0
	github.com/go-openapi/swag v0.23.0
	github.com/go-openapi/validate v0.24.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/gophercloud/gophercloud/v2 v2.4.0
	github.com/gophercloud/utils/v2 v2.0.0-20250106080359-c8704d65b4fa
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/iancoleman/strcase v0.3.0
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438
	github.com/jackc/pgx/v5 v5.7.2
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jessevdk/go-flags v1.6.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
	github.com/mcuadros/go-defaults v1.2.0
	github.com/nats-io/nats.go v1.38.0
	github.com/prometheus/client_golang v1.20.5
	github.com/rs/cors v1.11.1
	github.com/sapcc/go-api-declarations v1.13.2
	github.com/sapcc/go-bits v0.0.0-20250204142118-dd50ca8a0a5f
	github.com/scottdware/go-bigip v0.0.0-20240809002616-deb9b0aff84a
	github.com/slok/go-http-metrics v0.13.0
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v2 v2.27.5
	github.com/xo/dburl v0.23.2
	golang.org/x/net v0.34.0
	google.golang.org/protobuf v1.36.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/databus23/goslo.policy v0.0.0-20210929125152-81bf2876dbdb // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/gobuffalo/logger v1.0.3 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr/v2 v2.8.0 // indirect
	github.com/gofrs/uuid/v5 v5.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/karrick/godirwalk v1.15.3 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/logrusorgru/aurora v0.0.0-20181002194514-a7b3b318ed4e // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/markbates/errx v1.1.0 // indirect
	github.com/markbates/oncer v1.0.0 // indirect
	github.com/markbates/safe v1.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/ztrue/tracerr v0.3.0 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/term v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)
