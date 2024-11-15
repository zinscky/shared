package shared

import (
	"net/http"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/zinscky/log"
)

// Function is the interface that we're exposing as a plugin.
type Function interface {
	Execute(args FunctionArgs) FunctionArgs
}

type FunctionArgs struct {
	Req        *http.Request
	Log        *log.Logger
	Resp       any
	Headers    map[string]string
	StatusCode int
}

type FunctionResp struct {
	Args FunctionArgs
}

// Here is an implementation that talks over RPC
type FunctionRPC struct{ client *rpc.Client }

func (g *FunctionRPC) Execute(args FunctionArgs) FunctionArgs {
	resp := FunctionResp{}
	err := g.client.Call("Plugin.Execute", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		// panic(err)
		return resp.Args
	}

	return resp.Args
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type FunctionRPCServer struct {
	// This is the real implementation
	Impl Function
}

func (s *FunctionRPCServer) Execute(args FunctionArgs, resp *FunctionResp) error {
	var err error
	resp.Args = s.Impl.Execute(args)
	return err
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type FunctionPlugin struct {
	// Impl Injection
	Impl Function
}

func (p *FunctionPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &FunctionRPCServer{Impl: p.Impl}, nil
}

func (FunctionPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FunctionRPC{client: c}, nil
}
