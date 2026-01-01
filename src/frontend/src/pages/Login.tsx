import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Auth.css';

export function LoginPage() {
    const navigate = useNavigate();
    const location = useLocation();
    const { login } = useAuth();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const from = (location.state as any)?.from || '/';

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);

        try {
            await login(email, password);
            navigate(from, { replace: true });
        } catch (err: any) {
            setError(err.message || 'Login failed');
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
                        <h1>Welcome Back</h1>
                        <p>Sign in to your account</p>
                    </div>

                    <form onSubmit={handleSubmit} className="auth-form">
                        {error && <div className="alert alert-error">{error}</div>}

                        <div className="form-group">
                            <label className="form-label">Email</label>
                            <input
                                type="email"
                                className="form-input"
                                placeholder="you@example.com"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label className="form-label">Password</label>
                            <input
                                type="password"
                                className="form-input"
                                placeholder="••••••••"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
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
                                    <span>Signing in...</span>
                                </>
                            ) : (
                                'Sign In'
                            )}
                        </button>
                    </form>

                    <div className="auth-footer">
                        <p>
                            Don't have an account?{' '}
                            <Link to="/register">Create one</Link>
                        </p>
                    </div>
                </div>

                <div className="auth-decoration">
                    <div className="decoration-content">
                        <h2>Rent with Confidence</h2>
                        <p>Access thousands of rental items from trusted owners</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
