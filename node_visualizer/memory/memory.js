const DATA_COUNT = 10;
const memoryChart = createChart('Memory Load', 'memoryChart', DATA_COUNT);
const CPUChart = createChart('CPU load', 'CPUChart', DATA_COUNT);
const COLOR_SET = ["#ea5545", "#f46a9b", "#ef9b20", "#edbf33", "#ede15b", "#bdcf32", "#87bc45", "#27aeef", "#b33dc6"];
let table;

const REFRESH_INTERVAL = 2000;

const updateData = async () => {

    try {
        let response = await fetch('http://localhost:8080/memory');
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }
        const memoryJson = await response.json();

        response = await fetch('http://localhost:8080/nodes');
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }
        const nodesJson = await response.json();

        response = await fetch('http://localhost:8080/cpu');
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }
        const cpuJson = await response.json();

        if (!nodesJson) throw new Error(`Nodes data not found.`);
        if (!memoryJson) throw new Error(`Memory data not found`);
        if (!cpuJson) throw new Error(`CPU data not found`);

        if (!table) {
            table = document.getElementById("page1Table");
            nodesJson.forEach((node, index) => {
                node.borderColor = COLOR_SET[Math.floor(Math.random() * COLOR_SET.length)];

                addGraph(memoryChart, node.name, node.borderColor);
                addGraph(CPUChart, node.name, node.borderColor);

                let row = document.createElement("tr");
                row.setAttribute("id", node.id + " " + node.name);

                row.innerHTML = `
                    <td>${node.id || undefined}</td>
                    <td>${node.name || undefined}</td>
                    <td>${node.address.IP || undefined}</td>
                    <td>${memoryJson[node.name].installed || undefined}</td>
                    <td>${memoryJson[node.name].Available || undefined}</td>
                    <td>${memoryJson[node.name].Memory_Load || undefined}</td>
                    <td>${cpuJson[node.name].usage || undefined}</td>
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
                    <td>${cpuJson[node.name].usage || undefined}</td>
                `;

                memoryChart.data.datasets[index].data.push(memoryJson[node.name].Memory_Load);
                memoryChart.data.datasets[index].data.shift();
                memoryChart.update("none");

                CPUChart.data.datasets[index].data.push(cpuJson[node.name].usage);
                CPUChart.data.datasets[index].data.shift();
                CPUChart.update("none");

            });
        }

    } catch (error) {
        console.log(error);
    }
}

updateData();

let loop = setInterval(updateData, REFRESH_INTERVAL);

