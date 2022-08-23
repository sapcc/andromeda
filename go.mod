module github.com/sapcc/andromeda

go 1.17

replace github.com/micro/protoc-gen-micro => go-micro.dev/v4/cmd/protoc-gen-micro/v4 v4.8.1

replace github.com/micro/go-micro => go-micro.dev/v4 v4.8.1

require (
	github.com/Boostport/migration v1.1.1
	github.com/Boostport/migration/driver/mysql v1.1.1
	github.com/Boostport/migration/driver/postgres v1.1.1
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 v2.15.0
	github.com/cockroachdb/cockroach-go/v2 v2.2.15
	github.com/didip/tollbooth v4.0.2+incompatible
	github.com/dre1080/recovr v1.0.3
	github.com/go-micro/plugins/v4/broker/nats v1.1.0
	github.com/go-micro/plugins/v4/registry/nats v1.1.0
	github.com/go-micro/plugins/v4/transport/nats v1.1.0
	github.com/go-openapi/errors v0.20.3
	github.com/go-openapi/loads v0.21.2
	github.com/go-openapi/runtime v0.24.1
	github.com/go-openapi/spec v0.20.7
	github.com/go-openapi/strfmt v0.21.3
	github.com/go-openapi/swag v0.22.3
	github.com/go-openapi/validate v0.22.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gophercloud/gophercloud v0.25.0
	github.com/gophercloud/utils v0.0.0-20220704184730-55bdbbaec4ba
	github.com/iancoleman/strcase v0.2.0
	github.com/jackc/pgconn v1.13.0
	github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa
	github.com/jackc/pgx/v4 v4.17.0
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jessevdk/go-flags v1.5.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/rs/cors v1.8.2
	github.com/sapcc/go-bits v0.0.0-20220823141731-fd7850ce8cb7
	github.com/scottdware/go-bigip v0.0.0-20220517145708-9fe3e2f9f005
	github.com/stretchr/testify v1.8.0
	github.com/xo/dburl v0.11.0
	go-micro.dev/v4 v4.8.1
	golang.org/x/net v0.0.0-20220812174116-3211cb980234
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20220407094043-a94812496cf5 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/apex/log v1.9.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/databus23/goslo.policy v0.0.0-20210929125152-81bf2876dbdb // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/evanphx/json-patch/v5 v5.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-acme/lego/v4 v4.4.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/gobuffalo/logger v1.0.6 // indirect
	github.com/gobuffalo/packd v1.0.1 // indirect
	github.com/gobuffalo/packr/v2 v2.8.3 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.0.4 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/karrick/godirwalk v1.17.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.15.3 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/markbates/errx v1.1.0 // indirect
	github.com/markbates/oncer v1.0.0 // indirect
	github.com/markbates/safe v1.0.1 // indirect
	github.com/matryer/is v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/miekg/dns v1.1.48 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/nats-io/nats-server/v2 v2.7.4 // indirect
	github.com/nats-io/nats.go v1.13.1-0.20220308171302-2f2f6968e98d // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/urfave/cli/v2 v2.6.0 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	github.com/ztrue/tracerr v0.3.0 // indirect
	go.mongodb.org/mongo-driver v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220411224347-583f2d630306 // indirect
	golang.org/x/tools v0.1.11-0.20220316014157-77aa08bb151a // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
