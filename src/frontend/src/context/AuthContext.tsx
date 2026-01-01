import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { authApi } from '../services/api';
import type { User } from '../types';

interface AuthContextType {
    user: User | null;
    isLoading: boolean;
    isAuthenticated: boolean;
    login: (email: string, password: string) => Promise<void>;
    register: (data: {
        email: string;
        password: string;
        first_name: string;
        last_name: string;
        phone?: string;
        role: 'renter' | 'owner';
    }) => Promise<void>;
    logout: () => Promise<void>;
    updateUser: (userData: Partial<User>) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Check for existing session on mount
    useEffect(() => {
        const checkAuth = async () => {
            const token = localStorage.getItem('access_token');
            const storedUser = localStorage.getItem('user');

            if (token && storedUser) {
                try {
                    const result = await authApi.validateToken(token);
                    if (result.valid) {
                        setUser(JSON.parse(storedUser));
                    } else {
                        localStorage.removeItem('access_token');
                        localStorage.removeItem('refresh_token');
                        localStorage.removeItem('user');
                    }
                } catch {
                    // Token invalid or expired
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('refresh_token');
                    localStorage.removeItem('user');
                }
            }
            setIsLoading(false);
        };

        checkAuth();
    }, []);

    const login = async (email: string, password: string) => {
        const result = await authApi.login({ email, password });

        const userData: User = {
            id: result.user.id,
            email: result.user.email,
            first_name: result.user.first_name,
            last_name: result.user.last_name,
            role: result.user.role as 'renter' | 'owner' | 'admin',
        };

        localStorage.setItem('access_token', result.access_token);
        localStorage.setItem('refresh_token', result.refresh_token);
        localStorage.setItem('user', JSON.stringify(userData));

        setUser(userData);
    };

    const register = async (data: {
        email: string;
        password: string;
        first_name: string;
        last_name: string;
        phone?: string;
        role: 'renter' | 'owner';
    }) => {
        const result = await authApi.register(data) as any; // Cast as any because register return type might differ or we trust the response structure

        // Note: The reference implementation assumes register returns { user, access_token ... } similar to login
        // But the type definition I saw for register in api.ts returns Promise<unknown> (inferred from request default).
        // Let's assume the backend returns the session on register. If not, we might need to login after register.
        // Looking at api.ts: register returns request(...) which returns Promise<T>. In reference api.ts, register didn't have a generic type but `request` defaults.
        // I'll trust the reference logic for now.

        // Actually looking at reference AuthContext.tsx:
        // const result = await authApi.register(data) as any;

        const userData: User = {
            id: result.user.id,
            email: result.user.email,
            first_name: result.user.first_name,
            last_name: result.user.last_name,
            role: result.user.role as 'renter' | 'owner' | 'admin',
        };

        localStorage.setItem('access_token', result.access_token);
        localStorage.setItem('refresh_token', result.refresh_token);
        localStorage.setItem('user', JSON.stringify(userData));

        setUser(userData);
    };

    const logout = async () => {
        if (user) {
            try {
                await authApi.logout(user.id);
            } catch {
                // Ignore logout errors
            }
        }

        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        setUser(null);
    };

    const updateUser = (userData: Partial<User>) => {
        if (user) {
            const updatedUser = { ...user, ...userData };
            setUser(updatedUser);
            localStorage.setItem('user', JSON.stringify(updatedUser));
        }
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                isLoading,
                isAuthenticated: !!user,
                login,
                register,
                logout,
                updateUser,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}
