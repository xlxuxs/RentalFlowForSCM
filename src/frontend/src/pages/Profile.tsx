import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { usersApi } from '../services/api';
import { uploadToCloudinary } from '../services/cloudinary';
import './Profile.css';

export function ProfilePage() {
    const { user, updateUser } = useAuth();
    const [isEditing, setIsEditing] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [uploadingAvatar, setUploadingAvatar] = useState(false);

    const [formData, setFormData] = useState({
        first_name: '',
        last_name: '',
        email: '',
        phone: '',
        bio: '',
    });

    const [passwordData, setPasswordData] = useState({
        current_password: '',
        new_password: '',
        confirm_password: '',
    });

    const [showPasswordChange, setShowPasswordChange] = useState(false);

    useEffect(() => {
        if (user) {
            setFormData({
                first_name: user.first_name || '',
                last_name: user.last_name || '',
                email: user.email || '',
                phone: user.phone || '',
                bio: user.bio || '',
            });
        }
    }, [user]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value,
        });
    };

    const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setPasswordData({
            ...passwordData,
            [e.target.name]: e.target.value,
        });
    };

    const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file || !user) return;

        setUploadingAvatar(true);
        setError('');

        try {
            const imageUrl = await uploadToCloudinary(file);
            await usersApi.updateAvatar(user.id, imageUrl);
            updateUser({ ...user, avatar_url: imageUrl });
            setSuccess('Profile picture updated successfully!');
        } catch (err: any) {
            setError(err.message || 'Failed to upload avatar');
        } finally {
            setUploadingAvatar(false);
        }
    };

    const handleProfileUpdate = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) return;

        setLoading(true);
        setError('');
        setSuccess('');

        try {
            await usersApi.updateProfile(user.id, formData);
            updateUser({ ...user, ...formData });
            setSuccess('Profile updated successfully!');
            setIsEditing(false);
        } catch (err: any) {
            setError(err.message || 'Failed to update profile');
        } finally {
            setLoading(false);
        }
    };

    const handlePasswordUpdate = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) return;

        if (passwordData.new_password !== passwordData.confirm_password) {
            setError('New passwords do not match');
            return;
        }

        if (passwordData.new_password.length < 6) {
            setError('Password must be at least 6 characters');
            return;
        }

        setLoading(true);
        setError('');
        setSuccess('');

        try {
            await usersApi.changePassword(user.id, passwordData.current_password, passwordData.new_password);
            setSuccess('Password changed successfully!');
            setPasswordData({ current_password: '', new_password: '', confirm_password: '' });
            setShowPasswordChange(false);
        } catch (err: any) {
            setError(err.message || 'Failed to change password');
        } finally {
            setLoading(false);
        }
    };

    if (!user) {
        return <div className="profile-page">Loading...</div>;
    }

    return (
        <div className="profile-page">
            <div className="profile-container">
                <h1>My Profile</h1>

                {error && <div className="alert alert-error">{error}</div>}
                {success && <div className="alert alert-success">{success}</div>}

                {/* Avatar Section */}
                <div className="profile-avatar-section">
                    <div className="avatar-display">
                        {user.avatar_url ? (
                            <img src={user.avatar_url} alt={user.first_name} />
                        ) : (
                            <div className="avatar-placeholder">
                                {user.first_name?.[0]}{user.last_name?.[0]}
                            </div>
                        )}
                    </div>
                    <div className="avatar-actions">
                        <label className="btn btn-secondary" htmlFor="avatar-upload">
                            {uploadingAvatar ? 'Uploading...' : 'Change Photo'}
                        </label>
                        <input
                            id="avatar-upload"
                            type="file"
                            accept="image/*"
                            onChange={handleAvatarUpload}
                            disabled={uploadingAvatar}
                            style={{ display: 'none' }}
                        />
                    </div>
                </div>

                {/* Profile Info */}
                <div className="profile-info-section">
                    <div className="section-header">
                        <h2>Personal Information</h2>
                        {!isEditing && (
                            <button className="btn btn-secondary" onClick={() => setIsEditing(true)}>
                                Edit Profile
                            </button>
                        )}
                    </div>

                    {isEditing ? (
                        <form onSubmit={handleProfileUpdate} className="profile-form">
                            <div className="form-row">
                                <div className="form-group">
                                    <label>First Name</label>
                                    <input
                                        type="text"
                                        name="first_name"
                                        value={formData.first_name}
                                        onChange={handleInputChange}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label>Last Name</label>
                                    <input
                                        type="text"
                                        name="last_name"
                                        value={formData.last_name}
                                        onChange={handleInputChange}
                                        required
                                    />
                                </div>
                            </div>

                            <div className="form-group">
                                <label>Email</label>
                                <input
                                    type="email"
                                    name="email"
                                    value={formData.email}
                                    onChange={handleInputChange}
                                    required
                                    disabled
                                />
                                <small>Email cannot be changed</small>
                            </div>

                            <div className="form-group">
                                <label>Phone (Optional)</label>
                                <input
                                    type="tel"
                                    name="phone"
                                    value={formData.phone}
                                    onChange={handleInputChange}
                                    placeholder="+251..."
                                />
                            </div>

                            <div className="form-group">
                                <label>Bio (Optional)</label>
                                <textarea
                                    name="bio"
                                    value={formData.bio}
                                    onChange={handleInputChange}
                                    rows={4}
                                    placeholder="Tell us about yourself..."
                                />
                            </div>

                            <div className="form-actions">
                                <button type="button" className="btn btn-secondary" onClick={() => setIsEditing(false)}>
                                    Cancel
                                </button>
                                <button type="submit" className="btn btn-primary" disabled={loading}>
                                    {loading ? 'Saving...' : 'Save Changes'}
                                </button>
                            </div>
                        </form>
                    ) : (
                        <div className="profile-details">
                            <div className="detail-row">
                                <span className="label">Name:</span>
                                <span className="value">{user.first_name} {user.last_name}</span>
                            </div>
                            <div className="detail-row">
                                <span className="label">Email:</span>
                                <span className="value">{user.email}</span>
                            </div>
                            <div className="detail-row">
                                <span className="label">Phone:</span>
                                <span className="value">{user.phone || 'Not provided'}</span>
                            </div>
                            <div className="detail-row">
                                <span className="label">Role:</span>
                                <span className="value badge">{user.role}</span>
                            </div>
                            {user.bio && (
                                <div className="detail-row">
                                    <span className="label">Bio:</span>
                                    <span className="value">{user.bio}</span>
                                </div>
                            )}
                        </div>
                    )}
                </div>

                {/* Password Change */}
                <div className="profile-security-section">
                    <div className="section-header">
                        <h2>Security</h2>
                        {!showPasswordChange && (
                            <button className="btn btn-secondary" onClick={() => setShowPasswordChange(true)}>
                                Change Password
                            </button>
                        )}
                    </div>

                    {showPasswordChange && (
                        <form onSubmit={handlePasswordUpdate} className="password-form">
                            <div className="form-group">
                                <label>Current Password</label>
                                <input
                                    type="password"
                                    name="current_password"
                                    value={passwordData.current_password}
                                    onChange={handlePasswordChange}
                                    required
                                />
                            </div>

                            <div className="form-group">
                                <label>New Password</label>
                                <input
                                    type="password"
                                    name="new_password"
                                    value={passwordData.new_password}
                                    onChange={handlePasswordChange}
                                    required
                                    minLength={6}
                                />
                            </div>

                            <div className="form-group">
                                <label>Confirm New Password</label>
                                <input
                                    type="password"
                                    name="confirm_password"
                                    value={passwordData.confirm_password}
                                    onChange={handlePasswordChange}
                                    required
                                />
                            </div>

                            <div className="form-actions">
                                <button type="button" className="btn btn-secondary" onClick={() => setShowPasswordChange(false)}>
                                    Cancel
                                </button>
                                <button type="submit" className="btn btn-primary" disabled={loading}>
                                    {loading ? 'Updating...' : 'Update Password'}
                                </button>
                            </div>
                        </form>
                    )}
                </div>
            </div>
        </div>
    );
}
