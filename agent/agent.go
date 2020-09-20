package agent

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/google/go-github/v32/github"
	srlndk "github.com/srl-wim/protos"
	"google.golang.org/grpc"
)

var grpcAddress = "localhost:50053"

var retryTimeout = 5 * time.Second

// HandleFunc function
type HandleFunc func(context.Context, *srlndk.NotificationStreamResponse)

// Agent type
type Agent struct {
	Name       string
	RetryTimer time.Duration

	AppID uint32

	GRPCConn *grpc.ClientConn

	SdkMgrService struct {
		Client srlndk.SdkMgrServiceClient
	}
	NotificationService struct {
		Client srlndk.SdkNotificationServiceClient
	}
	TelemetryService struct {
		Client srlndk.SdkMgrTelemetryServiceClient
	}
	RouteService struct {
		Client srlndk.SdkMgrRouteServiceClient
	}
	MPLSRouteService struct {
		Client srlndk.SdkMgrMplsRouteServiceClient
	}
	NextHopGroupService struct {
		Client srlndk.SdkMgrNextHopGroupServiceClient
	}
	Config struct {
		cfgTranxMap map[string][]cfgTranxEntry
		YangConfig  *yangGit
	}

	Github struct {
		ctx           context.Context
		client        *github.Client
		token         *string
		commitMessage *string
		baseBranch    *string
		prSubject     *string
		prDescription *string
		file          *string
		Ref           *github.Reference
		Tree          *github.Tree
		state         *gitClientState
	}
}

// NewAgent initializes the agent
func NewAgent(ctx context.Context, name string) (*Agent, error) {
	a := new(Agent)
	a.Name = name
	a.RetryTimer = retryTimeout

	var err error
	a.GRPCConn, err = grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial failed: %v", err)
		return nil, err
	}
	a.SdkMgrService.Client = srlndk.NewSdkMgrServiceClient(a.GRPCConn)

	nctx, cancel := context.WithTimeout(ctx, a.RetryTimer)
	defer cancel()
	r, err := a.SdkMgrService.Client.AgentRegister(nctx, &srlndk.AgentRegistrationRequest{})
	if err != nil {
		return nil, fmt.Errorf("agent %s registration failed: %v", a.Name, err)
	}
	a.AppID = r.GetAppId()
	log.Printf("agent %s: registration status: %v", a.Name, r.GetStatus())
	log.Printf("agent %s: registration appID: %v", a.Name, r.GetAppId())
	// create telemetry and notifications Clients
	a.TelemetryService.Client = srlndk.NewSdkMgrTelemetryServiceClient(a.GRPCConn)
	a.NotificationService.Client = srlndk.NewSdkNotificationServiceClient(a.GRPCConn)
	a.Config.cfgTranxMap = make(map[string][]cfgTranxEntry)
	a.Config.YangConfig = new(yangGit)
	return a, nil
}

// KeepAlive provides readiness of the Agent
func (a *Agent) KeepAlive(ctx context.Context, period time.Duration) {
	newTicker := time.NewTicker(period)
	for {
		select {
		case <-newTicker.C:
			keepAliveResponse, err := a.SdkMgrService.Client.KeepAlive(ctx, &srlndk.KeepAliveRequest{})
			if err != nil {
				log.Printf("agent %s: failed to send keep alive request: %v", a.Name, err)
				continue
			}
			log.Printf("agent %s: received keepAliveResponse, status=%v", a.Name, keepAliveResponse.Status)
		case <-ctx.Done():
			log.Printf("agent %s: received %v, shutting down keepAlives", a.Name, ctx.Err())
			newTicker.Stop()
			return
		}
	}
}

