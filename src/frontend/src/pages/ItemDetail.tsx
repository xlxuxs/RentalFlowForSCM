import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { itemsApi, bookingsApi, reviewsApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { RatingStars } from '../components/features';
import './ItemDetail.css';

export function ItemDetailPage() {

    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { user, isAuthenticated } = useAuth();

    const [item, setItem] = useState<any>(null);
    const [reviews, setReviews] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [isBooking, setIsBooking] = useState(false);
    const [bookingError, setBookingError] = useState('');
    const [bookingSuccess, setBookingSuccess] = useState(false);

    const [startDate, setStartDate] = useState('');
    const [endDate, setEndDate] = useState('');

    const [reviewRating, setReviewRating] = useState(5);
    const [reviewComment, setReviewComment] = useState('');
    const [isSubmittingReview, setIsSubmittingReview] = useState(false);
    const [reviewError, setReviewError] = useState('');
    const [reviewSuccess, setReviewSuccess] = useState(false);

    const isValidBookingDates = () => {
        if (!startDate || !endDate) return false;
        const start = new Date(startDate);
        const end = new Date(endDate);
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        return start >= today && end > start;
    };

    const categoryEmoji: Record<string, string> = {
        vehicle: '🚗',
        equipment: '🔧',
        property: '🏢',
        other: '📦'
    };

    useEffect(() => {
        const loadItem = async () => {
            if (!id) return;

            try {
                const [itemData, reviewsData] = await Promise.all([
                    itemsApi.get(id),
                    reviewsApi.getItemReviews(id, 1, 5).catch(() => ({ reviews: [], total: 0 })),
                ]);
                setItem(itemData);
                setReviews(reviewsData.reviews || []);
            } catch (error) {
                console.error('Failed to load item:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadItem();
    }, [id]);

    const handleBooking = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!isAuthenticated) {
            navigate('/login', { state: { from: `/items/${id}` } });
            return;
        }

        if (!startDate || !endDate || !item || !user) {
            setBookingError('Please select start and end dates');
            return;
        }

        setIsBooking(true);
        setBookingError('');

        try {
            // Convert date strings to ISO format with time
            const start = new Date(startDate);
            const end = new Date(endDate);
            start.setHours(0, 0, 0, 0);
            end.setHours(23, 59, 59, 999);

            await bookingsApi.create({
                renter_id: user.id,
                owner_id: item.owner_id,
                rental_item_id: item.id,
                start_date: start.toISOString(),
                end_date: end.toISOString(),
                daily_rate: item.daily_rate,
                security_deposit: item.security_deposit || 0,
            });
            setBookingSuccess(true);
        } catch (error: any) {
            setBookingError(error.message || 'Booking failed');
        } finally {
            setIsBooking(false);
        }
    };

    const calculateTotal = () => {
        if (!startDate || !endDate || !item) return 0;
        const start = new Date(startDate);
        const end = new Date(endDate);
        const days = Math.max(1, Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)));
        const subtotal = days * item.daily_rate;
        const serviceFee = subtotal * 0.1;
        return subtotal + serviceFee + (item.security_deposit || 0);
    };

    const handleReviewSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!isAuthenticated || !user || !id) return;

        setIsSubmittingReview(true);
        setReviewError('');
        setReviewSuccess(false);

        try {
            await reviewsApi.create({
                item_id: id,
                reviewer_id: user.id,
                rating: reviewRating,
                comment: reviewComment,
                review_type: 'renter_to_item',
            });
            setReviewSuccess(true);
            setReviewComment('');
            setReviewRating(5);
            // Refresh reviews
            const reviewsData = await reviewsApi.getItemReviews(id, 1, 5);
            setReviews(reviewsData.reviews || []);
        } catch (error: any) {
            setReviewError(error.message || 'Failed to submit review');
        } finally {
            setIsSubmittingReview(false);
        }
    };



    if (isLoading) {
        return (
            <div className="page-loading">
                <div className="spinner"></div>
            </div>
        );
    }

    if (!item) {
        return (
            <div className="container p-8">
                <div className="empty-state">
                    <div className="empty-state-icon">😕</div>
                    <h3 className="empty-state-title">Item not found</h3>
                    <p className="empty-state-text">This item may have been removed</p>
                    <button className="btn btn-primary" onClick={() => navigate('/browse')}>
                        Browse Items
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="item-detail-page">
            <div className="container">
                <div className="item-detail-grid">
                    {/* Images */}
                    <div className="item-images">
                        <div className="main-image">
                            {item.images && item.images.length > 0 ? (
                                <img src={item.images[0]} alt={item.title} className="main-image" />
                            ) : (
                                <div className="no-image">📦 No image available</div>
                            )}
                        </div>

                        {/* Item Info */}
                        <div className="item-info">
                            <div className="item-header">
                                <h1>{item.title}</h1>
                                <div className="item-rating">
                                    {item.average_rating > 0 && (
                                        <RatingStars rating={item.average_rating || 0} size="md" />
                                    )}
                                </div>
                            </div>

                            <p className="item-category">
                                {categoryEmoji[item.category.toLowerCase()] || '📦'} {item.category}
                            </p>

                            <p className="item-description">{item.description}</p>

                            <div className="item-pricing">
                                <div className="price-item">
                                    <span className="price-label">Daily Rate:</span>
                                    <span className="price-value">{item.daily_rate} ETB</span>
                                </div>
                                {item.weekly_rate && (
                                    <div className="price-item">
                                        <span className="price-label">Weekly Rate:</span>
                                        <span className="price-value">{item.weekly_rate} ETB</span>
                                    </div>
                                )}
                                {item.monthly_rate && (
                                    <div className="price-item">
                                        <span className="price-label">Monthly Rate:</span>
                                        <span className="price-value">{item.monthly_rate} ETB</span>
                                    </div>
                                )}
                                <div className="price-item">
                                    <span className="price-label">Security Deposit:</span>
                                    <span className="price-value">{item.security_deposit} ETB</span>
                                </div>
                            </div>

                            <div className="item-location">
                                <strong>📍 Location:</strong> {item.city}
                                {item.address && `, ${item.address}`}
                            </div>

                            {item.specifications && Object.keys(item.specifications).length > 0 && (
                                <div className="item-specifications">
                                    <h3>Specifications</h3>
                                    <ul>
                                        {Object.entries(item.specifications).map(([key, value]) => (
                                            <li key={key}>
                                                <strong>{key}:</strong> {value as string}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            )}

                            {/* Booking Form */}
                            {user?.id !== item.owner_id && item.is_active && (
                                <form onSubmit={handleBooking} className="booking-form">
                                    <h3>Book this item</h3>

                                    {bookingError && <div className="alert alert-error">{bookingError}</div>}
                                    {bookingSuccess && (
                                        <div className="alert alert-success">
                                            Booking created successfully! Redirecting to payment...
                                        </div>
                                    )}

                                    <div className="form-group">
                                        <label>Start Date *</label>
                                        <input
                                            type="date"
                                            value={startDate}
                                            onChange={(e) => setStartDate(e.target.value)}
                                            required
                                            min={new Date().toISOString().split('T')[0]}
                                        />
                                    </div>

                                    <div className="form-group">
                                        <label>End Date *</label>
                                        <input
                                            type="date"
                                            value={endDate}
                                            onChange={(e) => setEndDate(e.target.value)}
                                            required
                                            min={startDate || new Date().toISOString().split('T')[0]}
                                        />
                                    </div>

                                    {isValidBookingDates() && (
                                        <div className="booking-summary">
                                            <div className="summary-row">
                                                <span>Daily Rate</span>
                                                <span>${item.daily_rate}</span>
                                            </div>
                                            <div className="summary-row">
                                                <span>Service Fee (10%)</span>
                                                <span>${(calculateTotal() * 0.1).toFixed(2)}</span>
                                            </div>
                                            {item.security_deposit > 0 && (
                                                <div className="summary-row">
                                                    <span>Security Deposit</span>
                                                    <span>${item.security_deposit}</span>
                                                </div>
                                            )}
                                            <div className="summary-total">
                                                <span>Total</span>
                                                <span>${calculateTotal().toFixed(2)}</span>
                                            </div>
                                        </div>
                                    )}

                                    {bookingError && (
                                        <div className="alert alert-error">{bookingError}</div>
                                    )}

                                    <button
                                        type="submit"
                                        className="btn btn-primary btn-lg"
                                        style={{ width: '100%' }}
                                        disabled={isBooking || !item.is_active || (user?.id === item.owner_id)}
                                    >
                                        {isBooking ? 'Booking...' :
                                            !item.is_active ? 'Not Available' :
                                                user?.id === item.owner_id ? 'This is your item' :
                                                    isAuthenticated ? 'Request Booking' : 'Login to Book'}
                                    </button>
                                </form>
                            )}
                        </div>
                    </div>

                    {/* Reviews Section */}
                    <div className="reviews-section">
                        <h2>Reviews</h2>
                        {reviews.length > 0 ? (
                            <div className="reviews-list">
                                {reviews.map((review: any) => (
                                    <div key={review.id} className="review-card">
                                        <div className="review-header">
                                            <div className="review-rating">
                                                {'⭐'.repeat(Math.round(review.rating))}
                                            </div>
                                            <span className="review-date">
                                                {new Date(review.created_at).toLocaleDateString()}
                                            </span>
                                        </div>
                                        <p className="review-comment">{review.comment}</p>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <p className="no-reviews">No reviews yet</p>
                        )}

                        {isAuthenticated && user?.id !== item.owner_id && (
                            <div className="add-review-section">
                                <h3>Leave a Review</h3>
                                <form onSubmit={handleReviewSubmit} className="review-form">
                                    {reviewError && <div className="alert alert-error">{reviewError}</div>}
                                    {reviewSuccess && <div className="alert alert-success">Review submitted successfully!</div>}

                                    <div className="form-group">
                                        <label>Rating</label>
                                        <div className="rating-input">
                                            {[1, 2, 3, 4, 5].map((star) => (
                                                <button
                                                    key={star}
                                                    type="button"
                                                    className={`star-btn ${reviewRating >= star ? 'active' : ''}`}
                                                    onClick={() => setReviewRating(star)}
                                                >
                                                    ⭐
                                                </button>
                                            ))}
                                        </div>
                                    </div>

                                    <div className="form-group">
                                        <label>Comment</label>
                                        <textarea
                                            value={reviewComment}
                                            onChange={(e) => setReviewComment(e.target.value)}
                                            placeholder="Tell others about your experience..."
                                            required
                                            rows={4}
                                        ></textarea>
                                    </div>

                                    <button
                                        type="submit"
                                        className="btn btn-primary"
                                        disabled={isSubmittingReview}
                                    >
                                        {isSubmittingReview ? 'Submitting...' : 'Submit Review'}
                                    </button>
                                </form>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}
