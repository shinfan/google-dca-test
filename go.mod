module github.com/shinfan/google-dca-test

go 1.11

replace google.golang.org/api => ../api

require (
	cloud.google.com/go v0.70.0
	cloud.google.com/go/bigquery v1.9.0
	cloud.google.com/go/bigtable v1.4.0
	cloud.google.com/go/datastore v1.2.0
	cloud.google.com/go/logging v1.0.0
	cloud.google.com/go/pubsub v1.6.1
	cloud.google.com/go/spanner v1.10.0
	cloud.google.com/go/storage v1.11.0
	github.com/aws/aws-sdk-go v1.34.18
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/fluent/fluent-logger-golang v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.2
	github.com/google/uuid v1.1.2
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/linkedin/goavro/v2 v2.9.8
	github.com/mailgun/mailgun-go/v3 v3.6.4
	github.com/mailjet/mailjet-apiv3-go v0.0.0-20190724151621-55e56f74078c
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/sendgrid/smtpapi-go v0.6.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	go.opencensus.io v0.22.5
	golang.org/x/net v0.0.0-20201020065357-d65d470038a5 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f
	golang.org/x/tools v0.0.0-20201020123448-f5c826d1900e // indirect
	google.golang.org/api v0.33.0
	google.golang.org/appengine v1.6.6
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.32.0
	gopkg.in/sendgrid/sendgrid-go.v2 v2.0.0
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
