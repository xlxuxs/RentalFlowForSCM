// TypeScript types for RentalFlow frontend

// ========== User & Auth Types ==========
export interface User {
    id: string;
    email: string;
    first_name: string;
    last_name: string;
    phone?: string;
    bio?: string;
    avatar_url?: string;
    role: 'renter' | 'owner' | 'admin';
    identity_verified?: boolean;
    verification_status?: 'pending' | 'verified' | 'rejected';
    created_at?: string;
}

export interface AuthResponse {
    user: User;
    access_token: string;
    refresh_token: string;
    expires_in: number;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    phone?: string;
    role: 'renter' | 'owner';
}

// ========== Rental Item Types ==========
export type ItemCategory = 'vehicle' | 'equipment' | 'property';

export interface RentalItem {
    id: string;
    owner_id: string;
    title: string;
    description: string;
    category: ItemCategory;
    subcategory?: string;
    daily_rate: number;
    weekly_rate?: number;
    monthly_rate?: number;
    security_deposit: number;
    address?: string;
    city: string;
    latitude?: number;
    longitude?: number;
    specifications?: Record<string, string>;
    images: string[];
    is_active: boolean;
    created_at: string;
}

export interface CreateItemRequest {
    owner_id: string;
    title: string;
    description: string;
    category: ItemCategory;
    subcategory?: string;
    daily_rate: number;
    weekly_rate?: number;
    monthly_rate?: number;
    security_deposit: number;
    address?: string;
    city: string;
    latitude?: number;
    longitude?: number;
    specifications?: Record<string, string>;
    images?: string[];
}

export interface ItemFilters {
    category?: ItemCategory;
    city?: string;
    page?: number;
    page_size?: number;
}

export interface ItemsResponse {
    items: RentalItem[];
    total: number;
    page: number;
}

// ========== Booking Types ==========
export type BookingStatus = 'pending' | 'confirmed' | 'active' | 'completed' | 'cancelled';

export interface Booking {
    id: string;
    booking_number: string;
    renter_id: string;
    owner_id: string;
    rental_item_id: string;
    status: BookingStatus;
    start_date: string;
    end_date: string;
    total_days: number;
    daily_rate: number;
    subtotal?: number;
    security_deposit: number;
    service_fee?: number;
    total_amount: number;
    pickup_address?: string;
    pickup_notes?: string;
    return_address?: string;
    return_notes?: string;
    agreement_signed: boolean;
    cancellation_reason?: string;
    payment_status?: string;
    created_at?: string;
}

export interface CreateBookingRequest {
    renter_id: string;
    owner_id: string;
    rental_item_id: string;
    start_date: string;
    end_date: string;
    daily_rate: number;
    security_deposit: number;
}

export interface BookingsResponse {
    bookings: Booking[];
    total: number;
}

// ========== Payment Types ==========
export type PaymentMethod = 'chapa' | 'cash';
export type PaymentStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'refunded';

export interface Payment {
    id: string;
    booking_id: string;
    user_id: string;
    amount: number;
    method: PaymentMethod;
    status: PaymentStatus;
    checkout_url?: string;
    transaction_id?: string;
    created_at: string;
}

export interface InitializePaymentRequest {
    booking_id: string;
    user_id: string;
    amount: number;
    method: PaymentMethod;
}

// ========== Review Types ==========
export type ReviewType = 'item_review' | 'renter_review' | 'owner_review';

export interface Review {
    id: string;
    booking_id: string;
    reviewer_id: string;
    review_type: ReviewType;
    rating: number;
    comment: string;
    created_at: string;
}

export interface CreateReviewRequest {
    booking_id: string;
    reviewer_id: string;
    review_type: ReviewType;
    rating: number;
    comment: string;
}

export interface ReviewsResponse {
    reviews: Review[];
    total: number;
}

// ========== Notification Types ==========
export type NotificationChannel = 'push' | 'email' | 'sms';

export interface Notification {
    id: string;
    user_id: string;
    type: string;
    title: string;
    message: string;
    channel: NotificationChannel;
    is_read: boolean;
    created_at: string;
}

export interface NotificationsResponse {
    notifications: Notification[];
    total: number;
}

// ========== Message Types ==========
export interface Message {
    id: string;
    booking_id: string;
    sender_id: string;
    receiver_id: string;
    content: string;
    is_read: boolean;
    created_at: string;
}

export interface MessagesResponse {
    messages: Message[];
    total: number;
}

// ========== API Response Types ==========
export interface ApiError {
    error: string;
}

export interface SuccessResponse {
    success: boolean;
}
