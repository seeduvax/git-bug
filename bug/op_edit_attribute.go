package bug

import (
	"encoding/json"
	"github.com/MichaelMure/git-bug/entity"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/MichaelMure/git-bug/util/timestamp"
)

var _ Operation = &EditAttributeOperation{}

type EditAttributeOperation struct {
	OpBase
	name string
	value string
	set bool // true set, false unset (erase attribute)
}

func (op *EditAttributeOperation) base() *OpBase {
	return &op.OpBase
}

func (op *EditAttributeOperation) Id() entity.Id {
	return idOperation(op)
}


type AttributeEditTimelineItem struct {
	id	entity.Id
	Author	identity.Interface
	UnixTime	timestamp.Timestamp
	name	string
	value	string
	set	bool
}

func (a *AttributeEditTimelineItem) Id() entity.Id {
	return a.id
}

func (a *AttributeEditTimelineItem) IsAuthored() {}

func (op *EditAttributeOperation) Apply(snapshot *Snapshot) {
	snapshot.addActor(op.Author)
	snapshot.addParticipant(op.Author)
	changed:=false;
	for i, attr := range snapshot.Attributes {
		if attr.name == op.name {
			if op.set {
				attr.value=op.value
				changed=true
			} else {
				snapshot.Attributes[i] = snapshot.Attributes[len(snapshot.Attributes)-1]
				snapshot.Attributes = snapshot.Attributes[:len(snapshot.Attributes)-1]
			}
		}
	}
	if op.set && !changed {
		attribute := Attribute{
			name: op.name,
			value: op.value,
		}
		snapshot.Attributes = append(snapshot.Attributes, attribute)
	}
	item:=&AttributeEditTimelineItem{
		id: op.Id(),
		Author: op.Author,
		UnixTime: timestamp.Timestamp(op.UnixTime),
		name: op.name,
		value: op.value,
		set: op.set,
	}
	snapshot.Timeline=append(snapshot.Timeline,item)
}

func (op *EditAttributeOperation) Validate() error {
		if err:=opBaseValidate(op, EditAttributeOp); err != nil {
				return err
		}
		// TODO what specific operation needed for attribute validation ????
		// May be if attribute name is starting with "link:", check the value
		// is a valid bug id.
		return nil
}

func (op* EditAttributeOperation) MarshalJSON() ([]byte, error) {
	base, err:=json.Marshal(op.OpBase)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err:=json.Unmarshal(base,&data); err!=nil {
		return nil, err
	}
	data["name"]=op.name
	data["value"]=op.value
	data["set"]=op.set
	return json.Marshal(data)
}

func (op* EditAttributeOperation) UnmarshalJSON(data []byte) error {
	base := OpBase{}
	err := json.Unmarshal(data, &base)
	if err != nil {
		return err
	}

	aux := struct {
		Name string `json:"name"`
		Value string `json:"value"`
		Set bool `json:"set"`
	}{}
	err = json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	op.OpBase = base
	op.name = aux.Name
	op.value = aux.Value
	op.set = aux.Set

	return nil
}

func NewEditAttributeOp(author identity.Interface, unixTime int64, name string, value string, set bool) *EditAttributeOperation {
		return &EditAttributeOperation {
				OpBase: newOpBase(EditAttributeOp, author, unixTime),
				name: name,
				value: value,
				set: set,
		}
}



