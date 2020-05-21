module github.com/iamtheyammer/canvascbl/backend

go 1.13

require (
	github.com/Masterminds/squirrel v1.1.0
	github.com/aws/aws-sdk-go v1.29.8
	github.com/getsentry/sentry-go v0.6.1
	github.com/iamtheyammer/cfjwt v0.1.3
	github.com/julienschmidt/httprouter v1.2.0
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v2.0.2+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sendgrid/rest v2.4.1+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.5.0+incompatible
	github.com/stripe/stripe-go v66.0.0+incompatible
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
)

// +heroku goVersion go1.13.3
