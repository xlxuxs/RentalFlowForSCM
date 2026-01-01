import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { itemsApi } from '../services/api';
import { uploadImages } from '../services/cloudinary';
import './CreateItem.css';

export function EditItemPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { user } = useAuth();
    const [loading, setLoading] = useState(false);
    const [fetchingItem, setFetchingItem] = useState(true);
    const [error, setError] = useState('');
    const [uploadingImages, setUploadingImages] = useState(false);

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        category: 'vehicle' as 'vehicle' | 'equipment' | 'property',
        subcategory: '',
        daily_rate: '',
        weekly_rate: '',
        monthly_rate: '',
        security_deposit: '',
        address: '',
        city: '',
        specifications: {} as Record<string, string>,
    });

    const [existingImages, setExistingImages] = useState<string[]>([]);
    const [newImageFiles, setNewImageFiles] = useState<File[]>([]);
    const [newImagePreviews, setNewImagePreviews] = useState<string[]>([]);

    useEffect(() => {
        if (!user || user.role !== 'owner') {
            navigate('/');
            return;
        }

        const fetchItem = async () => {
            if (!id) return;

            try {
                const item = await itemsApi.get(id);

                // Verify ownership
                if (item.owner_id !== user.id) {
                    setError('You do not have permission to edit this item');
                    setTimeout(() => navigate('/owner/items'), 2000);
                    return;
                }

                setFormData({
                    title: item.title || '',
                    description: item.description || '',
                    category: (item.category as any) || 'vehicle',
                    subcategory: item.subcategory || '',
                    daily_rate: item.daily_rate?.toString() || '',
                    weekly_rate: item.weekly_rate?.toString() || '',
                    monthly_rate: item.monthly_rate?.toString() || '',
                    security_deposit: item.security_deposit?.toString() || '',
                    address: item.address || '',
                    city: item.city || '',
                    specifications: item.specifications || {},
                });
                setExistingImages(item.images || []);
            } catch (err: any) {
                setError(err.message || 'Failed to load item');
                setTimeout(() => navigate('/owner/items'), 2000);
            } finally {
                setFetchingItem(false);
            }
        };

        fetchItem();
    }, [id, user, navigate]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(e.target.files || []);
        setNewImageFiles(files);

        // Create previews
        const previews = files.map(file => URL.createObjectURL(file));
        setNewImagePreviews(previews);
    };

    const removeExistingImage = (index: number) => {
        setExistingImages(existingImages.filter((_, i) => i !== index));
    };

    const removeNewImage = (index: number) => {
        setNewImageFiles(newImageFiles.filter((_, i) => i !== index));
        setNewImagePreviews(newImagePreviews.filter((_, i) => i !== index));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user || !id) return;

        setLoading(true);
        setError('');

        try {
            // Upload new images if any
            let newImageUrls: string[] = [];
            if (newImageFiles.length > 0) {
                setUploadingImages(true);
                newImageUrls = await uploadImages(newImageFiles);
            }

            // Combine existing and new images
            const allImages = [...existingImages, ...newImageUrls];

            const itemData = {
                title: formData.title,
                description: formData.description,
                category: formData.category,
                subcategory: formData.subcategory || undefined,
                daily_rate: parseFloat(formData.daily_rate),
                weekly_rate: formData.weekly_rate ? parseFloat(formData.weekly_rate) : undefined,
                monthly_rate: formData.monthly_rate ? parseFloat(formData.monthly_rate) : undefined,
                security_deposit: parseFloat(formData.security_deposit),
                address: formData.address || undefined,
                city: formData.city,
                specifications: Object.keys(formData.specifications).length > 0 ? formData.specifications : undefined,
                images: allImages,
            };

            await itemsApi.update(id, itemData);
            navigate('/owner/items');
        } catch (err: any) {
            setError(err.message || 'Failed to update item. Please try again.');
        } finally {
            setLoading(false);
            setUploadingImages(false);
        }
    };

    if (fetchingItem) {
        return <div className="create-item-page"><div className="loading">Loading item...</div></div>;
    }

    return (
        <div className="create-item-page">
            <div className="create-item-form">
                <h1>Edit Rental Item</h1>

                {error && <div className="alert alert-error">{error}</div>}

                <form onSubmit={handleSubmit} className="item-form">
                    <div className="form-group">
                        <label className="form-label">Title *</label>
                        <input
                            type="text"
                            name="title"
                            className="form-input"
                            value={formData.title}
                            onChange={handleInputChange}
                            required
                            placeholder="e.g., 2020 Toyota Camry"
                        />
                    </div>

                    <div className="form-group">
                        <label className="form-label">Description *</label>
                        <textarea
                            name="description"
                            className="form-input"
                            value={formData.description}
                            onChange={handleInputChange}
                            required
                            rows={4}
                            placeholder="Describe your item in detail..."
                        />
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label className="form-label">Category *</label>
                            <select
                                name="category"
                                className="form-input"
                                value={formData.category}
                                onChange={handleInputChange}
                                required
                            >
                                <option value="vehicle">Vehicle</option>
                                <option value="equipment">Equipment</option>
                                <option value="property">Property</option>
                            </select>
                        </div>

                        <div className="form-group">
                            <label className="form-label">Subcategory</label>
                            <input
                                type="text"
                                name="subcategory"
                                className="form-input"
                                value={formData.subcategory}
                                onChange={handleInputChange}
                                placeholder="e.g., Sedan, Power Tools"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label className="form-label">Daily Rate (ETB) *</label>
                            <input
                                type="number"
                                name="daily_rate"
                                className="form-input"
                                value={formData.daily_rate}
                                onChange={handleInputChange}
                                required
                                min="0"
                                step="0.01"
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Security Deposit (ETB) *</label>
                            <input
                                type="number"
                                name="security_deposit"
                                className="form-input"
                                value={formData.security_deposit}
                                onChange={handleInputChange}
                                required
                                min="0"
                                step="0.01"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label className="form-label">Weekly Rate (ETB)</label>
                            <input
                                type="number"
                                name="weekly_rate"
                                className="form-input"
                                value={formData.weekly_rate}
                                onChange={handleInputChange}
                                min="0"
                                step="0.01"
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Monthly Rate (ETB)</label>
                            <input
                                type="number"
                                name="monthly_rate"
                                className="form-input"
                                value={formData.monthly_rate}
                                onChange={handleInputChange}
                                min="0"
                                step="0.01"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label className="form-label">City *</label>
                            <input
                                type="text"
                                name="city"
                                className="form-input"
                                value={formData.city}
                                onChange={handleInputChange}
                                required
                                placeholder="e.g., Addis Ababa"
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Address</label>
                            <input
                                type="text"
                                name="address"
                                className="form-input"
                                value={formData.address}
                                onChange={handleInputChange}
                                placeholder="Street address (optional)"
                            />
                        </div>
                    </div>

                    <div className="form-group">
                        <label className="form-label">Images</label>

                        <div className="image-upload-area">
                            {/* Existing images */}
                            {existingImages.map((url, index) => (
                                <div key={`existing-${index}`} className="image-preview">
                                    <img src={url} alt={`Item ${index + 1}`} />
                                    <button
                                        type="button"
                                        className="remove-image"
                                        onClick={() => removeExistingImage(index)}
                                    >
                                        ✕
                                    </button>
                                </div>
                            ))}

                            {/* New images */}
                            {newImagePreviews.map((url, index) => (
                                <div key={`new-${index}`} className="image-preview">
                                    <img src={url} alt={`New ${index + 1}`} />
                                    <button
                                        type="button"
                                        className="remove-image"
                                        onClick={() => removeNewImage(index)}
                                    >
                                        ✕
                                    </button>
                                </div>
                            ))}

                            <div className="image-add-btn" onClick={() => document.getElementById('edit-file-input')?.click()}>
                                <span className="add-icon">📷</span>
                                <span>Add New</span>
                            </div>
                        </div>

                        <input
                            id="edit-file-input"
                            type="file"
                            accept="image/*"
                            multiple
                            onChange={handleImageSelect}
                            style={{ display: 'none' }}
                        />
                    </div>

                    <div className="form-actions">
                        <button
                            type="button"
                            className="btn btn-ghost"
                            onClick={() => navigate('/owner/items')}
                            disabled={loading}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="btn btn-primary btn-lg"
                            disabled={loading || uploadingImages}
                        >
                            {uploadingImages ? 'Uploading Images...' : loading ? 'Updating...' : 'Update Item'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
