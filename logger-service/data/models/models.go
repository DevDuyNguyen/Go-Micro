package models

import (
	"logger-service/internal"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client
var database *mongo.Database
var collection *mongo.Collection

func New(mongoDBClient *mongo.Client, databaseName string, collectionName string) Model{
	client=mongoDBClient
	database=client.Database(databaseName)
	collection=database.Collection(collectionName)
	return Model{
		LogEntry: LogEntry{},
	}
}

type Model struct{
	LogEntry LogEntry
}

type LogEntry struct{
	ID string `bson:"_id,omitempty" json:"id,omitempty"` //this field is usually left empty so the driver can automatically generate unique ObjectId values. 
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (*LogEntry) Insert(entry *LogEntry) error{
	ctx, cancel:= internal.GetTimeOutContext()
	defer cancel()

	entry.CreatedAt=time.Now()
	entry.UpdatedAt=time.Now()

	_,err:=collection.InsertOne(ctx, *entry)
	if err!=nil{
		return err
	}

	return nil
}

func (*LogEntry) All() ([]LogEntry, error){
	queryOpts:=options.Find().SetSort(bson.D{{"created_at", -1}})
	ctx,cancel:=internal.GetTimeOutContext()
	defer cancel()

	cursor, err:= collection.Find(ctx, bson.D{}, queryOpts)
	if err!=nil{
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var res [] LogEntry
	if err:=cursor.All(ctx, &res); err!=nil{
		return nil, err
	}

	return res, nil

}

func (*LogEntry) GetOne(id string) (*LogEntry, error){
	ctx, cancel:= internal.GetTimeOutContext()
	defer cancel()

	_id, err:=bson.ObjectIDFromHex(id)
	if err!=nil{
		return nil, err
	}

	var logEntry LogEntry
	err=collection.FindOne(ctx, bson.D{{"_id", _id}}).Decode(&logEntry)
	if err!=nil{
		return nil, err
	}

	return &logEntry, nil
	
}

func (*LogEntry) DropCollection() error{
	ctx,cancel:= internal.GetTimeOutContext()
	defer cancel()

	err:=collection.Drop(ctx)
	if err!=nil{
		return err
	}
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error){
	ctx,cancel:= internal.GetTimeOutContext()
	defer cancel()

	updateCollection:=bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	}

	_id,err:=bson.ObjectIDFromHex(l.ID)
	if err!=nil{
		return nil, err
	}

	res,err:=collection.UpdateByID(ctx, _id, updateCollection)
	if err!=nil{
		return nil, err
	}
	
	return res, nil
}