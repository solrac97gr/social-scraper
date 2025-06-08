// Dashboard Authentication and JWT token management
const DashboardAuthManager = {
    getToken() {
        return localStorage.getItem('authToken');
    },
    
    getUserInfo() {
        const userInfo = localStorage.getItem('userInfo');
        return userInfo ? JSON.parse(userInfo) : null;
    },
    
    isLoggedIn() {
        return !!this.getToken();
    },
    
    async logout() {
        try {
            // Call logout API if user is logged in
            if (this.isLoggedIn()) {
                await authenticatedFetch('/api/v1/users/logout', {
                    method: 'POST'
                });
            }
        } catch (error) {
            console.warn('Logout API call failed:', error);
        } finally {
            // Always clear local storage regardless of API call success
            localStorage.removeItem('authToken');
            localStorage.removeItem('userInfo');
            localStorage.removeItem('rememberUser');
            this.updateAuthUI();
            window.location.reload();
        }
    },
    
    updateAuthUI() {
        const loginSection = document.getElementById('loginSection');
        
        if (this.isLoggedIn()) {
            const userInfo = this.getUserInfo();
            if (loginSection) {
                loginSection.style.display = 'none';
            }
            
            if (userInfo) {
                // Update navbar user info if elements exist
                const navUserEmail = document.getElementById('nav-user-email');
                const userEmail = document.getElementById('user-email');
                const userRole = document.getElementById('user-role');
                const userSubscription = document.getElementById('user-subscription');
                
                if (navUserEmail) navUserEmail.textContent = userInfo.email;
                if (userEmail) userEmail.textContent = userInfo.email;
                if (userRole) userRole.textContent = userInfo.role || 'User';
                if (userSubscription) userSubscription.textContent = userInfo.subscription || 'Free';
            }
        } else {
            if (loginSection) {
                loginSection.style.display = 'block';
            }
        }
    },
    
    async verifyToken() {
        const token = this.getToken();
        if (!token) return false;
        
        try {
            const response = await fetch('/api/v1/influencers/health', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (!response.ok) {
                this.logout();
                return false;
            }
            
            return true;
        } catch (error) {
            console.error('Token verification failed:', error);
            this.logout();
            return false;
        }
    },
    
    getAuthHeaders() {
        const token = this.getToken();
        return token ? {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        } : {
            'Content-Type': 'application/json'
        };
    }
};

// Protected fetch function that includes JWT token
async function authenticatedFetch(url, options = {}) {
    if (!DashboardAuthManager.isLoggedIn()) {
        throw new Error('Authentication required');
    }
    
    // Check if we're sending FormData (for file uploads)
    const isFormData = options.body instanceof FormData;
    
    const headers = {
        ...DashboardAuthManager.getAuthHeaders(),
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
        DashboardAuthManager.logout();
        window.location.href = 'login.html';
        throw new Error('Authentication expired');
    }
    
    if (!response.ok) {
        const errorText = await response.text();
        console.error('Server error:', response.status, errorText);
        throw new Error(`Server error (${response.status}): ${errorText}`);
    }
    
    return response;
}

// Show authentication required modal
function showAuthModal() {
    const modalHtml = `
        <div class="modal fade" id="authModal" tabindex="-1" role="dialog">
            <div class="modal-dialog modal-dialog-centered" role="document">
                <div class="modal-content">
                    <div class="modal-header bg-primary text-white">
                        <h5 class="modal-title">
                            <i class="fas fa-lock"></i> Authentication Required
                        </h5>
                        <button type="button" class="close text-white" data-dismiss="modal">
                            <span>&times;</span>
                        </button>
                    </div>
                    <div class="modal-body text-center">
                        <i class="fas fa-user-shield fa-3x text-primary mb-3"></i>
                        <h5>Please log in to continue</h5>
                        <p class="text-muted">You need to be logged in to use the Social Scraper features.</p>
                    </div>
                    <div class="modal-footer">
                        <a href="login.html" class="btn btn-primary">
                            <i class="fas fa-sign-in-alt"></i> Login
                        </a>
                        <a href="register.html" class="btn btn-outline-primary">
                            <i class="fas fa-user-plus"></i> Register
                        </a>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    document.body.insertAdjacentHTML('beforeend', modalHtml);
    $('#authModal').modal('show');
    
    $('#authModal').on('hidden.bs.modal', function() {
        document.getElementById('authModal').remove();
    });
}

const translations = {
    en: {
        title: "Social Scraper 🚀",
        subtitle: "Extract information from your favorite social media channels",
        uploadTitle: "Upload File 📂",
        uploadButton: "Upload ⬆️",
        loaderText: "Processing your file... ⏳",
        languageLabel: "Language",
        resultsTitle: "Analysis Results",
        downloadButton: "Download Results ⬇️",
        headers: ["Channel Name", "Followers Count", "Original Link", "Platform", "Registration Status"],
        minFollowersInput: "Min followers",
        maxFollowersInput: "Max followers",
        nameFilterInput: "Search by name",
        platformFilterDropdown: "Filter by Platform",
        all: "All",
        telegram: "Telegram",
        rutube: "Rutube",
        vk: "VK",
        instagram: "Instagram",
        youtube: "YouTube"
    },
    es: {
        title: "Extractor Social 🚀",
        subtitle: "Extrae información de tus canales de redes sociales favoritos",
        uploadTitle: "Subir Archivo 📂",
        uploadButton: "Subir ⬆️",
        loaderText: "Procesando tu archivo... ⏳",
        languageLabel: "Idioma",
        resultsTitle: "Resultados del Análisis",
        downloadButton: "Descargar Resultados ⬇️",
        headers: ["Nombre del Canal", "Cantidad de Seguidores", "Enlace Original", "Plataforma", "Estado de Registro"],
        minFollowersInput: "Mín. seguidores",
        maxFollowersInput: "Máx. seguidores",
        nameFilterInput: "Buscar por nombre",
        platformFilterDropdown: "Filtrar por plataforma",
        all: "Todos",
        telegram: "Telegram",
        rutube: "Rutube",
        vk: "VK",
        instagram: "Instagram",
        youtube: "YouTube"
    },
    ru: {
        title: "Социальный Скрапер 🚀",
        subtitle: "Извлекайте информацию из ваших любимых социальных сетей",
        uploadTitle: "Загрузить файл 📂",
        uploadButton: "Загрузить ⬆️",
        loaderText: "Обработка вашего файла... ⏳",
        languageLabel: "Язык",
        resultsTitle: "Результаты анализа",
        downloadButton: "Скачать результаты ⬇️",
        headers: ["Название канала", "Количество подписчиков", "Оригинальная ссылка", "Платформа", "Статус регистрации"],
        minFollowersInput: "Мин. подписчиков",
        maxFollowersInput: "Макс. подписчиков",
        nameFilterInput: "Поиск по имени",
        platformFilterDropdown: "Фильтр по платформе",
        all: "Все",
        telegram: "Telegram",
        rutube: "Rutube",
        vk: "VK",
        instagram: "Instagram",
        youtube: "YouTube"
    }
};

function changeLanguage(lang) {
    const elements = translations[lang];
    document.getElementById('title').innerText = elements.title;
    document.getElementById('subtitle').innerText = elements.subtitle;
    document.getElementById('uploadTitle').innerText = elements.uploadTitle;
    document.getElementById('uploadButton').innerText = elements.uploadButton;
    document.getElementById('loaderText').innerText = elements.loaderText;
    document.getElementById('languageLabel').innerText = elements.languageLabel;
    document.getElementById('resultsTitle').innerText = elements.resultsTitle;
    document.getElementById('downloadButton').innerText = elements.downloadButton;

    const headers = document.querySelectorAll('#resultsTable thead th');
    headers.forEach((header, index) => {
        header.innerText = elements.headers[index];
    });

    document.getElementById('minFollowersInput').placeholder = elements.minFollowersInput;
    document.getElementById('maxFollowersInput').placeholder = elements.maxFollowersInput;
    document.getElementById('nameFilterInput').placeholder = elements.nameFilterInput;
    document.getElementById('platformFilterDropdown').innerText = elements.platformFilterDropdown;
    const platformItems = document.querySelectorAll('#platformFilterDropdown + .dropdown-menu .dropdown-item');
    platformItems[0].innerText = elements.all;
    platformItems[1].innerText = elements.telegram;
    platformItems[2].innerText = elements.rutube;
    platformItems[3].innerText = elements.vk;
    platformItems[4].innerText = elements.instagram;
    platformItems[5].innerText = elements.youtube;
}

function updateLoaderText(estimatedTime) {
    const loaderText = document.getElementById('loaderText');
    loaderText.innerText = `Processing your file... ⏳ Estimated time: ${estimatedTime} seconds`;
}

function startCountdown(estimatedTime) {
    const loaderText = document.getElementById('loaderText');
    const interval = setInterval(() => {
        if (estimatedTime > 0) {
            estimatedTime--;
            loaderText.innerText = `Processing your file... ⏳ Estimated time: ${estimatedTime} seconds`;
        } else {
            clearInterval(interval);
        }
    }, 1000);
}

function applyFilters() {
    const platform = document.querySelector('#platformFilterDropdown').innerText.toLowerCase();
    const nameFilter = document.getElementById('nameFilterInput').value.toLowerCase();
    const minFollowers = document.getElementById('minFollowersInput').value;
    const maxFollowers = document.getElementById('maxFollowersInput').value;
    const rows = document.querySelectorAll('#resultsTableBody tr');

    rows.forEach(row => {
        const platformCell = row.querySelector('td:nth-child(4) span');
        const nameCell = row.querySelector('td:nth-child(1)');
        const followersCount = parseInt(row.querySelector('td:nth-child(2)').innerText);

        const matchesPlatform = platform === 'filter by platform' || platformCell.classList.contains(`badge-${platform}`);
        const matchesName = nameCell.innerText.toLowerCase().includes(nameFilter);
        const matchesFollowers = (minFollowers === '' || followersCount >= minFollowers) && (maxFollowers === '' || followersCount <= maxFollowers);

        if (matchesPlatform && matchesName && matchesFollowers) {
            row.style.display = '';
        } else {
            row.style.display = 'none';
        }
    });
}

function filterByPlatform(platform) {
    document.querySelector('#platformFilterDropdown').innerText = platform.charAt(0).toUpperCase() + platform.slice(1);
    applyFilters();
}

function filterByName() {
    applyFilters();
}

function filterByFollowers() {
    applyFilters();
}

function clearFilters() {
    document.getElementById('minFollowersInput').value = '';
    document.getElementById('maxFollowersInput').value = '';
    document.getElementById('nameFilterInput').value = '';
    document.querySelector('#platformFilterDropdown').innerText = 'Filter by Platform';
    applyFilters();
}

document.getElementById('uploadForm').addEventListener('submit', async function(event) {
    event.preventDefault();
    
    // Check authentication first
    if (!DashboardAuthManager.isLoggedIn()) {
        showAuthModal();
        return;
    }
    
    // Verify token is still valid
    if (!(await DashboardAuthManager.verifyToken())) {
        showAuthModal();
        return;
    }
    
    document.getElementById('loader').style.display = 'block';
    document.getElementById('uploadForm').classList.add('animate__animated', 'animate__bounceOut');

    var formData = new FormData();
    
    // Add input based on current mode
    if (currentInputMode === 'file') {
        const fileInput = document.getElementById('fileInput');
        if (fileInput.files.length === 0) {
            alert('Please select a file first');
            document.getElementById('loader').style.display = 'none';
            return;
        }
        formData.append('file', fileInput.files[0]);
    } else if (currentInputMode === 'text') {
        const textInput = document.getElementById('linkTextInput');
        if (!textInput.value.trim()) {
            alert('Please enter some links');
            document.getElementById('loader').style.display = 'none';
            return;
        }
        formData.append('textInput', textInput.value.trim());
    } else {
        alert('Please select an input method');
        document.getElementById('loader').style.display = 'none';
        return;
    }

    try {
        // Call the new endpoint to get the estimated time
        const estimateResponse = await authenticatedFetch('/api/v1/influencers/estimate-time', {
            method: 'POST',
            body: formData
        });
        
        if (!estimateResponse.ok) {
            throw new Error('Failed to estimate processing time');
        }
        
        const estimateData = await estimateResponse.json();
        updateLoaderText(estimateData.estimatedTime);
        startCountdown(estimateData.estimatedTime);

        // Proceed with the upload/analysis
        const uploadResponse = await authenticatedFetch('/api/v1/influencers/upload', {
            method: 'POST',
            body: formData
        });
        
        if (!uploadResponse.ok) {
            throw new Error('Upload/analysis failed');
        }
        
        const uploadData = await uploadResponse.json();
        document.getElementById('loader').style.display = 'none';
        document.getElementById('resultsContainer').style.display = 'block';

        const resultsTableBody = document.getElementById('resultsTableBody');
        resultsTableBody.innerHTML = '';

        // remove the first row (header)
        uploadData.results.shift();

        uploadData.results.forEach(result => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${result[0]}</td>
                <td class="count-up" data-count="${result[1]}">${result[1]}</td>
                <td><a href="${result[2]}" target="_blank">${result[2]}</a></td>
                <td><span class="badge badge-${result[3].toLowerCase()}">${result[3]}</span></td>
                <td>${result[4]}</td>
            `;
            resultsTableBody.appendChild(row);
        });

        document.getElementById('downloadButton').addEventListener('click', function() {
            const token = DashboardAuthManager.getToken();
            const downloadUrl = `/api/v1/influencers/download?filename=${uploadData.outputFile}`;
            
            // Create a temporary link with auth header for download
            fetch(downloadUrl, {
                headers: DashboardAuthManager.getAuthHeaders()
            })
            .then(response => response.blob())
            .then(blob => {
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = uploadData.outputFile;
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                document.body.removeChild(a);
            })
            .catch(error => {
                console.error('Download failed:', error);
                alert('Download failed. Please try again.');
            });
        });
        
    } catch (error) {
        console.error('Upload error:', error);
        if (error.message === 'Authentication required' || error.message === 'Authentication expired') {
            showAuthModal();
        } else {
            alert('There was an error processing your file: ' + error.message);
        }
        document.getElementById('loader').style.display = 'none';
    }
});

