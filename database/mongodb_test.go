package database

import (
	"context"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)


//dados mockados
func createMockUser() *User {
	return &User{
		Username: fake.UserName(),
		Password: fake.Password(8, 14, true, true, true),
		Email:    fake.EmailAddress(),
		CPF:      fake.DigitsN(11),

	}
}

func TestMongoDBClient_CheckCPFsExistInDB(t *testing.T) {

	client := &MongoDBClient{}


	err := client.Connect()
	assert.NoError(t, err, "Erro ao conectar ao banco de dados")

	// Crie um usuário mock
	mockUser := createMockUser()


	collection := client.client.Database("users").Collection("user")
	_, err = collection.InsertOne(context.TODO(), mockUser)
	assert.NoError(t, err, "Erro ao inserir usuário mock no banco de dados")


	exists, err := client.CheckCPFsExistInDB(mockUser.CPF)
	assert.NoError(t, err, "Erro ao verificar CPF existente")
	assert.True(t, exists, "O CPF deveria existir no banco de dados")
}

func TestMongoDBClient_GetEmailByCPF(t *testing.T) {

	client := &MongoDBClient{}


	err := client.Connect()
	assert.NoError(t, err, "Erro ao conectar ao banco de dados")


	mockUser := createMockUser()


	collection := client.client.Database("users").Collection("user")
	_, err = collection.InsertOne(context.TODO(), mockUser)
	assert.NoError(t, err, "Erro ao inserir usuário mock no banco de dados")


	email, err := client.GetEmailByCPF(mockUser.CPF)
	assert.NoError(t, err, "Erro ao obter email por CPF existente")
	assert.NotEmpty(t, email, "O email não deveria estar vazio")
}


func TestMongoDBClient_CheckCPFsExistInDB_NonExistentCPF(t *testing.T) {
    client := &MongoDBClient{}
    err := client.Connect()
    assert.NoError(t, err, "Erro ao conectar ao banco de dados")



    exists, err := client.CheckCPFsExistInDB("CPFInexistente")
    assert.NoError(t, err, "Erro ao verificar CPF inexistente")
    assert.False(t, exists, "O CPF não deveria existir no banco de dados")
}

func TestMongoDBClient_GetEmailByCPF_NonExistentCPF(t *testing.T) {
    client := &MongoDBClient{}
    err := client.Connect()
    assert.NoError(t, err, "Erro ao conectar ao banco de dados")



    email, err := client.GetEmailByCPF("CPFInexistente")
    assert.Error(t, err, "Erro esperado ao obter email por CPF inexistente")
    assert.Empty(t, email, "O email deveria estar vazio")
}

