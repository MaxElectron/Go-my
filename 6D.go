//go:build !solution

package pubsub

import (
	"context"
	"errors"
	"sync"
)

type Subscriber struct {
	id                 int
	messageHandler     MsgHandler
	lastProcessedIndex int
	topic              *PubSubSystem

	unsubscribeSignal chan struct{}
	closeSignal       chan struct{}
	notifySignal      chan struct{}
	finishedSignal    chan struct{}
}

func (sub *Subscriber) Unsubscribe() {
	close(sub.topic.unsubscribeChannels[sub.id])

	sub.topic.mutex.Lock()
	delete(sub.topic.subscriberChannels, sub.id)
	sub.topic.mutex.Unlock()
}

var _ PubSub = (*PubSubSystem)(nil)

type PubSubSystem struct {
	isClosed             bool
	isClosedTopic        bool
	topics               map[string]*PubSubSystem
	messages             []interface{}
	lastPublishedIndex   int
	nextSubscriberID     int
	subscriberChannels   map[int]chan struct{}
	unsubscribeChannels  map[int]chan struct{}
	closeChannels        map[int]chan struct{}
	notificationChannels map[int]chan struct{}
	mutex                sync.RWMutex
}

func NewPubSub() PubSub {
	return &PubSubSystem{
		topics:               make(map[string]*PubSubSystem),
		messages:             []interface{}{},
		lastPublishedIndex:   -1,
		nextSubscriberID:     0,
		subscriberChannels:   make(map[int]chan struct{}),
		unsubscribeChannels:  make(map[int]chan struct{}),
		closeChannels:        make(map[int]chan struct{}),
		notificationChannels: make(map[int]chan struct{}),
	}
}

func (ps *PubSubSystem) Subscribe(subject string, handler MsgHandler) (Subscription, error) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if ps.isClosed {
		return nil, errors.New("pubsub system is closed")
	}

	closeChan := make(chan struct{})
	unsubscribeChan := make(chan struct{})
	notifyChan := make(chan struct{}, 1)
	finishedChan := make(chan struct{}, 1)

	topic, exists := ps.topics[subject]
	if !exists {
		topic = &PubSubSystem{
			messages:             []interface{}{},
			lastPublishedIndex:   -1,
			nextSubscriberID:     0,
			subscriberChannels:   make(map[int]chan struct{}),
			unsubscribeChannels:  make(map[int]chan struct{}),
			closeChannels:        make(map[int]chan struct{}),
			notificationChannels: make(map[int]chan struct{}),
		}
		ps.topics[subject] = topic
	}

	topic.mutex.Lock()
	defer topic.mutex.Unlock()

	subID := topic.nextSubscriberID
	topic.nextSubscriberID++

	topic.subscriberChannels[subID] = notifyChan
	topic.unsubscribeChannels[subID] = unsubscribeChan
	topic.closeChannels[subID] = closeChan
	topic.notificationChannels[subID] = finishedChan

	newSubscriber := &Subscriber{
		id:                 subID,
		messageHandler:     handler,
		lastProcessedIndex: topic.lastPublishedIndex,
		topic:              topic,
		notifySignal:       notifyChan,
		unsubscribeSignal:  unsubscribeChan,
		closeSignal:        closeChan,
		finishedSignal:     finishedChan,
	}

	go newSubscriber.listen()

	return newSubscriber, nil
}

func (sub *Subscriber) listen() {
	for {
		select {
		case <-sub.unsubscribeSignal:
			return
		case <-sub.notifySignal:
			for {
				sub.topic.mutex.RLock()
				if sub.lastProcessedIndex == sub.topic.lastPublishedIndex {
					sub.topic.mutex.RUnlock()
					break
				}
				sub.lastProcessedIndex++
				message := sub.topic.messages[sub.lastProcessedIndex]
				sub.topic.mutex.RUnlock()

				sub.messageHandler(message)
			}
		case <-sub.closeSignal:
			for {
				sub.topic.mutex.RLock()
				if sub.lastProcessedIndex == sub.topic.lastPublishedIndex {
					sub.topic.mutex.RUnlock()
					break
				}
				sub.lastProcessedIndex++
				message := sub.topic.messages[sub.lastProcessedIndex]
				sub.topic.mutex.RUnlock()

				sub.messageHandler(message)
			}
			sub.finishedSignal <- struct{}{}
			return
		}
	}
}

func (ps *PubSubSystem) Publish(subject string, message interface{}) error {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if ps.isClosed {
		return errors.New("pubsub system is closed")
	}

	topic, exists := ps.topics[subject]
	if !exists {
		return errors.New("topic does not exist")
	}

	topic.mutex.Lock()
	defer topic.mutex.Unlock()

	topic.messages = append(topic.messages, message)
	topic.lastPublishedIndex++

	for _, notifyChan := range topic.subscriberChannels {
		select {
		case notifyChan <- struct{}{}:
		default:
		}
	}

	return nil
}

func (ps *PubSubSystem) Close(ctx context.Context) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.isClosed = true
	for _, topic := range ps.topics {
		topic.mutex.Lock()
		topic.isClosedTopic = true
		topic.mutex.Unlock()

		topic.mutex.RLock()
		for id := range topic.subscriberChannels {
			close(topic.closeChannels[id])
		}
		topic.mutex.RUnlock()

		for id := range topic.subscriberChannels {
			select {
			case <-ctx.Done():
			case <-topic.notificationChannels[id]:
			}
		}
	}

	return nil
}