// StartConfigNotificationStream function
func (a *Agent) StartConfigNotificationStream(ctx context.Context) chan *srlndk.NotificationStreamResponse {
CREATESUB:
	// get subscription and streamID
	notificationResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx,
		&srlndk.NotificationRegisterRequest{
			Op: srlndk.NotificationRegisterRequest_Create,
		})
	if err != nil {
		log.Printf("agent %s could not register for config notifications: %v", a.Name, err)
		log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	log.Printf("config notification registration status : %s streamID %d", notificationResponse.Status, notificationResponse.GetStreamId())

	notificationRegisterRequest := &srlndk.NotificationRegisterRequest{
		Op:       srlndk.NotificationRegisterRequest_AddSubscription,
		StreamId: notificationResponse.GetStreamId(),
		SubscriptionTypes: &srlndk.NotificationRegisterRequest_Config{ // config
			Config: &srlndk.ConfigSubscriptionRequest{},
		},
	}
	return a.startNotificationStream(ctx, notificationRegisterRequest, notificationResponse.GetSubId())
}

// StartNwInstNotificationStream function
func (a *Agent) StartNwInstNotificationStream(ctx context.Context) chan *srlndk.NotificationStreamResponse {
CREATESUB:
	// get subscription and streamID
	notificationResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx,
		&srlndk.NotificationRegisterRequest{
			Op: srlndk.NotificationRegisterRequest_Create,
		})
	if err != nil {
		log.Printf("agent %s could not register for NwInst notifications: %v", a.Name, err)
		log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	log.Printf("NwInst notification registration status: %s, subscriptionID=%d, streamID=%d",
		notificationResponse.Status, notificationResponse.GetSubId(), notificationResponse.GetStreamId())

	notificationRegisterRequest := &srlndk.NotificationRegisterRequest{
		Op:       srlndk.NotificationRegisterRequest_AddSubscription,
		StreamId: notificationResponse.GetStreamId(),
		SubscriptionTypes: &srlndk.NotificationRegisterRequest_NwInst{ // NwInst
			NwInst: &srlndk.NetworkInstanceSubscriptionRequest{},
		},
	}
	return a.startNotificationStream(ctx, notificationRegisterRequest, notificationResponse.GetSubId())
}

