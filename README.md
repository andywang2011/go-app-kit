# go-app-kit

Shared HTTP application kit for Gin services (context envelope, middleware, bootstrap, logging, field encryption).

Packages: `def` (constants), `core` (request context), `errcode`, `middleware`, `bootstrap`, `config`, `logging`, `telemetry`, `singleflight`, `fieldcrypt`; root `appkit` defines `Module` / `Closer`.

## Install

```bash
go get github.com/andywang2011/go-app-kit@v0.1.0
```

## Quick start

```go
gin.SetMode(gin.TestMode)
r := gin.New()
r.Use(kitmw.Recovery)
r.GET("/ping", core.Handle(func(c *core.Context) {
	c.Success(gin.H{"ok": true})
}))

w := httptest.NewRecorder()
req := httptest.NewRequest("GET", "/ping", nil)
r.ServeHTTP(w, req)
fmt.Println(w.Code)
```

## Examples

```bash
cd examples
go run .
```

## License

MIT — see [LICENSE](LICENSE).
