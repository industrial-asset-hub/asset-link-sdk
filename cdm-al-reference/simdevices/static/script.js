//TODO: refactor the code below and reuse the code for the devices also for sub-devices

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

    // Store open devices and subdevices
    const openDevices = new Set();
    const openSubDevices = new Set();
    document.querySelectorAll('.device.active').forEach(el => {
        const name = el.querySelector('.device-header .header-left span:last-child')?.textContent;
        if (name) openDevices.add(name);
    });
    document.querySelectorAll('.subdevice.open').forEach(el => {
        const name = el.querySelector('.subdevice-header .header-left span:last-child')?.textContent;
        if (name) openSubDevices.add(name);
    });
    deviceList.innerHTML = '';

    devices.forEach((device, deviceIdx) => {
        const deviceElement = document.createElement('div');
        deviceElement.className = `device`;

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

        // Sub-devices (collapsible)
        if (device.sub_devices && device.sub_devices.length > 0) {
            device.sub_devices.forEach((subDevice, subIdx) => {
                const subDeviceElement = document.createElement('div');
                subDeviceElement.className = 'subdevice';

                // Sub-device Header
                const subDeviceHeader = document.createElement('div');
                subDeviceHeader.className = 'subdevice-header';
                subDeviceHeader.innerHTML = `
                    <span class=\"header-left\">
                        <span class=\"device-icon\" aria-hidden=\"true\" style=\"display:inline-flex;align-items:center;margin-right:12px;\">
                            <img src=\"images/subdevice.svg\" width=\"26\" height=\"26\" alt=\"Subdevice Icon\" />
                        </span>
                        <span style=\"white-space:nowrap;overflow:hidden;text-overflow:ellipsis;\">${subDevice.device_name}</span>
                    </span>
                    <span class=\"state-tag ${subDevice.device_state.toLowerCase()}\">${subDevice.device_state.toLowerCase()}</span>
                    <button class=\"toggle-btn\" aria-label=\"Toggle subdevice details\"><img class=\"plusminus-icon\" src=\"images/plus.svg\" width=\"24\" height=\"24\" alt=\"Expand\" /></button>
                `;
                subDeviceElement.appendChild(subDeviceHeader);

                // Sub-device Content
                const subDeviceContent = document.createElement('div');
                subDeviceContent.className = 'subdevice-content';
                subDeviceContent.innerHTML = `
                    <dl>
                        <dt>Device Name:</dt>
                        <dd>${subDevice.device_name}</dd>
                        <dt>Product Designation:</dt>
                        <dd>${subDevice.product_designation}</dd>
                        <dt>Article Number:</dt>
                        <dd>${subDevice.article_number}</dd>
                        <dt>Manufacturer:</dt>
                        <dd>${subDevice.manufacturer}</dd>
                        <dt>Serial Number:</dt>
                        <dd>${subDevice.serial_number}</dd>
                        <dt>Firmware Version:</dt>
                        <dd>${subDevice.firmware_version}</dd>
                        <dt>Hardware Version:</dt>
                        <dd>${subDevice.hardware_version}</dd>
                        <dt>Device State:</dt>
                        <dd>${subDevice.device_state.toLowerCase()}</dd>
                        <dt>IP Address:</dt>
                        <dd>${subDevice.ip_device != "" ? subDevice.ip_device : "-"}</dd>
                        <dt>MAC Address:</dt>
                        <dd>${subDevice.mac_address != "" ? subDevice.mac_address : "-"}</dd>
                    </dl>
                `;
                subDeviceElement.appendChild(subDeviceContent);

                // Sub-device collapse logic
                subDeviceHeader.addEventListener('click', function(e) {
                    e.stopPropagation();
                    const isOpen = subDeviceElement.classList.toggle('open');
                    const btnImg = subDeviceHeader.querySelector('.toggle-btn img');
                    btnImg.src = isOpen ? 'images/minus.svg' : 'images/plus.svg';
                    btnImg.alt = isOpen ? 'Collapse' : 'Expand';
                });

                // Open sub-device if it was previously open
                if (openSubDevices.has(subDevice.device_name)) {
                    subDeviceElement.classList.add('open');
                    const btnImg = subDeviceHeader.querySelector('.toggle-btn img');
                    btnImg.src = 'images/minus.svg';
                    btnImg.alt = 'Collapse';
                }

                deviceContent.appendChild(subDeviceElement);
            });
        }

        deviceElement.appendChild(deviceContent);

        // Device collapse logic
        deviceHeader.addEventListener('click', function() {
            const isActive = deviceElement.classList.toggle('active');
            const btnImg = deviceHeader.querySelector('.toggle-btn img');
            btnImg.src = isActive ? 'images/minus.svg' : 'images/plus.svg';
            btnImg.alt = isActive ? 'Collapse' : 'Expand';
        });

        // Open device if it was previously open
        if (openDevices.has(device.device_name)) {
            deviceElement.classList.add('active');
            const btnImg = deviceHeader.querySelector('.toggle-btn img');
            btnImg.src = 'images/minus.svg';
            btnImg.alt = 'Collapse';
        }

        deviceList.appendChild(deviceElement);
    });
}

// Connect when page loads
connectWebSocket();
