import { useState } from 'react';
import './SearchFilters.css';

export interface FilterState {
    category: string;
    city: string;
    minPrice: string;
    maxPrice: string;
    sortBy: string;
}

interface SearchFiltersProps {
    filters: FilterState;
    onFilterChange: (filters: FilterState) => void;
}

export function SearchFilters({ filters, onFilterChange }: SearchFiltersProps) {
    const [isExpanded, setIsExpanded] = useState(true);

    const handleChange = (field: keyof FilterState, value: string) => {
        onFilterChange({
            ...filters,
            [field]: value,
        });
    };

    const clearFilters = () => {
        onFilterChange({
            category: '',
            city: '',
            minPrice: '',
            maxPrice: '',
            sortBy: 'newest',
        });
    };

    const hasActiveFilters = filters.category || filters.city || filters.minPrice || filters.maxPrice;

    return (
        <div className="search-filters">
            <div className="filters-header">
                <h3>Filters</h3>
                <button
                    className="btn-toggle"
                    onClick={() => setIsExpanded(!isExpanded)}
                    aria-label={isExpanded ? 'Collapse filters' : 'Expand filters'}
                >
                    {isExpanded ? '−' : '+'}
                </button>
            </div>

            {isExpanded && (
                <div className="filters-body">
                    {/* Category Filter */}
                    <div className="filter-group">
                        <label>Category</label>
                        <select
                            value={filters.category}
                            onChange={(e) => handleChange('category', e.target.value)}
                        >
                            <option value="">All Categories</option>
                            <option value="vehicle">🚗 Vehicles</option>
                            <option value="equipment">🔧 Equipment</option>
                            <option value="property">🏢 Property</option>
                        </select>
                    </div>

                    {/* Location Filter */}
                    <div className="filter-group">
                        <label>Location</label>
                        <select
                            value={filters.city}
                            onChange={(e) => handleChange('city', e.target.value)}
                        >
                            <option value="">All Cities</option>
                            <option value="Addis Ababa">Addis Ababa</option>
                            <option value="Dire Dawa">Dire Dawa</option>
                            <option value="Mekelle">Mekelle</option>
                            <option value="Gondar">Gondar</option>
                            <option value="Bahir Dar">Bahir Dar</option>
                            <option value="Hawassa">Hawassa</option>
                            <option value="Jimma">Jimma</option>
                        </select>
                    </div>

                    {/* Price Range */}
                    <div className="filter-group">
                        <label>Price Range (ETB/day)</label>
                        <div className="price-inputs">
                            <input
                                type="number"
                                placeholder="Min"
                                value={filters.minPrice}
                                onChange={(e) => handleChange('minPrice', e.target.value)}
                                min="0"
                            />
                            <input
                                type="number"
                                placeholder="Max"
                                value={filters.maxPrice}
                                onChange={(e) => handleChange('maxPrice', e.target.value)}
                                min="0"
                            />
                        </div>
                    </div>

                    {/* Sort By */}
                    <div className="filter-group">
                        <label>Sort By</label>
                        <select
                            value={filters.sortBy}
                            onChange={(e) => handleChange('sortBy', e.target.value)}
                        >
                            <option value="newest">Newest First</option>
                            <option value="price_low">Price: Low to High</option>
                            <option value="price_high">Price: High to Low</option>
                            <option value="popular">Most Popular</option>
                        </select>
                    </div>

                    {/* Clear Filters */}
                    {hasActiveFilters && (
                        <button className="btn btn-secondary btn-block" onClick={clearFilters}>
                            Clear All Filters
                        </button>
                    )}
                </div>
            )}
        </div>
    );
}
