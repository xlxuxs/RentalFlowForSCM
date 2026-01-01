import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { Layout } from './components/layout';
import { ProtectedRoute } from './components/ProtectedRoute';
import {
  HomePage,
  BrowsePage,
  ItemDetailPage,
  LoginPage,
  RegisterPage,
  MyBookingsPage,
  OwnerBookingsPage,
  OwnerItemsPage,
  CreateItemPage,
  NotificationsPage,
  BookingDetailPage,
  PaymentCallbackPage,
  ProfilePage,
  EditItemPage,
} from './pages';
import './styles/global.css';

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Auth pages (no layout) */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />

          {/* Main layout */}
          <Route element={<Layout />}>
            {/* Public routes */}
            <Route path="/" element={<HomePage />} />
            <Route path="/browse" element={<BrowsePage />} />
            <Route path="/items/:id" element={<ItemDetailPage />} />

            {/* Protected routes - any authenticated user */}
            <Route
              path="/my-bookings"
              element={
                <ProtectedRoute>
                  <MyBookingsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/bookings/:id"
              element={
                <ProtectedRoute>
                  <BookingDetailPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/payment/callback"
              element={
                <ProtectedRoute>
                  <PaymentCallbackPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/notifications"
              element={
                <ProtectedRoute>
                  <NotificationsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/profile"
              element={
                <ProtectedRoute>
                  <ProfilePage />
                </ProtectedRoute>
              }
            />

            {/* Owner only routes */}
            <Route
              path="/owner/items"
              element={
                <ProtectedRoute requiredRole="owner">
                  <OwnerItemsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/owner/items/new"
              element={
                <ProtectedRoute>
                  <CreateItemPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/owner/items/:id/edit"
              element={
                <ProtectedRoute>
                  <EditItemPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/owner/bookings"
              element={
                <ProtectedRoute requiredRole="owner">
                  <OwnerBookingsPage />
                </ProtectedRoute>
              }
            />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
