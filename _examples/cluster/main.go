package main

import (
	"github.com/anson-xcloud/xdp-demo/api"
	"github.com/anson-xcloud/xdp-demo/server"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

func hello(svr server.Server, req *server.Request) {
	svr.Reply(req, []byte("hello"))
}

func onHandshake(svr server.Server, req *server.Request) {
	var evt api.EventHandshake
	if err := proto.Unmarshal(req.Data.Data, &evt); err != nil {
		return
	}

	server.GetLogger().Debug("service %s(%d) connected", evt.Appid, evt.Tempid)
}

func main() {
	server.SetEnv("debug")

	eg := new(errgroup.Group)
	eg.Go(clusterServiceOne)
	eg.Go(clusterServiceTwo("two1.yaml"))
	eg.Go(clusterServiceTwo("two2.yaml"))
	eg.Wait()
}

func clusterServiceOne() error {
	sm := server.NewServeMux()
	sm.HandleFunc(server.HandlerRemoteAll, "one", hello)
	sm.HandleFunc(server.HandlerRemoteXcloud, "on_handshake", onHandshake)

	svr := server.NewServer(server.WithHandler(sm))
	return svr.Serve("appcluster:appkey", server.WithConfigFile("one.yaml"))
}

func clusterServiceTwo(fp string) func() error {
	return func() error {
		sm := server.NewServeMux()
		sm.HandleFunc(server.HandlerRemoteAll, "two", hello)
		sm.HandleFunc(server.HandlerRemoteXcloud, "on_handshake", onHandshake)

		svr := server.NewServer(server.WithHandler(sm))
		return svr.Serve("appcluster:appkey", server.WithConfigFile(fp))
	}
}
