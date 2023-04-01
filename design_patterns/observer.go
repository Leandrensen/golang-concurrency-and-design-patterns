package main

import "fmt"

type Topic interface {
	register(observer Observer)
	broadcast()
}

type Observer interface {
	getId() string
	updateValue(string)
}

// Item -> No disponible
// Item -> Disponible -> Avise que hay item
type Item struct {
	observers []Observer
	name      string
	available bool
}

// Constructor del nuevo item
func NewItem(name string) *Item {
	return &Item{
		name: name,
	}
}

// Creamos estas 2 funciones para que Item sea
// un Topic (Implementamos la interfaz)
func (i *Item) UpdateAvailable() {
	fmt.Printf("Item %s is available\n", i.name)
	i.available = true
	i.broadcast()
}

func (i *Item) broadcast() {
	for _, observer := range i.observers {
		observer.updateValue(i.name)
	}
}

func (i *Item) register(observer Observer) {
	i.observers = append(i.observers, observer)
}

type EmailClient struct {
	id string
}

// Creamos estas 2 funciones para que EmailClient
// sea un Observer (Implement su interfaz)
func (eC *EmailClient) updateValue(value string) {
	fmt.Printf("Sending Email - %s available from client %s\n", value, eC.id)
}

func (eC *EmailClient) getId() string {
	return eC.id
}

func main() {
	nvidiaItem := NewItem("RTX 3080")
	firstObserver := &EmailClient{
		id: "12ab",
	}
	secondObserver := &EmailClient{
		id: "34cd",
	}

	nvidiaItem.register(firstObserver)
	nvidiaItem.register(secondObserver)
	nvidiaItem.UpdateAvailable()
}
