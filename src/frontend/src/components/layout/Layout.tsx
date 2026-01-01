import { Outlet } from 'react-router-dom';
import { Navbar } from './Navbar';
import './Layout.css';

export function Layout() {
    return (
        <div className="layout">
            <Navbar />
            <main className="main-content">
                <Outlet />
            </main>
            <footer className="footer">
                <div className="footer-content">
                    <div className="footer-brand">
                        <span className="brand-icon">🏠</span>
                        <span>RentalFlow</span>
                    </div>
                    <p className="footer-text">
                        Rent equipment, vehicles, or property easily and securely.
                    </p>
                    <p className="footer-copyright">
                        © 2024 RentalFlow. All rights reserved.
                    </p>
                </div>
            </footer>
        </div>
    );
}
