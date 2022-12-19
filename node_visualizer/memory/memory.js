const DATA_COUNT = 10;
const memoryChart = createChart('myChart', DATA_COUNT);
let table;
let memoryFetchPromise;
let nodeFetchPromise;

const REFRESH_INTERVAL = 2000;

let loop = setInterval(() => {
    memoryFetchPromise = fetch('http://localhost:8080/memory');
    memoryFetchPromise.then(response => {
        if (!response.ok) {
            console.log(`HTTP error: ${response.status}`);
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();

    }).then(memoryJson => {

        nodeFetchPromise = fetch('http://localhost:8080/nodes');
        nodeFetchPromise.then(response => {
            if (!response.ok) {
                console.log(`HTTP error: ${response.status}`);
                throw new Error(`HTTP error: ${response.status}`);
            }
            return response.json();

        }).then(nodesJson => {
            if (!nodesJson || !memoryJson) {
                return;
            }

            if (!table) {
                table = document.getElementById("page1Table");
                nodesJson.forEach((node, index) => {
                    let borderColor = `rgba(${Math.floor(Math.random() * 255)}, ${Math.floor(Math.random() * 255)}, ${Math.floor(Math.random() * 255)}, 1)`;

                    addGraph(memoryChart, node.name, borderColor);

                    let row = document.createElement("tr");
                    row.setAttribute("id", node.id + " " + node.name);

                    row.innerHTML = `
                    <td>${node.id || undefined}</td>
                    <td>${node.name || undefined}</td>
                    <td>${node.address.IP || undefined}</td>
                    <td>${memoryJson[node.name].installed || undefined}</td>
                    <td>${memoryJson[node.name].Available || undefined}</td>
                    <td>${memoryJson[node.name].Memory_Load || undefined}</td>
                `;

                    table.append(row);

                });
            } else {
                nodesJson.forEach((node, index) => {
                    let row = document.getElementById(node.id + " " + node.name);
                    row.innerHTML = `
                    <td>${node.id || undefined}</td>
                    <td>${node.name || undefined}</td>
                    <td>${node.address.IP || undefined}</td>
                    <td>${memoryJson[node.name].installed || undefined}</td>
                    <td>${memoryJson[node.name].Available || undefined}</td>
                    <td>${memoryJson[node.name].Memory_Load || undefined}</td>
                `;

                    memoryChart.data.datasets[index].data.push(memoryJson[node.name].Memory_Load);
                    memoryChart.data.datasets[index].data.shift();
                    memoryChart.update("none");
                });
            }
        }).catch(err => {
            console.log(err);
        });
    });
}, REFRESH_INTERVAL);
