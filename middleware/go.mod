module github.com/PlatformCore/libpackage/middleware

go 1.23.0

require (
	github.com/PlatformCore/libpackage/messaging v0.0.0
	github.com/PlatformCore/libpackage/observability v0.0.0
	github.com/PlatformCore/libpackage/ratelimit v0.0.0
	github.com/PlatformCore/libpackage/resilience v0.0.0
	github.com/PlatformCore/libpackage/transport v0.0.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/prometheus/client_golang v1.23.2
	go.opentelemetry.io/otel v1.43.0
	go.opentelemetry.io/otel/trace v1.43.0
	go.opentelemetry.io/otel/metric v1.43.0
	go.uber.org/zap v1.28.0
	golang.org/x/time v0.14.0
	google.golang.org/grpc v1.80.0
)

replace github.com/PlatformCore/libpackage/transport => ../transport
replace github.com/PlatformCore/libpackage/observability => ../observability
replace github.com/PlatformCore/libpackage/ratelimit => ../ratelimit
replace github.com/PlatformCore/libpackage/resilience => ../resilience
replace github.com/PlatformCore/libpackage/messaging => ../messaging
