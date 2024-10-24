package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Destination is the interface that we're exposing as a plugin.
type Destination interface {
	Setup(args Args) error
	Execute(args Args) error
	Teardown(args Args) error
}

// Here is an implementation that talks over RPC
type DestinationRPC struct{ client *rpc.Client }

func (d *DestinationRPC) Execute(args Args) error {
	var resp string
	return d.client.Call("Plugin.Execute", args, &resp)
}

func (d *DestinationRPC) Setup(args Args) error {
	var resp string
	err := d.client.Call("Plugin.Setup", args, &resp)
	return err
}

func (d *DestinationRPC) Teardown(args Args) error {
	var resp string
	err := d.client.Call("Plugin.Teardown", args, &resp)
	return err
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type DestinationRPCServer struct {
	// This is the real implementation
	Impl Destination
}

func (s *DestinationRPCServer) Execute(args Args, resp *string) error {
	return s.Impl.Execute(args)
}

func (s *DestinationRPCServer) Setup(args Args, resp *string) error {
	return s.Impl.Setup(args)
}

func (s *DestinationRPCServer) Teardown(args Args, resp *string) error {
	return s.Impl.Teardown(args)
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
