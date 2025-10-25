package terror

import "fmt"

const (
	TypeObjectNotFound  = "ObjectNotFound"
)

type (
	ObjectNotFound struct {
		Base
		Id        int32
		ObjectType string
	}
)

func NewObjectNotFound(id int32, objecType string) *ObjectNotFound {
	msg := fmt.Sprintf("Object %d of type %s is not found", id, objecType)
	return &ObjectNotFound{
		Base: Base{
			Message: msg,
			Type:    TypeObjectNotFound,
		},
		Id:        id,
		ObjectType: objecType,
	}
}

var errInstances = map[string]interface{}{
	TypeObjectNotFound: ObjectNotFound{},
}

func Unmarshal(body []byte) (result error) {
	return unmarshal(body, errInstances)
}