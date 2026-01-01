import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import './Navbar.css';

export function Navbar() {
    const { user, isAuthenticated, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = async () => {
        await logout();
        navigate('/');
    };

    return (
        <nav className="navbar">
            <div className="navbar-container">
                <Link to="/" className="navbar-brand">
                    <span className="brand-icon">🏠</span>
                    <span className="brand-text">RentalFlow</span>
                </Link>

                <div className="navbar-links">
                    <Link to="/browse" className="nav-link">Browse</Link>

                    {isAuthenticated ? (
                        <>
                            {user?.role === 'owner' && (
                                <>
                                    <Link to="/owner/items" className="nav-link">My Items</Link>
                                    <Link to="/owner/bookings" className="nav-link">Bookings Received</Link>
                                </>
                            )}
                            <Link to="/my-bookings" className="nav-link">My Bookings</Link>
                            <Link to="/notifications" className="nav-link nav-link-icon">
                                🔔
                            </Link>

                            <div className="nav-user">
                                <span className="nav-user-name">
                                    {user?.first_name}
                                </span>
                                <span className="nav-user-role badge badge-primary">
                                    {user?.role}
                                </span>
                                <button onClick={handleLogout} className="btn btn-ghost btn-sm">
                                    Logout
                                </button>
                            </div>
                        </>
                    ) : (
                        <div className="nav-auth">
                            <Link to="/login" className="btn btn-ghost">Login</Link>
                            <Link to="/register" className="btn btn-primary">Sign Up</Link>
                        </div>
                    )}
                </div>
            </div>
        </nav>
    );
}
