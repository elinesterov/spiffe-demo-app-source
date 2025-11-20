const getJwtButton = document.getElementById('getJwtButton');
const getX509Button = document.getElementById('getX509Button');
const getX509TrustBundleButton = document.getElementById('getX509TrustBundleButton');
const getJwtTrustBundleButton = document.getElementById('getJwtTrustBundleButton');
let output = document.getElementById('output');
let parsedCert = document.getElementById('parsed-cert');

function clearOutput() {
  output.textContent = '';
  parsedCert.innerHTML = '';
}

function setLoading(button, isLoading) {
  if (!button) return;

  const buttonText = button.querySelector('.button-text');
  const buttonLoader = button.querySelector('.button-loader');

  if (!buttonText || !buttonLoader) return;

  if (isLoading) {
    buttonText.classList.add('hidden');
    buttonLoader.classList.remove('hidden');
    button.disabled = true;
    button.classList.add('button-loading');
    button.setAttribute('aria-busy', 'true');
  } else {
    buttonText.classList.remove('hidden');
    buttonLoader.classList.add('hidden');
    button.disabled = false;
    button.classList.remove('button-loading');
    button.setAttribute('aria-busy', 'false');
  }
}

getJwtButton.addEventListener('click', async () => {
  try {
    clearOutput();
    setLoading(getJwtButton, true);

    const response = await fetch('/api/getjwtsvid');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    const jsonResponse = await response.json();
    const token = jsonResponse.token;

    // Create container elements for encoded and decoded JWT
    const encodedContainer = document.createElement('div');
    const decodedContainer = document.createElement('div');

    // Display encoded JWT in first column
    const encodedHeader = document.createElement('h3');
    encodedHeader.textContent = 'Encoded';

    const encodedText = document.createElement('div');
    encodedText.classList.add('token');
    encodedText.textContent = token;

    encodedContainer.appendChild(encodedHeader);
    encodedContainer.appendChild(encodedText);

    // Decode JWT and display decoded contents in second column
    const decodedHeader = document.createElement('h3');
    decodedHeader.textContent = 'Decoded';

    const decodedText = document.createElement('pre');
    decodedText.classList.add('token');
    decodedText.textContent = JSON.stringify(JSON.parse(atob(token.split('.')[1])), null, 2);
    decodedContainer.appendChild(decodedHeader);
    decodedContainer.appendChild(decodedText);

    // Add both columns to parsed-cert element
    let parsedCert = document.getElementById('parsed-cert');
    parsedCert.appendChild(encodedContainer);
    parsedCert.appendChild(decodedContainer);

  } catch (error) {
    console.error(error);
    output.textContent = `Error: ${error.message}`;
  } finally {
    setLoading(getJwtButton, false);
  }
});

getX509Button.addEventListener('click', async () => {
  try {
    clearOutput();
    setLoading(getX509Button, true);

    const response = await fetch('/api/getx509svid');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }

    const res = await response.text();

    parsedCert.innerHTML = '';
    const certData = JSON.parse(res).cert;
    console.log("-----BEGIN PRIVATE KEY-----\n"+JSON.parse(res).key+"\n-----END PRIVATE KEY-----\n");
    const parsedContainer = document.getElementById('parsed-cert');
    // we use pv-cert-viewer to display the certificate
    // https://github.com/PeculiarVentures/pv-certificates-viewer/blob/master/packages/webcomponents/README.md
    const certViewer = document.createElement('peculiar-certificate-viewer');

    // Set the certificate property to the certificate data
    certViewer.setAttribute('certificate', certData);
    parsedContainer.appendChild(certViewer);

  } catch (error) {
    console.error(error);
    output.textContent = `Error: ${error.message}`;
  } finally {
    setLoading(getX509Button, false);
  }
});

getX509TrustBundleButton.addEventListener('click', async () => {
  try {
    clearOutput();
    setLoading(getX509TrustBundleButton, true);

    const response = await fetch('/api/getx509trustbundle');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    const trustBundle = await response.text();
    // output.textContent = trustBundle;
    parsedCert.innerHTML = '';
    const bundles = JSON.parse(trustBundle).bundles;
    const parsedContainer = document.getElementById('parsed-cert');
    // we use pv-cert-viewer to display the certificate
    // https://github.com/PeculiarVentures/pv-certificates-viewer/blob/master/packages/webcomponents/README.md
    const certViewer = document.createElement('peculiar-certificates-viewer');
    const certData = [];

    // Iterate over each trust domain in the bundles object
    for (const trustDomain in bundles) {
      if (Object.hasOwnProperty.call(bundles, trustDomain)) {
        // Create a new object for each trust domain
        for (const authority of bundles[trustDomain]) {
          let bundle = {};
          bundle.name = trustDomain;
          bundle.value = authority;
          certData.push(bundle);
        }
      }
    }

    // Set the certificate property to the certificate data
    certViewer.certificates = certData;
    parsedContainer.appendChild(certViewer);

  } catch (error) {
    console.error(error);
    output.textContent = `Error: ${error.message}`;
  } finally {
    setLoading(getX509TrustBundleButton, false);
  }
});

getJwtTrustBundleButton.addEventListener('click', async () => {
  try {
    clearOutput();
    setLoading(getJwtTrustBundleButton, true);

    const response = await fetch('/api/getjwttrustbundle');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    const jsonResponse = await response.json();

    const bundleTable = document.createElement('table');
    bundleTable.style.cssText = 'width: 100%;';

    const headerRow = bundleTable.insertRow();
    const tdHeader = document.createElement('th');
    tdHeader.style.cssText = 'max-width: 30%; white-space: nowrap;';
    tdHeader.textContent = 'Trust Domain';
    const jwksHeader = document.createElement('th');
    jwksHeader.textContent = 'JWKS';
    headerRow.appendChild(tdHeader)
    headerRow.appendChild(jwksHeader)

    // Iterate over each trust domain in the bundles object
    for (const trustDomain in jsonResponse.bundles) {
      const tr = bundleTable.insertRow();
      td = tr.insertCell()
      td.textContent = trustDomain;
      td.style.cssText = 'max-width: 30%; white-space: nowrap';

      const jwks = document.createElement('pre');
      jwks.textContent = JSON.stringify(jsonResponse.bundles[trustDomain], null, 2);
      jwks.classList.add('token');
      jwks.style.cssText = 'word-break: break-all; white-space: break-spaces; overflow-wrap: break-word;';
      tr.insertCell().appendChild(jwks);
    }

    // Add to parsed-cert element
    let parsedCert = document.getElementById('parsed-cert');
    parsedCert.appendChild(bundleTable);
  } catch (error) {
    console.error(error);
    output.textContent = `Error: ${error.message}`;
  } finally {
    setLoading(getJwtTrustBundleButton, false);
  }
});