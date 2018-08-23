// mautrix-whatsapp - A Matrix-WhatsApp puppeting bridge.
// Copyright (C) 2018 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package whatsapp_ext

import (
	"github.com/Rhymen/go-whatsapp"
	"encoding/json"
	"strings"
)

type PresenceType string

const (
	PresenceUnavailable PresenceType = "unavailable"
	PresenceAvailable   PresenceType = "available"
)

type Presence struct {
	JID       string       `json:"id"`
	Status    PresenceType `json:"type"`
	Timestamp int64        `json:"t"`
}

type PresenceHandler interface {
	whatsapp.Handler
	HandlePresence(Presence)
}

func (ext *ExtendedConn) handleMessagePresence(message []byte) {
	var event Presence
	err := json.Unmarshal(message, &event)
	if err != nil {
		ext.jsonParseError(err)
		return
	}
	event.JID = strings.Replace(event.JID, OldUserSuffix, NewUserSuffix, 1)
	for _, handler := range ext.handlers {
		presenceHandler, ok := handler.(PresenceHandler)
		if !ok {
			continue
		}
		presenceHandler.HandlePresence(event)
	}
}