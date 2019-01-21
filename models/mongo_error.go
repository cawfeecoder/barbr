package models

type MongoError struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

func GetErrorFromMongo(err error) string {
	switch (err.Error()) {
	case "mongo: no documents in result":
		return "ID does not reference any documents"
	case "the provided hex string is not a valid ObjectID":
		return "Provided ID is invalid because it is not a valid ObjectID"
	default:
		return err.Error()
	}
}
