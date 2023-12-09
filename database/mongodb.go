package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoURI = "mongodb://localhost:27017/users"
)

// User representa a estrutura do documento na coleção de usuários
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
	Role     string `bson:"role"`
	CPF      string `bson:"cpf"`
	Class    string `bson:"_class"`
}

// DatabaseOperations interface defines the common operations for the database
type DatabaseOperations interface {
	CheckCPFsExistInDB(veiculeOwneCPF string) (bool, error)
	GetEmailByCPF(veiculeOwnerCPF string) (string, error)
}

// MongoDBProxy struct acts as a proxy for the real MongoDBClient
type MongoDBProxy struct {
	RealDB *MongoDBClient
}

// NewMongoDBProxy creates a new instance of MongoDBProxy
func NewMongoDBProxy() *MongoDBProxy {
	// The real database client is not initialized here
	return &MongoDBProxy{}
}

// CheckCPFsExistInDB implements the CheckCPFsExistInDB method of DatabaseOperations interface
func (p *MongoDBProxy) CheckCPFsExistInDB(veiculeOwneCPF string) (bool, error) {
	// Lazy initialize the real database client if not done yet
	if p.RealDB == nil {
		p.RealDB = &MongoDBClient{}
		if err := p.RealDB.Connect(); err != nil {
			return false, err
		}
	}

	// Delegate the operation to the real database client
	return p.RealDB.CheckCPFsExistInDB(veiculeOwneCPF)
}

// GetEmailByCPF implements the GetEmailByCPF method of DatabaseOperations interface
func (p *MongoDBProxy) GetEmailByCPF(veiculeOwnerCPF string) (string, error) {
	// Lazy initialize the real database client if not done yet
	if p.RealDB == nil {
		p.RealDB = &MongoDBClient{}
		if err := p.RealDB.Connect(); err != nil {
			return "", err
		}
	}

	// Delegate the operation to the real database client
	return p.RealDB.GetEmailByCPF(veiculeOwnerCPF)
}

// MongoDBClient struct implements the DatabaseOperations interface
type MongoDBClient struct {
	client *mongo.Client
}

// Connect method initializes the MongoDB client
func (c *MongoDBClient) Connect() error {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c.client = client
	fmt.Println("Connected to MongoDB!")

	return nil
}

// CheckCPFsExistInDB implements the CheckCPFsExistInDB method of DatabaseOperations interface
func (c *MongoDBClient) CheckCPFsExistInDB(veiculeOwneCPF string) (bool, error) {
	collection := c.client.Database("users").Collection("user")
	filter := bson.M{"cpf": veiculeOwneCPF}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		log.Printf("CPF %s does not exist in the collection.", veiculeOwneCPF)
		return false, nil
	} else if err != nil {
		log.Fatal(err)
		return false, err
	}

	log.Printf("CPF %s exists in the collection.", veiculeOwneCPF)
	return true, nil
}

// GetEmailByCPF implements the GetEmailByCPF method of DatabaseOperations interface
func (c *MongoDBClient) GetEmailByCPF(veiculeOwnerCPF string) (string, error) {
	collection := c.client.Database("users").Collection("user")
	filter := bson.M{"cpf": veiculeOwnerCPF}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return "", fmt.Errorf("CPF %s not found in the database", veiculeOwnerCPF)
	} else if err != nil {
		log.Fatal(err)
		return "", err
	}

	return user.Email, nil
}
