import { defineStore } from 'pinia';
import { notificationService } from '../services/notification.service';
import { authStore } from './auth.store';

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [],
    loading: false,
    pollIntervalId: null,
  }),

  getters: {
    unreadCount: (state) => state.notifications.filter(n => !n.is_read).length,
  },

  actions: {
    async fetch() {
      // Avoid overlapping requests if already loading
      if (this.loading) return;
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

    startPolling() {
      if (this.pollIntervalId) return;
      // Fetch immediately
      this.fetch();
      // Then poll every 30 seconds
      this.pollIntervalId = setInterval(() => {
        if (authStore.isAuthenticated) {
          this.fetch();
        } else {
          this.stopPolling();
        }
      }, 30000);
    },

    stopPolling() {
      if (this.pollIntervalId) {
        clearInterval(this.pollIntervalId);
        this.pollIntervalId = null;
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
