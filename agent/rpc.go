package agent

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// Server struct
type Server struct {
	Agent *Agent
}

// Args struct
type Args struct {
	Comment string
}

// rpcRequest represents a RPC request.
// rpcRequest implements the io.ReadWriteCloser interface.
type rpcRequest struct {
	r    io.Reader     // holds the JSON formated RPC request
	rw   io.ReadWriter // holds the JSON formated RPC response
	done chan bool     // signals then end of the RPC request
}

// NewRPCRequest returns a new rpcRequest.
func NewRPCRequest(r io.Reader) *rpcRequest {
	var buf bytes.Buffer
	done := make(chan bool)
	return &rpcRequest{r, &buf, done}
}

// Read implements the io.ReadWriteCloser Read method.
func (r *rpcRequest) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// Write implements the io.ReadWriteCloser Write method.
func (r *rpcRequest) Write(p []byte) (n int, err error) {
	return r.rw.Write(p)
}

// Close implements the io.ReadWriteCloser Close method.
func (r *rpcRequest) Close() error {
	r.done <- true
	return nil
}

// Call invokes the RPC request, waits for it to complete, and returns the results.
func (r *rpcRequest) Call() io.Reader {
	go jsonrpc.ServeConn(r)
	<-r.done
	return r.rw
}

// StartJSONRPCServer start a JSONRPC server and waits for connection
func StartJSONRPCServer(a *Agent) {
	s := new(Server)
	s.Agent = a
	rpc.Register(s)

	//rpc.Register(arith)
	rpc.Register(s)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(req.Body).Call()
		io.Copy(w, res)
	})
	log.Fatal(http.ListenAndServe(":7777", nil))

}

// Branch function
func (s *Server) Branch(args *Args, reply *string) error {
	log.Print("Git branch action")
	a := s.Agent
	if err := a.GetRef(&a.Config.YangConfig.Branch.Value); err != nil {
		log.Printf("Error: Unable to get/create the commit reference: %s\n", err)
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error Unable to get/create the commit reference: %s\n", err)
		return nil
	}
	if a.Github.Ref == nil {
		log.Printf("Error: No error where returned but the reference is nil")
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error: No error where returned but the reference is nil")
		return nil
	}
	a.Github.state.Success++
	a.updateConfigTelemetry()

	*reply = fmt.Sprintf("success")
	return nil
}

// Commit function
func (s *Server) Commit(args *Args, reply *string) error {
	log.Print("Git commit action")
	a := s.Agent
	if err := a.GetRef(&a.Config.YangConfig.Branch.Value); err != nil {
		log.Printf("Error Unable to get/create the commit reference: %s\n", err)
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error Unable to get/create the commit reference: %s\n", err)
		return nil
	}
	if a.Github.Ref == nil {
		log.Printf("Error: No error where returned but the reference is nil")
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error: No error where returned but the reference is nil")
		return nil
	}
	if err := a.GetTree(); err != nil {
		log.Printf("Error Unable to create the tree based on the provided files: %s\n", err)
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error Unable to create the tree based on the provided files: %s\n", err)
		return nil
	}
	if err := a.PushCommit(a.Github.Ref, a.Github.Tree); err != nil {
		log.Printf("Error Unable to create the commit: %s\n", err)
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error Unable to create the commit: %s\n", err)
		return nil
	}
	a.Github.state.Success++
	a.updateConfigTelemetry()

	*reply = fmt.Sprintf("success")
	return nil
}

// PullRequest function
func (s *Server) PullRequest(args *Args, reply *string) error {
	log.Print("Git pull-request action")
	a := s.Agent
	if err := a.GetRef(&a.Config.YangConfig.Branch.Value); err != nil {
		log.Printf("Error Unable to get/create the commit reference: %s\n", err)
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error Unable to get/create the commit reference: %s\n", err)
		return nil
	}
	if a.Github.Ref == nil {
		log.Printf("Error: No error where returned but the reference is nil")
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error: No error where returned but the reference is nil")
		return nil
	}
	if err := a.GetTree(); err != nil {
		log.Printf("Error: Unable to create the tree based on the provided files: %s\n", err)
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error: Unable to create the tree based on the provided files: %s\n", err)
		return nil
	}
	if err := a.CreatePR(&a.Config.YangConfig.Branch.Value); err != nil {
		log.Printf("Error while creating the pull request: %s", err)
		a.Github.state.Failure++
		a.updateConfigTelemetry()
		*reply = fmt.Sprintf("Error while creating the pull request: %s", err)
		return nil
	}
	a.Github.state.Success++
	a.updateConfigTelemetry()

	*reply = fmt.Sprintf("success")

	return nil
}
