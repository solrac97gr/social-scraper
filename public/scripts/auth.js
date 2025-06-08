// Authentication utility functions
class AuthManager {
    constructor() {
        this.token = localStorage.getItem('authToken');
        this.user = this.getUserFromToken();
    }

    // Check if user is authenticated
    isAuthenticated() {
        if (!this.token) return false;
        
        try {
            const tokenPayload = this.decodeToken(this.token);
            const currentTime = Date.now() / 1000;
            
            // Check if token is expired
            if (tokenPayload.exp < currentTime) {
                this.logout();
                return false;
            }
            
            return true;
        } catch (error) {
            console.error('Token validation error:', error);
            this.logout();
            return false;
        }
    }

    // Decode JWT token (simple base64 decode)
    decodeToken(token) {
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));
            
            return JSON.parse(jsonPayload);
        } catch (error) {
            throw new Error('Invalid token format');
        }
    }

    // Get user information from token
    getUserFromToken() {
        if (!this.token) return null;
        
        try {
            const payload = this.decodeToken(this.token);
            return {
                id: payload.user_id,
                email: payload.email,
                role: payload.role,
                subscription: payload.subscription,
                username: payload.username || null // Optional field
            };
        } catch (error) {
            return null;
        }
    }

    // Set authentication token
    setToken(token) {
        this.token = token;
        localStorage.setItem('authToken', token);
        this.user = this.getUserFromToken();
    }

    // Get authentication token
    getToken() {
        return this.token;
    }

    // Get user information
    getUser() {
        return this.user;
    }

    // Logout user
    logout() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('authToken');
        window.location.href = '/';
    }

    // Redirect to login if not authenticated
    requireAuth() {
        if (!this.isAuthenticated()) {
            window.location.href = '/';
            return false;
        }
        return true;
    }

    // Make authenticated API request
    async apiRequest(url, options = {}) {
        if (!this.isAuthenticated()) {
            throw new Error('Not authenticated');
        }

        const defaultHeaders = {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.token}`
        };

        const requestOptions = {
            ...options,
            headers: {
                ...defaultHeaders,
                ...options.headers
            }
        };

        try {
            const response = await fetch(url, requestOptions);
            
            if (response.status === 401) {
                // Token expired or invalid
                this.logout();
                throw new Error('Authentication expired');
            }
            
            return response;
        } catch (error) {
            if (error.message === 'Authentication expired') {
                throw error;
            }
            console.error('API request error:', error);
            throw error;
        }
    }

    // Get authentication headers
    getAuthHeaders() {
        return {
            'Authorization': `Bearer ${this.token}`
        };
    }
}

// Create global auth manager instance
window.authManager = new AuthManager();

// Global authenticatedFetch function for API requests
window.authenticatedFetch = async function(url, options = {}) {
    if (!window.authManager.isAuthenticated()) {
        window.location.href = 'login.html';
        return;
    }
    
    // Check if we're sending FormData (for file uploads)
    const isFormData = options.body instanceof FormData;
    
    const headers = {
        ...window.authManager.getAuthHeaders(),
        ...(options.headers || {})
    };
    
    // Remove Content-Type header for FormData to let browser set it with boundary
    if (isFormData) {
        delete headers['Content-Type'];
    }
    
    const response = await fetch(url, {
        ...options,
        headers
    });
    
    if (response.status === 401) {
        window.authManager.logout();
        return;
    }
    
    if (!response.ok) {
        const errorText = await response.text();
        console.error('Server error:', response.status, errorText);
        throw new Error(`Server error (${response.status}): ${errorText}`);
    }
    
    return response;
};

// Utility function to show user info in the UI
function updateUserInterface() {
    const user = window.authManager.getUser();
    if (user) {
        // Update user email in the UI if element exists
        const userEmailElement = document.getElementById('user-email');
        if (userEmailElement) {
            userEmailElement.textContent = user.email;
        }

        // Update user role in the UI if element exists
        const userRoleElement = document.getElementById('user-role');
        if (userRoleElement) {
            userRoleElement.textContent = user.role;
        }

        // Update subscription info if element exists
        const userSubscriptionElement = document.getElementById('user-subscription');
        if (userSubscriptionElement) {
            userSubscriptionElement.textContent = user.subscription;
        }
    }
}

// Initialize authentication check when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Check if we're on a protected route
    const protectedRoutes = ['/dashboard.html', '/database.html'];
    const currentPath = window.location.pathname;
    
    if (protectedRoutes.some(route => currentPath.endsWith(route))) {
        // This is a protected route, require authentication
        if (!window.authManager.requireAuth()) {
            return; // Will redirect to login
        }
        
        // Update UI with user information
        updateUserInterface();
    }
    
    // If we're on the landing page and user is authenticated, redirect to dashboard
    if (currentPath === '/' || currentPath.endsWith('/index.html')) {
        if (window.authManager.isAuthenticated()) {
            window.location.href = '/dashboard.html';
        }
    }
});

// Function to handle logout
function handleLogout() {
    if (confirm('Are you sure you want to logout?')) {
        // Make logout API call
        window.authManager.apiRequest('/api/v1/users/logout', {
            method: 'POST'
        }).then(() => {
            window.authManager.logout();
        }).catch((error) => {
            console.error('Logout error:', error);
            // Even if API call fails, logout locally
            window.authManager.logout();
        });
    }
}
