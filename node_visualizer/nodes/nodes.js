
const getNodes = () => {
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
        }).catch(error => {
            console.log(error);

        });
}


// POST request to http://localhost:8080/connect

const connectForm = document.getElementById("connect-form");

connectForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const formData = new FormData(connectForm);
    const data = {
        "ip": formData.get("ip"),
        "port": formData.get("port")
    }

    fetch('http://localhost:8080/connect', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    }).then(response => {
        if (!response.ok) {
            console.log(`HTTP error: ${response.status}`);
            throw new Error(`HTTP error: ${response.status}`);
        }
        getNodes();
    }).catch(error => {
        console.log(error);
    })

});
