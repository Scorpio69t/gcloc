async function analyze() {
    const resultTable = document.getElementById('resultTable');
    const treePre = document.getElementById('tree');
    resultTable.innerHTML = '<tr><td class="text-cyan-300">Loading...</td></tr>';
    treePre.textContent = '';

    // Upload zip file
    let uploadId = '';
    const zipInput = document.getElementById('zipFile');
    if (zipInput.files.length > 0) {
        const form = new FormData();
        form.append('file', zipInput.files[0]);
        const resp = await fetch('/upload', {method: 'POST', body: form});
        const data = await resp.json();
        uploadId = data.id;
    }

    const paths = document.getElementById('paths').value
        .split(/\n+/).map(p => p.trim()).filter(Boolean);

    const body = {
        paths,
        uploadId,
        byFile: document.getElementById('byFile').checked,
        sort: document.getElementById('sort').value,
        excludeExt: document.getElementById('excludeExt').value,
        excludeLang: document.getElementById('excludeLang').value,
        includeLang: document.getElementById('includeLang').value,
        match: document.getElementById('match').value,
        notMatch: document.getElementById('notMatch').value,
        matchDir: document.getElementById('matchDir').value,
        notMatchDir: document.getElementById('notMatchDir').value,
        debug: document.getElementById('debug').checked,
        skipDuplicated: document.getElementById('skipDuplicated').checked
    };

    const res = await fetch('/analyze', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(body)
    });

    const json = await res.json();

    // Render table
    const headers = Object.keys(json[0] || {});
    resultTable.innerHTML = `
    <thead><tr>${headers.map(h => `<th>${h}</th>`).join('')}</tr></thead>
    <tbody>
      ${json.map(row => `
        <tr>${headers.map(h => `<td>${row[h]}</td>`).join('')}</tr>
      `).join('')}
    </tbody>
  `;

    // Chart data
    const labels = json.map(item => item.name || item.lang || item.file || 'N/A');
    const values = json.map(item => item.codes || 0);

    new Chart(document.getElementById('barChart'), {
        type: 'bar',
        data: {
            labels,
            datasets: [{
                label: 'Code Lines',
                data: values,
                backgroundColor: 'rgba(6,182,212,0.6)',
                borderColor: 'rgba(6,182,212,1)',
                borderWidth: 1
            }]
        },
        options: {responsive: true, scales: {y: {beginAtZero: true}}}
    });

    new Chart(document.getElementById('pieChart'), {
        type: 'pie',
        data: {
            labels,
            datasets: [{
                label: 'Code Distribution',
                data: values,
                backgroundColor: [
                    '#0ea5e9', '#06b6d4', '#22d3ee', '#67e8f9', '#a5f3fc', '#cffafe',
                ]
            }]
        }
    });

    if (uploadId) {
        const treeRes = await fetch('/tree?id=' + encodeURIComponent(uploadId));
        const tree = await treeRes.json();
        treePre.textContent = JSON.stringify(tree, null, 2);
    }
}

document.getElementById('analyzeBtn').addEventListener('click', analyze);