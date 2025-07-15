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

function updateDeviceList(devices) {
    const deviceList = document.getElementById('deviceList');

    // Store open devices
    const openDevices = new Set();
    document.querySelectorAll('.device.open').forEach(el => {
        openDevices.add(el.id);
    });
    deviceList.innerHTML = '';

    devices.forEach((device, deviceIdx) => {
        const deviceElement = document.createElement('div');
        deviceElement.className = `device`;
        deviceElement.id = device.unique_device_id;

        // Device Header (collapsible)
        const deviceHeader = document.createElement('div');
        deviceHeader.className = 'device-header';
        deviceHeader.innerHTML = `
            <span class="header-left">
                <span class="device-icon" aria-hidden="true" style="display:inline-flex;align-items:center;margin-right:16px;">
                    <img src="images/device.svg" width="32" height="32" alt="Device Icon" />
                </span>
                <span style="white-space:nowrap;overflow:hidden;text-overflow:ellipsis;">${device.device_name}</span>
            </span>
            <span class="state-tag ${device.device_state.toLowerCase()}">${device.device_state.toLowerCase()}</span>
            <button class="toggle-btn" aria-label="Toggle device details">
                <img class="plusminus-icon" src="images/plus.svg" width="24" height="24" alt="Expand" />
            </button>
        `;
        deviceElement.appendChild(deviceHeader);

        // Device Content
        const deviceContent = document.createElement('div');
        deviceContent.className = 'device-content';
        deviceContent.innerHTML = `
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
                <dd>${device.ip_device}</dd>
                <dt>MAC Address:</dt>
                <dd>${device.mac_address}</dd>
            </dl>
        `;

        deviceElement.appendChild(deviceContent);

        // Device collapse logic
        deviceHeader.addEventListener('click', function() {
            const isOpen = deviceElement.classList.toggle('open');
            const btnImg = deviceHeader.querySelector('.toggle-btn img');
            btnImg.src = isOpen ? 'images/minus.svg' : 'images/plus.svg';
            btnImg.alt = isOpen ? 'Collapse' : 'Expand';
        });

        // Open device if it was previously open
        if (openDevices.has(device.unique_device_id)) {
            deviceElement.classList.add('open');
            const btnImg = deviceHeader.querySelector('.toggle-btn img');
            btnImg.src = 'images/minus.svg';
            btnImg.alt = 'Collapse';
        }

        deviceList.appendChild(deviceElement);
    });
}

// Connect when page loads
connectWebSocket();
