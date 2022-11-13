// const fetchPromise = fetch('http://localhost:8080/');

// fetchPromise
//     .then(response => {
//         if (!response.ok) {
//             console.log(`HTTP error: ${response.status}`);
//             throw new Error(`HTTP error: ${response.status}`);
//         }
//         return response.json();
//     })
//     .then(json => {
//         document.getElementById("id").innerHTML = json.id
//         document.getElementById("IP").innerHTML = json.address.IP
//         document.getElementById("port").innerHTML = json.address.Port
//         document.getElementById("zone").innerHTML = json.address.Zone
//     });