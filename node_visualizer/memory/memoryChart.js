const createChart = (container, DATA_COUNT) => {
    const ctx = document.getElementById(container);
    const Utils = Chart.helpers;

    let datapoints = [];

    const labels = [];
    for (let i = 0; i < DATA_COUNT; ++i) {
        labels.push(i.toString());
    }

    const data = {
        labels: labels,

    };

    let memoryChart = new Chart(ctx, {
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
                    suggestedMin: 0,
                    suggestedMax: 100
                }
            }
        },
    });

    return memoryChart;
}

const addGraph = (chart, label, borderColor) => {
    let newData = {
        label: label,
        data: new Array(DATA_COUNT).fill(0),
        borderColor: borderColor,
        fill: true,
        borderWidth: 1,
        pointHoverRadius: 4,
        pointHoverBorderWidth: 1,
        pointHoverBackgroundColor: 'rgba(255, 255, 255, 1)',
        pointHoverBorderColor: 'rgba(0, 0, 0, 1)',
        cubicInterpolationMode: 'monotone',
        tension: 0.4
    }

    chart.data.datasets.push(newData);
    chart.update();
}


