package framework

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"sort"
	"strconv"
)

type MusicQuizEntry struct {
	UserId        string
	MusicScore    int
	TotalAttempts int
}

type PrefixEntry struct {
	GuildId string
	Prefix  string
}

type MusicQuizDatabase interface {
	GetScore(string) (int, int)
	UpdateScore(string, int, int)
	CreateScore(string, int, int)
	GetScores() []MusicQuizEntry
}

type PrefixDatabase interface {
	CreateGuild(string, string)
	UpdateGuild(string, string)
	RemoveGuild(string)
	GetPrefix(string) string
	SetGuilds()
	GetGuilds() map[string]string
}

type DynamoDBMusicQuizDatabase struct {
	TableName string
}

type DynamoDBPrefixDatabase struct {
	TableName   string
	PrefixCache map[string]string
}

var AwsSession *session.Session
var DynamoDBInstance *dynamodb.DynamoDB
var MQDB MusicQuizDatabase
var PDB PrefixDatabase

func init() {
	// Initialize session
	AwsSession = session.Must(session.NewSession())

	// Initialize instance
	DynamoDBInstance = dynamodb.New(AwsSession)

	// Setup music quiz database struct
	MQDB = DynamoDBMusicQuizDatabase{TableName: "music_quiz"}

	PDB = &DynamoDBPrefixDatabase{
		TableName:   "prefix_table",
		PrefixCache: make(map[string]string),
	}
}

func (db DynamoDBMusicQuizDatabase) GetScore(id string) (int, int) {
	result, _ := DynamoDBInstance.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(id),
			},
		},
	})

	item := MusicQuizEntry{}
	_ = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if item.UserId == "" {
		return 0, 0
	}

	return item.MusicScore, item.TotalAttempts
}

func (db DynamoDBMusicQuizDatabase) UpdateScore(id string, score, attempts int) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				N: aws.String(strconv.Itoa(score)),
			},
			":t": {
				N: aws.String(strconv.Itoa(attempts)),
			},
		},
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set MusicScore = :s, TotalAttempts = :t"),
	}

	_, _ = DynamoDBInstance.UpdateItem(input)
}

func (db DynamoDBMusicQuizDatabase) CreateScore(id string, score, attempts int) {
	item := MusicQuizEntry{
		UserId:        id,
		MusicScore:    score,
		TotalAttempts: attempts,
	}

	val, _ := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      val,
		TableName: aws.String(db.TableName),
	}

	_, _ = DynamoDBInstance.PutItem(input)
}

func (db DynamoDBMusicQuizDatabase) GetScores() []MusicQuizEntry {
	filt := expression.Name("MusicScore").GreaterThan(expression.Value(0))
	proj := expression.NamesList(expression.Name("UserId"),
		expression.Name("MusicScore"),
		expression.Name("TotalAttempts"))
	expr, _ := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.TableName),
	}
	result, _ := DynamoDBInstance.Scan(params)

	entries := make([]MusicQuizEntry, 0)
	for _, i := range result.Items {
		item := MusicQuizEntry{}
		_ = dynamodbattribute.UnmarshalMap(i, &item)
		entries = append(entries, item)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].MusicScore > entries[j].MusicScore
	})

	return entries
}

func (db *DynamoDBPrefixDatabase) CreateGuild(id, prefix string) {
	item := PrefixEntry{
		GuildId: id,
		Prefix:  prefix,
	}

	val, _ := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      val,
		TableName: aws.String(db.TableName),
	}

	_, _ = DynamoDBInstance.PutItem(input)
	db.PrefixCache[id] = prefix
}

func (db *DynamoDBPrefixDatabase) UpdateGuild(id, prefix string) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(prefix),
			},
		},
		TableName: aws.String(db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"GuildId": {
				S: aws.String(id),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Prefix = :s"),
	}

	_, _ = DynamoDBInstance.UpdateItem(input)
	db.PrefixCache[id] = prefix
}

func (db *DynamoDBPrefixDatabase) RemoveGuild(id string) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"GuildId": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.TableName),
	}

	_, _ = DynamoDBInstance.DeleteItem(input)
	delete(db.PrefixCache, id)
}

func (db *DynamoDBPrefixDatabase) GetPrefix(id string) string {
	return db.PrefixCache[id]
}

func (db *DynamoDBPrefixDatabase) SetGuilds() {
	proj := expression.NamesList(expression.Name("GuildId"),
		expression.Name("Prefix"))
	expr, _ := expression.NewBuilder().WithProjection(proj).Build()
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.TableName),
	}
	result, _ := DynamoDBInstance.Scan(params)

	for _, i := range result.Items {
		item := PrefixEntry{}
		_ = dynamodbattribute.UnmarshalMap(i, &item)
		db.PrefixCache[item.GuildId] = item.Prefix
	}
}

func (db *DynamoDBPrefixDatabase) GetGuilds() map[string]string {
	return db.PrefixCache
}
