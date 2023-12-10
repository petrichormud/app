package routes

import "fmt"

const (
	Characters           = "/characters"
	CharacterApplication = "/character/application"
)

func NewCharacterApplicationPath() string {
	return fmt.Sprintf("%s/new", CharacterApplication)
}

func CharacterApplicationPath(id string) string {
	return fmt.Sprintf("%s/%s", CharacterApplication, id)
}

func CharacterApplicationNamePath(id string) string {
	return fmt.Sprintf("%s/%s/name", CharacterApplication, id)
}

func CharacterApplicationGenderPath(id string) string {
	return fmt.Sprintf("%s/%s/gender", CharacterApplication, id)
}

func CharacterApplicationShortDescriptionPath(id string) string {
	return fmt.Sprintf("%s/%s/shortdescription", CharacterApplication, id)
}

func CharacterApplicationDescriptionPath(id string) string {
	return fmt.Sprintf("%s/%s/description", CharacterApplication, id)
}

func CharacterApplicationBackstoryPath(id string) string {
	return fmt.Sprintf("%s/%s/backstory", CharacterApplication, id)
}