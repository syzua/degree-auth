package services

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	sdk       *fabsdk.FabricSDK
	channelCtx context.Channel
)

const (
	channelName   = "mychannel"
	chaincodeName = "education"
	orgName      = "Org1"
	userName     = "User1"
	configPath   = "config/connection-profile.yaml"
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

func GetClient(org, user string) (*channel.Client, error) {
	clientChannelContext := fabsdk.ChannelUserContext{
		OrgID:        org,
		ChannelID:    channelName,
		User:         user,
		CryptoPath:   os.Getenv("CRYPTO_PATH"),
		ConfigPath:   configPath,
	}
	client, err := channel.New(sdk.ChannelContext(clientChannelContext))
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