// File upload functionality
function setupFileUpload() {
    const fileInput = document.getElementById('fileInput');
    const uploadArea = document.getElementById('uploadArea');
    const fileInfo = document.getElementById('fileInfo');
    const fileName = document.getElementById('fileName');
    const fileSize = document.getElementById('fileSize');
    const uploadText = document.getElementById('uploadText');

    if (!fileInput || !uploadArea) return;

    // File input change handler
    fileInput.addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
            showFileInfo(file);
        }
    });

    // Drag and drop handlers
    uploadArea.addEventListener('click', function() {
        fileInput.click();
    });

    uploadArea.addEventListener('dragover', function(e) {
        e.preventDefault();
        uploadArea.classList.add('dragover');
    });

    uploadArea.addEventListener('dragleave', function(e) {
        e.preventDefault();
        uploadArea.classList.remove('dragover');
    });

    uploadArea.addEventListener('drop', function(e) {
        e.preventDefault();
        uploadArea.classList.remove('dragover');
        
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            const file = files[0];
            // Validate file type
            if (file.type.includes('sheet') || file.name.endsWith('.xlsx') || file.name.endsWith('.xls')) {
                fileInput.files = files;
                showFileInfo(file);
            } else {
                alert('Please select a valid Excel file (.xlsx or .xls)');
            }
        }
    });

    function showFileInfo(file) {
        if (fileName && fileSize && fileInfo && uploadText) {
            fileName.textContent = file.name;
            fileSize.textContent = formatFileSize(file.size);
            fileInfo.style.display = 'block';
            uploadText.textContent = 'File selected successfully';
            uploadArea.classList.add('file-selected');
        }
    }

    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }
}

