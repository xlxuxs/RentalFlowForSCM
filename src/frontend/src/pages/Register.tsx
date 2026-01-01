import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Auth.css';

export function RegisterPage() {
    const navigate = useNavigate();
    const { register } = useAuth();

    const [formData, setFormData] = useState({
        email: '',
        password: '',
        confirmPassword: '',
        first_name: '',
        last_name: '',
        phone: '',
        role: 'renter' as 'renter' | 'owner',
    });
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        setFormData(prev => ({
            ...prev,
            [e.target.name]: e.target.value,
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (formData.password !== formData.confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        if (formData.password.length < 6) {
            setError('Password must be at least 6 characters');
            return;
        }

        setIsLoading(true);

        try {
            await register({
                email: formData.email,
                password: formData.password,
                first_name: formData.first_name,
                last_name: formData.last_name,
                phone: formData.phone || undefined,
                role: formData.role,
            });
            navigate('/');
        } catch (err: any) {
            setError(err.message || 'Registration failed');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-container">
                <div className="auth-card">
                    <div className="auth-header">
                        <Link to="/" className="auth-logo">
                            <span>🏠</span>
                            <span>RentalFlow</span>
                        </Link>
                        <h1>Create Account</h1>
                        <p>Join thousands of renters and owners</p>
                    </div>

                    <form onSubmit={handleSubmit} className="auth-form">
                        {error && <div className="alert alert-error">{error}</div>}

                        <div className="form-row">
                            <div className="form-group">
                                <label className="form-label">First Name</label>
                                <input
                                    type="text"
                                    name="first_name"
                                    className="form-input"
                                    placeholder="John"
                                    value={formData.first_name}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Last Name</label>
                                <input
                                    type="text"
                                    name="last_name"
                                    className="form-input"
                                    placeholder="Doe"
                                    value={formData.last_name}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        <div className="form-group">
                            <label className="form-label">Email</label>
                            <input
                                type="email"
                                name="email"
                                className="form-input"
                                placeholder="you@example.com"
                                value={formData.email}
                                onChange={handleChange}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Phone (optional)</label>
                            <input
                                type="tel"
                                name="phone"
                                className="form-input"
                                placeholder="+1 234 567 8900"
                                value={formData.phone}
                                onChange={handleChange}
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">I want to</label>
                            <div className="role-selector">
                                <label className={`role-option ${formData.role === 'renter' ? 'active' : ''}`}>
                                    <input
                                        type="radio"
                                        name="role"
                                        value="renter"
                                        checked={formData.role === 'renter'}
                                        onChange={handleChange}
                                    />
                                    <span className="role-icon">🛒</span>
                                    <span className="role-text">Rent Items</span>
                                </label>
                                <label className={`role-option ${formData.role === 'owner' ? 'active' : ''}`}>
                                    <input
                                        type="radio"
                                        name="role"
                                        value="owner"
                                        checked={formData.role === 'owner'}
                                        onChange={handleChange}
                                    />
                                    <span className="role-icon">💼</span>
                                    <span className="role-text">List Items</span>
                                </label>
                            </div>
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label className="form-label">Password</label>
                                <input
                                    type="password"
                                    name="password"
                                    className="form-input"
                                    placeholder="••••••••"
                                    value={formData.password}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Confirm Password</label>
                                <input
                                    type="password"
                                    name="confirmPassword"
                                    className="form-input"
                                    placeholder="••••••••"
                                    value={formData.confirmPassword}
                                    onChange={handleChange}
                                    required
                                />
                            </div>
                        </div>

                        <button
                            type="submit"
                            className="btn btn-primary btn-lg"
                            style={{ width: '100%' }}
                            disabled={isLoading}
                        >
                            {isLoading ? (
                                <>
                                    <div className="spinner spinner-sm"></div>
                                    <span>Creating account...</span>
                                </>
                            ) : (
                                'Create Account'
                            )}
                        </button>
                    </form>

                    <div className="auth-footer">
                        <p>
                            Already have an account?{' '}
                            <Link to="/login">Sign in</Link>
                        </p>
                    </div>
                </div>

                <div className="auth-decoration">
                    <div className="decoration-content">
                        <h2>Start Your Journey</h2>
                        <p>Whether you want to rent or list items, we've got you covered</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
