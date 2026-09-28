package entity

const (
	FieldTitle                     = "title"
	FieldDescription               = "description"
	FieldStatus                    = "status"
	FieldVersion                   = "version"
	FieldCardID                    = "cardId"
	ValidationMessageCanonicalUUID = "must be a canonical UUID"
)

const (
	OperationListCards  = "list_cards"
	OperationCreateCard = "create_card"
	OperationUpdateCard = "update_card"
	OperationDeleteCard = "delete_card"
	FieldLogCardID      = "card_id"
	DecisionConflict    = "conflict"
)
