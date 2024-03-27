document.getElementById('shortenBtn').addEventListener('click', function () {
    var urlToShorten = document.getElementById('urlToShorten').value;
    fetch('/shorten', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: urlToShorten }),
    })
        .then(response => response.json())
        .then(data => {
            document.getElementById('shortUrl').textContent = 'https://myshl.ru/' + data.shortUrl;
        })
        .catch((error) => {
            console.error('Error:', error);
        });
});