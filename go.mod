module github.com/clint456/north-sdk-go

go 1.25.0

replace (
	github.com/edgexfoundry/go-mod-bootstrap/v4 => ../go-mod-bootstrap
	github.com/edgexfoundry/go-mod-configuration/v4 => ../go-mod-configuration
	github.com/edgexfoundry/go-mod-core-contracts/v4 => ../go-mod-core-contracts
	github.com/edgexfoundry/go-mod-messaging/v4 => ../go-mod-messaging
	github.com/edgexfoundry/go-mod-registry/v4 => ../go-mod-registry
	github.com/edgexfoundry/go-mod-secrets/v4 => ../go-mod-secrets
)

require (
	github.com/bytedance/gopkg v0.1.4
	github.com/edgexfoundry/go-mod-core-contracts/v4 v4.0.0-00010101000000-000000000000
)

require (
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
)
