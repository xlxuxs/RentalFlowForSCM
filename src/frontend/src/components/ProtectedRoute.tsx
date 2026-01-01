import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import React from 'react';

interface ProtectedRouteProps {
    children: React.ReactNode;
    requiredRole?: 'owner' | 'admin';
}

export function ProtectedRoute({ children, requiredRole }: ProtectedRouteProps) {
    const { isAuthenticated, isLoading, user } = useAuth();
    const location = useLocation();

    if (isLoading) {
        return (
            <div className="page-loading" style={{ minHeight: '100vh', display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <div className="spinner"></div>
            </div>
        );
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" state={{ from: location.pathname }} replace />;
    }

    if (requiredRole && user?.role !== requiredRole && user?.role !== 'admin') {
        return (
            <div className="container" style={{ padding: '2rem', textAlign: 'center' }}>
                <div className="empty-state">
                    <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🚫</div>
                    <h3 style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>Access Denied</h3>
                    <p style={{ color: '#666' }}>
                        You don't have permission to access this page.
                        {requiredRole === 'owner' && ' Only owners can access this page.'}
                    </p>
                </div>
            </div>
        );
    }

    return <>{children}</>;
}
