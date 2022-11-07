const fetchPromise = fetch('http://localhost:8080/');

fetchPromise
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }
        return response.json();
    })
    .then(json => {
        document.getElementById("attr1").innerHTML = json.attr1
        document.getElementById("attr2").innerHTML = json.attr2
        document.getElementById("attr3").innerHTML = json.attr3
    });