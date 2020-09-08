package serializer

import (
	funk "github.com/thoas/go-funk"
)

type StringSubscription map[string]Subscription

type Subscriptions struct {
	ByChannelID map[string]StringSubscription // store the list of subscriptions for a channelID
	ByKey       map[string][]string           // stores the list of channelIDs to which the message needs to be posted for a subscription
}

func NewSubscriptions() *Subscriptions {
	return &Subscriptions{
		ByChannelID: map[string]StringSubscription{},
		ByKey:       map[string][]string{},
	}
}

// Add adds a new subscription to the list of all subscriptions
func (l *Subscriptions) Add(s Subscription) {
	key := s.GetKey()
	if _, contains := l.ByKey[key]; !contains {
		l.ByKey[key] = make([]string, 0)
	}

	if !funk.Contains(l.ByKey[key], s.ChannelID) {
		l.ByKey[key] = append(l.ByKey[key], s.ChannelID)
	}

	if _, found := l.ByChannelID[s.ChannelID]; !found {
		l.ByChannelID[s.ChannelID] = make(StringSubscription)
	}

	l.ByChannelID[s.ChannelID][key] = s
}

// Remove removes a subscription from the list of all subscriptions
func (l *Subscriptions) Remove(s Subscription) {
	key := s.GetKey()
	delete(l.ByChannelID[s.ChannelID], key)
	l.ByKey[key] = funk.FilterString(l.ByKey[key], func(el string) bool {
		return el == s.ChannelID
	})
}

// GetChannelID returns the channelID to which the message for a subscription should be posted to
func (l *Subscriptions) GetChannelIDs(s Subscription) []string {
	return l.ByKey[s.GetKey()]
}

// List returns the list for a particular channel as a formatted mattermost message
func (l *Subscriptions) List(channelID string) []Subscription {
	values := make([]Subscription, 0, len(l.ByChannelID[channelID]))
	for _, v := range l.ByChannelID[channelID] {
		values = append(values, v)
	}

	return values
}