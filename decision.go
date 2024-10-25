package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Decision is the interface that we're exposing as a plugin.
type Decision interface {
	Execute(args Args) (Args, error)
}

// Here is an implementation that talks over RPC
type DecisionRPC struct{ client *rpc.Client }

func (g *DecisionRPC) Execute(args Args) (Args, error) {
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
type DecisionRPCServer struct {
	// This is the real implementation
	Impl Decision
}

func (s *DecisionRPCServer) Execute(args Args, resp *Resp) error {
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
type DecisionPlugin struct {
	// Impl Injection
	Impl Decision
}

func (p *DecisionPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DecisionRPCServer{Impl: p.Impl}, nil
}

func (DecisionPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DecisionRPC{client: c}, nil
}
