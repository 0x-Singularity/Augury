document.addEventListener("DOMContentLoaded", function () {
    const searchBtn = document.getElementById("search-btn");
    const iocInput = document.getElementById("ioc-input");
    const tabButtons = document.getElementById("tab-buttons");
    const tabContent = document.getElementById("tab-content");

    searchBtn.addEventListener("click", function () {
        const ioc = iocInput.value.trim();
        if (!ioc) {
            alert("Please enter an IOC.");
            return;
        }

        // API Request URL
        const apiUrl = `http://localhost:8080/api/ioc/fakeula?ioc=${encodeURIComponent(ioc)}`;

        fetch(apiUrl)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                createTab(ioc, data);
            })
            .catch(error => {
                console.error("Error fetching Fakeula data:", error);
                alert("Failed to fetch data from Fakeula.");
            });
    });

    function createTab(ioc, data) {
        console.log(`Creating tab for: ${ioc}`);
    
        // Generate a unique tab ID
        const tabId = `tab-${ioc.replace(/[^a-zA-Z0-9]/g, "_")}`;
    
        // Check if the tab already exists
        let existingTab = document.getElementById(tabId);
        if (existingTab) {
            console.log(`Tab for ${ioc} already exists, switching to it.`);
            setActiveTab(tabId);
            return;
        }
    
        console.log(`Creating new tab: ${tabId}`);
    
        // Create new tab button
        const tabButton = document.createElement("button");
        tabButton.innerText = ioc;
        tabButton.classList.add("tab-button");
        tabButton.onclick = function () {
            console.log(`Switching to tab: ${tabId}`);
            setActiveTab(tabId);
        };
    
        // Create new tab content
        const tabPanel = document.createElement("div");
        tabPanel.id = tabId;
        tabPanel.classList.add("tab-panel");
        tabPanel.innerHTML = `
            <h3>Results for: ${ioc}</h3>
            <pre>${JSON.stringify(data, null, 2)}</pre>
        `;
    
        // Append to UI
        document.getElementById("tab-buttons").appendChild(tabButton);
        document.getElementById("tab-content").appendChild(tabPanel);
    
        console.log(`Tab ${tabId} created successfully!`);
    
        // Activate the new tab
        setActiveTab(tabId);
    }

    function setActiveTab(tabId) {
        document.querySelectorAll(".tab-panel").forEach(panel => {
            panel.style.display = panel.id === tabId ? "block" : "none";
        });
        document.querySelectorAll(".tab-button").forEach(btn => {
            btn.classList.toggle("active", btn.innerText === tabId.replace("tab-", ""));
        });
    }
});