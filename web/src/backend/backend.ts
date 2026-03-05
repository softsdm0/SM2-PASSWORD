const serverUrl = ''

function getAccount() {
  return fetch(`${serverUrl}/api/userinfo`, {
    method: 'GET',
    credentials: 'include',
  }).then((res) => {
    console.log(res)
    return res.json()
  })
}

function logOut() {
  return fetch(`${serverUrl}/api/signout`, {
    method: 'POST',
    credentials: 'include',
  }).then((res) => {
    return res.json()
  })
}

export default {
  getAccount,
  logOut,
}
