package single_test

import (
	"context"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/internal/handler"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/entities"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/repository"
	"github.com/ricardojonathanromero/lambda-golang-example/create-user-lambda/pkg/service"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/models"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/tests"
	"net/http"
	"time"
)

const (
	startPort = 8001
	tableName = "users"
)

var (
	port             int
	dynamodbTestConn tests.DBSuite
	conn             *dynamodb.Client
	log              logger.Logger
)

var _ = Describe("Single Record", func() {
	var hdl handler.Handle
	var lambdaCtx *lambdacontext.LambdaContext
	var ctx context.Context

	BeforeEach(func() {
		repo := repository.New(tableName, conn, log)
		srv := service.New(repo, log)
		hdl = handler.New(srv, log)

		lambdaCtx = &lambdacontext.LambdaContext{
			AwsRequestID:       "awsRequestId1234",
			InvokedFunctionArn: "arn:aws:lambda:xxx",
			Identity:           lambdacontext.CognitoIdentity{},
			ClientContext:      lambdacontext.ClientContext{},
		}
	})

	// init connection
	Describe("save item in db", func() {

		When("deadline is up to 10 secs", func() {
			var cancel context.CancelFunc

			BeforeEach(func() {
				// configure
				ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
				ctx = lambdacontext.NewContext(ctx, lambdaCtx)
			})

			It("save record successfully", func() {
				defer cancel()

				req := entities.UserReq{
					Name:     "john",
					Lastname: "smith",
					Age:      23,
					Email:    "john.smith@test.com",
				}

				log.Debug("start send handler request")
				res, errRes := hdl.HandleCreateUser(ctx, req)
				Expect(errRes).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.StatusCode).To(Equal(http.StatusCreated))

				log.Debug("check record exists in db")
				out, errScan := conn.Scan(context.Background(), &dynamodb.ScanInput{TableName: aws.String(tableName)})
				Expect(errScan).To(BeNil())
				Expect(out).NotTo(BeNil())
				Expect(out.Count).To(Equal(int32(1)))
				Expect(out.Items).NotTo(BeNil())

				var items []*models.UserDB
				log.Debug("unmarshal response")
				errUnmarshal := attributevalue.UnmarshalListOfMaps(out.Items, &items)
				Expect(errUnmarshal).To(BeNil())
				Expect(items).To(HaveLen(1))
				Expect(items[0].Name).To(Equal("john"))
				log.Debug("item exists as expected")
			})
		})
	})
})
