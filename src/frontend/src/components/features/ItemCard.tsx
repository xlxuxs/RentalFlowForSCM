import { Link } from 'react-router-dom';
import { getThumbnailUrl } from '../../services/cloudinary';
import './ItemCard.css';

interface ItemCardProps {
    id: string;
    title: string;
    category: string;
    city: string;
    daily_rate: number;
    images?: string[];
    is_active?: boolean;
}

export function ItemCard({ id, title, category, city, daily_rate, images, is_active = true }: ItemCardProps) {
    const emojiMap: Record<string, string> = {
        vehicle: '🚗',
        equipment: '🔧',
        property: '🏢',
    };

    const categoryEmoji = emojiMap[category] || '📦';

    return (
        <Link to={`/items/${id}`} className="item-card">
            <div className="item-card-image">
                {images && images.length > 0 ? (
                    <img src={getThumbnailUrl(images[0])} alt={title} loading="lazy" />
                ) : (
                    <div className="item-card-placeholder">
                        <span>{categoryEmoji}</span>
                    </div>
                )}
                <span className={`item-card-badge badge ${is_active ? 'badge-success' : 'badge-gray'}`}>
                    {is_active ? 'Available' : 'Unavailable'}
                </span>
            </div>
            <div className="item-card-body">
                <div className="item-card-category">
                    <span>{categoryEmoji}</span>
                    <span>{category}</span>
                </div>
                <h3 className="item-card-title">{title}</h3>
                <p className="item-card-location">📍 {city}</p>
                <div className="item-card-price">
                    <span className="price-amount">${daily_rate}</span>
                    <span className="price-period">/ day</span>
                </div>
            </div>
        </Link>
    );
}
