import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { itemsApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { ItemCard } from '../components/features';
import './Dashboard.css';
import './OwnerItems.css';

export function OwnerItemsPage() {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [items, setItems] = useState<any[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadItems = async () => {
            if (!user) return;

            try {
                const result = await itemsApi.getOwnerItems(user.id, 1, 50);
                setItems(result.items || []);
            } catch (error) {
                console.error('Failed to load items:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadItems();
    }, [user]);

    const handleDelete = async (id: string) => {
        if (!user || !window.confirm('Are you sure you want to delete this item?')) return;

        try {
            await itemsApi.delete(id, user.id);
            setItems(prev => prev.filter(item => (item.ID || item.id) !== id));
        } catch (error) {
            console.error('Failed to delete item:', error);
            alert('Failed to delete item. It might have active bookings.');
        }
    };

    const handleToggleStatus = async (item: any) => {
        try {
            const newStatus = !item.is_active;
            await itemsApi.update(item.ID || item.id, { is_active: newStatus });
            setItems(prev => prev.map(i =>
                (i.ID || i.id) === (item.ID || item.id) ? { ...i, is_active: newStatus } : i
            ));
        } catch (error) {
            console.error('Failed to toggle status:', error);
        }
    };

    return (
        <div className="owner-items-page dashboard-page">
            <div className="dashboard-header">
                <div className="container">
                    <div className="header-content" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div>
                            <h1>My Rental Items</h1>
                            <p>Manage and track the availability of your items</p>
                        </div>
                        <button className="btn btn-primary" onClick={() => navigate('/owner/items/new')}>
                            + Add New Item
                        </button>
                    </div>
                </div>
            </div>

            <div className="container">
                {isLoading ? (
                    <div className="page-loading">
                        <div className="spinner"></div>
                    </div>
                ) : items.length > 0 ? (
                    <div className="items-grid">
                        {items.map((item, index) => (
                            <div key={item.ID || item.id || `item-${index}`} style={{ position: 'relative' }}>
                                <div className="item-actions-overlay">
                                    <button
                                        className="item-action-btn edit"
                                        onClick={() => navigate(`/owner/items/${item.ID || item.id}/edit`)}
                                        title="Edit Item"
                                    >
                                        ✏️
                                    </button>
                                    <button
                                        className="item-action-btn delete"
                                        onClick={() => handleDelete(item.ID || item.id)}
                                        title="Delete Item"
                                    >
                                        🗑️
                                    </button>
                                </div>
                                <div className={`item-status-badge ${item.is_active ? 'active' : 'inactive'}`} onClick={() => handleToggleStatus(item)} style={{ cursor: 'pointer' }}>
                                    {item.is_active ? 'Active' : 'Inactive'}
                                </div>
                                <ItemCard
                                    id={item.ID || item.id}
                                    title={item.title}
                                    daily_rate={item.daily_rate}
                                    city={item.city}
                                    images={item.images}
                                    category={item.category}
                                />
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="empty-state">
                        <div className="empty-state-icon">📦</div>
                        <h3 className="empty-state-title">No items listed</h3>
                        <p className="empty-state-text">
                            You haven't listed any items for rent yet.
                        </p>
                        <button className="btn btn-primary" onClick={() => navigate('/owner/items/new')}>
                            List Your First Item
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}
