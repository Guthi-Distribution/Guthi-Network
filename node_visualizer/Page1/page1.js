
const fetchPromise = fetch('http://localhost:8080/');

fetchPromise
    .then(response => {
        if (!response.ok) {
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
                <td>${e.id}</td>
                <td>${e.address.IP}</td>
                <td>${e.address.Port}</td>
                <td>${e.address.Zone}</td>
            `;
            table.append(row);
        });
    });
