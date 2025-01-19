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
	var title = []string{"Captain", "Admiral", "Apprentice", "Pirate", "Skipper", "Commander", "Boatswain", "Officer",
		"Traitor", "Ghostly", "Commodore", "Agent", "Seaman", "Rebel", "Privateer", "First Mate", "Buccaneer", "Sir"}
	var name = []string{"Gleeson", "Orvin", "Ripley", "Preston", "Eldon", "Dorset", "Falk", "Jorge", "Frederick",
		"Hunter", "Jasper", "Salvodor", "Hailey", "Rackham", "Crowther", "Lucifer", "Woolworth", "Dunstan", "Claire",
		"Kaiser", "Kerwin", "Morris", "Ulrik", "Asema", "Storm", "Gladwin", "Morse", "Zell", "Penny", "Janice",
		"Barbara", "Ironia", "Shauntelle", "Elvira", "Esmeralda", "Bob", "Trixie", "Wendy", "Franz", "Peggy", "Anous",
		"Dick", "Gaylord", "Angus", "Pud", "Bruce", "Marty", "Wolfgang", "Hyacinth", "Zimoslav", "Rufulus", "Nolif",
		"Lollie", "Malvina", "Stella", "Xensor", "Bentley", "Cordelia", "Johnson", "Muff", "Titus", "Anthony"}
	var flair = []string{"Dishonest", "Soft Heart", "Balding", "Rum Lover", "Two Toes", "Hair", "Gloomy",
		"Cutthroat", "Dastardly", "Vile", "Ripe", "Pungent", "Piggy", "Pleasant", "Crazy", "Weasel", "Squealer", "Feral",
		"Snake", "Slayer", "Ghostly", "Traitor", "Coxswain", "One-tooth", "Windy", "Butter", "Betrayer", "Foxy"}
	var flairWithThe = []string{"Dishonest", "Soft Heart", "Balding", "Rum Lover", "Hair", "Gloomy",
		"Cutthroat", "Dastardly", "Vile", "Ripe", "Pungent", "Piggy", "Pleasant", "Crazy", "Weasel", "Squealer",
		"Snake", "Slayer", "Traitor", "Coxswain", "One-tooth", "Windy", "Butter", "Cozy", "Tide Turner", "Bear", "Savage"}
	var placePrefix = []string{"Port", "Isle of", "Saint", "South", "North", "East", "West", "Mt"}
	var place = []string{"Coxswain", "Rackham", "Seezley", "Salty", "Briller", "Dunstan", "Cordith", "Firth", "Barbady",
		"Yorben", "Nillith", "Sanctitly", "Laction", "Derzley", "Jitterham", "Milktown", "Appleton", "Greently", "Asstin",
		"Hoplonton", "Welgadin", "Klappertown", "Windville", "Folkenwald", "Dids", "Munkton", "Shallows", "Plaqard", "Oiltown", "Willows", "Quellton"}
	var placeSuffix = []string{"Bay", "Island", "Falls", "Harbour", "Lake", "River", "Way", "Rock", "Springs", "Bend", "Beach", "Point"}

	fullName := []string{}
	last := "none"
	if roll() {
		fullName = append(fullName, grab(title))
		last = "title"
	}
	if roll() {
		fullName = append(fullName, grab(name))
		last = "firstName"
	}
	if roll() && roll() && roll() && last != "title" {
		fullName = append(fullName, "von")
		fullName = append(fullName, grab(name))
		last = "von"
	} else if roll() && last != "title" {
		fullName = append(fullName, "the")
		fullName = append(fullName, grab(flairWithThe))
		last = "theflair"
	} else if roll() {
		fullName = append(fullName, "\""+grab(flair)+"\"")
		last = "flair"
	} else {
		fullName = append(fullName, grab(name))
		last = "lastName"
	}

	if last == "none" || last == "title" || last == "flair" {
		fullName = append(fullName, grab(name))
	}

	if roll() && roll() {
		fullName = append(fullName, "of")
		if roll() && roll() {
			fullName = append(fullName, grab(placePrefix))
		}
		fullName = append(fullName, grab(place))
		if roll() && roll() {
			fullName = append(fullName, grab(placeSuffix))
		}
	} else if roll() && roll() && roll() && roll() && roll() {
		fullName = append(fullName, "yon")
		fullName = append(fullName, grab(place))
	}
	return cases.Title(language.English).String(strings.Join(fullName, " "))
}

func GetRandomFlag() string {
	return Flags[rand.Intn(len(Flags))]
}
