package objid

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ChangeObjectiveId2StringId(oid primitive.ObjectID) string {
	return string(oid.Hex())
}

func ChangeStringId2ObjectId(oid string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(oid)
}
