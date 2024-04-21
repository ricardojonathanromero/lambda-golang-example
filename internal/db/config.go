package db

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/ricardojonathanromero/go-utilities/logger"
	"time"
)

type DB interface {
	ConfigureTable(tableName string) error
}

type dbInfra struct {
	conn *dynamodb.Client
	log  logger.Logger
}

func New(conn *dynamodb.Client, log logger.Logger) DB {
	return &dbInfra{conn: conn, log: log}
}

func (check *dbInfra) ConfigureTable(tableName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	check.log.Debug("describe table")
	out, err := check.conn.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
	if err != nil {
		var nfErr *types.ResourceNotFoundException
		if errors.As(err, &nfErr) {
			check.log.Debug("resource not exists, create")
			return check.configureTable(tableName)
		}

		var ae smithy.APIError
		if errors.As(err, &ae) {
			check.log.Errorf("code: %s, message: %s, fault: %s", ae.ErrorCode(), ae.ErrorMessage(), ae.ErrorFault().String())
			return err
		}

		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			check.log.Errorf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Service(), oe.Unwrap())
			return err
		}

		check.log.Errorf("error describing table: %v", err)
		return err
	}

	check.log.Debugf("check if table name is expected: %v", *out.Table)
	if *out.Table.TableName == tableName {
		check.log.Info("table already configured")
		return nil
	}

	return nil
}

func (check *dbInfra) configureTable(tableName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// table not created
	check.log.Debug("creating table")
	input := getTableDefinition(tableName)
	_, err := check.conn.CreateTable(ctx, input)
	if err != nil {
		check.log.Errorf("error creating table: %v", err)
		return err
	}

	check.log.Info("table configured")
	return nil
}
