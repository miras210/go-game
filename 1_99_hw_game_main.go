package main

import (
	"strings"
)

type Location struct {
	name      string
	neighbors []Location
	places    []Place
}

func (l *Location) addNeighbour(n Location) {
	l.neighbors = append(l.neighbors, n)
}

func (l *Location) getNeighbours() string {
	length := len(l.neighbors)
	ans := "можно пройти - "
	for i := 0; i < length; {
		ans += l.neighbors[i].name
		if i != length-1 {
			ans += ", "
		}
	}
	return ans
}

func (l *Location) addPlace(p Place) {
	l.places = append(l.places, p)
}

func (l *Location) getPlaces() string {
	ans := ""
	length := len(l.places)
	for i := 0; i < length; {
		ans += l.places[i].getObjects()
		if i != length-1 {
			ans += ", "
		}
	}
	return ans
}

type Place struct {
	name    string
	objects []Object
}

func (p *Place) getObjects() string {
	length := len(p.objects)
	ans := "на " + p.name + ": "
	for i := 0; i < length; {
		ans += p.objects[i].getName()
		if i != length-1 {
			ans += ", "
		}
	}
	return ans
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

func main() {
	/*

		сделать построчный ввод команд тут

	*/
}

func initGame() {
	/*
		эта функция инициализирует игровой мир - все команты
		если что-то было - оно корректно перезатирается
	*/

}

func handleCommand(command string) string {
	s := strings.Split(command, " ")
	switch s[0] {
	case "осмотреться":
		return "осматриваюсь"
	case "идти":
		return "иду"
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