// StartInterfaceNotificationStream function
func (a *Agent) StartInterfaceNotificationStream(ctx context.Context, ifName string) chan *srlndk.NotificationStreamResponse {
CREATESUB:
	// get subscription and streamID
	notificationResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx,
		&srlndk.NotificationRegisterRequest{
			Op: srlndk.NotificationRegisterRequest_Create,
		})
	if err != nil {
		log.Printf("agent %s could not register for Intf notifications: %v", a.Name, err)
		log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	log.Printf("interface notification registration status: %s, subscriptionID=%d, streamID=%d",
		notificationResponse.Status, notificationResponse.GetSubId(), notificationResponse.GetStreamId())
	if notificationResponse.Status == srlndk.SdkMgrStatus_kSdkMgrFailed {
		log.Printf("interface notification subscribe failed")
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	key := new(srlndk.InterfaceKey)
	if ifName != "" {
		key = &srlndk.InterfaceKey{
			IfName: ifName,
		}
	}
	notificationRegisterRequest := &srlndk.NotificationRegisterRequest{
		Op:       srlndk.NotificationRegisterRequest_AddSubscription,
		StreamId: notificationResponse.GetStreamId(),
		SubscriptionTypes: &srlndk.NotificationRegisterRequest_Intf{ // Intf
			Intf: &srlndk.InterfaceSubscriptionRequest{
				Key: key,
			},
		},
	}
	return a.startNotificationStream(ctx, notificationRegisterRequest, notificationResponse.GetSubId())
}

// StartLLDPNeighNotificationStream function
func (a *Agent) StartLLDPNeighNotificationStream(ctx context.Context, ifName, chassisType, chassisID string) chan *srlndk.NotificationStreamResponse {
CREATESUB:
	// get subscription and streamID
	notificationResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx,
		&srlndk.NotificationRegisterRequest{
			Op: srlndk.NotificationRegisterRequest_Create,
		})
	if err != nil {
		log.Printf("agent %s could not register for Intf notifications: %v", a.Name, err)
		log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	log.Printf("LLDPNeighbor notification registration status: %s, subscriptionID=%d, streamID=%d",
		notificationResponse.Status, notificationResponse.GetSubId(), notificationResponse.GetStreamId())
	if notificationResponse.Status == srlndk.SdkMgrStatus_kSdkMgrFailed {
		log.Printf("LLDPNeighbor notification subscribe failed")
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	key := new(srlndk.LldpNeighborKeyPb)
	if ifName != "" || chassisID != "" || chassisType != "" {
		key = &srlndk.LldpNeighborKeyPb{
			InterfaceName: ifName,
			// ChassisId:     chassisID,
			// ChassisType: ndk.LldpNeighborKeyPb_CHASSIS_COMPONENT,
		}
	}
	notificationRegisterRequest := &srlndk.NotificationRegisterRequest{
		Op:       srlndk.NotificationRegisterRequest_AddSubscription,
		StreamId: notificationResponse.GetStreamId(),
		SubscriptionTypes: &srlndk.NotificationRegisterRequest_LldpNeighbor{ // LLDPNeigh
			LldpNeighbor: &srlndk.LldpNeighborSubscriptionRequest{
				Key: key,
			},
		},
	}
	return a.startNotificationStream(ctx, notificationRegisterRequest, notificationResponse.GetSubId())
}

// StartBFDSessionNotificationStream function
func (a *Agent) StartBFDSessionNotificationStream(ctx context.Context, srcIP, dstIP net.IP, instance uint32) chan *srlndk.NotificationStreamResponse {
CREATESUB:
	// get subscription and streamID
	notificationResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx,
		&srlndk.NotificationRegisterRequest{
			Op: srlndk.NotificationRegisterRequest_Create,
		})
	if err != nil {
		log.Printf("agent %s could not register for Intf notifications: %v", a.Name, err)
		log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	log.Printf("BFDSession notification registration status: %s, subscriptionID=%d, streamID=%d",
		notificationResponse.Status, notificationResponse.GetSubId(), notificationResponse.GetStreamId())
	if notificationResponse.Status == srlndk.SdkMgrStatus_kSdkMgrFailed {
		log.Printf("BFDSession notification subscribe failed")
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	bfdSession := &srlndk.BfdSessionSubscriptionRequest{
		Key: &srlndk.BfdmgrGeneralSessionKeyPb{},
	}
	if srcIP != nil {
		bfdSession.Key.SrcIpAddr = &srlndk.IpAddressPb{Addr: srcIP}
	}
	if dstIP != nil {
		bfdSession.Key.DstIpAddr = &srlndk.IpAddressPb{Addr: dstIP}
	}
	bfdSession.Key.InstanceId = instance
	// bfdSession.Key.Type =
	notificationRegisterRequest := &srlndk.NotificationRegisterRequest{
		Op:       srlndk.NotificationRegisterRequest_AddSubscription,
		StreamId: notificationResponse.GetStreamId(),
		SubscriptionTypes: &srlndk.NotificationRegisterRequest_BfdSession{ // BFDSession
			BfdSession: bfdSession,
		},
	}
	return a.startNotificationStream(ctx, notificationRegisterRequest, notificationResponse.GetSubId())
}

// StartRouteNotificationStream function
func (a *Agent) StartRouteNotificationStream(ctx context.Context, netInstance string, ipAddr net.IP, prefixLen uint32) chan *srlndk.NotificationStreamResponse {
CREATESUB:
	// get subscription and streamID
	notificationResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx,
		&srlndk.NotificationRegisterRequest{
			Op: srlndk.NotificationRegisterRequest_Create,
		})
	if err != nil {
		log.Printf("agent %s could not register for Intf notifications: %v", a.Name, err)
		log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	log.Printf("Route notification registration status: %s, subscriptionID=%d, streamID=%d",
		notificationResponse.Status, notificationResponse.GetSubId(), notificationResponse.GetStreamId())
	if notificationResponse.Status == srlndk.SdkMgrStatus_kSdkMgrFailed {
		log.Printf("Route notification subscribe failed")
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	key := new(srlndk.RouteKeyPb)
	if netInstance != "" {
		key.NetInstName = netInstance
	}
	if ipAddr != nil {
		key.IpPrefix = &srlndk.IpAddrPrefLenPb{
			IpAddr:       &srlndk.IpAddressPb{Addr: ipAddr},
			PrefixLength: prefixLen,
		}
	}
	notificationRegisterRequest := &srlndk.NotificationRegisterRequest{
		Op:       srlndk.NotificationRegisterRequest_AddSubscription,
		StreamId: notificationResponse.GetStreamId(),
		SubscriptionTypes: &srlndk.NotificationRegisterRequest_Route{ // route
			Route: &srlndk.IpRouteSubscriptionRequest{
				Key: key,
			},
		},
	}
	return a.startNotificationStream(ctx, notificationRegisterRequest, notificationResponse.GetSubId())
}

// StartAppIDNotificationStream function
func (a *Agent) StartAppIDNotificationStream(ctx context.Context, id uint32) chan *srlndk.NotificationStreamResponse {
CREATESUB:
	// get subscription and streamID
	notificationResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx,
		&srlndk.NotificationRegisterRequest{
			Op: srlndk.NotificationRegisterRequest_Create,
		})
	if err != nil {
		log.Printf("agent %s could not register for Intf notifications: %v", a.Name, err)
		log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	log.Printf("AppId notification registration status: %s, subscriptionID=%d, streamID=%d",
		notificationResponse.Status, notificationResponse.GetSubId(), notificationResponse.GetStreamId())
	if notificationResponse.Status == srlndk.SdkMgrStatus_kSdkMgrFailed {
		log.Printf("AppId notification subscribe failed")
		time.Sleep(a.RetryTimer)
		goto CREATESUB
	}
	notificationRegisterRequest := &srlndk.NotificationRegisterRequest{
		Op:       srlndk.NotificationRegisterRequest_AddSubscription,
		StreamId: notificationResponse.GetStreamId(),
		SubscriptionTypes: &srlndk.NotificationRegisterRequest_Appid{ // AppId
			Appid: &srlndk.AppIdentSubscriptionRequest{
				Key: &srlndk.AppIdentKey{Id: id},
			},
		},
	}
	return a.startNotificationStream(ctx, notificationRegisterRequest, notificationResponse.GetSubId())
}

// startNotificationStream function
func (a *Agent) startNotificationStream(ctx context.Context, req *srlndk.NotificationRegisterRequest, subID uint64) chan *srlndk.NotificationStreamResponse {
	streamChan := make(chan *srlndk.NotificationStreamResponse)
	log.Printf("starting stream with req=%+v", req)
	go func() {
		defer close(streamChan)
		defer func() {
			log.Printf("agent %s deleting subscription %d", a.Name, subID)
			a.SdkMgrService.Client.NotificationRegister(context.TODO(), &srlndk.NotificationRegisterRequest{
				Op:    srlndk.NotificationRegisterRequest_DeleteSubscription,
				SubId: subID,
			})
		}()
	GETSTREAM:
		registerResponse, err := a.SdkMgrService.Client.NotificationRegister(ctx, req)
		if err != nil {
			log.Printf("agent %s failed registering to notification with req=%+v: %v", a.Name, req, err)
			log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
			time.Sleep(a.RetryTimer)
			goto GETSTREAM
		}
		if registerResponse.GetStatus() == srlndk.SdkMgrStatus_kSdkMgrFailed {
			log.Printf("failed to get stream with req: %v", req)
			log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
			time.Sleep(a.RetryTimer)
			goto GETSTREAM
		}
		stream, err := a.NotificationService.Client.NotificationStream(ctx,
			&srlndk.NotificationStreamRequest{
				StreamId: req.GetStreamId(),
			})
		if err != nil {
			log.Printf("agent %s failed creating stream client with req=%+v: %v", a.Name, req, err)
			log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
			time.Sleep(a.RetryTimer)
			goto GETSTREAM
		}

		for {
			// select {
			// case <-ctx.Done():
			// 	return
			// default:
			ev, err := stream.Recv()
			if err == io.EOF {
				log.Printf("agent %s received EOF for stream %v", a.Name, req.GetSubscriptionTypes())
				log.Printf("agent %s retrying in %s", a.Name, a.RetryTimer)
				time.Sleep(a.RetryTimer)
				goto GETSTREAM
			}
			if err != nil {
				log.Printf("agent %s failed to receive notification: %v", a.Name, err)
				continue
			}
			streamChan <- ev
		}
		//}
	}()
	return streamChan
}
