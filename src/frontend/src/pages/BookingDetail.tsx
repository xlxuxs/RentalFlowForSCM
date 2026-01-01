import { useState, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import { bookingsApi, itemsApi, paymentsApi, reviewsApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import './BookingDetail.css';

export function BookingDetailPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const { user } = useAuth();

    const [booking, setBooking] = useState<any>(null);
    const [item, setItem] = useState<any>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState('');
    const [error, setError] = useState('');
    const [showReviewForm, setShowReviewForm] = useState(false);
    const [reviewData, setReviewData] = useState({ rating: 5, comment: '' });

    useEffect(() => {
        const loadBooking = async () => {
            if (!id) return;

            try {
                const bookingData = await bookingsApi.get(id);
                setBooking(bookingData);

                // Load the item details
                if (bookingData.rental_item_id) {
                    const itemData = await itemsApi.get(bookingData.rental_item_id);
                    setItem(itemData);
                }
            } catch (err: any) {
                setError(err.message || 'Failed to load booking');
            } finally {
                setIsLoading(false);
            }
        };

        loadBooking();
    }, [id]);

    const handleConfirm = async () => {
        if (!booking || !user) return;
        setActionLoading('confirm');
        try {
            await bookingsApi.confirm(booking.id, user.id);
            setBooking((prev: any) => ({ ...prev, status: 'confirmed' }));
        } catch (err: any) {
            setError(err.message);
        } finally {
            setActionLoading('');
        }
    };

    const handleCancel = async () => {
        if (!booking || !user || !window.confirm('Are you sure you want to cancel this booking?')) return;
        setActionLoading('cancel');
        try {
            await bookingsApi.cancel(booking.id, user.id, 'Cancelled by user');
            setBooking((prev: any) => ({ ...prev, status: 'cancelled' }));
        } catch (err: any) {
            setError(err.message);
        } finally {
            setActionLoading('');
        }
    };

    const handlePayment = async () => {
        if (!booking || !user) return;
        setActionLoading('payment');
        try {
            const result = await paymentsApi.initialize({
                booking_id: booking.id,
                user_id: user.id,
                amount: booking.total_amount,
                method: 'chapa',
            });
            if (result.checkout_url) {
                window.open(result.checkout_url, '_blank');
            }
        } catch (err: any) {
            setError(err.message);
        } finally {
            setActionLoading('');
        }
    };

    const handleReviewSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!booking || !user) return;
        setActionLoading('review');
        try {
            await reviewsApi.create({
                item_id: booking.rental_item_id,
                booking_id: booking.id,
                reviewer_id: user.id,
                rating: reviewData.rating,
                comment: reviewData.comment,
            });
            setShowReviewForm(false);
            setReviewData({ rating: 5, comment: '' });
        } catch (err: any) {
            setError(err.message);
        } finally {
            setActionLoading('');
        }
    };

    const formatDate = (dateStr: string) => {
        return new Date(dateStr).toLocaleDateString('en-US', {
            weekday: 'short',
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        });
    };

    const statusColors: Record<string, string> = {
        pending: 'badge-warning',
        confirmed: 'badge-primary',
        active: 'badge-success',
        completed: 'badge-gray',
        cancelled: 'badge-error',
    };

    const isOwner = user?.id === booking?.owner_id;
    const isRenter = user?.id === booking?.renter_id;
    const isPaymentSuccess = searchParams.get('payment_status') === 'success';

    if (isLoading) {
        return (
            <div className="page-loading">
                <div className="spinner"></div>
            </div>
        );
    }

    if (error || !booking) {
        return (
            <div className="container p-8">
                <div className="empty-state">
                    <div className="empty-state-icon">😕</div>
                    <h3 className="empty-state-title">Booking not found</h3>
                    <p className="empty-state-text">{error || 'This booking may have been removed'}</p>
                    <button className="btn btn-primary" onClick={() => navigate(-1)}>
                        Go Back
                    </button>
                </div>
            </div>
        );
    }

    const isPaid = booking.payment_status === 'success' || booking.payment_status === 'completed';

    return (
        <div className="booking-detail-page">
            <div className="container">
                <button className="back-button" onClick={() => navigate(-1)}>
                    ← Back
                </button>

                {isPaymentSuccess && (
                    <div className="alert alert-success m-b-6">
                        <span className="alert-icon">✅</span>
                        <div className="alert-content">
                            <h4>Payment Successful!</h4>
                            <p>Your payment has been verified. You can now coordinate with the owner for pickup.</p>
                        </div>
                    </div>
                )}

                <div className="booking-detail-grid">
                    {/* Main Details */}
                    <div className="booking-main">
                        <div className="booking-header">
                            <div>
                                <span className="booking-number">Booking #{booking.booking_number}</span>
                                <h1>Booking Details</h1>
                            </div>
                            <span className={`badge ${statusColors[booking.status] || 'badge-gray'}`}>
                                {booking.status}
                            </span>
                        </div>

                        {/* Item Info */}
                        {item && (
                            <div className="booking-item-card" onClick={() => navigate(`/items/${item.id}`)}>
                                <div className="booking-item-image">
                                    {item.images && item.images.length > 0 ? (
                                        <img src={item.images[0]} alt={item.title} />
                                    ) : (
                                        <span className="placeholder-icon">📦</span>
                                    )}
                                </div>
                                <div className="booking-item-info">
                                    <h3>{item.title}</h3>
                                    <p>📍 {item.city}</p>
                                </div>
                            </div>
                        )}

                        {/* Dates */}
                        <div className="booking-section">
                            <h3>Rental Period</h3>
                            <div className="date-range">
                                <div className="date-box">
                                    <span className="date-label">Start Date</span>
                                    <span className="date-value">{formatDate(booking.start_date)}</span>
                                </div>
                                <span className="date-arrow">→</span>
                                <div className="date-box">
                                    <span className="date-label">End Date</span>
                                    <span className="date-value">{formatDate(booking.end_date)}</span>
                                </div>
                            </div>
                            <p className="date-duration">{booking.total_days} day{booking.total_days > 1 ? 's' : ''}</p>
                        </div>

                        {/* Pricing Breakdown */}
                        <div className="booking-section">
                            <h3>Pricing Breakdown</h3>
                            <div className="price-breakdown">
                                <div className="price-row">
                                    <span>${booking.daily_rate} × {booking.total_days} days</span>
                                    <span>${(booking.daily_rate * booking.total_days).toFixed(2)}</span>
                                </div>
                                {booking.service_fee > 0 && (
                                    <div className="price-row">
                                        <span>Service Fee</span>
                                        <span>${booking.service_fee?.toFixed(2)}</span>
                                    </div>
                                )}
                                {booking.security_deposit > 0 && (
                                    <div className="price-row">
                                        <span>Security Deposit</span>
                                        <span>${booking.security_deposit?.toFixed(2)}</span>
                                    </div>
                                )}
                                <div className="price-total">
                                    <span>Total</span>
                                    <span>${booking.total_amount?.toFixed(2)}</span>
                                </div>
                            </div>
                        </div>

                        {/* Review Form */}
                        {showReviewForm && (
                            <div className="booking-section">
                                <h3>Leave a Review</h3>
                                <form onSubmit={handleReviewSubmit} className="review-form">
                                    <div className="rating-select">
                                        {[1, 2, 3, 4, 5].map(star => (
                                            <button
                                                key={star}
                                                type="button"
                                                className={`star-btn ${reviewData.rating >= star ? 'active' : ''}`}
                                                onClick={() => setReviewData(prev => ({ ...prev, rating: star }))}
                                            >
                                                ⭐
                                            </button>
                                        ))}
                                    </div>
                                    <textarea
                                        className="form-input"
                                        placeholder="Share your experience..."
                                        rows={3}
                                        value={reviewData.comment}
                                        onChange={e => setReviewData(prev => ({ ...prev, comment: e.target.value }))}
                                        required
                                    />
                                    <div className="review-actions">
                                        <button type="button" className="btn btn-ghost" onClick={() => setShowReviewForm(false)}>
                                            Cancel
                                        </button>
                                        <button type="submit" className="btn btn-primary" disabled={actionLoading === 'review'}>
                                            {actionLoading === 'review' ? 'Submitting...' : 'Submit Review'}
                                        </button>
                                    </div>
                                </form>
                            </div>
                        )}
                    </div>

                    {/* Actions Sidebar */}
                    <div className="booking-actions-card">
                        <h3>Actions</h3>

                        {error && <div className="alert alert-error">{error}</div>}

                        {/* Owner Actions */}
                        {isOwner && booking.status === 'pending' && (
                            <>
                                <button
                                    className="btn btn-primary"
                                    style={{ width: '100%' }}
                                    onClick={handleConfirm}
                                    disabled={actionLoading === 'confirm'}
                                >
                                    {actionLoading === 'confirm' ? 'Confirming...' : '✓ Confirm Booking'}
                                </button>
                                <button
                                    className="btn btn-outline"
                                    style={{ width: '100%', marginTop: 'var(--spacing-3)' }}
                                    onClick={handleCancel}
                                    disabled={actionLoading === 'cancel'}
                                >
                                    {actionLoading === 'cancel' ? 'Cancelling...' : '✕ Reject Booking'}
                                </button>
                            </>
                        )}

                        {/* Renter Actions */}
                        {isRenter && booking.status === 'confirmed' && !isPaid && (
                            <button
                                className="btn btn-primary"
                                style={{ width: '100%' }}
                                onClick={handlePayment}
                                disabled={actionLoading === 'payment'}
                            >
                                {actionLoading === 'payment' ? 'Processing...' : '💳 Pay Now'}
                            </button>
                        )}
                        {isRenter && isPaid && booking.status === 'confirmed' && (
                            <div className="payment-status-badge success">
                                <span>💳 Payment Completed</span>
                            </div>
                        )}

                        {isRenter && booking.status === 'pending' && (
                            <button
                                className="btn btn-outline text-error"
                                style={{ width: '100%' }}
                                onClick={handleCancel}
                                disabled={actionLoading === 'cancel'}
                            >
                                {actionLoading === 'cancel' ? 'Cancelling...' : '✕ Cancel Booking'}
                            </button>
                        )}

                        {isRenter && booking.status === 'completed' && !showReviewForm && (
                            <button
                                className="btn btn-primary"
                                style={{ width: '100%' }}
                                onClick={() => setShowReviewForm(true)}
                            >
                                ⭐ Leave a Review
                            </button>
                        )}

                        {/* Status Info */}
                        <div className="status-info">
                            {booking.status === 'pending' && (
                                <p>⏳ Waiting for {isOwner ? 'your confirmation' : 'owner confirmation'}</p>
                            )}
                            {booking.status === 'confirmed' && (
                                <p>✅ Booking confirmed. {isRenter && 'Complete payment to proceed.'}</p>
                            )}
                            {booking.status === 'active' && (
                                <p>🚀 Rental is currently active</p>
                            )}
                            {booking.status === 'completed' && (
                                <p>✓ This rental has been completed</p>
                            )}
                            {booking.status === 'cancelled' && (
                                <p>❌ This booking was cancelled</p>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
