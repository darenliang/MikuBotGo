package framework

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/darenliang/MikuBotGo/config"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type MusicQuizEntry struct {
	UserId        string
	MusicScore    int
	TotalAttempts int
}

type MusicQuizEntryTuple struct {
	MusicScore    int
	TotalAttempts int
}

type PrefixEntry struct {
	GuildId string
	Prefix  string
}

type GifItemList struct {
	ID      string
	Data    []GifImage `json:"data"`
	Success bool       `json:"success"`
}

type GifImage struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type GifItem struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
	Success bool `json:"success"`
}

type GifUpload struct {
	Data    GifImage `json:"data"`
	Success bool     `json:"success"`
}

type MusicQuizDatabase interface {
	GetScore(string) (int, int)
	UpdateScore(string, int, int)
	CreateScore(string, int, int)
	SetScores()
	GetScores() map[string]MusicQuizEntryTuple
}

type PrefixDatabase interface {
	CreateGuild(string, string)
	UpdateGuild(string, string)
	RemoveGuild(string)
	GetPrefix(string) string
	SetGuilds()
	GetGuilds() map[string]string
}

type GifDatabase interface {
	GetGif(string) (string, string)
	UploadGif(string, string, string, string) error
	SetAlbums()
	CheckDup(string, string) bool
}

type DynamoDBMusicQuizDatabase struct {
	TableName      string
	MusicQuizCache *sync.Map
}

type DynamoDBPrefixDatabase struct {
	TableName   string
	PrefixCache *sync.Map
}

type GifCacheDatabase struct {
	GifCache *sync.Map
}

var (
	AwsSession       *session.Session
	DynamoDBInstance *dynamodb.DynamoDB
	MQDB             MusicQuizDatabase
	PDB              PrefixDatabase
	GBD              GifDatabase
)

func init() {
	// Initialize session
	AwsSession = session.Must(session.NewSession())

	// Initialize instance
	DynamoDBInstance = dynamodb.New(AwsSession)

	// Setup music quiz database struct
	MQDB = &DynamoDBMusicQuizDatabase{
		TableName:      "music_quiz",
		MusicQuizCache: &sync.Map{},
	}

	// Setup prefix database struct
	PDB = &DynamoDBPrefixDatabase{
		TableName:   "prefix_table",
		PrefixCache: &sync.Map{},
	}

	// Setup gif album cache struct
	GBD = &GifCacheDatabase{
		GifCache: &sync.Map{},
	}
}

func (db *DynamoDBMusicQuizDatabase) GetScore(id string) (int, int) {
	res, ok := db.MusicQuizCache.Load(id)
	if ok {
		entry := res.(MusicQuizEntryTuple)
		return entry.MusicScore, entry.TotalAttempts
	}
	return 0, 0
}

func (db *DynamoDBMusicQuizDatabase) UpdateScore(id string, score, attempts int) {
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
	db.MusicQuizCache.Store(id, MusicQuizEntryTuple{
		MusicScore:    score,
		TotalAttempts: attempts,
	})
}

func (db *DynamoDBMusicQuizDatabase) CreateScore(id string, score, attempts int) {
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
	db.MusicQuizCache.Store(id, MusicQuizEntryTuple{
		MusicScore:    score,
		TotalAttempts: attempts,
	})
}

func (db *DynamoDBMusicQuizDatabase) SetScores() {
	proj := expression.NamesList(expression.Name("UserId"),
		expression.Name("MusicScore"),
		expression.Name("TotalAttempts"))
	expr, _ := expression.NewBuilder().WithProjection(proj).Build()
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.TableName),
	}
	result, _ := DynamoDBInstance.Scan(params)

	for _, i := range result.Items {
		item := MusicQuizEntry{}
		_ = dynamodbattribute.UnmarshalMap(i, &item)
		db.MusicQuizCache.Store(item.UserId, MusicQuizEntryTuple{
			MusicScore:    item.MusicScore,
			TotalAttempts: item.TotalAttempts,
		})
	}
}

func (db *DynamoDBMusicQuizDatabase) GetScores() map[string]MusicQuizEntryTuple {
	res := make(map[string]MusicQuizEntryTuple)
	db.MusicQuizCache.Range(func(k, v interface{}) bool {
		res[k.(string)] = v.(MusicQuizEntryTuple)
		return true
	})
	return res
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
	db.PrefixCache.Store(id, prefix)
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
	db.PrefixCache.Store(id, prefix)
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
	db.PrefixCache.Delete(id)
}

