import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { itemsApi } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { uploadImages } from '../services/cloudinary';
import './CreateItem.css';

export function CreateItemPage() {
    const { user } = useAuth();
    const navigate = useNavigate();
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        category: 'equipment',
        subcategory: '',
        daily_rate: '',
        weekly_rate: '',
        monthly_rate: '',
        security_deposit: '',
        city: '',
        address: '',
    });
    const [images, setImages] = useState<File[]>([]);
    const [imagePreviews, setImagePreviews] = useState<string[]>([]);
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState('');

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        setFormData(prev => ({
            ...prev,
            [e.target.name]: e.target.value,
        }));
    };

    const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(e.target.files || []);
        if (files.length === 0) return;

        // Limit to 5 images
        const newFiles = files.slice(0, 5 - images.length);

        setImages(prev => [...prev, ...newFiles]);

        // Create previews
        newFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                setImagePreviews(prev => [...prev, reader.result as string]);
            };
            reader.readAsDataURL(file);
        });
    };

    const removeImage = (index: number) => {
        setImages(prev => prev.filter((_, i) => i !== index));
        setImagePreviews(prev => prev.filter((_, i) => i !== index));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) return;

        setError('');
        setIsLoading(true);

        try {
            // Upload images to Cloudinary first
            let imageUrls: string[] = [];
            if (images.length > 0) {
                setUploadProgress('Uploading images...');
                imageUrls = await uploadImages(images);
                setUploadProgress('');
            }

            // Create item with image URLs
            await itemsApi.create({
                owner_id: user.id,
                title: formData.title,
                description: formData.description,
                category: formData.category,
                subcategory: formData.subcategory || undefined,
                daily_rate: parseFloat(formData.daily_rate),
                weekly_rate: formData.weekly_rate ? parseFloat(formData.weekly_rate) : undefined,
                monthly_rate: formData.monthly_rate ? parseFloat(formData.monthly_rate) : undefined,
                security_deposit: formData.security_deposit ? parseFloat(formData.security_deposit) : 0,
                city: formData.city,
                address: formData.address || undefined,
                images: imageUrls,
            });
            navigate('/owner/items');
        } catch (err: any) {
            setError(err.message || 'Failed to create item');
            setUploadProgress('');
        } finally {
            setIsLoading(false);
        }
    };

    const categories = [
        { id: 'vehicle', name: 'Vehicle', icon: '🚗' },
        { id: 'equipment', name: 'Equipment', icon: '🔧' },
        { id: 'property', name: 'Property', icon: '🏢' },
    ];

    return (
        <div className="create-item-page">
            <div className="container">
                <div className="create-item-header">
                    <h1>Add New Item</h1>
                    <p>List your item for rent</p>
                </div>

                <form onSubmit={handleSubmit} className="create-item-form">
                    {error && <div className="alert alert-error">{error}</div>}

                    <div className="form-section">
                        <h3>Basic Information</h3>

                        <div className="form-group">
                            <label className="form-label">Title</label>
                            <input
                                type="text"
                                name="title"
                                className="form-input"
                                placeholder="e.g., Power Drill Kit"
                                value={formData.title}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Description</label>
                            <textarea
                                name="description"
                                className="form-input"
                                placeholder="Describe your item, its condition, and what's included..."
                                rows={4}
                                value={formData.description}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Category</label>
                            <div className="category-select">
                                {categories.map(cat => (
                                    <label
                                        key={cat.id}
                                        className={`category-option ${formData.category === cat.id ? 'active' : ''}`}
                                    >
                                        <input
                                            type="radio"
                                            name="category"
                                            value={cat.id}
                                            checked={formData.category === cat.id}
                                            onChange={handleChange}
                                        />
                                        <span className="category-icon">{cat.icon}</span>
                                        <span className="category-name">{cat.name}</span>
                                    </label>
                                ))}
                            </div>
                        </div>

                        <div className="form-group">
                            <label className="form-label">Subcategory (optional)</label>
                            <input
                                type="text"
                                name="subcategory"
                                className="form-input"
                                placeholder="e.g., Power Tools"
                                value={formData.subcategory}
                                onChange={handleChange}
                            />
                        </div>
                    </div>

                    {/* Image Upload Section */}
                    <div className="form-section">
                        <h3>Photos</h3>
                        <p className="form-helper">Add up to 5 photos of your item</p>

                        <div className="image-upload-area">
                            {imagePreviews.map((preview, index) => (
                                <div key={index} className="image-preview">
                                    <img src={preview} alt={`Preview ${index + 1}`} />
                                    <button
                                        type="button"
                                        className="image-remove-btn"
                                        onClick={() => removeImage(index)}
                                    >
                                        ✕
                                    </button>
                                </div>
                            ))}

                            {images.length < 5 && (
                                <button
                                    type="button"
                                    className="image-add-btn"
                                    onClick={() => fileInputRef.current?.click()}
                                >
                                    <span className="add-icon">📷</span>
                                    <span>Add Photo</span>
                                </button>
                            )}
                        </div>

                        <input
                            ref={fileInputRef}
                            type="file"
                            accept="image/*"
                            multiple
                            onChange={handleImageSelect}
                            style={{ display: 'none' }}
                        />
                    </div>

                    <div className="form-section">
                        <h3>Pricing</h3>

                        <div className="form-row">
                            <div className="form-group">
                                <label className="form-label">Daily Rate ($) *</label>
                                <input
                                    type="number"
                                    name="daily_rate"
                                    className="form-input"
                                    placeholder="25"
                                    min="1"
                                    step="0.01"
                                    value={formData.daily_rate}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Weekly Rate ($)</label>
                                <input
                                    type="number"
                                    name="weekly_rate"
                                    className="form-input"
                                    placeholder="150"
                                    min="1"
                                    step="0.01"
                                    value={formData.weekly_rate}
                                    onChange={handleChange}
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Monthly Rate ($)</label>
                                <input
                                    type="number"
                                    name="monthly_rate"
                                    className="form-input"
                                    placeholder="500"
                                    min="1"
                                    step="0.01"
                                    value={formData.monthly_rate}
                                    onChange={handleChange}
                                />
                            </div>
                        </div>

                        <div className="form-group">
                            <label className="form-label">Security Deposit ($)</label>
                            <input
                                type="number"
                                name="security_deposit"
                                className="form-input"
                                placeholder="100"
                                min="0"
                                step="0.01"
                                value={formData.security_deposit}
                                onChange={handleChange}
                            />
                        </div>
                    </div>

                    <div className="form-section">
                        <h3>Location</h3>

                        <div className="form-row">
                            <div className="form-group">
                                <label className="form-label">City *</label>
                                <input
                                    type="text"
                                    name="city"
                                    className="form-input"
                                    placeholder="Addis Ababa"
                                    value={formData.city}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Address (optional)</label>
                                <input
                                    type="text"
                                    name="address"
                                    className="form-input"
                                    placeholder="123 Main St"
                                    value={formData.address}
                                    onChange={handleChange}
                                />
                            </div>
                        </div>
                    </div>

                    <div className="form-actions">
                        <button
                            type="button"
                            className="btn btn-ghost"
                            onClick={() => navigate('/owner/items')}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="btn btn-primary btn-lg"
                            disabled={isLoading}
                        >
                            {isLoading ? (uploadProgress || 'Creating...') : 'Create Item'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
