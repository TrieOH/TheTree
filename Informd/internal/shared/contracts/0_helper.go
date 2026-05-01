package contracts

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type BulkGetParams struct {
	FilterKey   string `fun_query:"filter_key,required"`
	FilterOp    string `fun_query:"filter_op,required"`
	FilterValue string `fun_query:"filter_value,required"`
	FilterOrder string `fun_query:"filter_order,default=desc"`
}

var AllowedFilterKeys = map[string]string{
	"id":           "id",
	"namespace_id": "namespace_id",
	"owner_id":     "owner_id",
	"name":         "name",
	"status":       "status",
	"opened_at":    "opened_at",
	"closed_at":    "closed_at",
	"archived_at":  "archived_at",
	"created_at":   "created_at",
	"updated_at":   "updated_at",
}

var AllowedOps = map[string]string{
	"eq":   "=",
	"neq":  "!=",
	"like": "LIKE",
}
