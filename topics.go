package main

import "strings"

type pubsubDiscordTopicAddr struct {
	Direction string // inbound or outbound
	Protocol  string // discord, always
	GuildID   string // guild/server ID
	ChannelID string // channel ID (ie, #general)
	UserID    string // user ID (ie, @user). This should be the user's snowflake ID.
}

func (p *pubsubDiscordTopicAddr) String() string {
	return p.Direction + "." + p.Protocol + "." + p.GuildID + "." + p.ChannelID + "." + p.UserID
}

func newPubsubDiscordTopicAddr(topic string) (*pubsubDiscordTopicAddr, error) {
	var p pubsubDiscordTopicAddr
	parts := strings.Split(topic, ".")
	p.Direction = parts[0]
	p.Protocol = parts[1]
	p.GuildID = parts[2]
	p.ChannelID = parts[3]
	p.UserID = parts[4]
	return &p, nil // TODO: error handling
}

// getReplyTopic returns the topic to which a reply should be sent.
func (p *pubsubDiscordTopicAddr) getReplyTopic() string {
	var dir string = "inbound"
	if p.Direction == "inbound" {
		dir = "outbound"
	}
	return dir + "." + p.Protocol + "." + p.GuildID + "." + p.ChannelID + "." + p.UserID
}
