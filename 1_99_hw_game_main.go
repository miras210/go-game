package main

import (
	"strings"
)

type Location struct {
	name      string
	neighbors []*Location
	places    []*Place
}

func newLocation(name string) *Location {

	return &Location{
		name:      name,
		neighbors: []*Location{},
		places:    []*Place{},
	}
}

func (l *Location) addNeighbour(n *Location) {
	l.neighbors = append(l.neighbors, n)
}

func (l *Location) getNeighbours() string {
	length := len(l.neighbors)
	ans := "можно пройти - "
	for i := 0; i < length; i++ {
		ans += l.neighbors[i].name
		if i != length-1 {
			ans += ", "
		}
	}
	return ans
}

func (l *Location) addPlace(p *Place) {
	l.places = append(l.places, p)
}

func (l *Location) getPlaces() string {
	ans := ""
	length := len(l.places)
	for i := 0; i < length; i++ {
		ans += l.places[i].getObjects()
		if i != length-1 {
			ans += ", "
		}
	}
	ans += ". "
	return ans
}

type Place struct {
	name    string
	objects []Object
}

func newPlace(name string) *Place {
	return &Place{
		name:    name,
		objects: []Object{},
	}
}

func (p *Place) getObjects() string {
	length := len(p.objects)
	ans := "на " + p.name + ": "
	for i := 0; i < length; i++ {
		ans += p.objects[i].getName()
		if i != length-1 {
			ans += ", "
		}
	}
	return ans
}

func (p *Place) addObject(o Object) {
	p.objects = append(p.objects, o)
}

type Object interface {
	getName() string
}

type Equipment struct {
	name string
}

func (e *Equipment) getName() string {
	return e.name
}

type Item struct {
	name string
}

func (i *Item) getName() string {
	return i.name
}

func (i *Item) useItem(p Place) {}

type Player struct {
	cur        *Location
	pEquipment Equipment
	pInventory []Item
}

func (p *Player) lookout() string {
	ans := p.cur.name + ", "
	ans += p.cur.getPlaces()
	ans += p.cur.getNeighbours()
	return ans
}

func (p *Player) moveto(location string) string {
	for _, loc := range p.cur.neighbors {
		if loc.name == location {
			p.cur = loc
			break
		}
	}
	ans := p.cur.name + ", "
	ans += p.cur.getNeighbours()
	return ans
}

var player *Player

func main() {
	/*

		сделать построчный ввод команд тут

	*/
}

func initGame() {
	kitchen := newLocation("кухня")
	kTable := newPlace("стол")
	tea := &Item{name: "чай"}
	kTable.addObject(tea)
	kitchen.addPlace(kTable)

	corridor := newLocation("коридор")
	door := newPlace("дверь")
	corridor.addPlace(door)

	room := newLocation("комната")
	rTable := newPlace("стол")
	keys := &Item{name: "ключи"}
	rTable.addObject(keys)
	notes := &Item{name: "конспекты"}
	rTable.addObject(notes)
	rChair := newPlace("стул")
	backpack := &Equipment{name: "рюкзак"}
	rChair.addObject(backpack)
	room.addPlace(rTable)
	room.addPlace(rChair)

	street := newLocation("улица")

	kitchen.addNeighbour(corridor)
	room.addNeighbour(corridor)
	street.addNeighbour(corridor)
	corridor.addNeighbour(kitchen)
	corridor.addNeighbour(room)
	corridor.addNeighbour(street)

	player = &Player{
		cur:        kitchen,
		pEquipment: Equipment{},
		pInventory: []Item{},
	}
}

func handleCommand(command string) string {
	s := strings.Split(command, " ")
	switch s[0] {
	case "осмотреться":
		return player.lookout()
	case "идти":
		return player.moveto(s[1])
	case "надеть":
		return "надеваю"
	case "взять":
		return "беру"
	case "применить":
		return "применяю"
	default:
		return "неизвестная команда"
	}
}
