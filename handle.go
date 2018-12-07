package main

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

func handle(cnn net.Conn) {

	var user User

	user.id = rand.Intn(1000)

	for userexist(user.id) != false {

		user.id = rand.Intn(1000)

	}

	fmt.Fprint(cnn, "Please Enter a Nickname --> ")

	user.cnn = cnn

	user.chatscanner = bufio.NewScanner(user.cnn)

	user.chatscanner.Scan()

	user.name = user.chatscanner.Text()

	users[user.id] = user

	log.Println("User joined ", user.id, user.name)

	fmt.Fprintln(cnn, "WELCOME ", user.name, " (You are online) USE cm for show COMMANDS.")

	commandhandle(user.id)
}

func commandhandle(userid int) {

	for {

		user := users[userid]

		user.chatscanner.Scan()

		fs := strings.Fields(user.chatscanner.Text())
		if len(fs) > 0 {

			switch fs[0] {
			case "cm":
				fmt.Fprintln(user.cnn, "Commands		                 	Actions")
				fmt.Fprintln(user.cnn, "myid				                show your ID")
				fmt.Fprintln(user.cnn, "accept				            accept chat requests")
				fmt.Fprintln(user.cnn, "reject			                reject chat requests")
				fmt.Fprintln(user.cnn, "start			                    start chat after your target accepted")
				fmt.Fprintln(user.cnn, "chat [User Target ID]				chat with some one")
			case "myid":
				fmt.Fprintln(user.cnn, "your id = ", user.id)
			case "chat":

				if len(fs) != 2 {
					fmt.Fprintln(user.cnn, "CM Error : WRONG Command use cm for show Commands.")
				} else {
					targetid, err := strconv.ParseInt(fs[1], 10, 64)

					if err != nil {
						fmt.Fprintln(user.cnn, "Target ID must be a number]\n")
					}

					t := int(targetid)

					if userexist(t) {

						chatrequest(userid, t)

					} else {

						fmt.Fprintln(user.cnn, "User Does Not Exist .")

					}

				}
			case "accept":

				user = users[userid]

				if user.chatrequest.hostid == 0 {
					fmt.Fprintln(user.cnn, "There is no chat request for you .")
				} else {
					fmt.Fprintln(users[user.chatrequest.hostid].cnn, user.name, "accepted your request type start command to start chat.")

					targetuser := users[user.chatrequest.hostid]
					targetuser.chatrequest.guestaccept = true
					users[user.chatrequest.hostid] = targetuser

					fmt.Fprintln(user.cnn, "chat Started with ", targetuser.name, ". /n Enter dc command to End the chat")

					fmt.Fprintln(user.cnn, "waiting for", targetuser.name)

					chathandle(userid, user.chatrequest.hostid)
				}

			case "reject":

				user = users[userid]

				if user.chatrequest.hostid == 0 {
					fmt.Fprintln(user.cnn, "There is no chat request for you .")
				} else {
					fmt.Fprintln(users[user.chatrequest.hostid].cnn, user.name, "Did not accept your chat request.")

					targetuser := users[user.chatrequest.hostid]
					targetuser.chatrequest.guestaccept = false
					targetuser.chatrequest.guestid = 0
					users[user.chatrequest.hostid] = targetuser

					user.chatrequest.hostid = 0
					users[userid] = user
				}

			case "start":

				user = users[userid]

				if user.chatrequest.guestid == 0 {

					fmt.Fprintln(user.cnn, "There is not chat connection")

				} else if user.chatrequest.guestid != 0 && !user.chatrequest.guestaccept {

					fmt.Fprintf(user.cnn, "%v did not reject or accept your chat request.\n", users[user.chatrequest.guestid].name)

				} else {

					fmt.Fprintf(users[user.chatrequest.guestid].cnn, "%v joined the chat.\n", user.name)
					chathandle(userid, user.chatrequest.guestid)

				}
			}

		}

	}
}

func chatrequest(u, t int) {

	targetuser := users[t]
	targetuser.chatrequest.hostid = u
	users[t] = targetuser

	hostuser := users[u]
	hostuser.chatrequest.guestid = t
	users[u] = hostuser

	fmt.Fprintf(users[t].cnn, "chat request from %v[%v] Do you accept ? (yes/no) \n", users[u].name, users[u].id)

	log.Printf("%v [%v] requested to %v [%v].", users[u].name, users[u].id, users[t].name, users[t].id)

}

func chathandle(u int, t int) {

	log.Printf("chat started between %v[%v] and %v[%v]\n", users[u].name, users[u].id, users[t].name, users[t].id)

	for users[u].chatscanner.Scan() {

		fs := strings.Fields(users[u].chatscanner.Text())

		if len(fs) == 1 && fs[0] == "dc" {

			if userexist(t) {
				fmt.Fprintf(users[t].cnn, "%v Left the chat .\n", users[u].name)
			}

			fmt.Fprintf(users[u].cnn, "You Left the chat.\n")

			if users[u].chatrequest.hostid != 0 {

				targetuser := users[t]
				targetuser.chatrequest.guestaccept = false
				targetuser.chatrequest.guestid = 0
				users[t] = targetuser

				user := users[u]
				user.chatrequest.hostid = 0
				users[u] = user

			} else {

				targetuser := users[u]
				targetuser.chatrequest.guestaccept = false
				targetuser.chatrequest.guestid = 0
				users[u] = targetuser

				user := users[t]
				user.chatrequest.hostid = 0
				users[t] = user
			}

			return
		} else {
			fmt.Fprintf(users[t].cnn, "%v : %v \n", users[u].name, users[u].chatscanner.Text())
		}

	}

}
