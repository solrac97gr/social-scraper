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

document.getElementById('uploadForm').addEventListener('submit', function(event) {
    event.preventDefault();
    document.getElementById('loader').style.display = 'block';
    document.getElementById('uploadForm').classList.add('animate__animated', 'animate__bounceOut');

    var formData = new FormData(this);

    // Call the new endpoint to get the estimated time
    fetch('/estimate-time', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        updateLoaderText(data.estimatedTime);
        startCountdown(data.estimatedTime);

        // Proceed with the file upload
        fetch('/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            document.getElementById('loader').style.display = 'none';
            document.getElementById('resultsContainer').style.display = 'block';

            const resultsTableBody = document.getElementById('resultsTableBody');
            resultsTableBody.innerHTML = '';

            // remove the first row (header)
            data.results.shift();

            data.results.forEach(result => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${result[0]}</td>
                    <td class="count-up" data-count="${result[1]}">${result[1]}</td>
                    <td><a href="${result[2]}" target="_blank">${result[2]}</a></td>
                    <td><span class="badge badge-${result[3].toLowerCase()}">${result[3]}</span></td>
                    <td>${result[4]}</td>
					<td>${result[5]}</td>
					<td>${result[6]}</td>
                `;
                resultsTableBody.appendChild(row);
            });

            document.getElementById('downloadButton').addEventListener('click', function() {
                window.location.href = `/download?filename=${data.outputFile}`;
            });
        })
        .catch(() => {
            alert('There was an error processing your file.');
            document.getElementById('loader').style.display = 'none';
        });
    })
    .catch(() => {
        alert('There was an error estimating the processing time.');
        document.getElementById('loader').style.display = 'none';
    });
});

