package Golang_TCP_ChatRoom

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

var users = make(map[int]User)

func userexist(id int) bool {

	_, ok := users[id]

	return ok

}

func allowchat(u int, t int) bool {

	fmt.Fprintf(users[t].cnn, "chat request from %v [%v] \n Do you accept? (y/n)", users[u].name, users[u].id)

	scaner := bufio.NewScanner(users[t].cnn)

	for scaner.Scan() {

		fs := strings.Fields(scaner.Text())

		switch strings.ToLower(fs[0]) {

		case "y":

			return true

		case "n":

			break

		default:
			fmt.Fprintln(users[t].cnn, "Wrong command please enter Y or N .")
		}

	}

	return false
}

func chathandle(u int, t int) {

	fmt.Fprintln(users[u].cnn, "chat Started with ", users[t].name, ".")

	scaner := bufio.NewScanner(users[u].cnn)

	for scaner.Scan() {

		fs := strings.Fields(scaner.Text())

		if len(fs) == 1 && fs[0] == "dc" {
			fmt.Fprintf(users[t].cnn, "%v has been disconnect", users[u].name)
			fmt.Fprintf(users[u].cnn, "You Ended the chat . with : ", users[t].name)
			break
		} else {
			fmt.Fprintf(users[t].cnn, "%v : %v", users[u].name, scaner.Text())
		}

	}

}

func handle(cnn net.Conn) {

	var user User

	user.id = rand.Intn(1000)

	for userexist(user.id) == false {

		user.id = rand.Intn(1000)

	}

	user.cnn = cnn

	scaner := bufio.NewScanner(user.cnn)

	fmt.Printf("Enter A NickName : ")

	scaner.Scan()

	user.name = scaner.Text()

	log.Println("User joined ", user.id, user.name)

	fmt.Fprintln(cnn, "WELCOME (You are online) USE cm for show COMMANDS.")

	for scaner.Scan() {
		fs := strings.Fields(scaner.Text())
		switch fs[0] {
		case "cm":
			fmt.Fprintln(cnn, "--------------CM--------------")
			fmt.Fprintln(cnn, "Commands			Actions")
			fmt.Fprintln(cnn, "myid				show your id")
		case "myid":
			fmt.Fprintln(cnn, "your id = ", user.id)
		case "chat":

			if len(fs) != 2 {
				fmt.Fprintln(cnn, "CM Error : WRONG Command use cm for show Commands.")
			} else {
				targetid, err := strconv.ParseInt(fs[1], 10, 64)

				if err != nil {
					fmt.Fprintln(cnn, "Target ID must be a number]\n")
				}

				if userexist(int(targetid)) {

					if allowchat(user.id, int(targetid)) {

						go chathandle(user.id, int(targetid))
						go chathandle(int(targetid), user.id)

					} else {

						fmt.Fprintln(user.cnn, "Target User Did Not Accept Your Request .")

					}
				} else {

					fmt.Fprintf(user.cnn, "User Does Not Exist .")

				}

			}
		}
	}
}
