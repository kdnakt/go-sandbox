package main

import (
	"fmt"
	"log"
	"context"
	"net/http"
	"io/ioutil"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var svc *s3.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	svc = s3.NewFromConfig(cfg)
}

func hello(ctx context.Context, event S3ObjectLambdaEvent) (string, error) {
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
	if len(bodyString) > 20 {
		fmt.Println("Object: ", bodyString[:20])
	} else {
		fmt.Println("Object: ", bodyString)
	}

	return "Hello λ!", nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
}

// Implementation for aws-lambda-go/events
type S3ObjectLambdaEvent struct {
	XAmzRequestId		string			`json:"xAmzrequestId"`
	GetObjectContext	GetObjectContext	`json:"getObjectContext"`
	ListObjectsContext	ObjectContext		`json:"listObjectsContext"`
	ListObjectsV2Context	ObjectContext		`json:"listObjectsV2Context"`
	HeadObjectContext	ObjectContext		`json:"headObjectContext"`
	Configuration		Configuration		`json:"configuration"`
	UserRequest		UserRequest		`json:"userRequest"`
	UserIdentity		UserIdentity		`json:"userIdentity"`
	ProtocolVersion		string			`json:"protocolVersion"`
}

type GetObjectContext struct {
	InputS3Url	string	`json:"inputS3Url"`
	OutputRoute	string	`json:"outputRoute"`
	OutputToken	string	`json:"outputToken"`
}

type ObjectContext struct {
	InputS3Url	string	`json:"inputS3Url"`
}

type Configuration struct {
	AccessPointArn			string	`json:"accessPointArn"`
	SupportingAccessPointArn	string	`json:"supportingAccessPointArn"`
	payload				string	`json:"payload"`
}

type UserRequest struct {
	url	string			`json:"url"`
	headers	map[string]string	`json:"headers"`
}

type UserIdentity struct {
	Type		string		`json:"type"`
	PrincipalId	string		`json:"principalId"`
	Arn		string		`json:"arn"`
	AccountId	string		`json:"accountId"`
	AccessKeyId	string		`json:"accessKeyId"`
	SessionContext	SessionContext	`json:"sessionContext"`
}

type SessionContext struct {
	Attributes	map[string]string	`json:"attributes"`
	SessionIssuer	SessionIssuer		`json:"sessionIssuer"`
}

type SessionIssuer struct {
	Type		string	`json:"type"`
	PrincipalId	string	`json:"principalId"`
	Arn		string	`json:"arn"`
	AccountId	string	`json:"accountId"`
	UserName	string	`json:"userName"`
}

