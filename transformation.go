package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Args struct {
	Event  string
	Config map[string]string
}

// Transformation is the interface that we're exposing as a plugin.
type Transformation interface {
	Execute(args Args) (string, error)
}

// Here is an implementation that talks over RPC
type TransformationRPC struct{ client *rpc.Client }

func (g *TransformationRPC) Execute(args Args) (string, error) {
	var resp string
	err := g.client.Call("Plugin.Execute", args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		// panic(err)
		return "", err
	}

	return resp, nil
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type TransformationRPCServer struct {
	// This is the real implementation
	Impl Transformation
}

func (s *TransformationRPCServer) Execute(args Args, resp *string) error {
	var err error
	*resp, err = s.Impl.Execute(args)
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
