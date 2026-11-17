package service

import (
	"strings"
	"uuid"
)

// profileURL is the attendee's public profile on the Univents front: what a
// badge's QR points at.
func (o *Operations) profileURL(userID uuid.UUID) string {
	return strings.TrimRight(o.AppURL, "/") + "/profile/" + userID.String()
}
