package tests

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	dynamodb2 "github.com/ricardojonathanromero/go-utilities/db/dynamodb"
	"github.com/stretchr/testify/suite"
	"os"
	"time"
)

const (
	dynamodbImageName string = "amazon/dynamodb-local:latest"
	containerName     string = "dynamodb-test-local"
)

type DBSuite interface {
	StartDynamoDB() error
	CreateTable(input *dynamodb.CreateTableInput) error
	Shutdown()
	GetLocalClient() (*dynamodb.Client, error)
	PutItem(tableName string, item any) error
	DeleteItem(tableName string, key string, value string) error
}

type dbSuiteImpl struct {
	suite.Suite
	url         string
	exposedPort string
}

func New(exposedPort string) DBSuite {
	return &dbSuiteImpl{exposedPort: exposedPort}
}

func (db *dbSuiteImpl) StartDynamoDB() error {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	ctx := context.Background()

	// create docker container
	res, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: dynamodbImageName,
			ExposedPorts: map[nat.Port]struct{}{
				nat.Port("8000"): {},
			},
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{"8000": {{HostIP: "0.0.0.0", HostPort: db.exposedPort}}},
		},
		nil,
		nil,
		containerName)

	if err != nil {
		return err
	}

	// start container
	err = cli.ContainerStart(ctx, res.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	var isRunning bool
	var count int
	maxRetries := 3

	for !isRunning && count <= maxRetries {
		containerInfo, errInspect := cli.ContainerInspect(ctx, res.ID)
		if errInspect != nil {
			return err
		}

		isRunning = containerInfo.State.Running

		if !isRunning {
			time.Sleep(time.Second * 5)
			count++
		}
	}

	_, err = cli.ContainerInspect(ctx, res.ID)
	if err != nil {
		return err
	}

	db.url = fmt.Sprintf("http://localhost:%s", db.exposedPort)

	return nil
}

func (db *dbSuiteImpl) CreateTable(input *dynamodb.CreateTableInput) error {
	conn, err := db.GetLocalClient()
	if err != nil {
		return err
	}

	_, err = conn.CreateTable(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

func (db *dbSuiteImpl) Shutdown() {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	ctx := context.Background()
	if err != nil {
		panic(err)
	}
	if err = cli.ContainerStop(ctx, containerName, container.StopOptions{}); err != nil {
		fmt.Printf("Unable to stop container %s: %s\n", containerName, err)
	}

	removeOptions := container.RemoveOptions{RemoveVolumes: true, Force: true}

	if err = cli.ContainerRemove(ctx, containerName, removeOptions); err != nil {
		fmt.Printf("Unable to remove container: %v\n", err)
	}
}

func (db *dbSuiteImpl) GetLocalClient() (*dynamodb.Client, error) {
	sess := dynamodb2.New()

	err := os.Setenv("DYNAMODB_URL", db.url)
	if err != nil {
		return nil, err
	}

	conn, err := sess.Connect()
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (db *dbSuiteImpl) PutItem(tableName string, item any) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	putInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	conn, err := db.GetLocalClient()
	if err != nil {
		return err
	}

	_, err = conn.PutItem(context.Background(), putInput)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}

func (db *dbSuiteImpl) DeleteItem(tableName string, key string, value string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			key: &types.AttributeValueMemberS{
				Value: value,
			},
		},
		TableName: aws.String(tableName),
	}

	conn, err := db.GetLocalClient()
	if err != nil {
		return err
	}

	_, err = conn.DeleteItem(context.Background(), input)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return nil
}
