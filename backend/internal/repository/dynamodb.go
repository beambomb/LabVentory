package repository

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"labventroy/backend/internal/apperror"
	"labventroy/backend/internal/model"
)

// DynamoDBRepository implements ComponentRepository backed by Amazon DynamoDB.
type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBRepository creates a DynamoDB-backed ComponentRepository.
func NewDynamoDBRepository(client *dynamodb.Client, tableName string) *DynamoDBRepository {
	return &DynamoDBRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBRepository) Create(ctx context.Context, item model.Component) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return apperror.Internal("failed to marshal component", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(item_id)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return apperror.Conflict("component with item_id " + item.ItemID + " already exists")
		}
		return apperror.Internal("failed to create component in dynamodb", err)
	}
	return nil
}

func (r *DynamoDBRepository) List(ctx context.Context) ([]model.Component, error) {
	var items []model.Component
	var exclusiveStartKey map[string]types.AttributeValue

	for {
		out, err := r.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         aws.String(r.tableName),
			ExclusiveStartKey: exclusiveStartKey,
		})
		if err != nil {
			return nil, apperror.Internal("failed to scan components from dynamodb", err)
		}

		var page []model.Component
		if err := attributevalue.UnmarshalListOfMaps(out.Items, &page); err != nil {
			return nil, apperror.Internal("failed to unmarshal components", err)
		}
		items = append(items, page...)

		if len(out.LastEvaluatedKey) == 0 {
			break
		}
		exclusiveStartKey = out.LastEvaluatedKey
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ItemID < items[j].ItemID
	})
	return items, nil
}

func (r *DynamoDBRepository) GetByID(ctx context.Context, id string) (*model.Component, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"item_id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, apperror.Internal("failed to get component from dynamodb", err)
	}
	if len(out.Item) == 0 {
		return nil, apperror.NotFound("component " + id + " not found")
	}

	var item model.Component
	if err := attributevalue.UnmarshalMap(out.Item, &item); err != nil {
		return nil, apperror.Internal("failed to unmarshal component", err)
	}
	return &item, nil
}

func (r *DynamoDBRepository) Update(ctx context.Context, id string, input model.UpdateComponentInput, updatedAt string) (*model.Component, error) {
	setClauses := []string{"#updated_at = :updated_at"}
	exprNames := map[string]string{
		"#updated_at": "updated_at",
	}
	exprValues := map[string]types.AttributeValue{
		":updated_at": &types.AttributeValueMemberS{Value: updatedAt},
	}

	if input.Name != nil {
		setClauses = append(setClauses, "#name = :name")
		exprNames["#name"] = "name"
		exprValues[":name"] = &types.AttributeValueMemberS{Value: *input.Name}
	}
	if input.Category != nil {
		setClauses = append(setClauses, "#category = :category")
		exprNames["#category"] = "category"
		exprValues[":category"] = &types.AttributeValueMemberS{Value: string(*input.Category)}
	}
	if input.StockQty != nil {
		setClauses = append(setClauses, "#stock_qty = :stock_qty")
		exprNames["#stock_qty"] = "stock_qty"
		v, _ := attributevalue.Marshal(*input.StockQty)
		exprValues[":stock_qty"] = v
	}
	if input.MinStock != nil {
		setClauses = append(setClauses, "#min_stock = :min_stock")
		exprNames["#min_stock"] = "min_stock"
		v, _ := attributevalue.Marshal(*input.MinStock)
		exprValues[":min_stock"] = v
	}
	if input.DrawerLocation != nil {
		setClauses = append(setClauses, "#drawer_location = :drawer_location")
		exprNames["#drawer_location"] = "drawer_location"
		exprValues[":drawer_location"] = &types.AttributeValueMemberS{Value: *input.DrawerLocation}
	}
	if input.Condition != nil {
		setClauses = append(setClauses, "#condition = :condition")
		exprNames["#condition"] = "condition"
		exprValues[":condition"] = &types.AttributeValueMemberS{Value: string(*input.Condition)}
	}

	out, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"item_id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          aws.String("SET " + strings.Join(setClauses, ", ")),
		ConditionExpression:       aws.String("attribute_exists(item_id)"),
		ExpressionAttributeNames:  exprNames,
		ExpressionAttributeValues: exprValues,
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return nil, apperror.NotFound("component " + id + " not found")
		}
		return nil, apperror.Internal("failed to update component in dynamodb", err)
	}

	var updated model.Component
	if err := attributevalue.UnmarshalMap(out.Attributes, &updated); err != nil {
		return nil, apperror.Internal("failed to unmarshal updated component", err)
	}
	return &updated, nil
}

func (r *DynamoDBRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"item_id": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(item_id)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return apperror.NotFound("component " + id + " not found")
		}
		return apperror.Internal("failed to delete component from dynamodb", err)
	}
	return nil
}
