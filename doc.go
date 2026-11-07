// Package appkit is the root of go-app-kit, a shared HTTP application kit that
// extracts the common web-server plumbing used by multiple Gin-based services
// (request Context, response envelope, error codes, middleware, process
// bootstrap, logging/telemetry helpers and field encryption).
//
// # Design constraints
//
//   - No global singletons: callers construct dependencies and pass Config
//     structs; middleware and helpers accept their dependencies explicitly.
//   - ConfigDefault: N/A at package root; subpackages (e.g. telemetry) expose ConfigDefault where applicable.
//   - No log.Fatal / os.Exit inside library code: failures are returned as
//     errors so the caller decides how to react.
//   - Business logic stays out of the kit: it never imports caller business
//     packages (controller/service/repository/model). The kit exposes framework
//     primitives only; callers compose them into their own layering.
//   - Context propagation: per-request work flows through core.Context, which
//     embeds *gin.Context so handlers keep full Gin access.
//
// # Known exceptions (documented, intentional)
//
//   - gin dependency: this is an HTTP application kit, not a pure SDK, so it
//     intentionally depends on github.com/gin-gonic/gin. core.Context embeds
//     *gin.Context and middleware are Gin handlers.
//   - viper in bootstrap/ and config/: the bootstrap and config packages may use
//     github.com/spf13/viper for optional Load* helpers and process startup
//     (loading conf.yaml, APP_* env overrides, Nacos sync). This mirrors the
//     documented go-nacos env exception: it is confined to process startup and
//     is never used by the request-path packages (core, errcode, middleware).
//
// # Packages
//
//   - def         framework constants (env names, language codes, header keys)
//   - core        request Context, Handle decorator, response envelope, Subject
//   - errcode     ErrCode type plus shared Lib*/Web* error codes
//   - middleware  reusable Gin middleware (recovery, logging, cors, ...)
//   - bootstrap   process startup: viper + optional Nacos sync + secret checks
//   - config      pure config structs and viper Load helpers (no service deps)
//   - logging     process-level logging setup helpers
//   - telemetry   optional tracing/metrics bootstrap helpers
//   - singleflight in-process request de-duplication
//   - fieldcrypt  field-level encryption / masking / normalization helpers
package appkit
