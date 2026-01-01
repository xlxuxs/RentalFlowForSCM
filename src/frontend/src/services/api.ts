// API Service Layer for RentalFlow

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000';
console.log('Using API Base URL:', API_BASE_URL);


// ========== Helper Functions ==========
async function request<T>(
    endpoint: string,
    options: RequestInit = {}
): Promise<T> {
    const token = localStorage.getItem('access_token');

    const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    if (token) {
        (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers,
    });

    if (!response.ok) {
        const error = await response.json().catch(() => ({ error: 'Unknown error' }));
        throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
}

// ========== Auth API ==========
export const authApi = {
    register: (data: {
        email: string;
        password: string;
        first_name: string;
        last_name: string;
        phone?: string;
        role: 'renter' | 'owner';
    }) => request('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify(data),
    }),

    login: (data: { email: string; password: string }) =>
        request<{
            user: { id: string; email: string; first_name: string; last_name: string; role: string };
            access_token: string;
            refresh_token: string;
            expires_in: number;
        }>('/api/auth/login', {
            method: 'POST',
            body: JSON.stringify(data),
        }),

    logout: (userId: string) =>
        request('/api/auth/logout', {
            method: 'POST',
            body: JSON.stringify({ user_id: userId }),
        }),

    getProfile: (userId: string) =>
        request(`/api/auth/profile?user_id=${userId}`),

    validateToken: (token: string) =>
        request<{ valid: boolean; user_id?: string; email?: string; role?: string }>(
            '/api/auth/validate',
            {
                method: 'POST',
                body: JSON.stringify({ token }),
            }
        ),
};

// ========== Users API ==========
export const usersApi = {
    updateProfile: (userId: string, data: {
        first_name?: string;
        last_name?: string;
        phone?: string;
        bio?: string;
    }) =>
        request('/api/auth/profile', {
            method: 'PUT',
            body: JSON.stringify({ user_id: userId, ...data }),
        }),

    updateAvatar: (userId: string, avatarUrl: string) =>
        request('/api/auth/avatar', {
            method: 'POST',
            body: JSON.stringify({ user_id: userId, avatar_url: avatarUrl }),
        }),

    changePassword: (userId: string, currentPassword: string, newPassword: string) =>
        request('/api/auth/change-password', {
            method: 'POST',
            body: JSON.stringify({
                user_id: userId,
                current_password: currentPassword,
                new_password: newPassword,
            }),
        }),
};

