package repository

import (
	"sync"
	"time"

	"recommendation-platform/model"
)

type NotificationRepository struct {
	lock  sync.RWMutex
	next  int
	items []model.Notification
}

var notifications = &NotificationRepository{
	next:  1,
	items: []model.Notification{},
}

func AddNotification(n model.Notification) model.Notification {
	notifications.lock.Lock()
	defer notifications.lock.Unlock()
	n.ID = notifications.next
	n.CreatedAt = time.Now().Unix()
	notifications.next++
	notifications.items = append(notifications.items, n)
	return n
}

func GetNotificationsByUser(userID string) []model.Notification {
	notifications.lock.RLock()
	defer notifications.lock.RUnlock()
	out := []model.Notification{}
	for i := len(notifications.items) - 1; i >= 0; i-- {
		n := notifications.items[i]
		if n.UserID == userID {
			out = append(out, n)
		}
	}
	return out
}

func MarkNotificationRead(userID string, id int) {
	notifications.lock.Lock()
	defer notifications.lock.Unlock()
	for i := range notifications.items {
		if notifications.items[i].UserID == userID && notifications.items[i].ID == id {
			notifications.items[i].Read = true
			return
		}
	}
}

func MarkAllRead(userID string) {
	notifications.lock.Lock()
	defer notifications.lock.Unlock()
	for i := range notifications.items {
		if notifications.items[i].UserID == userID {
			notifications.items[i].Read = true
		}
	}
}
