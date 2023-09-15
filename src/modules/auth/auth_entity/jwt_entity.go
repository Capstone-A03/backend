package authentity

import (
	uc "capstonea03/be/src/modules/user/user_constant"

	"github.com/google/uuid"
)

type JWTPayload struct {
	ID         *uuid.UUID `json:"id,omitempty"`
	Role       *uc.Role   `json:"role,omitempty"`
	Expiration *int64     `json:"exp,omitempty"`
}
