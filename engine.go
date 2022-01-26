package XWeb

import (
	"github.com/valyala/fasthttp"
)

type XWeb struct {
	router Router
}

func NewEngine() XWeb {
	return XWeb{}
}
func (z *XWeb) SetRouter(router Router) {
	z.router = router
}
func (z *XWeb) RunAndServe(address string) error {
	return fasthttp.ListenAndServe(address, z.router.Handler)
}
func (z *XWeb) RunAndServerTLS(address string, certData []byte, keyData []byte) error {
	return fasthttp.ListenAndServeTLSEmbed(address, certData, keyData, z.router.Handler)
}
