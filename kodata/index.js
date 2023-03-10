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
    output.value = jwt;
  } catch (error) {
    console.error(error);
    output.value = `Error: ${error.message}`;
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
    output.value = x509;
  } catch (error) {
    console.error(error);
    output.value = `Error: ${error.message}`;
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
    output.value = trustBundle;
  } catch (error) {
    console.error(error);
    output.value = `Error: ${error.message}`;
  }
});