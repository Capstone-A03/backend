package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Where []primitive.E
type Order []primitive.E

type FindAllWhere struct {
	Where          Where
	IncludeInCount bool
}

type CountOptions struct {
	Where *[]Where
}

type FindOneOptions struct {
	Where *[]Where
	Order *[]Order
}

type FindAllOptions struct {
	Where  *[]FindAllWhere
	Order  *[]Order
	Limit  *int
	Offset *int
}

type DistinctOptions struct {
	FieldName string
	Where     *[]FindAllWhere
}

type UpdateOptions struct {
	Where *[]Where
}

type ReplaceOptions struct {
	Where *[]Where
}

type DestroyOptions struct {
	Where *[]Where
}

type Pagination struct {
	Limit int
	Count int
	Total int
}
