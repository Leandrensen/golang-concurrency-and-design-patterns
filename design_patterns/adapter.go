package main

import "fmt"

type Payment interface {
	Pay()
}

type CashPayment struct{}

func (CashPayment) Pay() {
	fmt.Println("Pay using Cash")
}

func ProcessPayment(p Payment) {
	p.Pay()
}

type BankPayment struct{}

// Bank es diferente en su funcion Pay
// Y necesita un bankAccount number para funcionar
func (BankPayment) Pay(bankAccount int) {
	fmt.Printf("Paying using Bank Account: %d\n", bankAccount)
}

// ProcessPayment no nos va a dejar usar bank de primera
// Porque es un type que no implementa correctamente la interface de Payment
// Entonces creamos un adaptador para solucionar este problema
type BankPaymentAdapter struct {
	BankPayment *BankPayment
	bankAccount int
}

// Y ahora hacemos que BankPaymentAdapter implemente correctamente
// la funcion Pay()
func (bpa *BankPaymentAdapter) Pay() {
	// Y dentro de la funcion implementamos el Pay
	// en la manera que lo hace el type BankPayment
	bpa.BankPayment.Pay(bpa.bankAccount)
}

func main() {
	cash := &CashPayment{}
	ProcessPayment(cash)
	// Esto no funciona
	// bank := &BankPayment{}
	// ProcessPayment(bank)
	//Esto si funciona
	bpa := &BankPaymentAdapter{
		bankAccount: 5,
		BankPayment: &BankPayment{},
	}
	ProcessPayment(bpa)
}
