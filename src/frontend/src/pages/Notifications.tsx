import { useState, useEffect } from 'react';
import { notificationsApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { formatDistanceToNow } from 'date-fns';
import './Dashboard.css';
import './Notifications.css';

export function NotificationsPage() {
    const { user } = useAuth();
    const [notifications, setNotifications] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadNotifications = async () => {
            if (!user) return;

            try {
                const result = await notificationsApi.getUserNotifications(user.id, 1, 50);
                setNotifications(result.notifications || []);
            } catch (error) {
                console.error('Failed to load notifications:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadNotifications();
    }, [user]);

    const handleMarkAsRead = async (id: string) => {
        try {
            await notificationsApi.markAsRead(id);
            setNotifications(prev =>
                prev.map(n => ((n.ID || n.id) === id ? { ...n, is_read: true } : n))
            );
        } catch (error) {
            console.error('Failed to mark notification as read:', error);
        }
    };

    const getIcon = (type: string) => {
        switch (type) {
            case 'booking_request': return '📅';
            case 'booking_confirmed': return '✅';
            case 'booking_cancelled': return '❌';
            case 'payment_received': return '💰';
            case 'message_received': return '💬';
            case 'item_review': return '⭐';
            default: return '🔔';
        }
    };

    return (
        <div className="notifications-page dashboard-page">
            <div className="dashboard-header">
                <div className="container">
                    <h1>Notifications</h1>
                    <p>Stay updated with your rentals and bookings</p>
                </div>
            </div>

            <div className="container">
                {isLoading ? (
                    <div className="page-loading">
                        <div className="spinner"></div>
                    </div>
                ) : notifications.length > 0 ? (
                    <div className="notifications-list">
                        {notifications.map((notif, index) => (
                            <div
                                key={notif.ID || notif.id || `notif-${index}`}
                                className={`notification-item ${!notif.is_read ? 'unread' : ''}`}
                                onClick={() => !notif.is_read && handleMarkAsRead(notif.ID || notif.id)}
                            >
                                <div className="notification-icon">
                                    {getIcon(notif.type)}
                                </div>
                                <div className="notification-content">
                                    <div className="notification-header">
                                        <h3 className="notification-title">{notif.title}</h3>
                                        <span className="notification-time">
                                            {notif.created_at ? formatDistanceToNow(new Date(notif.created_at), { addSuffix: true }) : ''}
                                        </span>
                                    </div>
                                    <p className="notification-message">{notif.message}</p>
                                    {!notif.is_read && (
                                        <div className="notification-actions">
                                            <button
                                                className="mark-read-btn"
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    handleMarkAsRead(notif.ID || notif.id);
                                                }}
                                            >
                                                Mark as read
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="empty-state" style={{ marginTop: '4rem' }}>
                        <div className="empty-state-icon">🔔</div>
                        <h3 className="empty-state-title">All caught up!</h3>
                        <p className="empty-state-text">
                            You don't have any notifications at the moment.
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}
