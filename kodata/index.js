const getJwtButton = document.getElementById('getJwtButton');
const getX509Button = document.getElementById('getX509Button');
const getTrustBundleButton = document.getElementById('getTrustBundleButton');
const output = document.getElementById('output');
const parsedCert = document.getElementById('parsed-cert');

function clearParsedCert() {
  parsedCert.innerHTML = '';
}

getJwtButton.addEventListener('click', async () => {
  try {
    clearParsedCert();
    const response = await fetch('/api/getjwtsvid');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    const jwt = await response.text();
    output.innerHTML = `<div>${jwt}</div>`;
  } catch (error) {
    console.error(error);
    output.innerHTML = `<div>Error: ${error.message}</div>`;
  }
});

getX509Button.addEventListener('click', async () => {
  try {
    clearParsedCert();
    const response = await fetch('/api/getx509svid');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    const x509 = await response.text();
    output.innerHTML = `<div>${x509}</div>`;
  } catch (error) {
    console.error(error);
    output.innerHTML = `<div>Error: ${error.message}</div>`;
  }
});

getTrustBundleButton.addEventListener('click', async () => {
  try {
    clearParsedCert();
    const response = await fetch('/api/gettrustbundle');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    const trustBundle = await response.text();
    output.textContent = trustBundle;
    parsedCert.innerHTML = '';
    const bundles = JSON.parse(trustBundle).bundles;
    const parsedContainer = document.getElementById('parsed-cert');
    // we use pv-cert-viewer to display the certificate
    // https://github.com/PeculiarVentures/pv-certificates-viewer/blob/master/packages/webcomponents/README.md
    const certViewer = document.createElement('peculiar-certificate-viewer');
    const certData = [];

    // Iterate over each trust domain in the bundles object
    for (const trustDomain in bundles) {
      if (Object.hasOwnProperty.call(bundles, trustDomain)) {
        certData.push(bundles[trustDomain]);
        console.log(`Trust Domain: ${trustDomain}`);
        console.log(`Certificate Data: ${certData}`);
      }
    }

    // Set the certificate property to the certificate data
    certViewer.setAttribute('certificate', certData);
    parsedContainer.appendChild(certViewer);

  } catch (error) {
    console.error(error);
    output.textContent = `Error: ${error.message}`;
  }
});
