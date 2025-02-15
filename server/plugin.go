package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"sync/atomic"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/pkg/errors"
)

const (
	botUsername    = "bside-bot"
	botDisplayName = "B-Side"
	botDescription = "A bot account created by the Welcomebot plugin."

	welcomebotChannelWelcomeKey = "chanmsg_"
)

// Plugin represents the welcome bot plugin
type Plugin struct {
	plugin.MattermostPlugin

	client *pluginapi.Client

	welcomeMessages atomic.Value

	// botUserID of the created bot account.
	botUserID string
}

// OnActivate ensure the bot account exists
func (p *Plugin) OnActivate() error {
	p.client = pluginapi.NewClient(p.API, p.Driver)

	existingBot, _ := p.client.Bot.Get(botUsername, false)
	if existingBot == nil {
		bot := &model.Bot{
			Username:    botUsername,
			DisplayName: botDisplayName,
			Description: botDescription,
		}
		botUserID, appErr := p.client.Bot.EnsureBot(bot)
		if appErr != nil {
			return errors.Wrap(appErr, "failed to ensure bot user")
		}
		p.botUserID = botUserID
	} else {
		p.botUserID = existingBot.UserId
	}

	err := p.API.RegisterCommand(getCommand())
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}