// ========== Items API ==========
export const itemsApi = {
    list: (params?: {
        category?: string;
        city?: string;
        min_price?: number;
        max_price?: number;
        sort?: string;
        page?: number;
        page_size?: number;
    }) => {
        const searchParams = new URLSearchParams();
        if (params?.category) searchParams.set('category', params.category);
        if (params?.city) searchParams.set('city', params.city);
        if (params?.min_price) searchParams.set('min_price', params.min_price.toString());
        if (params?.max_price) searchParams.set('max_price', params.max_price.toString());
        if (params?.sort) searchParams.set('sort', params.sort);
        if (params?.page) searchParams.set('page', params.page.toString());
        if (params?.page_size) searchParams.set('page_size', params.page_size.toString());

        const query = searchParams.toString();
        return request<{ items: any[]; total: number; page: number }>(
            `/api/items${query ? `?${query}` : ''}`
        );
    },

    get: (id: string) =>
        request<any>(`/api/items?id=${id}`),

    create: (data: any) =>
        request('/api/items', {
            method: 'POST',
            body: JSON.stringify(data),
        }),

    update: (id: string, data: any) =>
        request(`/api/items?id=${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        }),

    delete: (id: string, ownerId: string) =>
        request(`/api/items?id=${id}&owner_id=${ownerId}`, {
            method: 'DELETE',
        }),

    getOwnerItems: (ownerId: string, page?: number, pageSize?: number) => {
        const searchParams = new URLSearchParams({ owner_id: ownerId });
        if (page) searchParams.set('page', page.toString());
        if (pageSize) searchParams.set('page_size', pageSize.toString());
        return request<{ items: any[]; total: number }>(`/api/items/owner?${searchParams}`);
    },

    search: (params?: { category?: string; city?: string; page?: number }) => {
        const searchParams = new URLSearchParams();
        if (params?.category) searchParams.set('category', params.category);
        if (params?.city) searchParams.set('city', params.city);
        if (params?.page) searchParams.set('page', params.page.toString());
        return request<{ items: any[]; total: number }>(`/api/items/search?${searchParams}`);
    },
};

// ========== Bookings API ==========
export const bookingsApi = {
    create: (data: {
        renter_id: string;
        owner_id: string;
        rental_item_id: string;
        start_date: string;
        end_date: string;
        daily_rate: number;
        security_deposit: number;
    }) => request('/api/bookings', {
        method: 'POST',
        body: JSON.stringify(data),
    }),

    get: (id: string) =>
        request<any>(`/api/bookings?id=${id}`),

    getRenterBookings: (renterId: string, page?: number, pageSize?: number) => {
        const params = new URLSearchParams({ renter_id: renterId });
        if (page) params.set('page', page.toString());
        if (pageSize) params.set('page_size', pageSize.toString());
        return request<{ bookings: any[]; total: number }>(`/api/bookings/renter?${params}`);
    },

    getOwnerBookings: (ownerId: string, page?: number, pageSize?: number) => {
        const params = new URLSearchParams({ owner_id: ownerId });
        if (page) params.set('page', page.toString());
        if (pageSize) params.set('page_size', pageSize.toString());
        return request<{ bookings: any[]; total: number }>(`/api/bookings/owner?${params}`);
    },

    confirm: (bookingId: string, ownerId: string) =>
        request('/api/bookings/confirm', {
            method: 'POST',
            body: JSON.stringify({ booking_id: bookingId, owner_id: ownerId }),
        }),

    cancel: (bookingId: string, userId: string, reason?: string) =>
        request('/api/bookings/cancel', {
            method: 'POST',
            body: JSON.stringify({ booking_id: bookingId, user_id: userId, reason: reason || '' }),
        }),
};

// ========== Payments API ==========
export const paymentsApi = {
    initialize: (data: {
        booking_id: string;
        user_id: string;
        amount: number;
        method: string;
    }) => request<{ payment_id: string; checkout_url: string; status: string }>(
        '/api/payments/initialize',
        {
            method: 'POST',
            body: JSON.stringify(data),
        }
    ),

    get: (id: string) =>
        request<any>(`/api/payments?id=${id}`),

    getBookingPayments: (bookingId: string) =>
        request<{ payments: any[]; count: number }>(`/api/payments/booking?booking_id=${bookingId}`),

    refund: (paymentId: string, amount: number) =>
        request('/api/payments/refund', {
            method: 'POST',
            body: JSON.stringify({ payment_id: paymentId, amount }),
        }),

    verify: (txRef: string) =>
        request<{ tx_ref: string; reference: string; amount: number; status: string; email: string }>(`/api/payments/verify?tx_ref=${txRef}`),
};

// ========== Reviews API ==========
export const reviewsApi = {
    create: (data: {
        item_id: string;
        booking_id?: string;
        reviewer_id: string;
        rating: number;
        comment: string;
        review_type?: string;
    }) => request('/api/reviews', {
        method: 'POST',
        body: JSON.stringify(data),
    }),

    get: (id: string) =>
        request<any>(`/api/reviews?id=${id}`),

    update: (reviewId: string, rating: number, comment: string) =>
        request('/api/reviews', {
            method: 'PUT',
            body: JSON.stringify({ review_id: reviewId, rating, comment }),
        }),

    delete: (id: string) =>
        request(`/api/reviews?id=${id}`, { method: 'DELETE' }),

    getItemReviews: (itemId: string, page?: number, pageSize?: number) => {
        const params = new URLSearchParams({ item_id: itemId });
        if (page) params.set('page', page.toString());
        if (pageSize) params.set('page_size', pageSize.toString());
        return request<{ reviews: any[]; total: number }>(`/api/reviews/item?${params}`);
    },

    getUserReviews: (userId: string, page?: number, pageSize?: number) => {
        const params = new URLSearchParams({ user_id: userId });
        if (page) params.set('page', page.toString());
        if (pageSize) params.set('page_size', pageSize.toString());
        return request<{ reviews: any[]; total: number }>(`/api/reviews/user?${params}`);
    },
};

// ========== Notifications API ==========
export const notificationsApi = {
    getUserNotifications: (userId: string, page?: number, pageSize?: number) => {
        const params = new URLSearchParams({ user_id: userId });
        if (page) params.set('page', page.toString());
        if (pageSize) params.set('page_size', pageSize.toString());
        return request<{ notifications: any[]; total: number }>(`/api/notifications/user?${params}`);
    },

    markAsRead: (notificationId: string) =>
        request('/api/notifications/mark-read', {
            method: 'POST',
            body: JSON.stringify({ notification_id: notificationId }),
        }),

    getUnreadCount: (userId: string) =>
        request<{ count: number }>(`/api/notifications/unread-count?user_id=${userId}`),
};

// ========== Messages API ==========
export const messagesApi = {
    send: (data: {
        booking_id: string;
        sender_id: string;
        receiver_id: string;
        content: string;
    }) => request('/api/messages/send', {
        method: 'POST',
        body: JSON.stringify(data),
    }),

    getBookingMessages: (bookingId: string, page?: number, pageSize?: number) => {
        const params = new URLSearchParams({ booking_id: bookingId });
        if (page) params.set('page', page.toString());
        if (pageSize) params.set('page_size', pageSize.toString());
        return request<{ messages: any[]; total: number }>(`/api/messages/booking?${params}`);
    },
};
