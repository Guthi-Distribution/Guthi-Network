
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
        nodes.forEach(node => {
            let row = document.createElement("tr");
            row.setAttribute("id", node.id);
            row.innerHTML = `
                <td>${node.id || undefined}</td>
                <td>${node.name || undefined}</td>
                <td>${node.address || undefined}</td>
                <td>${memoryJson[node.name].installed || undefined}</td>
                <td>${memoryJson[node.name].Available || undefined}</td>
                <td>${memoryJson[node.name].Memory_Load || undefined}</td>
            `;

            table.append(row);
        });
    });
});
