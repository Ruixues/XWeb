package XWeb

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type Context struct {
	raw *fasthttp.RequestCtx
}

func newContext(raw *fasthttp.RequestCtx) *Context {
	return &Context{
		raw: raw,
	}
}

func (z *Context) UserValue(key string) any {
	return z.raw.UserValue(key)
}
func (z *Context) WriteJson(object any) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	_, err = z.raw.Write(data)
	return err
}
func (z *Context) WriteString(str string) error {
	_, err := z.raw.WriteString(str)
	return err
}
