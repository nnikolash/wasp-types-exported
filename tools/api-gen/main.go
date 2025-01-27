package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"

	"github.com/iotaledger/hive.go/logger"
	"github.com/nnikolash/wasp-types-exported/components/app"
	"github.com/nnikolash/wasp-types-exported/components/webapi"
	"github.com/nnikolash/wasp-types-exported/packages/authentication"
	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/evm/jsonrpc"
	v2 "github.com/nnikolash/wasp-types-exported/packages/webapi"
)

type NodeIdentityProviderMock struct{}

func (n *NodeIdentityProviderMock) NodeIdentity() *cryptolib.KeyPair {
	return cryptolib.NewKeyPair()
}

func (n *NodeIdentityProviderMock) NodePublicKey() *cryptolib.PublicKey {
	return cryptolib.NewEmptyPublicKey()
}

func main() {
	mockLog := logger.NewNopLogger()
	e := echo.New()

	if app.Version == "" {
		app.Version = "0"
	}

	swagger := webapi.CreateEchoSwagger(e, app.Version)
	v2.Init(
		mockLog,
		swagger,
		app.Version,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		&NodeIdentityProviderMock{},
		nil,
		nil,
		nil,
		nil,
		authentication.AuthConfiguration{Scheme: authentication.AuthJWT},
		time.Second,
		nil,
		"",
		"",
		nil,
		jsonrpc.ParametersDefault(),
	)

	root, ok := swagger.(*echoswagger.Root)
	if !ok {
		panic("failed to get swagger root")
	}

	schema, err := root.GetSpec(nil, "/doc")
	if err != nil {
		panic(err.Error())
	}

	jsonSchema, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	fmt.Print(string(jsonSchema))
}
