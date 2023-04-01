package main

import "fmt"

// Interface o SuperClase
type IProduct interface {
	setStock(stock int)
	getStock() int
	setName(name string)
	getName() string
}

// Vamos a asumir que nuestra empresa vende computadoras
type Computer struct {
	name  string
	stock int
}

func (c *Computer) setStock(stock int) {
	c.stock = stock
}

func (c *Computer) setName(name string) {
	c.name = name
}

func (c *Computer) getStock() int {
	return c.stock
}

func (c *Computer) getName() string {
	return c.name
}

// Golang no maneja la herencia directa
// Sino que maneja la composicion por herencia
// Al poner COmputer dentro del struct Laptop
// Laptop va a pasar a tener todas las propiedades
// de un struct Computer y sus funciones
type Laptop struct {
	Computer
}

// Creamos un constructor para la laptop
// Pero devuelvo IProduct en vez de Laptop
// De esta manera yo puedo crear una: Laptop,
// una PC Escritorio, un monitor, (Zapatos, ropa, camiones)
// que van a tener
// todas la mismas propiedades de IProduct
// Entonces no nos tenemos que preocupar sobre los
// SubTipos especificos que estan siendo creados
// Ese es el poder del patron Factory
func newLaptop() IProduct {
	return &Laptop{
		Computer: Computer{
			name:  "Laptop Computer",
			stock: 25,
		},
	}
}

type Desktop struct {
	Computer
}

// Lo mimso para una desktop PC
func newDesktop() IProduct {
	return &Desktop{
		Computer: Computer{
			name:  "Desktop Computer",
			stock: 35,
		},
	}
}

func GetComputerFactory(computerType string) (IProduct, error) {
	if computerType == "laptop" {
		return newLaptop(), nil
	}

	if computerType == "desktop" {
		return newDesktop(), nil
	}

	return nil, fmt.Errorf("Invalid computer type")
}

func printNameAndStock(p IProduct) {
	fmt.Printf("Product name: %s, with stock %d\n", p.getName(), p.getStock())
}

func main() {
	laptop, _ := GetComputerFactory("laptop")
	desktop, _ := GetComputerFactory("desktop")

	printNameAndStock(laptop)
	printNameAndStock(desktop)
}