function clearFile() {
    const fileInput = document.getElementById('fileInput');
    const fileInfo = document.getElementById('fileInfo');
    const uploadText = document.getElementById('uploadText');
    const uploadArea = document.getElementById('uploadArea');

    if (fileInput) fileInput.value = '';
    if (fileInfo) fileInfo.style.display = 'none';
    if (uploadText) uploadText.textContent = 'Drag and drop your Excel file here, or click to browse';
    if (uploadArea) uploadArea.classList.remove('file-selected');
}

// Input mode management
let currentInputMode = 'file'; // 'file' or 'text'

function switchToFileUpload() {
    currentInputMode = 'file';
    document.getElementById('fileUploadBtn').classList.add('active');
    document.getElementById('textInputBtn').classList.remove('active');
    document.getElementById('fileUploadMode').style.display = 'block';
    document.getElementById('textInputMode').style.display = 'none';
    
    // Manage required attributes
    const fileInput = document.getElementById('fileInput');
    const textInput = document.getElementById('linkTextInput');
    if (fileInput) fileInput.required = true;
    if (textInput) textInput.required = false;
    
    // Update submit button state
    updateSubmitButtonState();
}

function switchToTextInput() {
    currentInputMode = 'text';
    document.getElementById('textInputBtn').classList.add('active');
    document.getElementById('fileUploadBtn').classList.remove('active');
    document.getElementById('fileUploadMode').style.display = 'none';
    document.getElementById('textInputMode').style.display = 'block';
    
    // Manage required attributes
    const fileInput = document.getElementById('fileInput');
    const textInput = document.getElementById('linkTextInput');
    if (fileInput) fileInput.required = false;
    if (textInput) textInput.required = true;
    
    // Update submit button state
    updateSubmitButtonState();
}

