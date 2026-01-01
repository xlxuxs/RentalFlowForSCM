import { getThumbnailUrl } from '../../services/cloudinary';
import './BookingCard.css';

interface BookingCardProps {
    id: string;
    booking_number: string;
    status: string;
    start_date: string;
    end_date: string;
    total_amount: number;
    item_title?: string;
    item_image?: string;
    onViewDetails?: () => void;
    onConfirm?: () => void;
    onCancel?: () => void;
    showActions?: boolean;
    isOwnerView?: boolean;
}

export function BookingCard({
    booking_number,
    status,
    start_date,
    end_date,
    total_amount,
    item_title,
    item_image,
    onViewDetails,
    onConfirm,
    onCancel,
    showActions = true,
    isOwnerView = false,
}: BookingCardProps) {
    const statusColors: Record<string, string> = {
        pending: 'badge-warning',
        confirmed: 'badge-primary',
        active: 'badge-success',
        completed: 'badge-gray',
        cancelled: 'badge-error',
    };

    const formatDate = (dateStr: string) => {
        if (!dateStr) return 'N/A';
        const date = new Date(dateStr);
        if (isNaN(date.getTime())) return 'Invalid';
        return date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        });
    };

    const imageUrl = item_image ? getThumbnailUrl(item_image) : null;

    return (
        <div className="booking-card">
            <div className="booking-card-header">
                {imageUrl && (
                    <div className="booking-card-image">
                        <img src={imageUrl} alt={item_title || 'Item'} />
                    </div>
                )}
                <div>
                    <span className="booking-number">#{booking_number}</span>
                    <span className={`badge ${statusColors[status] || 'badge-gray'}`}>
                        {status}
                    </span>
                </div>
            </div>

            {item_title && (
                <h4 className="booking-item-title">{item_title}</h4>
            )}

            <div className="booking-dates">
                <div className="booking-date">
                    <span className="date-label">Start</span>
                    <span className="date-value">{formatDate(start_date)}</span>
                </div>
                <span className="date-separator">→</span>
                <div className="booking-date">
                    <span className="date-label">End</span>
                    <span className="date-value">{formatDate(end_date)}</span>
                </div>
            </div>

            <div className="booking-total">
                <span className="total-label">Total Amount</span>
                <span className="total-value">${(total_amount || 0).toFixed(2)}</span>
            </div>

            {showActions && (
                <div className="booking-actions">
                    {onViewDetails && (
                        <button className="btn btn-primary btn-sm" onClick={onViewDetails}>
                            View Details
                        </button>
                    )}
                    {isOwnerView && status === 'pending' && onConfirm && (
                        <button className="btn btn-primary btn-sm" onClick={onConfirm}>
                            Confirm
                        </button>
                    )}
                    {status === 'pending' && onCancel && (
                        <button className="btn btn-ghost btn-sm text-error" onClick={onCancel}>
                            Cancel
                        </button>
                    )}
                </div>
            )}
        </div>
    );
}
