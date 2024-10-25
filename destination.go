package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Destination is the interface that we're exposing as a plugin.
type Destination interface {
	Execute(args Args) (Args, error)
	Setup(args Args) (Args, error)
	Teardown(args Args) (Args, error)
}

// Here is an implementation that talks over RPC
type DestinationRPC struct{ client *rpc.Client }

func (g *DestinationRPC) Execute(args Args) (Args, error) {
	resp := Resp{}
	err := g.client.Call("Plugin.Execute", args, &resp)
	if err != nil {
		return resp.Args, err
	}
	return resp.Args, nil
}

func (g *DestinationRPC) Setup(args Args) (Args, error) {
	resp := Resp{}
	err := g.client.Call("Plugin.Setup", args, &resp)
	if err != nil {
		return resp.Args, err
	}
	return resp.Args, nil
}

func (g *DestinationRPC) TearDown(args Args) (Args, error) {
	resp := Resp{}
	err := g.client.Call("Plugin.Teardown", args, &resp)
	if err != nil {
		return resp.Args, err
	}
	return resp.Args, nil
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type DestinationRPCServer struct {
	// This is the real implementation
	Impl Destination
}

func (s *DestinationRPCServer) Execute(args Args, resp *Resp) error {
	var err error
	resp.Args, err = s.Impl.Execute(args)
	return err
}

func (s *DestinationRPCServer) Setup(args Args, resp *Resp) error {
	var err error
	resp.Args, err = s.Impl.Setup(args)
	return err
}

func (s *DestinationRPCServer) TearDown(args Args, resp *Resp) error {
	var err error
	resp.Args, err = s.Impl.Teardown(args)
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
type DestinationPlugin struct {
	// Impl Injection
	Impl Destination
}

func (p *DestinationPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &DestinationRPCServer{Impl: p.Impl}, nil
}

func (DestinationPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &DestinationRPC{client: c}, nil
}