function updateSubmitButtonState() {
    const submitBtn = document.getElementById('submitBtn');
    const fileInput = document.getElementById('fileInput');
    const textInput = document.getElementById('linkTextInput');
    
    if (currentInputMode === 'file') {
        submitBtn.disabled = !fileInput.files.length;
    } else {
        submitBtn.disabled = !textInput.value.trim();
    }
}

// Add event listener for text input
document.addEventListener('DOMContentLoaded', function() {
    const textInput = document.getElementById('linkTextInput');
    if (textInput) {
        textInput.addEventListener('input', updateSubmitButtonState);
    }
    
    // Initialize required attributes based on default mode
    const fileInput = document.getElementById('fileInput');
    if (fileInput) fileInput.required = true; // Default is file mode
    if (textInput) textInput.required = false;
});

// Initialize authentication UI and event handlers
document.addEventListener('DOMContentLoaded', function() {
    // Update authentication UI
    DashboardAuthManager.updateAuthUI();
    
    // Verify token on page load
    DashboardAuthManager.verifyToken();
    
    // Set up logout handler
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', function(e) {
            e.preventDefault();
            DashboardAuthManager.logout();
        });
    }
    
    // Set up file upload functionality
    setupFileUpload();
    
    // Load user analyses if logged in
    if (DashboardAuthManager.isLoggedIn()) {
        loadUserAnalyses();
    }
});

// Function to load user analyses
async function loadUserAnalyses() {
    try {
        const response = await authenticatedFetch('/api/v1/influencers/analyses?page=1&limit=5');
        if (response.ok) {
            const data = await response.json();
        }
    } catch (error) {
        console.error('Failed to load user analyses:', error);
    }
}

