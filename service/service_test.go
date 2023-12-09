package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestValidateCPF(t *testing.T) {
    // Crie um objeto MessageContent para usar nos testes
    message := MessageContent{
        VehicleOwnerCPF: "206.376.746-13",
    }

    // Chame a função que você deseja testar
    result := ValidateCPF(message)

    // Avalie o resultado
    if !result {
        t.Errorf("A validação do CPF falhou para o CPF %s", message.VehicleOwnerCPF)
    }

	assert.Equal(t, true, result, "The result should be true")
}

func TestValidateFalseCPF(t *testing.T) {
	// Crie um objeto MessageContent para usar nos testes
	message := MessageContent{
		VehicleOwnerCPF: "206.376.746-1",
	}

	// Chame a função que você deseja testar
	result := ValidateCPF(message)

	// Avalie o resultado
	if result {
		t.Errorf("A validação do CPF falhou para o CPF %s", message.VehicleOwnerCPF)
	}

	assert.Equal(t, false, result, "The result should be false")


}

func TestEmptyCPF(t *testing.T) {
	// Crie um objeto MessageContent para usar nos testes
	message := MessageContent{
		VehicleOwnerCPF: "",
	}

	// Chame a função que você deseja testar
	result := ValidateCPF(message)

	// Avalie o resultado
	if result {
		t.Errorf("A validação do CPF falhou para o CPF %s", message.VehicleOwnerCPF)
	}



	assert.Equal(t, false, result, "The result should be false")

}



func TestEmptyEmail(t *testing.T) {
	// Crie um objeto MessageContent para usar nos testes
	message := MessageContent{
		VehicleOwnerCPF: "",
	}

	// Chame a função que você deseja testar
	result := GetEmail(message)

	// Avalie o resultado
	if result != "" {
		t.Errorf("A validação do CPF falhou para o CPF %s", message.VehicleOwnerCPF)
	}

	assert.Equal(t, "", result, "The result should be false")

}





func TestGetEmail(t *testing.T) {
	// Crie um objeto MessageContent para usar nos testes
	message := MessageContent{
		VehicleOwnerCPF: "206.376.746-13",
	}

	// Chame a função que você deseja testar
	result := GetEmail(message)

	// Avalie o resultado
	if result == "" {
		t.Errorf("A validação do CPF falhou para o CPF %s", message.VehicleOwnerCPF)
	}

	assert.Equal(t, "andremiranda.aluno@unipampa.edu.br", result, "The result should be true")
}


func TestGetEmailFalse(t *testing.T) {
	// Crie um objeto MessageContent para usar nos testes
	message := MessageContent{
		VehicleOwnerCPF: "206.376.746-1",
	}

	// Chame a função que você deseja testar
	result := GetEmail(message)


	assert.Equal(t, "", result, "The result should be false")
}



