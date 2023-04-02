// Chat Server
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

// Canal que va a transmitir string
// Este tipo es el que vamos a usar para enviar los diferentes mensajes a traves del chat
type Client chan<- string

var (
	// Clientes que se estan conectando a nuestro server chat
	// Es un canal de canales
	incomingClients = make(chan Client)
	// Clientes que se estan yendo del Server
	// Es un canal de canales
	leavingClients = make(chan Client)
	// Mensajes del chat
	messages = make(chan string)
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 3090, "port")
)

// Va escribiendo/imprimiendo todos los mensajes que se van recibiendo
// Recibimos la conexion y un canal para manejar los strings
func MessageWrite(conn net.Conn, messages <-chan string) {
	for message := range messages {
		fmt.Fprintln(conn, message)
	}
}

// Con esta funcion manejamos la conexion de 1 cliente al chat server
// Client1 -> Server -> HandleConnection(Client1)
// Cuando un cliente se conecta a nuestro servidor delegamos la coneccion a esta funcion
// Esta funcion se va a encargar de estar enviando todos los mensajes que ese cliente este escribiendo
// Y tambien de los mensajes que los otros clientes esten enviando
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	// Creamos un canal donde se van a enviar los mensajes del Cliente que se ha conectado
	message := make(chan string)
	// Esto es para que se esten imprimiendo/enviando los mensajes a traves de la conexion (conn)
	go MessageWrite(conn, message)
	// Asignamos un clientName a la conexion, como muestra arriba
	// Client1:2560 -> host:Platzi.com, port:38
	// clientName == platzi.com:38
	clientName := conn.RemoteAddr().String()
	// Con esto le enviamos un mensaje al Cliente. Al usuario que esta usando esta instancia
	message <- fmt.Sprintf("Welcome to the server, yout name is %s\n", clientName)
	// Con esto le enviamos un mensaje al canal de mensajes global de nuestro sistema
	// Esto le va a llegar a todos los clientes que esten conectados al chat server
	messages <- fmt.Sprintf("New client is here, his name is %s\n", clientName)
	// Con esto le avisamos al server que hay un nuevo cliente y que tiene que procesarlo
	// Le enviamos el message, que es el channel que creamos mas arriba dentro de esta funcion
	// Que el channel que este cliente va a utilizar para transmitir sus mensajes
	incomingClients <- message
	// Con esto vamos a leer todo lo que el cliente escriba en la terminal
	inputMessage := bufio.NewScanner(conn)
	// Este ciclo se va a romper cuando el cliente aprete "Ctrl+C" y cierre la terminal o algo similar
	// Y lo que va a hacer es enviar todo el texto que se escriba en la terminal al canal global de mensajes del servidor
	for inputMessage.Scan() {
		messages <- fmt.Sprintf("%s: %s\n", clientName, inputMessage.Text())
	}
	// Despues te este for, el cliente ha abandonado el chat
	// Entonces agregamos su canal message al canal de leavingClients
	leavingClients <- message
	// Avisamos que alguien se ha ido del chat
	messages <- fmt.Sprintf("%s said goodbye!", clientName)
}

// Logica para manejar los diferentes clientes
// Y como se van a manejar los unos con los otros
func Broadcast() {
	// Map que nos dice cuales clientes estan conectados y cuales no
	clients := make(map[Client]bool)
	for {
		// select: palabra reservada para multiplexacion de canales
		select {
		// Cuando un cliente manda un mensaje
		case message := <-messages:
			for client := range clients {
				client <- message

			}
		// Cuando se conecta un nuevo cliente
		case newClient := <-incomingClients:
			clients[newClient] = true
		// Cuando un cliente se desconecta
		case leavingClient := <-leavingClients:
			delete(clients, leavingClient)
			close(leavingClient)
		}

	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	go Broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go HandleConnection(conn)
	}
}
