import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { paymentsApi } from '../services/api';
import './PaymentCallback.css';

export function PaymentCallbackPage() {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const [status, setStatus] = useState<'loading' | 'success' | 'failed'>('loading');
    const [message, setMessage] = useState('');
    const [bookingId, setBookingId] = useState<string | null>(null);

    useEffect(() => {
        const verifyPayment = async () => {
            const txRef = searchParams.get('tx_ref');
            const bId = searchParams.get('booking_id');
            setBookingId(bId);

            if (!txRef) {
                setStatus('failed');
                setMessage('Invalid payment reference');
                return;
            }

            try {
                const result = await paymentsApi.verify(txRef);

                if (result.status === 'success') {
                    setStatus('success');
                    setMessage('Your payment has been successfully processed and verified.');

                    // Redirect to booking details after 5 seconds
                    setTimeout(() => {
                        if (bId) {
                            navigate(`/bookings/${bId}?payment_status=success`);
                        } else {
                            navigate('/my-bookings');
                        }
                    }, 5000);
                } else {
                    setStatus('failed');
                    setMessage('We couldn\'t verify your payment. If you believe this is an error, please contact support.');
                }
            } catch (error: any) {
                setStatus('failed');
                setMessage(error.message || 'Payment verification failed');
            }
        };

        verifyPayment();
    }, [searchParams, navigate]);

    return (
        <div className="payment-callback-page">
            <div className="callback-container">
                {status === 'loading' && (
                    <div className="callback-loading">
                        <div className="spinner"></div>
                        <h2>Verifying Payment...</h2>
                        <p>Please wait while we confirm your payment</p>
                    </div>
                )}

                {status === 'success' && (
                    <div className="callback-success">
                        <div className="success-icon">✅</div>
                        <h2>Payment Successful!</h2>
                        <p>{message}</p>
                        <div className="callback-actions">
                            <button
                                className="btn btn-primary"
                                onClick={() => navigate(bookingId ? `/bookings/${bookingId}?payment_status=success` : '/my-bookings')}
                            >
                                View Booking Now
                            </button>
                        </div>
                        <p className="redirect-text">Redirecting automatically in a few seconds...</p>
                    </div>
                )}

                {status === 'failed' && (
                    <div className="callback-failed">
                        <div className="error-icon">❌</div>
                        <h2>Payment Failed</h2>
                        <p>{message}</p>
                        <button
                            className="btn btn-primary"
                            onClick={() => navigate('/my-bookings')}
                        >
                            Go to My Bookings
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}
