import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { itemsApi } from '../services/api';
import { ItemCard } from '../components/features';
import './Home.css';

export function HomePage() {
    const [featuredItems, setFeaturedItems] = useState<any[]>([]);
    const [totalItems, setTotalItems] = useState(0);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const loadFeaturedItems = async () => {
            try {
                const result = await itemsApi.list({ page: 1, page_size: 6 });
                setFeaturedItems(result.items || []);
                setTotalItems(result.total || 0);
            } catch (error) {
                console.error('Failed to load items:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadFeaturedItems();
    }, []);

    const categories = [
        { id: 'vehicle', name: 'Vehicles', icon: '🚗', description: 'Cars, bikes, and more' },
        { id: 'equipment', name: 'Equipment', icon: '🔧', description: 'Tools and machinery' },
        { id: 'property', name: 'Property', icon: '🏢', description: 'Spaces and venues' },
    ];

    return (
        <div className="home-page">
            {/* Hero Section */}
            <section className="hero">
                <div className="hero-content">
                    <h1 className="hero-title">
                        Rent Anything, <span className="gradient-text">Anywhere</span>
                    </h1>
                    <p className="hero-subtitle">
                        Find the perfect equipment, vehicle, or property for your next project.
                        Easy booking, secure payments, trusted rentals.
                    </p>
                    <div className="hero-actions">
                        <Link to="/browse" className="btn btn-primary btn-lg">
                            Browse Rentals
                        </Link>
                        <Link to="/register" className="btn btn-accent btn-lg">
                            List Your Item
                        </Link>
                    </div>

                    {totalItems > 0 && (
                        <div className="hero-stats">
                            <div className="stat">
                                <span className="stat-value">{totalItems}</span>
                                <span className="stat-label">Items Available</span>
                            </div>
                        </div>
                    )}
                </div>
                <div className="hero-decoration">
                    <div className="decoration-circle circle-1"></div>
                    <div className="decoration-circle circle-2"></div>
                    <div className="decoration-circle circle-3"></div>
                </div>
            </section>

            {/* Categories Section */}
            <section className="section categories-section">
                <div className="container">
                    <h2 className="section-title">Browse by Category</h2>
                    <div className="categories-grid">
                        {categories.map(cat => (
                            <Link
                                to={`/browse?category=${cat.id}`}
                                key={cat.id}
                                className="category-card"
                            >
                                <span className="category-icon">{cat.icon}</span>
                                <h3 className="category-name">{cat.name}</h3>
                                <p className="category-desc">{cat.description}</p>
                            </Link>
                        ))}
                    </div>
                </div>
            </section>

            {/* Featured Items Section */}
            <section className="section featured-section">
                <div className="container">
                    <div className="section-header">
                        <h2 className="section-title">Featured Rentals</h2>
                        <Link to="/browse" className="btn btn-ghost">
                            View All →
                        </Link>
                    </div>

                    {isLoading ? (
                        <div className="page-loading">
                            <div className="spinner"></div>
                        </div>
                    ) : featuredItems.length > 0 ? (
                        <div className="items-grid grid grid-cols-3 gap-6">
                            {featuredItems.map(item => (
                                <ItemCard
                                    key={item.id}
                                    id={item.id}
                                    title={item.title}
                                    category={item.category}
                                    city={item.city}
                                    daily_rate={item.daily_rate}
                                    images={item.images}
                                    is_active={item.is_active}
                                />
                            ))}
                        </div>
                    ) : (
                        <div className="empty-state">
                            <div className="empty-state-icon">📦</div>
                            <h3 className="empty-state-title">No items available yet</h3>
                            <p className="empty-state-text">Be the first to list an item for rent!</p>
                            <Link to="/register" className="btn btn-accent">
                                Get Started
                            </Link>
                        </div>
                    )}
                </div>
            </section>

            {/* How It Works Section */}
            <section className="section how-it-works-section">
                <div className="container">
                    <h2 className="section-title text-center">How It Works</h2>
                    <div className="steps-grid">
                        <div className="step-card">
                            <div className="step-number">1</div>
                            <h3 className="step-title">Browse & Find</h3>
                            <p className="step-desc">Search through our catalog of rental items by category or location</p>
                        </div>
                        <div className="step-card">
                            <div className="step-number">2</div>
                            <h3 className="step-title">Book & Pay</h3>
                            <p className="step-desc">Select your dates and complete secure payment online</p>
                        </div>
                        <div className="step-card">
                            <div className="step-number">3</div>
                            <h3 className="step-title">Pickup & Enjoy</h3>
                            <p className="step-desc">Coordinate pickup with the owner and enjoy your rental</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* CTA Section */}
            <section className="section cta-section">
                <div className="container">
                    <div className="cta-content">
                        <h2 className="cta-title">Ready to start renting?</h2>
                        <p className="cta-text">Join thousands of renters and owners on RentalFlow</p>
                        <Link to="/register" className="btn btn-primary btn-lg">
                            Create Free Account
                        </Link>
                    </div>
                </div>
            </section>
        </div>
    );
}
