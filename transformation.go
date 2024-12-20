package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/zinscky/log"
)

type Args struct {
	Event        string
	Config       map[string]string
	Log          log.Logger
	DecisionResp bool
}

type Resp struct {
	Args Args
}

// Transformation is the interface that we're exposing as a plugin.
type Transformation interface {
	Execute(args Args) (Args, error)
}

// Here is an implementation that talks over RPC
type TransformationRPC struct{ client *rpc.Client }

func (g *TransformationRPC) Execute(args Args) (Args, error) {
	resp := Resp{}
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
type TransformationRPCServer struct {
	// This is the real implementation
	Impl Transformation
}

func (s *TransformationRPCServer) Execute(args Args, resp *Resp) error {
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
type TransformationPlugin struct {
	// Impl Injection
	Impl Transformation
}

func (p *TransformationPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &TransformationRPCServer{Impl: p.Impl}, nil
}

func (TransformationPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &TransformationRPC{client: c}, nil
}
