package serializer

import (
	"encoding/json"
	"fmt"
	url2 "net/url"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const (
	CommentCreatedEvent = "comment_created"
	CommentUpdatedEvent = "comment_updated"
	CommentRemovedEvent = "comment_removed"
	PageCreatedEvent    = "page_created"
	PageUpdatedEvent    = "page_updated"
	PageTrashedEvent    = "page_trashed"
	PageRestoredEvent   = "page_restored"
	PageRemovedEvent    = "page_removed"
	SubscriptionTypeSpace = "space_subscription"
	SubscriptionTypePage =	"page_subscription"
)

var eventDisplayName = map[string]string{
	CommentCreatedEvent: "Comment Create",
	CommentUpdatedEvent: "Comment Update",
	CommentRemovedEvent: "Comment Remove",
	PageCreatedEvent:    "Page Create",
	PageUpdatedEvent:    "Page Update",
	PageTrashedEvent:    "Page Trash",
	PageRestoredEvent:   "Page Restore",
	PageRemovedEvent:    "Page Remove",
}

type Subscription interface {
	Add(s *Subscriptions)
	Remove(s *Subscriptions)
	Edit(s *Subscriptions)
	Name() string
	GetFormattedSubscription() string
}

type BaseSubscription struct {
	Alias     string   `json:"alias"`
	BaseURL   string   `json:"baseURL"`
	Events    []string `json:"events"`
	ChannelID string   `json:"channelID"`
	Type      string   `json:"subscriptionType"`
}

type PageSubscription struct {
	PageID string `json:"pageID"`
	BaseSubscription
}

type SpaceSubscription struct {
	SpaceKey  string   `json:"spaceKey"`
	BaseSubscription
}

type StringSubscription map[string]Subscription
type StringArrayMap map[string][]string

type Subscriptions struct {
	ByChannelID map[string]StringSubscription //
	ByURLPagID map[string]StringArrayMap
	ByURLSpaceKey map[string]StringArrayMap
}

func NewSubscriptions() *Subscriptions {
	return &Subscriptions{
		ByChannelID:   map[string]StringSubscription{},
		ByURLPagID:    map[string]StringArrayMap{},
		ByURLSpaceKey: map[string]StringArrayMap{},
	}
}

func (ps PageSubscription) Add(s *Subscriptions) {
	fmt.Println("done done=", ps)
	_, valid := s.ByChannelID[ps.ChannelID]
	if !valid {
		s.ByChannelID[ps.ChannelID] = make(StringSubscription)
	}
	s.ByChannelID[ps.ChannelID][ps.Alias] = ps
	key := store.GetURLPageIDCombinationKey(ps.BaseURL, ps.PageID)
	_, ok :=  s.ByURLPagID[key]
	if !ok {
		s.ByURLPagID[key] = make(map[string][]string)
	}
	fmt.Println("url=", s.ByURLPagID[key])
	s.ByURLPagID[key][ps.ChannelID] = ps.Events
}

func (ps PageSubscription) Remove (s *Subscriptions) {
	delete(s.ByChannelID[ps.ChannelID], ps.Alias)
	key := store.GetURLPageIDCombinationKey(ps.BaseURL, ps.PageID)
	delete(s.ByURLPagID[key], ps.ChannelID)
}

func (ps PageSubscription) Edit (s *Subscriptions) {
	ps.Remove(s)
	ps.Add(s)
}

func (ps PageSubscription) Name() string {
	return SubscriptionTypePage
}

func (ps PageSubscription) GetFormattedSubscription() string {
	var events []string
	for _, event := range ps.Events {
		events = append(events, eventDisplayName[event])
	}
	return fmt.Sprintf("\n|%s|%s|%s|%s|", ps.Alias, ps.BaseURL, ps.PageID, strings.Join(events, ", "))
}

func (ss SpaceSubscription) Add (s *Subscriptions) {
	if s.ByChannelID[ss.ChannelID] == nil {
		s.ByChannelID[ss.ChannelID] = make(StringSubscription)
	}
	s.ByChannelID[ss.ChannelID][ss.Alias] = ss
	key := store.GetURLSpaceKeyCombinationKey(ss.BaseURL, ss.SpaceKey)
	if s.ByURLSpaceKey[key] == nil {
		s.ByURLSpaceKey[key] = make(map[string][]string)
	}
	s.ByURLSpaceKey[key][ss.ChannelID] = ss.Events
}

func (ss SpaceSubscription) Remove (s *Subscriptions) {
	delete(s.ByChannelID[ss.ChannelID], ss.Alias)
	key := store.GetURLSpaceKeyCombinationKey(ss.BaseURL, ss.SpaceKey)
	delete(s.ByURLSpaceKey[key], ss.ChannelID)
}

func (ss SpaceSubscription) Edit (s *Subscriptions) {
	ss.Add(s)
	ss.Remove(s)
}

func (ss SpaceSubscription) Name() string {
	return SubscriptionTypeSpace
}

func (ss SpaceSubscription) GetFormattedSubscription() string {
	var events []string
	for _, event := range ss.Events {
		events = append(events, eventDisplayName[event])
	}
	return fmt.Sprintf("\n|%s|%s|%s|%s|", ss.Alias, ss.BaseURL, ss.SpaceKey, strings.Join(events, ", "))
}

func (ps *PageSubscription) IsValid() error {
	// TODO : Clean subscription data
	if ps.Alias == "" {
		return errors.New("alias can not be empty")
	}
	if ps.BaseURL == "" {
		return errors.New("base url can not be empty")
	}
	if _, err := url2.Parse(ps.BaseURL); err != nil {
		return errors.New("enter a valid url")
	}
	if ps.PageID == "" {
		return errors.New("page id can not be empty")
	}
	return nil
}

// func (s *StringSubscription) UnmarshalJSON(data []byte) error {
// 	subs := make(map[string]map[string]interface{})
// 	err := json.Unmarshal(data, &subs)
// 	if err != nil {
// 		return err
// 	}
// 	result := make(StringSubscription)
// 	for k, v := range subs {
//
// 		switch k {
// 		case SubscriptionTypePage :
// 			ps := PageSubscription{}
// 			err := json.Unmarshal(v, &ps)
// 			if err != nil {
// 				return err
// 			}
// 			result[k] = ps
// 		case SubscriptionTypeSpace:
// 			ss := SpaceSubscription{}
// 			err := json.Unmarshal(v, &ss)
// 			if err != nil {
// 				return err
// 			}
// 			result[k] = ss
// 		default:
// 			return errors.New("Unrecognized subscription")
// 		}
// 	}
// 	*s = result
// 	return nil
// }


func (s *StringSubscription) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	result := make(StringSubscription)
	for k, v := range m {
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		value, err := UnmarshalCustomSubscription(bytes, "subscriptionType", map[string]reflect.Type{
				SubscriptionTypePage: reflect.TypeOf(PageSubscription{}),
				SubscriptionTypeSpace: reflect.TypeOf(SpaceSubscription{}),
		})
		if err != nil {
			return err
		}
		result[k] = value.(Subscription)
	}

	*s = result
	return nil
}

