module github.com/shinfan/google-dca-test

go 1.11

replace google.golang.org/api => ./api

replace cloud.google.com/go => github.com/googleapis/google-cloud-go v0.72.1-0.20201112055843-869bd24e213e

replace cloud.google.com/go/pubsub => github.com/googleapis/google-cloud-go/pubsub v1.8.4-0.20201112055843-869bd24e213e

replace cloud.google.com/go/bigquery => github.com/googleapis/google-cloud-go/bigquery v1.13.1-0.20201112055843-869bd24e213e

replace cloud.google.com/go/spanner => github.com/googleapis/google-cloud-go/spanner v1.12.1-0.20201111204810-23b69042c584

replace cloud.google.com/go/bigtable => github.com/googleapis/google-cloud-go/bigtable v1.6.1-0.20201112125007-99f5ac76814f

replace cloud.google.com/go/logging => github.com/googleapis/google-cloud-go/logging v1.1.3-0.20201112055843-869bd24e213e

replace cloud.google.com/go/storage => github.com/googleapis/google-cloud-go/storage v1.1.2-0.20201027002940-a0982f79d646

require (
	cloud.google.com/go v0.71.0
	cloud.google.com/go/bigquery v1.9.0
	cloud.google.com/go/bigtable v1.4.0
	cloud.google.com/go/datastore v1.2.0
	cloud.google.com/go/logging v1.0.0
	cloud.google.com/go/pubsub v1.6.1
	cloud.google.com/go/spanner v1.10.0
	cloud.google.com/go/storage v1.11.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/fluent/fluent-logger-golang v1.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.2
	github.com/google/uuid v1.1.2
	github.com/kr/pretty v0.2.1 // indirect
	github.com/mailgun/mailgun-go/v3 v3.6.4
	github.com/mailjet/mailjet-apiv3-go v0.0.0-20190724151621-55e56f74078c
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sendgrid/smtpapi-go v0.6.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
	google.golang.org/api v0.35.0
	google.golang.org/appengine v1.6.7
	google.golang.org/genproto v0.0.0-20201111145450-ac7456db90a6
	google.golang.org/grpc v1.33.2
	gopkg.in/sendgrid/sendgrid-go.v2 v2.0.0
)

// https://github.com/jstemmer/go-junit-report/issues/107
replace github.com/jstemmer/go-junit-report => github.com/tbpg/go-junit-report v0.9.2-0.20200506144438-50086c54f894
