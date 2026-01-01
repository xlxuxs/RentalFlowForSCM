import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { bookingsApi, itemsApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { BookingCard } from '../components/features';
import './Dashboard.css';

export function MyBookingsPage() {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [bookings, setBookings] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadBookings = async () => {
            if (!user) return;

            try {
                const result = await bookingsApi.getRenterBookings(user.id, 1, 20);
                const bookingsList = result.bookings || [];

                // Fetch item details for each booking to get images
                const enrichedBookings = await Promise.all(
                    bookingsList.map(async (booking) => {
                        try {
                            const itemId = booking.RentalItemID || booking.rental_item_id;
                            if (itemId) {
                                const item = await itemsApi.get(itemId);
                                return { ...booking, rental_item: item };
                            }
                        } catch (e) {
                            console.error(`Failed to load item for booking ${booking.ID || booking.id}`, e);
                        }
                        return booking;
                    })
                );

                setBookings(enrichedBookings);
            } catch (error) {
                console.error('Failed to load bookings:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadBookings();
    }, [user]);

    const handleCancel = async (bookingId: string) => {
        if (!user || !window.confirm('Are you sure you want to cancel this booking?')) return;

        try {
            await bookingsApi.cancel(bookingId, user.id, 'Cancelled by renter');
            setBookings(prev =>
                prev.map(b => ((b.ID || b.id) === bookingId ? { ...b, Status: 'cancelled', status: 'cancelled' } : b))
            );
        } catch (error) {
            console.error('Failed to cancel booking:', error);
        }
    };

    const statusTabs = ['pending', 'confirmed', 'active', 'completed', 'cancelled'];

    return (
        <div className="dashboard-page">
            <div className="dashboard-header">
                <div className="container">
                    <h1>My Bookings</h1>
                    <p>Track and manage your rental bookings in Kanban view</p>
                </div>
            </div>

            <div className="kanban-container">
                {isLoading ? (
                    <div className="page-loading">
                        <div className="spinner"></div>
                    </div>
                ) : (
                    <div className="kanban-board">
                        {statusTabs.map(status => (
                            <div key={status} className="kanban-column">
                                <div className="kanban-column-header">
                                    <h3 className="status-title">{status}</h3>
                                    <span className="count-badge">
                                        {bookings.filter(b => (b.Status || b.status) === status).length}
                                    </span>
                                </div>
                                <div className="kanban-column-content">
                                    {bookings
                                        .filter(b => (b.Status || b.status) === status)
                                        .map((booking, index) => (
                                            <BookingCard
                                                key={booking.ID || booking.id || `booking-${index}`}
                                                id={booking.ID || booking.id}
                                                booking_number={booking.BookingNumber || booking.booking_number}
                                                status={booking.Status || booking.status}
                                                start_date={booking.StartDate || booking.start_date}
                                                end_date={booking.EndDate || booking.end_date}
                                                total_amount={booking.TotalAmount || booking.total_amount}
                                                item_title={booking.rental_item?.title}
                                                item_image={booking.rental_item?.images?.[0]}
                                                onViewDetails={() => navigate(`/bookings/${booking.ID || booking.id}`)}
                                                onCancel={
                                                    (booking.Status || booking.status) === 'pending'
                                                        ? () => handleCancel(booking.ID || booking.id)
                                                        : undefined
                                                }
                                            />
                                        ))}
                                    {bookings.filter(b => (b.Status || b.status) === status).length === 0 && (
                                        <div className="column-empty">No {status} bookings</div>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}
