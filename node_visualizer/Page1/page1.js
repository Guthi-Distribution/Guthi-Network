
const fetchPromise = fetch('http://localhost:8080/');

fetchPromise
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();
    })
    .then(json => {
        let data = json;

        let tblBody = document.getElementById("page1Table");
        data.forEach(e => {
            let row = document.createElement("tr");
            row.setAttribute("id", e.id);
            // row.setAttribute("class", "tableRow");
            row.innerHTML = `
                <td>${e.attrb1}</td>
                <td>${e.attrb2}</td>
                <td>${e.attrb3}</td>
                <td>${e.attrb4}</td>
                <td>${e.attrb5}</td>
            `;
            tblBody.append(row);
        });
    });
