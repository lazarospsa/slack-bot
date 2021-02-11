package mocks

import (
	"encoding/json"
	"fmt"
	"github.com/innogames/slack-bot/bot/msg"
	"github.com/innogames/slack-bot/bot/util"
	"github.com/innogames/slack-bot/client"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"testing"
)

func AssertSlackMessage(slackClient *SlackClient, ref msg.Ref, text string) {
	slackClient.On("SendMessage", ref, text).Once().Return("")
}

func AssertReaction(slackClient *SlackClient, reaction string, ref msg.Ref) {
	slackClient.On("AddReaction", util.Reaction(reaction), ref).Once()
}

func AssertRemoveReaction(slackClient *SlackClient, reaction string, ref msg.Ref) {
	slackClient.On("RemoveReaction", util.Reaction(reaction), ref).Once()
}

func AssertQueuedMessage(t *testing.T, expected msg.Message) {
	t.Helper()

	actual := <-client.InternalMessages
	assert.Equal(t, actual, expected)
}

// AssertSlackJSON is a test helper to assert full slack attachments
func AssertSlackJSON(t *testing.T, slackClient *SlackClient, message msg.Ref, expected url.Values) {
	t.Helper()

	slackClient.On("SendMessage", message, "", mock.MatchedBy(func(option slack.MsgOption) bool {
		_, values, _ := slack.UnsafeApplyMsgOptions(
			"token",
			"channel",
			"apiUrl",
			option,
		)

		expected.Add("token", "token")
		expected.Add("channel", "channel")

		assert.Equal(t, expected, values)

		return true
	})).Once().Return("")
}

// AssertSlackBlocks test helper to assert a given JSON representation of "Blocks"
func AssertSlackBlocks(t *testing.T, slackClient *SlackClient, message msg.Ref, expectedJSON string) {
	t.Helper()

	slackClient.On("SendBlockMessage", message, mock.MatchedBy(func(givenBlocks []slack.Block) bool {
		// replace the random tokens to fixed ones for easier mocking
		for i := range givenBlocks {
			if actionBlock, ok := givenBlocks[i].(*slack.ActionBlock); ok {
				if button, ok := actionBlock.Elements.ElementSet[0].(*slack.ButtonBlockElement); ok {
					button.Value = fmt.Sprintf("token-%d", i)
				}
			}
		}
		givenJSON, err := json.Marshal(givenBlocks)
		assert.Nil(t, err)

		fmt.Println(string(givenJSON))

		return expectedJSON == string(givenJSON)
	})).Once().Return("")
}
