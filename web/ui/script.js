async function analyze() {
  const resultPre = document.getElementById('result');
  const treePre = document.getElementById('tree');
  resultPre.textContent = 'Running...';
  treePre.textContent = '';

  let uploadId = '';
  const zipInput = document.getElementById('zipFile');
  if (zipInput.files.length > 0) {
    const form = new FormData();
    form.append('file', zipInput.files[0]);
    const resp = await fetch('/upload', { method: 'POST', body: form });
    const data = await resp.json();
    uploadId = data.id;
  }

  const paths = document.getElementById('paths').value
    .split(/\n+/)
    .map(p => p.trim())
    .filter(p => p);

  const body = {
    paths: paths,
    uploadId: uploadId,
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
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  });
  const json = await res.json();
  resultPre.textContent = JSON.stringify(json, null, 2);

  if (uploadId) {
    const tRes = await fetch('/tree?id=' + encodeURIComponent(uploadId));
    const tree = await tRes.json();
    treePre.textContent = JSON.stringify(tree, null, 2);
  }
}

document.getElementById('analyzeBtn').addEventListener('click', analyze);
