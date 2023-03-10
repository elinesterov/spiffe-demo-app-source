const getJwtButton = document.getElementById('getJwtButton');
const getX509Button = document.getElementById('getX509Button');
const getTrustBundleButton = document.getElementById('getTrustBundleButton');
const output = document.getElementById('output');

getJwtButton.addEventListener('click', async () => {
  try {
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
    const response = await fetch('/api/gettrustbundle');
    if (!response.ok) {
      const error = await response.text();
      throw new Error(error);
    }
    const trustBundle = await response.text();
    output.innerHTML = `<div>${trustBundle}</div>`;
  } catch (error) {
    console.error(error);
    output.innerHTML = `<div>Error: ${error.message}</div>`;
  }
});
