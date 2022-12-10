
const ctx = document.getElementById('myChart');
const Utils = Chart.helpers;

const DATA_COUNT = 12;
const labels = [];
for (let i = 0; i < DATA_COUNT; ++i) {
    labels.push(i.toString());
}
const datapoints = [0, 20, 20, 60, 60, 120, 20, 180, 120, 125, 105, 110, 170];
const data = {
    labels: labels,
    datasets: [
        {
            label: '1',
            data: datapoints,
            borderColor: "red",
            fill: false,
            cubicInterpolationMode: 'monotone',
            tension: 0.4
        }
    ]
};

new Chart(ctx, {
    type: 'line',
    data: data,
    options: {
        responsive: true,
        plugins: {
            title: {
                display: true,
                text: 'Memory_Load'
            },
        },
        interaction: {
            intersect: false,
        },
        scales: {
            x: {
                display: true,
                title: {
                    display: true
                }
            },
            y: {
                display: true,
                title: {
                    display: true,
                    text: 'Value'
                },
                suggestedMin: -10,
                suggestedMax: 200
            }
        }
    },
});



const memoryFetchPromise = fetch('http://localhost:8080/memory');

memoryFetchPromise.then(response => {
    if (!response.ok) {
        console.log(`HTTP error: ${response.status}`);
        throw new Error(`HTTP error: ${response.status}`);
    }
    return response.json();

}).then(memoryJson => {
    const nodeFetchPromise = fetch('http://localhost:8080/nodes');
    let nodes = [];

    nodeFetchPromise.then(response => {
        if (!response.ok) {
            console.log(`HTTP error: ${response.status}`);
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();

    }).then(nodesJson => {
        nodesJson.forEach(node => {
            nodes.push({
                id: node.id,
                name: node.name,
                address: node.address.IP
            });
        })

        let table = document.getElementById("page1Table");
        nodes.forEach((node, index) => {
            let row = document.createElement("tr");
            row.setAttribute("id", node.id);

            row.innerHTML = `
                <td>${node.id || undefined}</td>
                <td>${node.name || undefined}</td>
                <td>${node.address || undefined}</td>
                <td>${memoryJson[node.name].installed || undefined}</td>
                <td>${memoryJson[node.name].Available || undefined}
                    <div class="progress-container">
                        <div class="progress-${index}"></div>
                    </div>
                </td>
                <td>${memoryJson[node.name].Memory_Load || undefined}</td>
            `;

            table.append(row);
        });

        nodes.forEach((node, index) => {
            let progressbar = document.querySelector(`.progress-${index}`);
            let progress = memoryJson[node.name].Memory_Load * 10;


            // create a variable to track 
            setTimeout(() => {
                progressbar.style.width = `${progress}%`;
            }, 1000);
        });

    });
});
