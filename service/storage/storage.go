package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

type Storage struct {
	BasePath string
}

func New(basePath string) *Storage {
	return &Storage{
		BasePath: basePath,
	}
}

func (s *Storage) SavePFP(data []byte, userID string) (string, error) {
	if len(data) == 0 {
		return "", database.ErrInvalidImage
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", database.ErrInvalidImage
	}

	if format != "jpeg" && format != "png" {
		return "", database.ErrInvalidImage
	}

	//SIZE CVHECK GOES HERE

	//img can be saved
	extension := ".png"
	if format == "jpeg" {
		extension = ".jpg"
	}
	filename := userID + extension
	dir := filepath.Join(s.BasePath, "pfp")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", database.ErrStorage
	}

	path := filepath.Join(dir, filename)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", database.ErrStorage
	}

	return "/pfp/" + filename, nil
}

// ask gpt for this one
func (s *Storage) DeletePFP(photoURI string) error {
	if photoURI == "" {
		return nil
	}

	filename := filepath.Base(photoURI)
	path := filepath.Join(s.BasePath, "pfp", filename)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return database.ErrStorage
	}
	return nil

}

func (s *Storage) SaveMSGPhoto(data []byte, messageID string, chatID string) (string, error) {
	if len(data) == 0 {
		return "", database.ErrInvalidImage
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", database.ErrInvalidImage
	}

	if format != "jpeg" && format != "png" {
		return "", database.ErrInvalidImage
	}

	extension := "." + "png"
	if format == "jpeg" {
		extension = ".jpg" //ASK AGAIN
	}

	filename := messageID + extension
	dir := filepath.Join(s.BasePath, "msgpics", chatID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", database.ErrStorage
	}

	path := filepath.Join(dir, filename)

	if err = os.WriteFile(path, data, 0644); err != nil {
		return "", database.ErrStorage
	}

	return "/msgpics/" + chatID + "/" + filename, nil
}

func (s *Storage) SaveMSGJson(message database.Message) error {
	dir := filepath.Join(s.BasePath, "msg", message.ChatID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return database.ErrStorage
	}

	data, err := json.MarshalIndent(message, "", "	")
	if err != nil {
		return database.ErrStorage
	}

	filename := message.MessageID + ".json"

	path := filepath.Join(dir, filename)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return database.ErrStorage
	}
	return nil

}

func (s *Storage) DeleteMSGFiles(messageID string, chatID string) error {

	jsonPath := filepath.Join(s.BasePath, "msg", chatID, messageID+".json")
	if err := os.Remove(jsonPath); err != nil && !os.IsNotExist(err) {
		return database.ErrStorage
	}

	//if present, delete msg photo
	for _, extension := range []string{".jpg", ".png"} {
		photopath := filepath.Join(s.BasePath, "msgpics", chatID, messageID+extension)
		if err := os.Remove(photopath); err != nil && !os.IsNotExist(err) {
			return database.ErrStorage
		}

	}
	//delete reactions
	reactionsdirpath := filepath.Join(s.BasePath, "reactions", chatID, messageID)
	if err := os.RemoveAll(reactionsdirpath); err != nil {
		return database.ErrStorage
	}

	return nil

}

// same as savemsgjson but with a reaction
func (s *Storage) AddReaction(reac database.Reaction, chatID string) error {

	dir := filepath.Join(s.BasePath, "reactions", chatID, reac.MessageID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return database.ErrStorage
	}

	data, err := json.MarshalIndent(reac, "", "	")
	if err != nil {
		return database.ErrStorage
	}

	filename := reac.ReactionID + ".json"

	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return database.ErrStorage
	}
	return nil
}

func (s *Storage) DeleteReaction(msgID string, userID string, chatID string, reactionID string) error {

	jsonPath := filepath.Join(s.BasePath, "reactions", chatID, msgID, reactionID+".json")
	if err := os.Remove(jsonPath); err != nil && !os.IsNotExist(err) {
		return database.ErrStorage
	}

	return nil

}

// func for deleting all messages, pictures and reactions when about
// to leave an orphaned chat
func (s *Storage) DeleteChatGroup(chatID string) error {
	msgdirpath := filepath.Join(s.BasePath, "msg", chatID)
	err := os.RemoveAll(msgdirpath)
	if err != nil {
		return database.ErrStorage
	}

	msgpicsdirpath := filepath.Join(s.BasePath, "msgpics", chatID)
	if err = os.RemoveAll(msgpicsdirpath); err != nil {
		return database.ErrStorage
	}

	reactionsdirpath := filepath.Join(s.BasePath, "reactions", chatID)
	if err = os.RemoveAll(reactionsdirpath); err != nil {
		return database.ErrStorage
	}

	return nil
}
