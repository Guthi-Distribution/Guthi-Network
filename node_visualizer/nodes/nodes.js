
const fetchPromise = fetch('http://localhost:8080/nodes');

fetchPromise
    .then(response => {
        if (!response.ok) {
            console.log(`HTTP error: ${response.status}`);
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();
    })
    .then(json => {
        let table = document.getElementById("page1Table");
        json.forEach(e => {
            let row = document.createElement("tr");
            row.setAttribute("id", e.id);
            // row.setAttribute("class", "tableRow");
            row.innerHTML = `
                <td>${e.id || undefined}</td>
                <td>${e.name || undefined}</td>
                <td>${e.address.IP || undefined}</td>
                <td>${e.address.Port || undefined}</td>
                <td>${e.address.Zone || undefined}</td>
            `;
            table.append(row);
        });
    });
