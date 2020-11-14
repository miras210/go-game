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

func (p *Place) deleteObject(name string) {
	for i, o := range p.objects {
		if o == nil {
			continue
		}
		if o.getName() == name {
			p.objects[i] = p.objects[len(p.objects)-1]
			p.objects[len(p.objects)-1] = nil
			p.objects = p.objects[:len(p.objects)-1]
		}
	}
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

func (i *Item) useItem() string {
	return "дверь открыта"
}

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

func (p *Player) equip(e string) string {
	for _, place := range p.cur.places {
		for _, equipment := range place.objects {
			if equipment == nil {
				continue
			}
			if equipment.getName() == e {
				switch equipment.(type) {
				case *Equipment:
					p.pEquipment = Equipment{name: e}
					place.deleteObject(e)
					return "вы надели: " + e
				}
			}
		}
	}
	return "нет такого"
}

func (p *Player) getItem(i string) string {
	for _, place := range p.cur.places {
		for _, equipment := range place.objects {
			if equipment == nil {
				continue
			}
			if equipment.getName() == i {
				switch equipment.(type) {
				case *Item:
					p.pInventory = append(p.pInventory, Item{name: i})
					place.deleteObject(i)
					return "предмет добавлен в инвентарь: " + i
				}
			}
		}
	}
	return "нет такого"
}

func (p *Player) useItem(item, place string) string {
	for _, i := range p.pInventory {
		if i.getName() == item {
			for _, p := range p.cur.places {
				if p.name == place {
					return i.useItem()
				}
			}
			return "не к чему применить"
		}
	}
	return "нет предмета в инвентаре - " + item
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
		return player.equip(s[1])
	case "взять":
		return player.getItem(s[1])
	case "применить":
		return player.useItem(s[1], s[2])
	default:
		return "неизвестная команда"
	}
}
