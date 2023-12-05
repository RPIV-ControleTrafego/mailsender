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

// CheckCPFsExistInDB verifica se o CPF existe na coleção de usuários
// CheckCPFsExistInDB verifica se o CPF existe na coleção de usuários
func CheckCPFsExistInDB(veiculeOwneCPF string) (bool, error) {
	client, err := getClient()
	if err != nil {
		return false, err
	}
	defer client.Disconnect(context.TODO())

	// Selecione a coleção de usuários
	collection := client.Database("users").Collection("user")

	// Crie um filtro para procurar pelo CPF
	filter := bson.M{"cpf": veiculeOwneCPF}

	// Execute a consulta
	var user User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		// CPF não encontrado
		log.Printf("CPF %s does not exist in the collection.", veiculeOwneCPF)
		return false, nil
	} else if err != nil {
		// Um erro ocorreu durante a consulta
		log.Fatal(err)
		return false, err
	}


	// O CPF existe na coleção de usuários
	log.Printf("CPF %s exists in the collection.", veiculeOwneCPF)
	return true, nil
}

func getClient() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to check if the connection is established
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}

type MongoDBProxy struct {
	RealDB *MongoDBClient
}

func NewMongoDBProxy() *MongoDBProxy {
	// Inicialize o MongoDBClient real
	realDB := &MongoDBClient{}

	// Retorne uma instância do proxy com o MongoDBClient real
	return &MongoDBProxy{
		RealDB: realDB,
	}
}


type MongoDBClient struct {
	
	client *mongo.Client
}

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

func (c *MongoDBClient) Close() {
	err := c.client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

// GetEmailByCPF obtém o e-mail associado a um CPF na coleção de usuários// GetEmailByCPF retorna o e-mail associado ao CPF do usuário
func GetEmailByCPF(veiculeOwnerCPF string) (string, error) {
    client, err := getClient()
    if err != nil {
        return "", err
    }
    defer client.Disconnect(context.TODO())

    // Selecione a coleção de usuários
    collection := client.Database("users").Collection("user")

    // Crie um filtro para procurar pelo CPF
    filter := bson.M{"cpf": veiculeOwnerCPF}

    // Execute a consulta
    var user User
    err = collection.FindOne(context.TODO(), filter).Decode(&user)

    if err == mongo.ErrNoDocuments {
        return "", fmt.Errorf("CPF %s not found in the database", veiculeOwnerCPF)
    } else if err != nil {
        // Um erro ocorreu durante a consulta
        log.Fatal(err)
        return "", err
    }

    // O CPF existe na coleção de usuários, retorne o e-mail
    return user.Email, nil
}