package single_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/db"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/allocate"
	"github.com/ricardojonathanromero/lambda-golang-example/internal/utils/tests"
	"testing"
	"time"
)

var _ = BeforeSuite(func() {
	log = logger.NewLoggerWithOptions(logger.Opts{
		AppName: "create-user-lambda-e2e-single-record",
		Level:   "debug",
	})

	port = startPort
	for port <= 10000 {
		isPortFree := allocate.IsPortFree(port)
		if isPortFree {
			break
		}

		log.Debugf("port %d not free", port)
		port++
	}

	log.Debugf("init db connection in port %d", port)
	// init docker cli
	dynamodbTestConn = tests.New(fmt.Sprintf("%d", port))

	// start container
	log.Debugf("starting docker dynamodb container")
	err := dynamodbTestConn.StartDynamoDB()
	Expect(err).To(BeNil())

	// init connection
	log.Debugf("init dynamodb client connection")
	conn, err = dynamodbTestConn.GetLocalClient()
	Expect(err).To(BeNil())
	Expect(conn).NotTo(BeNil())

	log.Debug("checking dynamodb connection is available")
	var isDBActive bool
	var count int
	maxRetries := 3

	for count <= maxRetries {
		_, err := conn.ListTables(context.Background(), &dynamodb.ListTablesInput{})
		if err != nil {
			log.Debug("not ready, sleeping by 5 seconds")
			time.Sleep(time.Second * 5)
			count++
			continue
		} else {
			log.Debug("connection ready")
			isDBActive = true
			break
		}
	}

	Expect(isDBActive).To(BeTrue())

	// configure table
	log.Debugf("configuring table")
	err = db.New(conn, log).ConfigureTable(tableName)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	if dynamodbTestConn != nil {
		// ends docker container
		dynamodbTestConn.Shutdown()
	}
})

func TestSingleRecord(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Single Record Suite")
}
