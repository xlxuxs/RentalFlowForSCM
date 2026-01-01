import { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { itemsApi } from '../services/api';
import { ItemCard, SearchFilters, type FilterState } from '../components/features';
import './Browse.css';

export function BrowsePage() {
    const [searchParams, setSearchParams] = useSearchParams();
    const [items, setItems] = useState<any[]>([]);
    const [total, setTotal] = useState(0);
    const [isLoading, setIsLoading] = useState(true);
    const [page, setPage] = useState(1);

    const [filters, setFilters] = useState<FilterState>({
        category: searchParams.get('category') || '',
        city: searchParams.get('city') || '',
        minPrice: searchParams.get('min_price') || '',
        maxPrice: searchParams.get('max_price') || '',
        sortBy: searchParams.get('sort') || 'newest',
    });

    useEffect(() => {
        const loadItems = async () => {
            setIsLoading(true);
            try {
                const result = await itemsApi.list({
                    category: filters.category || undefined,
                    city: filters.city || undefined,
                    min_price: filters.minPrice ? parseFloat(filters.minPrice) : undefined,
                    max_price: filters.maxPrice ? parseFloat(filters.maxPrice) : undefined,
                    sort: filters.sortBy || undefined,
                    page,
                    page_size: 12,
                });
                setItems(result.items || []);
                setTotal(result.total || 0);
            } catch (error) {
                console.error('Failed to load items:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadItems();
    }, [filters, page]);

    const handleFilterChange = (newFilters: FilterState) => {
        setFilters(newFilters);

        // Update URL params
        const params = new URLSearchParams();
        if (newFilters.category) params.set('category', newFilters.category);
        if (newFilters.city) params.set('city', newFilters.city);
        if (newFilters.minPrice) params.set('min_price', newFilters.minPrice);
        if (newFilters.maxPrice) params.set('max_price', newFilters.maxPrice);
        if (newFilters.sortBy) params.set('sort', newFilters.sortBy);

        setSearchParams(params);
        setPage(1); // Reset to first page when filters change
    };

    return (
        <div className="browse-page">
            <div className="browse-header">
                <div className="container">
                    <h1>Browse Rentals</h1>
                    <p className="browse-subtitle">Find the perfect item for your needs</p>
                </div>
            </div>

            <div className="browse-content container">
                {/* Filters Sidebar */}
                <aside className="browse-sidebar">
                    <SearchFilters filters={filters} onFilterChange={handleFilterChange} />
                </aside>

                {/* Items Grid */}
                <main className="browse-main">
                    <div className="browse-results-header">
                        <span className="results-count">
                            {total} {total === 1 ? 'item' : 'items'} found
                        </span>
                    </div>

                    {isLoading ? (
                        <div className="page-loading">
                            <div className="spinner"></div>
                        </div>
                    ) : items.length > 0 ? (
                        <>
                            <div className="items-grid">
                                {items.map(item => (
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

                            {/* Pagination */}
                            {total > 12 && (
                                <div className="pagination">
                                    <button
                                        className="btn btn-ghost"
                                        disabled={page === 1}
                                        onClick={() => setPage(p => p - 1)}
                                    >
                                        ← Prev
                                    </button>
                                    <span className="page-info">
                                        Page {page} of {Math.ceil(total / 12)}
                                    </span>
                                    <button
                                        className="btn btn-ghost"
                                        disabled={page >= Math.ceil(total / 12)}
                                        onClick={() => setPage(p => p + 1)}
                                    >
                                        Next →
                                    </button>
                                </div>
                            )}
                        </>
                    ) : (
                        <div className="empty-state">
                            <div className="empty-state-icon">🔍</div>
                            <h3 className="empty-state-title">No items found</h3>
                            <p className="empty-state-text">
                                Try adjusting your filters or check back later
                            </p>
                        </div>
                    )}
                </main>
            </div>
        </div>
    );
}
