package XWeb

import (
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler func(ctx *Context)
type Node struct {
	children  map[string]*Node
	param     bool   //儿子是否为参数节点 如果是参数节点 那么就可以匹配参数了
	paramName string //参数节点名称
	handler   Handler
}
type Router interface {
	Bind(url string, method Method, handler Handler)
	Handler(ctx *fasthttp.RequestCtx)
}
type DefaultRouter struct {
	root []*Node
}
type Method int

const (
	MethodPost = iota
	MethodGet
)

var methodLink = map[string]int{
	"POST": MethodPost,
	"GET":  MethodGet,
}
var methodRLink map[int]string = func() map[int]string {
	var ret = make(map[int]string)
	for k, v := range methodLink {
		ret[v] = k
	}
	return ret
}()

func NewDefaultRouter() Router {
	ret := new(DefaultRouter)
	ret.root = make([]*Node, 2)
	for i := 0; i < len(ret.root); i++ {
		ret.root[i] = newNode()
	}
	return ret
}
func newNode() *Node {
	node := new(Node)
	node.children = make(map[string]*Node, 0)
	return node
}
func (z *DefaultRouter) Bind(url string, method Method, handler Handler) {
	res := strings.Split(url, "\\")
	nowNode := z.root[method]
	for i, v := range res {
		if v[0] == ':' { //参数 那就必须是最后一个了
			if i != len(res)-1 {
				panic("error")
			}
			if len(nowNode.children) != 0 {
				panic("wrong children of node")
			}
			nowNode.param = true
			nowNode.paramName = v[1:]
		} else {
			_, ok := nowNode.children[v]
			if !ok {
				nowNode.children[v] = newNode()
			}
			nowNode = nowNode.children[v]
		}
	}
	nowNode.handler = handler
}
func (z *DefaultRouter) Handler(ctx *fasthttp.RequestCtx) {
	res := strings.Split(ctx.Request.URI().String(), "\\")
	nowNode := z.root[methodLink[string(ctx.Method())]]
	rctx := newContext(ctx)
	for _, v := range res {
		if nowNode.param {
			ctx.SetUserValue(nowNode.paramName, v)
			break
		}
		node, ok := nowNode.children[v]
		if !ok { // 未知请求路径 返回404
			ctx.SetStatusCode(404)
			return
		}
		nowNode = node
	}
	nowNode.handler(rctx)
}

// 推荐使用这个
type FasthttpRouter struct {
	raw fasthttprouter.Router
}

func (z *FasthttpRouter) newWrapper(handler Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		handler(newContext(ctx))
	}
}
func (z *FasthttpRouter) Bind(address string, method Method, handler Handler) {
	handle := z.newWrapper(handler)
	switch method {
	case MethodGet:
		z.raw.GET(address, handle)
	case MethodPost:
		z.raw.POST(address, handle)
	default:
		panic("XWEB doesn't implement this method yet")
	}
}
func (z *FasthttpRouter) Handler(ctx *fasthttp.RequestCtx) {
	z.raw.Handler(ctx)
}