// UnmarshalCustomSubscription returns subscription of type page or space from bytes.
func UnmarshalCustomSubscription(data []byte, typeJsonField string, customTypes map[string]reflect.Type) (interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	typeName := m[typeJsonField].(string)
	var value Subscription
	if ty, found := customTypes[typeName]; found {
		value = reflect.New(ty).Interface().(Subscription)
	}

	valueBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(valueBytes, &value); err != nil {
		return nil, err
	}

	return value, nil
}

func SubscriptionsFromJson(bytes []byte) (*Subscriptions, error) {

	var subs *Subscriptions
	if len(bytes) != 0 {
		fmt.Println("bytes=", string(bytes))
		unmarshalErr := json.Unmarshal(bytes, &subs)
		fmt.Println("unmarshal error=", unmarshalErr)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
	} else {
		subs = NewSubscriptions()
	}
	fmt.Println("subsssssss=", subs)
	return subs, nil
}

// func getSubscription(i interface{}) Subscription {
// 	if s, ok := i.(Subscription); ok {
// 		return s
// 	}
//
// 	bytes, err := json.Marshal(i)
// 	if err != nil {
// 		fmt.Println("111111111111111", err)
// 		return nil
// 	}
//
// 	msi, ok := i.(map[string]interface{})
// 	if !ok {
// 		return nil
// 	}
// 	switch msi["subscriptionType"] {
// 	case SubscriptionTypePage:
// 		var ps PageSubscription
// 		if ps, ok := i.(PageSubscription); ok {
// 			return ps
// 		}
// 		if err := json.Unmarshal(bytes, &ps); err != nil {
// 			return ps
// 		}
//
// 	case SubscriptionTypeSpace:
// 		var ss SpaceSubscription
// 		if ss, ok := i.(SpaceSubscription); ok {
// 			return ss
// 		}
// 		if err := json.Unmarshal(bytes, &ss); err != nil {
// 			return ss
// 		}
// 	}
//
// 	return nil
// }

// func SubscriptionsFromJson(data io.Reader) (*Subscriptions, error) {
// 	var subs Subscriptions
// 	if err := json.NewDecoder(data).Decode(&subs); err != nil {
// 		fmt.Println("unmarshal error=", err)
// 		return nil, err
// 	}
// 	fmt.Println("subsssssss=", subs)
// 	return &subs, nil
// }

func FormattedSubscriptionList(channelSubscriptions StringSubscription) string {
	var pageSubscriptions, spaceSubscription, list string
	pageSubscriptionsHeader := fmt.Sprintf("| Alias | Base Url | Page Id | Events|\n| :----|:--------| :--------| :-----|")
	spaceSubscriptionsHeader := fmt.Sprintf("| Alias | Base Url | Space Key | Events|\n| :----|:--------| :--------| :-----|")
	for _, sub := range channelSubscriptions {
		if sub.Name() == SubscriptionTypeSpace {
			pageSubscriptions += sub.GetFormattedSubscription()
		} else if sub.Name() == SubscriptionTypePage {
			spaceSubscription += sub.GetFormattedSubscription()
		}
	}
	if spaceSubscription != "" {
		list = spaceSubscriptionsHeader + "\n" + spaceSubscription + "\n"
	}
	if pageSubscriptions != "" {
		list = pageSubscriptionsHeader + "\n" + pageSubscriptions
	}
	return list
}

// func (s *Subscription) ToJSON() string {
// 	b, err := json.Marshal(s)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(b)
// }
