package main

import (
	"fmt"
	"strings"
)

type Level struct {
	locations []*Location
	name      string
	neighbour *Location
}

type Location struct {
	name             string
	description      string
	neighbors        []*Location
	places           []*Place
	lvlNeighbour     *Level
	specialCondition func() string
}

func newLocation(name, description string) *Location {
	return &Location{
		name:             name,
		neighbors:        []*Location{},
		description:      description,
		places:           []*Place{},
		specialCondition: nil,
	}
}

func (l *Location) addNeighbour(n *Location) {
	l.neighbors = append(l.neighbors, n)
}

func (l *Location) getNeighbours() string {
	length := len(l.neighbors)
	ans := "можно пройти - "
	if l.lvlNeighbour != nil {
		return ans + l.lvlNeighbour.name
	}
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
		if l.places[i].isEmpty() {
			continue
		}
		ans += l.places[i].getObjects()
		if i != length-1 && !l.places[i+1].isEmpty() {
			ans += ", "
		}
	}
	if ans == "" {
		return "пустая комната. "
	}
	ans += ". "
	return ans
}

type Place struct {
	name             string
	objects          []Object
	specialCondition func() (bool, string)
	special          bool
}

func newPlace(name string) *Place {
	return &Place{
		name:    name,
		objects: []Object{},
	}
}

func (p *Place) isEmpty() bool {
	if len(p.objects) == 0 {
		return true
	}
	return false
}

func (p *Place) getObjects() string {
	if p.isEmpty() {
		return ""
	}
	length := len(p.objects)
	ans := "на " + p.name + "е: "
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

type Mission struct {
	completed func() bool
	mission   string
}

type Player struct {
	cur         *Location
	mission     *[]Mission
	missionText string
	fLocDesc    string
	fLocName    string
	pEquipment  Equipment
	pInventory  []Item
}

func (p *Player) lookout() string {
	if p.fLocName == p.cur.name {
		ans := p.fLocDesc
		ans += p.cur.getPlaces()
		ans = ans[:len(ans)-2]
		newMissionText := ", надо "
		for i, m := range *p.mission {
			if !m.completed() {
				newMissionText += m.mission
			} else {
				continue
			}
			if i != len(*p.mission)-1 {
				newMissionText += " и "
			}
		}
		newMissionText += ". "
		ans += newMissionText
		ans += p.cur.getNeighbours()
		return ans
	} else {
		ans := p.cur.getPlaces()
		ans += p.cur.getNeighbours()
		return ans
	}
}

func (p *Player) moveto(location string) string {
	moved := false
	for _, loc := range p.cur.neighbors {
		if loc.name == location {
			if loc.specialCondition != nil {
				if loc.specialCondition() != "" {
					return loc.specialCondition()
				}
			}
			p.cur = loc
			moved = true
			break
		}
	}
	if moved {
		ans := p.cur.description
		ans += p.cur.getNeighbours()
		return ans
	}
	return "нет пути в " + location
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
				if p.pEquipment.getName() == "" {
					return "некуда класть"
				}
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
					if p.specialCondition != nil {
						msg := ""
						p.special, msg = p.specialCondition()
						if p.special {
							return msg
						}
					}
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
	initGame()
	var s string
	for player.cur.name != "улица" {
		//TODO SCAN COMMANDS WITH WHITESPACE
		//just wtf, why it does not work. I hate my life :(
		_, _ = fmt.Scanln(&s)
		fmt.Println(handleCommand(s))
	}
	fmt.Println("Поздравляем! Вы прошли эту игру!")
}

func initGame() {
	kitchen := newLocation("кухня", "кухня, ничего интересного. ")
	kTable := newPlace("стол")
	tea := &Item{name: "чай"}
	kTable.addObject(tea)
	kitchen.addPlace(kTable)

	corridor := newLocation("коридор", "ничего интересного. ")
	door := newPlace("дверь")
	corridor.addPlace(door)

	room := newLocation("комната", "ты в своей комнате. ")
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

	street := newLocation("улица", "на улице весна. ")

	kitchen.addNeighbour(corridor)
	room.addNeighbour(corridor)
	street.addNeighbour(corridor)
	corridor.addNeighbour(kitchen)
	corridor.addNeighbour(room)
	corridor.addNeighbour(street)

	street.specialCondition = func() string {
		if door.special == false {
			return "дверь закрыта"
		}
		return ""
	}

	door.specialCondition = func() (bool, string) {
		for _, item := range player.pInventory {
			if item.name == "ключи" {
				return true, "дверь открыта"
			}
		}
		return false, "дверь закрыта"
	}

	home := Level{
		locations: []*Location{},
		name:      "домой",
		neighbour: street,
	}

	street.lvlNeighbour = &home

	mission1 := Mission{completed: func() bool {
		for _, item := range player.pInventory {
			if item.name == "конспекты" {
				return true
			}
		}
		return false
	},
		mission: "собрать рюкзак"}
	mission2 := Mission{completed: func() bool {
		for _, item := range player.pInventory {
			if item.name == "ключи" && player.cur.name == "улица" {
				return true
			}
		}
		return false
	},
		mission: "идти в универ"}

	player = &Player{
		cur:         kitchen,
		pEquipment:  Equipment{name: ""},
		pInventory:  []Item{},
		fLocDesc:    "ты находишься на кухне, ",
		fLocName:    kitchen.name,
		missionText: "надо " + mission1.mission + " и " + mission2.mission + ". ",
		mission:     &[]Mission{mission1, mission2},
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