func (db *DynamoDBPrefixDatabase) GetPrefix(id string) string {
	res, ok := db.PrefixCache.Load(id)
	if ok {
		return res.(string)
	}
	return config.Prefix
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
		db.PrefixCache.Store(item.GuildId, item.Prefix)
	}
}

func (db *DynamoDBPrefixDatabase) GetGuilds() map[string]string {
	res := make(map[string]string)
	db.PrefixCache.Range(func(k, v interface{}) bool {
		res[k.(string)] = v.(string)
		return true
	})
	return res
}

func (db *GifCacheDatabase) GetGif(guildId string) (string, string) {
	res, ok := db.GifCache.Load(guildId)

	if !ok {
		return "", ""
	}

	images := res.(GifItemList)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/album/%s/images",
		config.ImgurEndpoint, images.ID), new(bytes.Buffer))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ImgurToken)

	resp, _ := HttpClient.Do(req)
	AlbumEntry := GifItemList{}

	err := json.NewDecoder(resp.Body).Decode(&AlbumEntry)

	if err != nil {
		return "", ""
	}

	if !AlbumEntry.Success {
		return "", ""
	}

	if len(AlbumEntry.Data) == 0 {
		return "", ""
	}

	img := AlbumEntry.Data[rand.Intn(len(AlbumEntry.Data))]

	return img.Title, img.Link
}

func (db *GifCacheDatabase) UploadGif(guildId, userId, imgUrl, hash string) error {
	res, ok := db.GifCache.Load(guildId)

	images := GifItemList{}

	if !ok {
		params := url.Values{}
		params.Set("title", guildId)
		params.Set("privacy", "secret")
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/album?%s",
			config.ImgurEndpoint, params.Encode()), new(bytes.Buffer))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.ImgurToken)
		resp, _ := HttpClient.Do(req)

		albumCreation := GifItem{}
		err := json.NewDecoder(resp.Body).Decode(&albumCreation)

		if err != nil {
			return err
		}

		if !albumCreation.Success {
			return errors.New("database error")
		}

		images = GifItemList{
			ID:      albumCreation.Data.ID,
			Data:    make([]GifImage, 0),
			Success: albumCreation.Success,
		}
	} else {
		images = res.(GifItemList)
	}

	params := url.Values{}
	params.Set("image", imgUrl)
	params.Set("album", images.ID)
	params.Set("type", "url")
	params.Set("title", userId)
	params.Set("description", hash)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/upload?%s",
		config.ImgurEndpoint, params.Encode()), new(bytes.Buffer))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ImgurToken)

	resp, _ := HttpClient.Do(req)
	status := GifUpload{}

	err := json.NewDecoder(resp.Body).Decode(&status)

	if err != nil {
		return err
	}

	if !status.Success {
		fmt.Println(status)
		return errors.New("gif cannot be added")
	}

	images.Data = append(images.Data, GifImage{
		ID:          status.Data.ID,
		Title:       status.Data.Title,
		Description: status.Data.Description,
		Link:        status.Data.Link,
	})

	db.GifCache.Store(guildId, images)

	return nil
}

func (db *GifCacheDatabase) SetAlbums() {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/account/%s/albums",
		config.ImgurEndpoint, config.ImgurUsername), new(bytes.Buffer))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.ImgurToken)
	resp, _ := HttpClient.Do(req)
	albums := GifItemList{}

	_ = json.NewDecoder(resp.Body).Decode(&albums)

	for _, i := range albums.Data {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/album/%s/images",
			config.ImgurEndpoint, i.ID), new(bytes.Buffer))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.ImgurToken)
		resp, _ := HttpClient.Do(req)
		images := GifItemList{}
		images.ID = i.ID
		_ = json.NewDecoder(resp.Body).Decode(&images)
		db.GifCache.Store(i.Title, images)
	}
}

func (db *GifCacheDatabase) CheckDup(guildId, hash string) bool {
	res, ok := db.GifCache.Load(guildId)

	if !ok {
		return false
	}

	images := res.(GifItemList)

	for _, i := range images.Data {
		if i.Description == hash {
			return true
		}
	}

	return false
}
