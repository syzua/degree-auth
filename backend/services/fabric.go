package services

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var sdk *fabsdk.FabricSDK

const (
	channelName   = "mychannel"
	chaincodeName = "education"
	configPath    = "config/connection-profile.yaml"
)

type FabricService struct {
	client *channel.Client
}

func InitSDK() error {
	var err error
	configBackend := config.FromFile(configPath)
	sdk, err = fabsdk.New(configBackend)
	if err != nil {
		return fmt.Errorf("failed to create SDK: %v", err)
	}
	return nil
}

// GetClient 用「组织 + 用户」构造一个 channel 客户端。
// fabric-sdk-go 的用法是先 sdk.ChannelContext 拿到上下文，再交给 channel.New。
func GetClient(org, user string) (*channel.Client, error) {
	if sdk == nil {
		return nil, fmt.Errorf("SDK 未初始化，请先调用 InitSDK")
	}
	ctx := sdk.ChannelContext(channelName,
		fabsdk.WithOrg(org),
		fabsdk.WithUser(user),
	)
	client, err := channel.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel client: %v", err)
	}
	return client, nil
}

func InvokeChaincode(function string, args [][]byte, org, user string) ([]byte, error) {
	client, err := GetClient(org, user)
	if err != nil {
		return nil, err
	}

	req := channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         function,
		Args:        args,
	}

	resp, err := client.Execute(req)
	if err != nil {
		return nil, fmt.Errorf("invoke failed: %v", err)
	}

	return resp.Payload, nil
}

func QueryChaincode(function string, args [][]byte, org, user string) ([]byte, error) {
	client, err := GetClient(org, user)
	if err != nil {
		return nil, err
	}

	req := channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         function,
		Args:        args,
	}

	resp, err := client.Query(req)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	return resp.Payload, nil
}

func CloseSDK() {
	if sdk != nil {
		sdk.Close()
	}
}
