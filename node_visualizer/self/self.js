
const fetchPromise = fetch('http://localhost:8080/self');

fetchPromise
    .then(response => {
        if (!response.ok) {
            console.log(`HTTP error: ${response.status}`);
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();
    })
    .then(json => {
        let table = document.getElementById("table");

        let tBody = document.createElement("tbody");
        tBody.innerHTML = `
            <tr>
                <th>Attribute</th>
                <th>Value</th>
            </tr>
            <tr>
                <td>ID</td>
                <td>${json.id || undefined}</td>
            </tr>
            <tr>
                <td>Name</td>            
                <td>${json.name || undefined}</td>
            </tr>
            <tr>
                <td>IP Address</td>
                <td>${json.address.IP || undefined}</td>
            </tr>
            <tr>
                <td>Port</td>
                <td>${json.address.Port || undefined}</td>
            </tr>
            <tr>
                <td>Zone</td>
                <td>${json.address.Zone || undefined}</td>
            </tr>
            `;

        table.append(tBody);
    });
