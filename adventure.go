package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// item's struct that contains item's fields
type item struct {
	name        string
	description string // look item will get this description
	visible     bool
	takeable    bool
}

// room's struct that contains room's fields
type room struct {
	directions  map[string]*room //map of direction that has key as string and room struct as value
	name        string
	description string
	inventory   []item
	visited     bool
}

// this function called look has room as a parameter, display room description
func (r room) look() {
	fmt.Println(r.description)
	if len(r.inventory) > 0 { // check if there is any items in a room, then display items
		fmt.Println("itmes here:")
		for i := range r.inventory {
			if r.inventory[i].visible {
				fmt.Println(" a", r.inventory[i].name)
			}
		}
	}
}

// this function called gamelooop that takes poiter to room struct as a parameter
func gameloop(start *room) {
	s := bufio.NewScanner(os.Stdin) // take input from user
	playing := true
	currentroom := start
	var inventory []item
	yourroom := room{name: "your room", description: "", directions: make(map[string]*room)}
	score := 0
	for fmt.Print(">>"); s.Scan() && playing; fmt.Print(">>") {
		switch argv := strings.Fields(strings.ToLower(s.Text())); argv[0] {
		case "look": //if user type "look" call look() that display room and items description
			if len(argv) == 1 {
				currentroom.look()
			} else {
				found := false
				for _, v := range inventory { //check if input match with item's name then display the description
					if v.name == argv[1] { //if name of item at element 1, print the description
						fmt.Println(v.description)
						found = true
					}
				}
				if !found {
					foundinroom := false
					for _, v := range currentroom.inventory { // if you found item in the current room, set to true
						if v.name == argv[1] { //if user type "look item's name" then display item's description
							fmt.Println(v.description)
							foundinroom = true
						}
					}
					if !foundinroom {
						fmt.Println("I don't know where that is.")
					}
				}
			}
		case "quit": // user type "quit" to quit the game loop
			fmt.Println("ok")
			playing = false
			break
		case "take": // user type take to take items that is in current room
			found := false
			for i, v := range currentroom.inventory {
				if v.name == argv[1] {
					if v.takeable { //if an item takeable, set to true and in inventory, user can take it
						currentroom.inventory[i] = currentroom.inventory[len(currentroom.inventory)-1]
						currentroom.inventory = currentroom.inventory[:len(currentroom.inventory)-1]
						fmt.Println("You take the ", argv[1])
						inventory = append(inventory, v) // add an item to inventory after you take
					} else {
						fmt.Println("You can't take that.") // items not takeable
					}
					found = true
				}
			}
			if !found {
				fmt.Println("That isn't here.") // the item is not in this room
			}
		case "drop":
			found := false
			for i, v := range inventory {
				if v.name == argv[1] { // if the item already in inventory, user can get rid by type "drop iem's name"
					if len(inventory) > 1 {
						inventory[i] = inventory[len(currentroom.inventory)-1]
					}
					inventory = inventory[:len(inventory)-1] //take out from inventory
					currentroom.inventory = append(currentroom.inventory, v)
					fmt.Println("You drop the ", argv[1])
					found = true
				}
			}
			if !found {
				fmt.Println("You don't have one.") // drop something that you don't have
			}
		case "use": // use item in the inventory
			found := false
			for _, v := range inventory {
				if v.name == argv[1] {
					if v.name == "tree" { // tree can't use anywhere, just take to collect bonus points
						fmt.Println("You can't use that here.")
					}
					if v.name == "key" { // key can be used to unlock room's door
						if currentroom.name == "your home" { //if you look and the current room is yor home then you can "use key"
							currentroom.directions["enter"] = &yourroom // direction point to yourroom room
							fmt.Println("You unlock the door!")
						} else {
							fmt.Println("You can't use that here.")
						}
					}
					found = true
				}
			}
			if !found {
				foundinroom := false
				for _, v := range currentroom.inventory { // found slot machine in casino
					if v.name == argv[1] {
						if v.name == "machine" { //play with slot machine that throw random interger 5 times you will win one time
							v := rand.Intn(5)
							if v == 1 {
								fmt.Println("you won", score/5, "points!")
								score = score + score/5
							} else {
								fmt.Println("You gave the slot machine", score/10, "points and you didn't win.")
								score = score - score/10
							}
						}
						foundinroom = true
					}
				}
				if !foundinroom {
					fmt.Println("I don't know where that is.")
				}
			}
		default:
			//default case when first start or direction of the room at [0] and boolean ok set to true
			//if the room has the direction that user type in the ok boolean is true then continue the loop
			val, ok := (currentroom.directions)[argv[0]]
			if ok {
				currentroom = val
				if !currentroom.visited { //if you have not vistited the current room set it to true and add 10 scores per each new room
					currentroom.visited = true
					score = score + 10
				}
				fmt.Println("You are now in ", currentroom.name)
			} else {
				fmt.Println("I don't understand.") //direction only north south west east
			}

		}
		if currentroom.name == "St. Mary's river" { //if current room at St. Mary's river, user die and exit the game loop
			fmt.Println("You drowned in the cold water.")
			break
		}
		if currentroom.name == "your room" { // if the current room is your room then user win and check if user inventory has tree then get 100 scores bonus + 200 scores for winner
			fmt.Println("You made it home!")
			for _, v := range inventory {
				if v.name == "tree" {
					score = score + 100
				}
			}
			score = score + 200
			break
		}
	}
	fmt.Println("Your score was: ", score)
}

