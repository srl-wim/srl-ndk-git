package agent

import (
	"context"
	"log"

	//srlNDK "github.com/henderiw/srl-ndk-agents/protobuf/sdk_protos/nokia.com/srlinux/sdk/protos"
	srlndk "github.com/srl-wim/protos"
	"google.golang.org/grpc/metadata"
)

func (a *Agent) updateTelemetry(jsPath *string, jsData *string) {
	ctx := context.Background()

	// Set up agent name
	ctx = metadata.AppendToOutgoingContext(ctx, "agent_name", a.Name)

	telClient := srlndk.NewSdkMgrTelemetryServiceClient(a.GRPCConn)

	key := &srlndk.TelemetryKey{JsPath: *jsPath}
	data := &srlndk.TelemetryData{JsonContent: *jsData}
	entry := &srlndk.TelemetryInfo{Key: key, Data: data}
	telReq := &srlndk.TelemetryUpdateRequest{}
	telReq.State = make([]*srlndk.TelemetryInfo, 0)
	telReq.State = append(telReq.State, entry)

	r1, err := telClient.TelemetryAddOrUpdate(ctx, telReq)
	if err != nil {
		log.Fatalf("Could not update telemetry for key : %s", *jsPath)
	}
	log.Printf("Telemetry add/update status: %s error_string: %s", r1.GetStatus(), r1.GetErrorStr())
}

func (a *Agent) deleteTelemetry(jsPath *string) {
	ctx := context.Background()

	// Set up agent name
	ctx = metadata.AppendToOutgoingContext(ctx, "agent_name", a.Name)

	telClient := srlndk.NewSdkMgrTelemetryServiceClient(a.GRPCConn)

	key := &srlndk.TelemetryKey{JsPath: *jsPath}
	telReq := &srlndk.TelemetryDeleteRequest{}
	telReq.Key = make([]*srlndk.TelemetryKey, 0)
	telReq.Key = append(telReq.Key, key)

	r1, err := telClient.TelemetryDelete(ctx, telReq)
	if err != nil {
		log.Fatalf("Could not delete telemetry for key : %s", *jsPath)
	}
	log.Printf("Telemetry delete status: %s error_string: %s", r1.GetStatus(), r1.GetErrorStr())
}
