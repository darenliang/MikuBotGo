package framework

// This code is really messy, major restructuring is needed

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"sort"
	"strconv"
)

type DatabaseEntry struct {
	UserId        string
	MusicScore    int
	TotalAttempts int
}

var AwsSession *session.Session
var DynamoDBInstance *dynamodb.DynamoDB

func init() {
	// Initialize session
	AwsSession = session.Must(session.NewSession())

	// Initialize instance
	DynamoDBInstance = dynamodb.New(AwsSession)
}

func GetDatabaseValue(id string) (int, int) {
	result, _ := DynamoDBInstance.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("music_quiz"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(id),
			},
		},
	})

	item := DatabaseEntry{}
	_ = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if item.UserId == "" {
		return 0, 0
	}

	return item.MusicScore, item.TotalAttempts
}

func UpdateDatabaseValue(id string, score int, attempts int) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				N: aws.String(strconv.Itoa(score)),
			},
			":t": {
				N: aws.String(strconv.Itoa(attempts)),
			},
		},
		TableName: aws.String("music_quiz"),
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set MusicScore = :s, TotalAttempts = :t"),
	}

	_, _ = DynamoDBInstance.UpdateItem(input)
	return
}

func CreateDatabaseEntry(id string, score int, attempts int) {
	item := DatabaseEntry{
		UserId:        id,
		MusicScore:    score,
		TotalAttempts: attempts,
	}

	val, _ := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      val,
		TableName: aws.String("music_quiz"),
	}

	_, err := DynamoDBInstance.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetHighscores() []DatabaseEntry {
	filt := expression.Name("MusicScore").GreaterThanEqual(expression.Value(0))
	proj := expression.NamesList(expression.Name("UserId"),
		expression.Name("MusicScore"),
		expression.Name("TotalAttempts"))
	expr, _ := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("music_quiz"),
	}
	result, _ := DynamoDBInstance.Scan(params)

	entries := make([]DatabaseEntry, 0)
	for _, i := range result.Items {
		item := DatabaseEntry{}
		_ = dynamodbattribute.UnmarshalMap(i, &item)
		entries = append(entries, item)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].MusicScore > entries[j].MusicScore
	})

	return entries
}
