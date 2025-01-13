package common

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/rand"
	"strings"
)

var Flags = []string{"French", "English", "Dutch", "Spanish", "Pirate"}

func roll() bool {
	return rand.Intn(2) == 0
}
func grab(s []string) string {
	return s[rand.Intn(len(s))]
}

func GenerateCaptainName() string {
	var title = []string{"Captian", "Admiral", "Apprentice", "Pirate", "Skipper", "Commander", "Boatswain", "Officer",
		"Traitor", "Ghostly", "Commodore", "Agent"}
	var name = []string{"Gleeson", "Orvin", "Ripley", "Preston", "Eldon", "Dorset", "Falk", "Jorge", "Frederick",
		"Hunter", "Jasper", "Salvodor", "Hailey", "Rackham", "Crowther", "Lucifer", "Woolworth", "Dunstan", "Clare",
		"Kaiser", "Kerwin", "Morris", "Ulrik", "Asema", "Storm", "Gladwin", "Morse", "Zell", "Penny"}
	var flair = []string{"Dishonest", "Soft Heart", "Balding", "Rum Lover", "Two Toes", "Hair", "Gloomy",
		"Cutthroat", "Dastardly", "Vile", "Ripe", "Pungent", "Piggy", "Pleasant", "Crazy", "Weasel", "Squealer", "Feral",
		"Snake", "Slayer", "Ghostly", "Traitor", "Coxswain", "One-tooth"}
	var flairWithThe = []string{"Dishonest", "Soft Heart", "Balding", "Rum Lover", "Hair", "Gloomy",
		"Cutthroat", "Dastardly", "Vile", "Ripe", "Pungent", "Piggy", "Pleasant", "Crazy", "Weasel", "Squealer",
		"Snake", "Slayer", "Traitor", "Coxswain", "One-tooth"}
	var place = []string{"Coxswain", "Rackham", "Seezley", "Salty", "Briller", "Dunstan", "Cordith", "Firth", "Barbady",
		"Yorben", "Nillith", "Salvador", "Lactipon", "Derzley", "Jitterham", "Milktown", "Appleton", "Greently", "Asstin",
		"Hoplonton", "Welgadin"}

	fullName := []string{}
	if roll() {
		fullName = append(fullName, grab(title))
	}
	if roll() {
		fullName = append(fullName, grab(name))
	}
	if roll() {
		fullName = append(fullName, "von")
		fullName = append(fullName, grab(name))
	} else if roll() {
		fullName = append(fullName, "the")
		fullName = append(fullName, grab(flairWithThe))
	} else if roll() {
		fullName = append(fullName, "\""+grab(flair)+"\"")
	} else {
		fullName = append(fullName, grab(name))
	}

	if roll() {
		fullName = append(fullName, "of")
		fullName = append(fullName, grab(place))
	}
	return cases.Title(language.English).String(strings.Join(fullName, " "))
}

func GetRandomFlag() string {
	return Flags[rand.Intn(len(Flags))]
}
