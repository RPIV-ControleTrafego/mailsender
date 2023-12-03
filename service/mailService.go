package service

import "fmt"

type MessageContent struct {
	CarPlate          string  `json:"carPlate"`
	Address           string  `json:"address"`
	Date              string  `json:"date"`
	Violation         string  `json:"violation"`
	CarType           string  `json:"carType"`
	CarColor          string  `json:"carColor"`
	CarBrand          string  `json:"carBrand"`
	VehicleOwnerName  string  `json:"vehicleOwnerName"`
	VehicleOwnerCPF   string  `json:"veiculeOwneCPF"`
	Speed             float64 `json:"speed"`
	MaxSpeed          int     `json:"maxSpeed"`
	FinePrice         float64 `json:"finePrice"`
	Sex               string  `json:"sex"`
	Age               int  `json:"age"`
}

func ShowInfraction(message MessageContent) {
	fmt.Println("Infraction details:")
	fmt.Printf("   Car Plate: %s\n", message.CarPlate)
	fmt.Printf("   Violation: %s\n", message.Violation)
	fmt.Printf("   Owner Name: %s\n", message.VehicleOwnerName)
	fmt.Printf("   Owner CPF: %s\n", message.VehicleOwnerCPF)
	// ... outras informações ...
}