package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type TransformedObject struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Length int    `json:"length"`
	Md5    string `json:"md5"`
}

var svc *s3.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	svc = s3.NewFromConfig(cfg)
}

func hello(ctx context.Context, event S3ObjectLambdaEvent) (interface{}, error) {
	fmt.Printf("[debug] %+v", event)
	c := event.GetObjectContext
	fmt.Printf("[%s] %s - %s", c.InputS3Url, c.OutputRoute, c.OutputToken)

	resp, err := http.Get(c.InputS3Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	fmt.Println("BodyString:")
	fmt.Println(bodyString)

	transformedObject := TransformedObject{
		Metadata: Metadata{
			Length: len(bodyString),
			Md5:    toMd5(bodyString),
		},
	}
	jsonData, err := json.Marshal(transformedObject)
	if err != nil {
		fmt.Println("JSON変換エラー:", err)
		return nil, err
	}

	input := s3.WriteGetObjectResponseInput{
		RequestRoute: &c.OutputRoute,
		RequestToken: &c.OutputToken,
		Body:         strings.NewReader(string(jsonData)),
	}
	return svc.WriteGetObjectResponse(ctx, &input)
}

func toMd5(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
}

// Implementation for aws-lambda-go/events
type S3ObjectLambdaEvent struct {
	XAmzRequestId        string            `json:"xAmzRequestId"`
	GetObjectContext     *GetObjectContext `json:"getObjectContext,omitempty"`
	ListObjectsContext   *ObjectContext    `json:"listObjectsContext,omitempty"`
	ListObjectsV2Context *ObjectContext    `json:"listObjectsV2Context,omitempty"`
	HeadObjectContext    *ObjectContext    `json:"headObjectContext,omitempty"`
	Configuration        Configuration     `json:"configuration"`
	UserRequest          UserRequest       `json:"userRequest"`
	UserIdentity         UserIdentity      `json:"userIdentity"`
	ProtocolVersion      string            `json:"protocolVersion"`
}

type GetObjectContext struct {
	InputS3Url  string `json:"inputS3Url"`
	OutputRoute string `json:"outputRoute"`
	OutputToken string `json:"outputToken"`
}

type ObjectContext struct {
	InputS3Url string `json:"inputS3Url"`
}

type Configuration struct {
	AccessPointArn           string `json:"accessPointArn"`
	SupportingAccessPointArn string `json:"supportingAccessPointArn"`
	payload                  string `json:"payload"`
}

type UserRequest struct {
	url     string            `json:"url"`
	headers map[string]string `json:"headers"`
}

type UserIdentity struct {
	Type           string          `json:"type"`
	PrincipalId    string          `json:"principalId"`
	Arn            string          `json:"arn"`
	AccountId      string          `json:"accountId"`
	AccessKeyId    string          `json:"accessKeyId"`
	SessionContext *SessionContext `json:"sessionContext,omitempty"`
}

type SessionContext struct {
	Attributes    map[string]string `json:"attributes"`
	SessionIssuer *SessionIssuer    `json:"sessionIssuer,omitempty"`
}

type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalId string `json:"principalId"`
	Arn         string `json:"arn"`
	AccountId   string `json:"accountId"`
	UserName    string `json:"userName"`
}