func main() {

	//create rooms for the game and each room has to follow room struct signatures
	startingroom := room{visited: false, name: "waterfront", description: "This is a long boardwalk along St. Marys river, to the North is station mall, to the South is St. Marys river", directions: make(map[string]*room), inventory: make([]item, 0)}
	waterroom := room{visited: false, name: "St. Marys's river", description: "You are drowning and dying soon", directions: make(map[string]*room), inventory: make([]item, 0)}
	mallroom := room{visited: false, name: "Station Mall", description: "It's a massive building with vaulted ceilings, a tiled floor, and tiny shops. To the North is the town of Saullt St. Marie, to the South is the waterfront.", directions: make(map[string]*room), inventory: make([]item, 0)}
	casinoroom := room{visited: false, name: "Casino", description: "You are in the casino, wanna bet? There is a big gate of the casino, To the West is a breakfast cafe. To the North is Downtown.", directions: make(map[string]*room), inventory: make([]item, 0)}
	cafe := room{visited: false, name: "Westside Cafe", description: "Your in front of a tiny cafe, it's closed today. To the East is a casino.", directions: make(map[string]*room), inventory: make([]item, 0)}
	downtown := room{visited: false, name: "downtown", description: "You're standing on Queens street. There are many buildings here, all of the shops are closed. To the North is Wacky Wings", directions: make(map[string]*room), inventory: make([]item, 0)}
	wackywings := room{visited: false, name: "Wacky Wings", description: "You're in front of a ridiculous looking bar. It's closed right now. Accross the street to the West is your house. To the North is Canadian Tire.", directions: make(map[string]*room), inventory: make([]item, 0)}
	home := room{visited: false, name: "your home", description: "You're standing in the hallway of a 5 story appartment complex. The door to your room is here", directions: make(map[string]*room), inventory: make([]item, 0)}
	canadiantire := room{visited: false, name: "Canadian Tire", description: "You're in a massive building full of hardware and sporting goods.", directions: make(map[string]*room), inventory: make([]item, 0)}
	college := room{visited: false, name: "Sault College", description: "You are in a large metal and glass building. There is an open classroom here. Everything is to the South.", directions: make(map[string]*room), inventory: make([]item, 0)}
	classroom := room{visited: false, name: "the classroom", description: "It's very boring here. There is a door to go back outside.", directions: make(map[string]*room), inventory: make([]item, 0)}
	/*
	   the town's map
	           classroom
	           college
	   youroom Canadiantire
	   home    wackywings
	           downtown
	   cafe     casino
	           mall
	           waterfront
	           =====
	*/
	//direction of each room indicate where the room pointer point to
	// the symbol & is a pointer that point current room direction to the room struct which player should be next
	startingroom.directions["south"] = &waterroom
	startingroom.directions["north"] = &mallroom
	mallroom.directions["south"] = &startingroom
	mallroom.directions["north"] = &casinoroom
	casinoroom.directions["north"] = &downtown
	casinoroom.directions["west"] = &cafe
	downtown.directions["north"] = &wackywings
	wackywings.directions["north"] = &canadiantire
	canadiantire.directions["north"] = &college
	downtown.directions["south"] = &casinoroom
	cafe.directions["east"] = &casinoroom
	wackywings.directions["south"] = &downtown
	canadiantire.directions["south"] = &wackywings
	college.directions["south"] = &canadiantire
	classroom.directions["door"] = &college
	college.directions["door"] = &classroom
	wackywings.directions["west"] = &home
	home.directions["east"] = &wackywings
	wackywings.directions["street"] = &home
	home.directions["street"] = &wackywings
	//create item using item struct
	slotmachine := item{name: "machine", description: "a beat up slot machine. You can gamble away your score with it.", visible: true, takeable: false}
	key := item{name: "key", description: "The key to your house. It's metal and has notches in it.", visible: true, takeable: true}
	tree := item{name: "tree", description: "It's a Christmass tree! You love those! It has lots of colored lights and ornaments", visible: true, takeable: true}
	//set item to each room
	casinoroom.inventory = append(casinoroom.inventory, slotmachine)
	classroom.inventory = append(classroom.inventory, key)
	canadiantire.inventory = append(canadiantire.inventory, tree)
	rand.Seed(time.Now().UnixNano())
	gameloop(&startingroom) //start at the startingroom
}
