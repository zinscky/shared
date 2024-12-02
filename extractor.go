package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/zinscky/log"
)

type ExtractorArgs struct {
	Data   []interface{}
	Config map[string]string
	Log    log.Logger
}

type ExtractorResp struct {
	Args ExtractorArgs
}

// Transformation is the interface that we're exposing as a plugin.
type Extractor interface {
	Execute(args ExtractorArgs) (ExtractorArgs, error)
}

// Here is an implementation that talks over RPC
type ExtractorRPC struct{ client *rpc.Client }

func (g *ExtractorRPC) Execute(args ExtractorArgs) (ExtractorArgs, error) {
	resp := ExtractorResp{}
	err := g.client.Call("Plugin.Execute", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		// panic(err)
		return resp.Args, err
	}

	return resp.Args, nil
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type ExtractorRPCServer struct {
	// This is the real implementation
	Impl Extractor
}

func (s *ExtractorRPCServer) Execute(args ExtractorArgs, resp *ExtractorResp) error {
	var err error
	resp.Args, err = s.Impl.Execute(args)
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
type ExtractorPlugin struct {
	// Impl Injection
	Impl Extractor
}

func (p *ExtractorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ExtractorRPCServer{Impl: p.Impl}, nil
}

func (ExtractorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ExtractorRPC{client: c}, nil
}
