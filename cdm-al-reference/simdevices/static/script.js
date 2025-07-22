let ws;

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
    const host = window.location.host;
    const wsURL = `${protocol}${host}/ws`;

    ws = new WebSocket(wsURL);

    ws.onopen = function() {
        console.log('Connected to WebSocket');
    };

    ws.onmessage = function(event) {
        const devices = JSON.parse(event.data);
        updateDeviceList(devices);
    };

    ws.onclose = function() {
        console.log('WebSocket connection closed');
        // Attempt to reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
    };
}

function addDevice(device, openSet) {
    const element = document.createElement('div');
    element.className = 'device';
    element.id = device.unique_device_id;

    // Header
    const header = document.createElement('div');
    header.className = 'header';
    header.innerHTML = `
        <span class="header-left">
            <span class="device-icon" aria-hidden="true">
                <img src="images/device.svg" alt="Device Icon" />
            </span>
            <span>${device.device_name}</span>
        </span>
        <span class="state-tag ${device.device_state.toLowerCase()}">${device.device_state.toLowerCase()}</span>
        <button class="toggle-btn" aria-label="Toggle device details">
            <img class="plusminus-icon" src="images/plus.svg" width="24" height="24" alt="Expand" />
        </button>
    `;
    element.appendChild(header);

    // Content
    const content = document.createElement('div');
    content.className = 'content';
    content.innerHTML = `
        <dl>
            <dt>Device Name:</dt>
            <dd>${device.device_name}</dd>
            <dt>Product Designation:</dt>
            <dd>${device.product_designation}</dd>
            <dt>Article Number:</dt>
            <dd>${device.article_number}</dd>
            <dt>Manufacturer:</dt>
            <dd>${device.manufacturer}</dd>
            <dt>Serial Number:</dt>
            <dd>${device.serial_number}</dd>
            <dt>Firmware Version:</dt>
            <dd>${device.firmware_version}</dd>
            <dt>Hardware Version:</dt>
            <dd>${device.hardware_version}</dd>
            <dt>Device State:</dt>
            <dd>${device.device_state.toLowerCase()}</dd>
            <dt>IP Address:</dt>
            <dd>${device.ip_device !== undefined && device.ip_device !== '' ? device.ip_device : '-'}</dd>
            <dt>MAC Address:</dt>
            <dd>${device.mac_address !== undefined && device.mac_address !== '' ? device.mac_address : '-'}</dd>
        </dl>
    `;
    element.appendChild(content);

    // Collapse logic
    header.addEventListener('click', function(e) {
        const isOpen = element.classList.toggle('open');
        const btnImg = header.querySelector('.toggle-btn img');
        btnImg.src = isOpen ? 'images/minus.svg' : 'images/plus.svg';
        btnImg.alt = isOpen ? 'Collapse' : 'Expand';
    });

    // Open if previously open
    if (openSet.has(device.unique_device_id)) {
        element.classList.add('open');
        const btnImg = header.querySelector('.toggle-btn img');
        btnImg.src = 'images/minus.svg';
        btnImg.alt = 'Collapse';
    }

    return element;
}

function updateDeviceList(devices) {
    const deviceList = document.getElementById('deviceList');

    // Store open devices
    const openDevices = new Set();
    document.querySelectorAll('.device.open').forEach(el => {
        openDevices.add(el.id);
    });
    deviceList.innerHTML = '';

    devices.forEach((device) => {
        const deviceElement = addDevice(device, openDevices);
        deviceList.appendChild(deviceElement);
    });
}

// Connect when page loads
connectWebSocket();
