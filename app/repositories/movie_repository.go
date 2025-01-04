package repositories

import (
	"fmt"
	"log"
	"revel-dynamodb-v2/app/models"

	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)


type DynamoDBRepository struct {
	Client *dynamodb.Client
	TableName string
}

func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		Client: client,
		TableName: tableName,
	}
}

func (a *DynamoDBRepository) GetMovie(id string) (*models.Movie, error) {
	var movie models.Movie
	result, err := a.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(a.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		log.Fatalf("failed to get movie: %v", err)
		return nil, fmt.Errorf("failed to get movie. here's why: %v", err)
	}

	err = attributevalue.UnmarshalMap(result.Item, &movie)
	if err != nil {
		log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		return nil, fmt.Errorf("couldn't unmarshal response. here's why: %v", err)
	}

	return &movie, nil
}