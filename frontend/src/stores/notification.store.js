import { defineStore } from 'pinia';
import { notificationService } from '../services/notification.service';

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [],
    loading: false,
  }),

  getters: {
    unreadCount: (state) => state.notifications.filter(n => !n.is_read).length,
  },

  actions: {
    async fetch() {
      this.loading = true;
      try {
        const data = await notificationService.getNotifications();
        this.notifications = Array.isArray(data) ? data : (data?.notifications || []);
      } catch (error) {
        console.error('Failed to fetch notifications:', error);
      } finally {
        this.loading = false;
      }
    },

    async markRead(notificationId) {
      try {
        await notificationService.markAsRead(notificationId);
        const notif = this.notifications.find(n => n.id === notificationId);
        if (notif) notif.is_read = true;
      } catch (error) {
        console.error('Failed to mark notification as read:', error);
      }
    },

    async markAllRead() {
      try {
        await notificationService.markAllAsRead();
        this.notifications = this.notifications.map(n => ({ ...n, is_read: true }));
      } catch (error) {
        console.error('Failed to mark all as read:', error);
      }
    },
  }
});
