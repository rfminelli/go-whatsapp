package whatsapp_connection

import (
	"strconv"
	"github.com/Rhymen/go-whatsapp/binary"
	"time"
	"fmt"
)

//TODO: filename? WhatsApp uses Store.Contacts for these functions
//TODO: functions probably shouldn't return a string, maybe build a struct / return json
//TODO: check for further queries
func (wac *conn) GetProfilePicThumb(jid string) (<-chan string, error) {
	data := []interface{}{"query", "ProfilePicThumb", jid}
	return wac.write(data)
}

func (wac *conn) GetStatus(jid string) (<-chan string, error) {
	data := []interface{}{"query", "Status", jid}
	return wac.write(data)
}

func (wac *conn) GetGroupMetaData(jid string) (<-chan string, error) {
	data := []interface{}{"query", "GroupMetadata", jid}
	return wac.write(data)
}

func (wac *conn) SubscribePresence(jid string) (<-chan string, error) {
	data := []interface{}{"action", "presence", "subscribe", jid}
	return wac.write(data)
}

func (wac *conn) CreateGroup(subject string, participants []string) (<-chan string, error) {
	return wac.setGroup("create", "", subject, participants)
}

func (wac *conn) UpdateGroupSubject(subject string, jid string) (<-chan string, error) {
	return wac.setGroup("subject", jid, subject, nil)
}

func (wac *conn) SetAdmin(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("promote", jid, "", participants)
}

func (wac *conn) RemoveAdmin(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("demote", jid, "", participants)
}

func (wac *conn) AddMember(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("add", jid, "", participants)
}

func (wac *conn) RemoveMember(jid string, participants []string) (<-chan string, error) {
	return wac.setGroup("remove", jid, "", participants)
}

func (wac *conn) setGroup(t, jid, subject string, participants []string) (<-chan string, error) {
	ts := time.Now().Unix()
	id := fmt.Sprintf("%d.--%d", ts, wac.msgCount)

	//TODO: get proto or improve encoder to handle []interface{}

	p := buildParticipantNodes(participants)

	g := binary.Node{
		Description: "group",
		Attributes: map[string]string{
			"author": wac.session.Wid,
			"id": id,
			"type": t,
		},
		Content: p,
	}

	if jid != "" {
		g.Attributes["jid"] = jid
	}

	if subject != "" {
		g.Attributes["subject"] = subject
	}

	n := binary.Node{
		Description: "action",
		Attributes: map[string]string{
			"type":  "set",
			"epoch": strconv.Itoa(wac.msgCount),
		},
		Content: []interface{}{g},
	}

	wac.msgCount++
	return wac.writeBinary(n, GROUP, IGNORE, id)
}

func buildParticipantNodes(participants []string) ([]binary.Node) {
	l := len(participants)
	if participants == nil || l == 0 {
		return nil
	}

	p := make([]binary.Node, len(participants))
	for i, participant := range participants {
		p[i] = binary.Node{
			Description: "participant",
			Attributes: map[string]string{
				"jid": participant,
			},
		}
	}
	return p
}